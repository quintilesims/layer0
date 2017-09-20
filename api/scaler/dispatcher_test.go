package scaler

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/provider/mock_provider"
	"github.com/quintilesims/layer0/common/models"
	"github.com/stretchr/testify/assert"
)

func setTimeMultiplier(m int) func() {
	timeMultiplier = time.Duration(m)
	return func() { timeMultiplier = 1 }
}

func TestDispatcherScheduleRun(t *testing.T) {
	defer setTimeMultiplier(0)()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEnvironmentProvider := mock_provider.NewMockEnvironmentProvider(ctrl)

	done := make(chan bool)
	scaler := ScalerFunc(func(environmentID string) error {
		assert.Equal(t, "eid", environmentID)
		done <- true
		return nil
	})

	dispatcher := NewDispatcher(mockEnvironmentProvider, scaler)
	dispatcher.ScheduleRun("eid")

	select {
	case <-done:
	case <-time.After(time.Second * 1):
		t.Fatalf("Failed to run after 1 second")
	}
}

func TestDispatcherRunAll(t *testing.T) {
	defer setTimeMultiplier(0)()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	environmentIDs := []string{"e1", "e2", "e3"}
	summaries := make([]models.EnvironmentSummary, len(environmentIDs))
	for i, environmentID := range environmentIDs {
		summaries[i] = models.EnvironmentSummary{EnvironmentID: environmentID}
	}

	mockEnvironmentProvider := mock_provider.NewMockEnvironmentProvider(ctrl)
	mockEnvironmentProvider.EXPECT().
		List().
		Return(summaries, nil)

	done := make(chan bool)
	scaler := ScalerFunc(func(environmentID string) error {
		assert.Contains(t, environmentIDs, environmentID)
		done <- true
		return nil
	})

	dispatcher := NewDispatcher(mockEnvironmentProvider, scaler)
	dispatcher.RunAll()

	for i := 0; i < len(environmentIDs); i++ {
		select {
		case <-done:
		case <-time.After(time.Second * 1):
			t.Fatalf("Failed to run after 1 second")
		}
	}
}
