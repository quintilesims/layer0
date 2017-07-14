package ecsbackend

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	awsecs "github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/api/backend/ecs/id"
	"github.com/quintilesims/layer0/common/aws/cloudwatchlogs"
	"github.com/quintilesims/layer0/common/aws/ecs"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
)

const MAX_TASK_IDS = 100

func boolp(b bool) *bool {
	return &b
}

func pbool(b *bool) bool {
	if b == nil {
		return false
	}

	return *b
}

func stringp(s string) *string {
	return &s
}

func intp(i int) *int {
	return &i
}

func int64p(i int64) *int64 {
	return &i
}

func pstring(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func pint(i *int) int {
	if i == nil {
		return 0
	}

	return *i
}

func pint64(i *int64) int64 {
	if i == nil {
		return 0
	}

	return *i
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

var CreateRenderedDeploy = func(body []byte) (*Deploy, error) {
	deploy, err := MarshalDeploy(body)
	if err != nil {
		return nil, err
	}

	for _, container := range deploy.ContainerDefinitions {
		if container.LogConfiguration == nil {
			container.LogConfiguration = &awsecs.LogConfiguration{
				LogDriver: stringp("awslogs"),
				Options: map[string]*string{
					"awslogs-group":         stringp(config.AWSLogGroupID()),
					"awslogs-region":        stringp(config.AWSRegion()),
					"awslogs-stream-prefix": stringp("l0"),
				},
			}
		}
	}

	return deploy, nil
}

var getTaskARNs = func(ecs ecs.Provider, ecsEnvironmentID id.ECSEnvironmentID, startedBy *string) ([]*string, error) {
	// we can only check each of the states individually, thus we must issue 2 API calls

	running := "RUNNING"
	tasks, err := ecs.ListTasks(ecsEnvironmentID.String(), nil, &running, startedBy, nil)
	if err != nil {
		return nil, err
	}

	stopped := "STOPPED"
	stoppedTasks, err := ecs.ListTasks(ecsEnvironmentID.String(), nil, &stopped, startedBy, nil)
	if err != nil {
		return nil, err
	}

	tasks = append(tasks, stoppedTasks...)

	return tasks, nil
}

var GetLogs = func(cloudWatchLogs cloudwatchlogs.Provider, taskARNs []*string, start, end string, tail int) ([]*models.LogFile, error) {
	taskIDCatalog := generateTaskIDCatalog(taskARNs)

	orderBy := "LastEventTime"
	logStreams, err := cloudWatchLogs.DescribeLogStreams(config.AWSLogGroupID(), orderBy)
	if err != nil {
		return nil, err
	}

	logFiles := []*models.LogFile{}
	for _, logStream := range logStreams {
		// filter by streams that have <prefix>/<container name>/<stream task id>
		streamNameSplit := strings.Split(*logStream.LogStreamName, "/")
		if len(streamNameSplit) != 3 {
			continue
		}

		streamTaskID := streamNameSplit[2]
		if _, ok := taskIDCatalog[streamTaskID]; !ok {
			continue
		}

		logFile := &models.LogFile{
			Name:  streamNameSplit[1],
			Lines: []string{},
		}

		// since the time range is exclusive, expand the range to get first/last events
		logEvents, err := cloudWatchLogs.GetLogEvents(
			config.AWSLogGroupID(),
			*logStream.LogStreamName,
			start,
			end,
			int64(tail))
		if err != nil {
			return nil, err
		}

		for _, logEvent := range logEvents {
			logFile.Lines = append(logFile.Lines, *logEvent.Message)
		}

		logFiles = append(logFiles, logFile)
	}

	return logFiles, nil
}

func generateTaskIDCatalog(taskARNs []*string) map[string]bool {
	catalog := map[string]bool{}
	for _, taskARN := range taskARNs {
		taskID := strings.Split(*taskARN, "/")[1]
		catalog[taskID] = true
	}

	return catalog
}
