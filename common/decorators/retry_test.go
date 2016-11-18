package decorators

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/golang/mock/gomock"
	"gitlab.imshealth.com/xfra/layer0/common/aws/ecs"
	"gitlab.imshealth.com/xfra/layer0/common/aws/ecs/mock_ecs"
	"gitlab.imshealth.com/xfra/layer0/common/testutils"
	"testing"
)

func TestRetry_successes(t *testing.T) {
	retryCount := []int{0, 1, 19}

	for _, count := range retryCount {
		ctrl := gomock.NewController(testutils.NewReporter(t, "TestRetry_successes"))
		defer ctrl.Finish()

		mockECS := mock_ecs.NewMockProvider(ctrl)
		if count > 0 {
			mockECS.EXPECT().CreateCluster(gomock.Any()).
				Return(nil, getRetryTrigger()).
				Times(count)
		}

		cluster := &ecs.Cluster{}
		mockECS.EXPECT().CreateCluster(gomock.Any()).
			Return(cluster, nil)

		wrap := prepareRetry(mockECS)

		obj, err := wrap.CreateCluster("test")
		if err != nil {
			t.Fatal(err)
		}

		if cluster != obj {
			t.Error("Cluster not equal to object on iteration %d", count)
		}
	}
}

func TestRetry_timeout(t *testing.T) {
	ctrl := gomock.NewController(testutils.NewReporter(t, "TestRetry_timeout"))
	defer ctrl.Finish()

	mockECS := mock_ecs.NewMockProvider(ctrl)

	mockECS.EXPECT().CreateCluster(gomock.Any()).
		Return(nil, getRetryTrigger()).
		Times(20)

	wrap := prepareRetry(mockECS)

	_, err := wrap.CreateCluster("test")
	if err == nil {
		t.Errorf("Error was unexpectedly nil")
	}
}

func TestRetry_otherError(t *testing.T) {
	ctrl := gomock.NewController(testutils.NewReporter(t, "TestRetry_otherError"))
	defer ctrl.Finish()

	mockECS := mock_ecs.NewMockProvider(ctrl)

	mockECS.EXPECT().CreateCluster(gomock.Any()).
		Return(nil, fmt.Errorf("Some error"))

	wrap := prepareRetry(mockECS)

	_, err := wrap.CreateCluster("test")
	if err == nil {
		t.Errorf("Error was unexpectedly nil")
	}
}

func prepareRetry(mockECS ecs.Provider) ecs.Provider {
	retry := &Retry{
		&testutils.StubClock{},
	}

	wrap := &ecs.ProviderDecorator{
		Inner:     mockECS,
		Decorator: retry.CallWithRetries,
	}

	return wrap
}

type testAwsError struct {
	code    string
	message string
}

func (this testAwsError) Code() string {
	return this.code
}

func (this testAwsError) Message() string {
	return this.message
}

func (this testAwsError) OrigErr() error {
	return nil
}

func (this testAwsError) Error() string {
	return fmt.Sprintf("%s -- %s", this.code, this.message)
}

func getRetryTrigger() awserr.Error {
	return testAwsError{code: "Throttling"}
}
