package cloudwatchlogs

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/quintilesims/layer0/common/aws/provider"
)

const (
	// DescribeLogStreams is throttled after five transactions per second.
	// With 50 streams/transaction, 1000 gives a reasonable streams:time ratio
	MAX_DESCRIBE_STREAMS_COUNT = 1000
<<<<<<< HEAD

	// 'YYYY-MM-DD HH:MM' time layout as described by https://golang.org/src/time/format.go
	TIME_LAYOUT = "2006-01-02 15:04"
=======
>>>>>>> remotes/origin/master
)

type Provider interface {
	CreateLogGroup(logGroupName string) error
	DeleteLogGroup(logGroupName string) error
	DescribeLogGroups(logGroupNamePrefix string, nextToken *string) ([]*LogGroup, error)
	DescribeLogStreams(logGroupName, orderBy string) ([]*LogStream, error)
	GetLogEvents(logGroupName, logStreamName, start, stop string, limit int64) ([]*OutputLogEvent, error)
	FilterLogEvents(filterPattern, logGroupName, nextToken *string, logStreamNames []*string, endTime, startTime *int64, interleaved *bool) ([]*FilteredLogEvent, []*SearchedLogStream, error)
}

type CloudWatchLogs struct {
	credProvider provider.CredProvider
	region       string
	Connect      func() (CloudWatchLogsInternal, error)
}

type CloudWatchLogsInternal interface {
	CreateLogGroup(input *cloudwatchlogs.CreateLogGroupInput) (*cloudwatchlogs.CreateLogGroupOutput, error)
	DeleteLogGroup(input *cloudwatchlogs.DeleteLogGroupInput) (*cloudwatchlogs.DeleteLogGroupOutput, error)
	DescribeLogGroups(input *cloudwatchlogs.DescribeLogGroupsInput) (*cloudwatchlogs.DescribeLogGroupsOutput, error)
	DescribeLogStreamsPages(input *cloudwatchlogs.DescribeLogStreamsInput, fn func(p *cloudwatchlogs.DescribeLogStreamsOutput, lastPage bool) (shouldContinue bool)) error
	GetLogEventsPages(input *cloudwatchlogs.GetLogEventsInput, fn func(p *cloudwatchlogs.GetLogEventsOutput, lastPage bool) (shouldContinue bool)) error
	FilterLogEvents(input *cloudwatchlogs.FilterLogEventsInput) (*cloudwatchlogs.FilterLogEventsOutput, error)
}

type LogGroup struct {
	*cloudwatchlogs.LogGroup
}

type LogStream struct {
	*cloudwatchlogs.LogStream
}

func NewLogStream(name string) *LogStream {
	return &LogStream{
		&cloudwatchlogs.LogStream{
			LogStreamName: aws.String(name),
		},
	}
}

type FilteredLogEvent struct {
	*cloudwatchlogs.FilteredLogEvent
}

type OutputLogEvent struct {
	*cloudwatchlogs.OutputLogEvent
}

func NewOutputLogEvent(message string) *OutputLogEvent {
	return &OutputLogEvent{
		&cloudwatchlogs.OutputLogEvent{
			Message: aws.String(message),
		},
	}
}

type SearchedLogStream struct {
	*cloudwatchlogs.SearchedLogStream
}

func NewCloudWatchLogs(credProvider provider.CredProvider, region string) (Provider, error) {
	cloudwatchlogs := CloudWatchLogs{
		credProvider,
		region,
		func() (CloudWatchLogsInternal, error) {
			return Connect(credProvider, region)
		},
	}

	if _, err := cloudwatchlogs.Connect(); err != nil {
		return nil, err
	}

	return &cloudwatchlogs, nil
}

func Connect(credProvider provider.CredProvider, region string) (CloudWatchLogsInternal, error) {
	connection, err := provider.GetCloudWatchLogsConnection(credProvider, region)
	if err != nil {
		return nil, err
	}

	return connection, nil
}

func (this *CloudWatchLogs) CreateLogGroup(logGroupName string) error {
	input := &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(logGroupName),
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}

	if _, err := connection.CreateLogGroup(input); err != nil {
		return err
	}

	return nil
}

func (this *CloudWatchLogs) DeleteLogGroup(logGroupName string) error {
	input := &cloudwatchlogs.DeleteLogGroupInput{
		LogGroupName: aws.String(logGroupName),
	}

	connection, err := this.Connect()
	if err != nil {
		return err
	}

	if _, err := connection.DeleteLogGroup(input); err != nil {
		return err
	}

	return nil
}

