package test_aws

// import (
// 	"testing"
// 	"time"
//
// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
// 	"github.com/aws/aws-sdk-go/service/ecs"
// 	"github.com/golang/mock/gomock"
// 	provider "github.com/quintilesims/layer0/api/provider/aws"
// 	"github.com/quintilesims/layer0/api/tag"
// 	awsc "github.com/quintilesims/layer0/common/aws"
// 	"github.com/quintilesims/layer0/common/config/mock_config"
// 	"github.com/quintilesims/layer0/common/models"
// 	"github.com/stretchr/testify/assert"
// )
//
// func TestServiceLogs(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
//
// 	mockAWS := awsc.NewMockClient(ctrl)
// 	tagStore := tag.NewMemoryStore()
// 	c := mock_config.NewMockAPIConfig(ctrl)
//
// 	c.EXPECT().Instance().Return("test").AnyTimes()
//
// 	tags := models.Tags{
// 		{
// 			EntityID:   "svc_id",
// 			EntityType: "service",
// 			Key:        "name",
// 			Value:      "svc_name",
// 		},
// 		{
// 			EntityID:   "svc_id",
// 			EntityType: "service",
// 			Key:        "environment_id",
// 			Value:      "env_id",
// 		},
// 	}
//
// 	for _, tag := range tags {
// 		if err := tagStore.Insert(tag); err != nil {
// 			t.Fatal(err)
// 		}
// 	}
//
// 	describeServicesInput := &ecs.DescribeServicesInput{}
// 	describeServicesInput.SetCluster("l0-test-env_id")
// 	describeServicesInput.SetServices([]*string{aws.String("l0-test-svc_id")})
//
// 	deployment1 := &ecs.Deployment{}
// 	deployment1.SetId("dpl_id:1")
// 	deployment1.SetStatus(ecs.DesiredStatusRunning)
// 	deployment1.SetTaskDefinition("arn:aws:ecs:region:012345678910:task-definition/dpl_id:1")
//
// 	deployment2 := &ecs.Deployment{}
// 	deployment2.SetId("dpl_id:2")
// 	deployment2.SetStatus(ecs.DesiredStatusStopped)
// 	deployment2.SetTaskDefinition("arn:aws:ecs:region:012345678910:task-definition/dpl_id:2")
//
// 	deployments := []*ecs.Deployment{
// 		deployment1,
// 		deployment2,
// 	}
//
// 	service := &ecs.Service{}
// 	service.SetDeployments(deployments)
// 	services := []*ecs.Service{
// 		service,
// 	}
//
// 	describeServicesOutput := &ecs.DescribeServicesOutput{}
// 	describeServicesOutput.SetServices(services)
//
// 	mockAWS.ECS.EXPECT().
// 		DescribeServices(describeServicesInput).
// 		Return(describeServicesOutput, nil)
//
// 	generateListTasksPagesFN := func(input *ecs.ListTasksInput) func(input *ecs.ListTasksInput, fn func(output *ecs.ListTasksOutput, lastPage bool) bool) error {
// 		var taskARN *string
// 		inputStatus := aws.StringValue(input.DesiredStatus)
// 		inputStartedBy := aws.StringValue(input.StartedBy)
//
// 		if inputStatus == ecs.DesiredStatusRunning && inputStartedBy == "dpl_id:1" {
// 			taskARN = aws.String("arn:aws:ecs:region:012345678910:task-definition/dpl_id:1")
// 		}
//
// 		if inputStatus == ecs.DesiredStatusStopped && inputStartedBy == "dpl_id:2" {
// 			taskARN = aws.String("arn:aws:ecs:region:012345678910:task-definition/dpl_id:2")
// 		}
//
// 		listTasksPagesFN := func(input *ecs.ListTasksInput, fn func(output *ecs.ListTasksOutput, lastPage bool) bool) error {
// 			output := &ecs.ListTasksOutput{}
// 			if taskARN != nil {
// 				taskARNs := []*string{
// 					taskARN,
// 				}
//
// 				output.SetTaskArns(taskARNs)
// 			}
//
// 			fn(output, true)
//
// 			return nil
// 		}
//
// 		return listTasksPagesFN
// 	}
//
// 	deploymentIDs := []string{
// 		"dpl_id:1",
// 		"dpl_id:2",
// 	}
//
// 	statuses := []string{
// 		ecs.DesiredStatusRunning,
// 		ecs.DesiredStatusStopped,
// 	}
//
// 	for _, deploymentID := range deploymentIDs {
// 		for _, status := range statuses {
// 			input := &ecs.ListTasksInput{}
// 			input.SetCluster("l0-test-env_id")
// 			input.SetDesiredStatus(status)
// 			input.SetStartedBy(deploymentID)
//
// 			mockAWS.ECS.EXPECT().
// 				ListTasksPages(input, gomock.Any()).
// 				Do(generateListTasksPagesFN(input)).
// 				Return(nil)
// 		}
// 	}
//
// 	c.EXPECT().
// 		LogGroupName().
// 		Return("l0-test")
//
// 	describeLogStreamsPagesFN := func(input *cloudwatchlogs.DescribeLogStreamsInput, fn func(ouput *cloudwatchlogs.DescribeLogStreamsOutput, lastPage bool) bool) error {
// 		logStream1 := &cloudwatchlogs.LogStream{}
// 		logStream1.SetLogStreamName("l0/container1/dpl_id:1")
//
// 		logStream2 := &cloudwatchlogs.LogStream{}
// 		logStream2.SetLogStreamName("l0/container2/dpl_id:2")
//
// 		logStreams := []*cloudwatchlogs.LogStream{
// 			logStream1,
// 			logStream2,
// 		}
//
// 		output := &cloudwatchlogs.DescribeLogStreamsOutput{}
// 		output.SetLogStreams(logStreams)
// 		fn(output, true)
//
// 		return nil
// 	}
//
// 	mockAWS.CloudWatchLogs.EXPECT().
// 		DescribeLogStreamsPages(gomock.Any(), gomock.Any()).
// 		Do(describeLogStreamsPagesFN).
// 		Return(nil)
//
// 	getLogEventsPagesFN := func(input *cloudwatchlogs.GetLogEventsInput, fn func(ouput *cloudwatchlogs.GetLogEventsOutput, lastPage bool) bool) error {
// 		event := &cloudwatchlogs.OutputLogEvent{}
// 		event.SetIngestionTime(int64(1234567890))
// 		event.SetMessage("log message")
// 		event.SetTimestamp(int64(1234567890))
//
// 		events := []*cloudwatchlogs.OutputLogEvent{
// 			event,
// 		}
//
// 		output := &cloudwatchlogs.GetLogEventsOutput{}
// 		output.SetEvents(events)
// 		output.SetNextForwardToken("next")
// 		fn(output, true)
//
// 		return nil
// 	}
//
// 	mockAWS.CloudWatchLogs.EXPECT().
// 		GetLogEventsPages(gomock.Any(), gomock.Any()).
// 		Do(getLogEventsPagesFN).
// 		Return(nil).
// 		Times(2)
//
// 	target := provider.NewServiceProvider(mockAWS.Client(), tagStore, c)
// 	result, err := target.Logs("svc_id", 0, time.Time{}, time.Time{})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	expected := []models.LogFile{
// 		{
// 			ContainerName: "container1",
// 			Lines:         []string{"log message"},
// 		},
// 		{
// 			ContainerName: "container2",
// 			Lines:         []string{"log message"},
// 		},
// 	}
//
// 	assert.Equal(t, expected, result)
// }
