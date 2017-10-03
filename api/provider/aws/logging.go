package aws

import (
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/quintilesims/layer0/common/models"
)

const (
	// DescribeLogStreams is throttled after five transactions per second.
	// With 50 streams/transaction, 1000 gives a reasonable streams:time ratio
	MAX_DESCRIBE_STREAMS_COUNT = 1000
)

func GetLogsFromTaskARNs(
	cloudWatchLogsAPI cloudwatchlogsiface.CloudWatchLogsAPI,
	logGroupName string,
	taskARNs []string,
	tail int,
	start,
	end time.Time,
) ([]models.LogFile, error) {
	taskIDMatch := generateTaskIDMap(taskARNs)
	logStreams, err := describeLogStreams(cloudWatchLogsAPI, logGroupName)
	if err != nil {
		return nil, err
	}

	logFiles := []models.LogFile{}
	for _, logStream := range logStreams {
		// ecs task log streams have format '<prefix>/<container name>/<task id>'
		streamNameSplit := strings.Split(aws.StringValue(logStream.LogStreamName), "/")
		if len(streamNameSplit) != 3 {
			continue
		}

		containerName := streamNameSplit[1]
		taskID := streamNameSplit[2]

		if !taskIDMatch[taskID] {
			continue
		}

		logStreamName := aws.StringValue(logStream.LogStreamName)
		logEvents, err := getLogEvents(cloudWatchLogsAPI, logGroupName, logStreamName, tail, start, end)
		if err != nil {
			return nil, err
		}

		logFile := models.LogFile{
			ContainerName: containerName,
			Lines:        make([]string, len(logEvents)),
		}

		for i, event := range logEvents {
			logFile.Lines[i] = aws.StringValue(event.Message)
		}

		logFiles = append(logFiles, logFile)
	}

	return logFiles, nil
}

func describeLogStreams(cloudWatchLogsAPI cloudwatchlogsiface.CloudWatchLogsAPI, logGroupName string) ([]*cloudwatchlogs.LogStream, error) {
	input := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(logGroupName),
		OrderBy:      aws.String(cloudwatchlogs.OrderByLastEventTime),
		Descending:   aws.Bool(true),
	}

	logStreams := []*cloudwatchlogs.LogStream{}
	fn := func(output *cloudwatchlogs.DescribeLogStreamsOutput, lastPage bool) bool {
		logStreams = append(logStreams, output.LogStreams...)

		if len(logStreams) >= MAX_DESCRIBE_STREAMS_COUNT {
			return false
		}

		return !lastPage
	}

	if err := cloudWatchLogsAPI.DescribeLogStreamsPages(input, fn); err != nil {
		return nil, err
	}

	return logStreams, nil
}

func getLogEvents(
	cloudWatchLogsAPI cloudwatchlogsiface.CloudWatchLogsAPI,
	logGroupName string,
	logStreamName string,
	tail int,
	start time.Time,
	end time.Time,
) ([]*cloudwatchlogs.OutputLogEvent, error) {
	input := &cloudwatchlogs.GetLogEventsInput{}
	input.SetLogGroupName(logGroupName)
	input.SetLogStreamName(logStreamName)

	if tail != 0 {
		input.SetLimit(int64(tail))
	}

	if !start.IsZero() {
		startMS := timeToMilliseconds(start)
		input.SetStartTime(startMS)
	}

	if !end.IsZero() {
		endMS := timeToMilliseconds(end)
		input.SetEndTime(endMS)
	}

	var previousToken string
	events := []*cloudwatchlogs.OutputLogEvent{}
	eventsFN := func(output *cloudwatchlogs.GetLogEventsOutput, lastPage bool) bool {
		defer func() { previousToken = aws.StringValue(output.NextForwardToken) }()
		events = append(events, output.Events...)

		// GetLogEvents re-uses the same NextToken when it is finished instead of returning nil
		return previousToken != aws.StringValue(output.NextForwardToken)
	}

	if err := cloudWatchLogsAPI.GetLogEventsPages(input, eventsFN); err != nil {
		return nil, err
	}

	return events, nil
}

func generateTaskIDMap(taskARNs []string) map[string]bool {
	catalog := map[string]bool{}
	for _, taskARN := range taskARNs {
		// task arn format: arn:aws:ecs:region:account_id:task/task_id
		taskID := strings.Split(taskARN, "/")[1]
		catalog[taskID] = true
	}

	return catalog
}

func timeToMilliseconds(t time.Time) int64 {
	date := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
	return date.UnixNano() / int64(time.Millisecond)
}