func (this *CloudWatchLogs) DescribeLogGroups(logGroupNamePrefix string, nextToken *string) ([]*LogGroup, error) {
	input := &cloudwatchlogs.DescribeLogGroupsInput{
		// The maximum number of items returned in the response. If you don't specify
		// a value, the request would return up to 50 items.
		// Limit: aws.Int64(limit),

		// Will only return log groups that match the provided logGroupNamePrefix. If
		// you don't specify a value, no prefix filter is applied.
		LogGroupNamePrefix: aws.String(logGroupNamePrefix),

		// A string token used for pagination that points to the next page of results.
		// It must be a value obtained from the response of the previous DescribeLogGroups
		// request.
		NextToken: nextToken,
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	resp, err := connection.DescribeLogGroups(input)
	if err != nil {
		return nil, err
	}

	result := []*LogGroup{}
	for _, group := range resp.LogGroups {
		result = append(result, &LogGroup{group})
	}

	return result, nil
}

func (this *CloudWatchLogs) DescribeLogStreams(logGroupName, orderBy string) ([]*LogStream, error) {
	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	input := &cloudwatchlogs.DescribeLogStreamsInput{
<<<<<<< HEAD
		LogGroupName: aws.String(logGroupName),
		OrderBy:      aws.String(orderBy),
		Descending:   aws.Bool(true),
=======
		LogGroupName:        aws.String(logGroupName),
		OrderBy:             aws.String(orderBy),
		Descending:          aws.Bool(true),
>>>>>>> remotes/origin/master
	}

	streams := []*LogStream{}
	pagef := func(output *cloudwatchlogs.DescribeLogStreamsOutput, lastPage bool) bool {
		for _, stream := range output.LogStreams {
			streams = append(streams, &LogStream{stream})
		}

		if len(streams) >= MAX_DESCRIBE_STREAMS_COUNT {
			return false
		}

		return !lastPage
	}

	if err := connection.DescribeLogStreamsPages(input, pagef); err != nil {
		return nil, err
	}

	return streams, nil
<<<<<<< HEAD
}

func timeToMilliseconds(v string) (int64, error) {
	t, err := time.Parse(TIME_LAYOUT, v)
	if err != nil {
		return 0, fmt.Errorf("Invalid time: must be in format YYYY-MM-DD HH:MM")
	}

	date := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)

	// convert ns to ms
	return date.UnixNano() / int64(time.Millisecond), nil
=======
>>>>>>> remotes/origin/master
}

func (this *CloudWatchLogs) GetLogEvents(
	logGroupName string,
	logStreamName string,
	start string,
	end string,
	limit int64,
) ([]*OutputLogEvent, error) {
	connection, err := this.Connect()
	if err != nil {
		return nil, err
	}

	input := &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
	}

	if limit > 0 {
		input.SetLimit(limit)
	}

	if start != "" {
		startTime, err := timeToMilliseconds(start)
		if err != nil {
			return nil, err
		}

		input.SetStartTime(startTime)
	}

	if end != "" {
		endTime, err := timeToMilliseconds(end)
		if err != nil {
			return nil, err
		}

		input.SetEndTime(endTime)
	}

	var previousToken string
	result := []*OutputLogEvent{}
	pagef := func(output *cloudwatchlogs.GetLogEventsOutput, lastPage bool) bool {
		for _, event := range output.Events {
			result = append(result, &OutputLogEvent{event})
		}

		// GetLogEvents re-uses the same NextToken when it is finished instead of returning nil
		if previousToken == *output.NextForwardToken {
			return false
		}

		previousToken = *output.NextForwardToken
		return true
	}

	if err := connection.GetLogEventsPages(input, pagef); err != nil {
		return nil, err
	}

	return result, nil
}

func (this *CloudWatchLogs) FilterLogEvents(
	filterPattern,
	logGroupName,
	nextToken *string,
	logStreamNames []*string,
	endTime,
	startTime *int64,
	interleaved *bool) ([]*FilteredLogEvent, []*SearchedLogStream, error) {

	if *nextToken == "" {
		nextToken = nil
	}

	input := &cloudwatchlogs.FilterLogEventsInput{
		// A point in time expressed as the number of milliseconds since Jan 1, 1970
		// 00:00:00 UTC. If provided, events with a timestamp later than this time are
		// not returned.
		EndTime: endTime,

		// A valid CloudWatch Logs filter pattern to use for filtering the response.
		// If not provided, all the events are matched.
		FilterPattern: filterPattern,

		// If provided, the API will make a best effort to provide responses that contain
		// events from multiple log streams within the log group interleaved in a single
		// response. If not provided, all the matched log events in the first log stream
		// will be searched first, then those in the next log stream, etc.
		Interleaved: interleaved,

		// The maximum number of events to return in a page of results. Default is 10,000
		// events.
		// Limit: limit,

		// The name of the log group to query.
		LogGroupName: logGroupName,

		// Optional list of log stream names within the specified log group to search.
		// Defaults to all the log streams in the log group.
		LogStreamNames: logStreamNames,

		// A pagination token obtained from a FilterLogEvents response to continue paginating
		// the FilterLogEvents results. This token is omitted from the response when
		// there are no other events to display.
		NextToken: nextToken,

		// A point in time expressed as the number of milliseconds since Jan 1, 1970
		// 00:00:00 UTC. If provided, events with a timestamp prior to this time are
		// not returned.
		StartTime: startTime,
	}

	connection, err := this.Connect()
	if err != nil {
		return nil, nil, err
	}
	resp, err := connection.FilterLogEvents(input)
	if err != nil {
		return nil, nil, err
	}

	resultFiltered := []*FilteredLogEvent{}
	for _, svc := range resp.Events {
		resultFiltered = append(resultFiltered, &FilteredLogEvent{svc})
	}

	resultSearched := []*SearchedLogStream{}
	for _, svc := range resp.SearchedLogStreams {
		resultSearched = append(resultSearched, &SearchedLogStream{svc})
	}

	return resultFiltered, resultSearched, nil
}
