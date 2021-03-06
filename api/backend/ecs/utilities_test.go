package ecsbackend

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/aws/cloudwatchlogs"
	"github.com/quintilesims/layer0/common/aws/cloudwatchlogs/mock_cloudwatchlogs"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
)

func testCreateDeploy(models.CreateDeployRequest) (*models.Deploy, error) {
	model := &models.Deploy{
		DeployID: "renderedid.2",
		Version:  "2",
	}

	return model, nil
}

func TestGetLogs(t *testing.T) {
	taskARN := "arn:aws:ecs:region:aws_account_id:task/taskARN"

	testCases := []testutils.TestCase{
		{
			Name: "Should add proper logs lines to models",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockCW := mock_cloudwatchlogs.NewMockProvider(ctrl)

				stream := cloudwatchlogs.NewLogStream("prefix/container_name/taskARN")

				mockCW.EXPECT().
					DescribeLogStreams(config.AWSLogGroupID(), "LastEventTime").
					Return([]*cloudwatchlogs.LogStream{stream}, nil)

				event := cloudwatchlogs.NewOutputLogEvent("some_message")

				mockCW.EXPECT().
					GetLogEvents(
						config.AWSLogGroupID(),
						*stream.LogStreamName,
						"start",
						"end",
						int64(30),
					).Return([]*cloudwatchlogs.OutputLogEvent{event}, nil)

				return mockCW
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				provider := target.(cloudwatchlogs.Provider)

				logs, err := GetLogs(provider, []*string{stringp(taskARN)}, "start", "end", 30)
				if err != nil {
					reporter.Fatal(err)
				}

				testutils.AssertEqual(t, len(logs), 1)
				testutils.AssertInSlice(t, "some_message", logs[0].Lines)
			},
		},
	}

	testutils.RunTests(t, testCases)
}
