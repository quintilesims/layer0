package ecsbackend

import (
	"github.com/golang/mock/gomock"
	"gitlab.imshealth.com/xfra/layer0/common/aws/cloudwatchlogs"
	"gitlab.imshealth.com/xfra/layer0/common/aws/cloudwatchlogs/mock_cloudwatchlogs"
	"gitlab.imshealth.com/xfra/layer0/common/models"
	"gitlab.imshealth.com/xfra/layer0/common/testutils"
	"testing"
)

func testCreateDeploy(models.CreateDeployRequest) (*models.Deploy, error) {
	model := &models.Deploy{
		DeployID: "renderedid.2",
		Version:  "2",
	}

	return model, nil
}

func TestGetLogs(t *testing.T) {
	testCases := []testutils.TestCase{
		testutils.TestCase{
			Name: "Should add proper logs lines to models",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockCW := mock_cloudwatchlogs.NewMockProvider(ctrl)

				stream := cloudwatchlogs.NewLogStream("some_name")
				stream.FirstEventTimestamp = int64p(int64(0))
				stream.LastEventTimestamp = int64p(int64(1))

				mockCW.EXPECT().
					DescribeLogStreams("logGroupID", "LogStreamName", nil).
					Return([]*cloudwatchlogs.LogStream{stream}, nil, nil)

				event := cloudwatchlogs.NewOutputLogEvent("some_message")

				mockCW.EXPECT().
					GetLogEvents(
						"logGroupID",
						"some_name",
						nil,
						gomock.Any(),
						gomock.Any(),
					).Return([]*cloudwatchlogs.OutputLogEvent{event}, nil, nil)

				return mockCW
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				provider := target.(cloudwatchlogs.Provider)

				logs, err := GetLogs(provider, "logGroupID", 0)
				if err != nil {
					reporter.Fatal(err)
				}

				reporter.AssertEqual(len(logs), 1)
				reporter.AssertInSlice("some_message", logs[0].Lines)
			},
		},
	}

	testutils.RunTests(t, testCases)
}
