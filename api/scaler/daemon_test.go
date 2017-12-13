package scaler

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/job/mock_job"
	"github.com/quintilesims/layer0/api/provider/mock_provider"
	"github.com/quintilesims/layer0/common/models"
)

func TestDaemonFN(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockJobStore := mock_job.NewMockStore(ctrl)
	mockEnvironmentProvider := mock_provider.NewMockEnvironmentProvider(ctrl)
	daemonFN := NewDaemonFN(mockJobStore, mockEnvironmentProvider)

	environmentSummaries := []models.EnvironmentSummary{
		{EnvironmentID: "env_id1"},
		{EnvironmentID: "env_id2"},
	}

	mockEnvironmentProvider.EXPECT().
		List().
		Return(environmentSummaries, nil)

	for _, es := range environmentSummaries {
		mockJobStore.EXPECT().
			Insert(models.ScaleEnvironmentJob, es.EnvironmentID).
			Return("", nil)
	}

	if err := daemonFN(); err != nil {
		t.Fatal(err)
	}
}
