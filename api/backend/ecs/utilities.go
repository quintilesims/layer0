package ecsbackend

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsecs "github.com/aws/aws-sdk-go/service/ecs"
	"gitlab.imshealth.com/xfra/layer0/api/backend"
	"gitlab.imshealth.com/xfra/layer0/api/backend/ecs/id"
	"gitlab.imshealth.com/xfra/layer0/common/aws/cloudwatchlogs"
	"gitlab.imshealth.com/xfra/layer0/common/aws/ecs"
	"gitlab.imshealth.com/xfra/layer0/common/config"
	"gitlab.imshealth.com/xfra/layer0/common/models"
	"strings"
)

const MAX_TASK_IDS = 100

func stringp(s string) *string {
	return &s
}

func int64p(i int64) *int64 {
	return &i
}

func boolp(b bool) *bool {
	return &b
}

func pstring(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func int64OrZero(i *int64) int64 {
	if i == nil {
		return 0
	}

	return *i
}

func ContainsErrCode(err error, code string) bool {
	if err == nil {
		return false
	}

	awsErr, ok := err.(awserr.Error)
	if !ok {
		return false
	}

	return strings.ToLower(awsErr.Code()) == strings.ToLower(code)
}

func ContainsErrMsg(err error, msg string) bool {
	if err == nil {
		return false
	}

	return strings.Contains(
		strings.ToLower(err.Error()),
		strings.ToLower(msg))
}

var CreateRenderedDeploy = func(
	backend backend.Backend,
	logGroupID string,
	task *ecs.TaskDefinition,
	createDeploy backend.CreateDeployf,
) (*models.Deploy, error) {
	for _, container := range task.ContainerDefinitions {
		if container.LogConfiguration == nil {
			container.LogConfiguration = &awsecs.LogConfiguration{
				LogDriver: stringp("awslogs"),
				Options: map[string]*string{
					"awslogs-group":  stringp(logGroupID),
					"awslogs-region": stringp(config.AWSRegion()),
				},
			}
		}
	}

	dockerrun, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	// *task.Family is in format 'l0-<prefix>-<deployName>'
	deployName := strings.TrimPrefix(*task.Family, id.PREFIX)
	request := models.CreateDeployRequest{
		DeployName: deployName,
		Dockerrun:  dockerrun,
	}

	return createDeploy(request)
}

// IteratePages performs a do-while loop on a paginatedf
// until nextToken is nil or an error is returned
type paginatedf func(*string) (*string, error)

func IteratePages(fn paginatedf) error {
	var err error
	var nextToken *string

	nextToken, err = fn(nextToken)
	if err != nil {
		return err
	}

	for nextToken != nil {
		nextToken, err = fn(nextToken)
		if err != nil {
			return err
		}
	}

	return nil
}

var GetLogs = func(cloudWatchLogs cloudwatchlogs.Provider, logGroupID string, tail int) ([]*models.LogFile, error) {
	logStreams := []*cloudwatchlogs.LogStream{}
	getStreams := func(nextToken *string) (*string, error) {
		orderBy := "LogStreamName"
		streams, updatedToken, err := cloudWatchLogs.DescribeLogStreams(logGroupID, orderBy, nextToken)
		if err != nil {
			return nil, err
		}

		logStreams = append(logStreams, streams...)
		return updatedToken, nil
	}

	if err := IteratePages(getStreams); err != nil {
		return nil, err
	}

	logFiles := []*models.LogFile{}
	for _, logStream := range logStreams {
		// each logStream corresponds to a container
		logFile := &models.LogFile{
			Name:  *logStream.LogStreamName,
			Lines: []string{},
		}

		// since the time range is exclusive, expand the range to get first/last events
		logStream.FirstEventTimestamp = int64p(*logStream.FirstEventTimestamp - 1)
		logStream.LastEventTimestamp = int64p(*logStream.LastEventTimestamp + 1)

		// todo: use tail to lookup from back to front
		getEvents := func(nextToken *string) (*string, error) {
			logEvents, updatedToken, err := cloudWatchLogs.GetLogEvents(
				logGroupID,
				*logStream.LogStreamName,
				nextToken,
				logStream.FirstEventTimestamp,
				logStream.LastEventTimestamp)
			if err != nil {
				return nil, err
			}

			for _, logEvent := range logEvents {
				logFile.Lines = append(logFile.Lines, *logEvent.Message)
			}

			// GetLogEvents re-uses the same NextToken when it is finished instead
			// of returning nil
			if nextToken != nil && *updatedToken == *nextToken {
				return nil, nil
			}

			return updatedToken, nil
		}

		if err := IteratePages(getEvents); err != nil {
			return nil, err
		}

		if numLines := len(logFile.Lines); tail != 0 && numLines > tail {
			logFile.Lines = logFile.Lines[numLines-tail:]
		}

		logFiles = append(logFiles, logFile)
	}

	return logFiles, nil
}
