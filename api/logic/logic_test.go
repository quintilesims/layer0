package logic

import (
	log "github.com/Sirupsen/logrus"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/backend/mock_backend"
	"github.com/quintilesims/layer0/api/scheduler/mock_scheduler"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/db/job_store"
	"github.com/quintilesims/layer0/common/db/tag_store"
	"github.com/quintilesims/layer0/common/models"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	config.SetTestConfig()
	log.SetLevel(log.FatalLevel)
	jobLogger.Level = log.FatalLevel
	retCode := m.Run()
	os.Exit(retCode)
}

type TestLogic struct {
	Backend  *mock_backend.MockBackend
	JobStore *job_store.MemoryJobStore
	TagStore *tag_store.MemoryTagStore
	Scaler   *mock_scheduler.MockEnvironmentScaler
}

func NewTestLogic(t *testing.T) (*TestLogic, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	logic := &TestLogic{
		Backend:  mock_backend.NewMockBackend(ctrl),
		JobStore: job_store.NewMemoryJobStore(),
		TagStore: tag_store.NewMemoryTagStore(),
		Scaler:   mock_scheduler.NewMockEnvironmentScaler(ctrl),
	}

	return logic, ctrl
}

func (l *TestLogic) AddTags(t *testing.T, tags []*models.Tag) {
	for _, tag := range tags {
		if err := l.TagStore.Insert(*tag); err != nil {
			t.Fatal(err)
		}
	}
}

func (l *TestLogic) AddJobs(t *testing.T, jobs []*models.Job) {
	for _, job := range jobs {
		if err := l.JobStore.Insert(job); err != nil {
			t.Fatal(err)
		}
	}
}

func (l *TestLogic) AssertTagExists(t *testing.T, tag models.Tag) {
	tags, err := l.TagStore.SelectByTypeAndID(tag.EntityType, tag.EntityID)
	if err != nil {
		t.Fatal(err)
	}

	exists := tags.Any(func(t models.Tag) bool {
		return t.Key == tag.Key && t.Value == tag.Value
	})

	if !exists {
		t.Fatalf("Tag '%#v' does not exist in JobStore", tag)
	}
}

func (l *TestLogic) Logic() Logic {
	return *NewLogic(l.TagStore, l.JobStore, l.Backend, l.Scaler)
}
