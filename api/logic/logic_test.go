package logic

import (
	log "github.com/Sirupsen/logrus"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/backend/mock_backend"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/db"
	"github.com/quintilesims/layer0/common/db/mock_data"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"os"
	"testing"
)

// Main test entrypoint
func TestMain(m *testing.M) {
	config.SetTestConfig()
	log.SetLevel(log.FatalLevel)
	jobLogger.Level = log.FatalLevel
	retCode := m.Run()
	os.Exit(retCode)
}

type MockLogic struct {
	Tag         *mock_data.MockTagStore
	Job         *mock_data.MockJobStore
	EbStackName *string
	Backend     *mock_backend.MockBackend
	SQLite      *tag_store.TagStoreStoreSQLite
}

func NewMockLogic(ctrl *gomock.Controller) *MockLogic {
	name := "Pre-set a StackName to avoid Beanstalk.ListAvailableSolutionStacks"
	return &MockLogic{
		Tag:         mock_data.NewMockTagStore(ctrl),
		Job:         mock_data.NewMockJobStore(ctrl),
		Backend:     mock_backend.NewMockBackend(ctrl),
		EbStackName: &name,
	}
}

func (this *MockLogic) StubTagMock() {
	this.Tag.EXPECT().
		Find(gomock.Any(), gomock.Any()).
		Return([]models.EntityWithTags{}, nil).
		AnyTimes()

	this.Tag.EXPECT().
		SelectByQuery(gomock.Any()).
		Return([]models.EntityWithTags{}, nil).
		AnyTimes()

	this.Tag.EXPECT().
		Make(gomock.Any()).
		Return(nil).
		AnyTimes()
}

func (this *MockLogic) Logic() Logic {
	var tagData tag_store.TagStore
	if this.SQLite != nil {
		tagData = data.NewTagLogicLayer(this.SQLite)
	} else {
		tagData = this.Tag
	}

	newLogic := NewLogic(
		nil, // SQLAdmin
		tagData,
		this.Job,
		this.Backend,
	)

	newLogic.ebStackName = this.EbStackName

	return *newLogic
}

func (this *MockLogic) UseSQLite(t *testing.T) {
	dataStore, err := data.NewTagSQLiteDataStore()
	if err != nil {
		t.Fatal(err)
	}

	this.SQLite = dataStore
}

func NewGoMock(t *testing.T, name string) *gomock.Controller {
	return gomock.NewController(testutils.NewReporter(t, name))
}

func newTag(tag_key, tag_value, eid, etype string) models.EntityTag {
	return models.EntityTag{
		EntityID:   eid,
		EntityType: etype,
		Key:        tag_key,
		Value:      tag_value,
	}
}

func addTag(t *testing.T, sqlite *tag_store.TagStoreStoreSQLite, tag models.EntityTag) {
	addTags(t, sqlite, []models.EntityTag{tag})
}

func addTags(t *testing.T, sqlite *tag_store.TagStoreStoreSQLite, tags []models.EntityTag) {
	for _, tag := range tags {
		if err := sqlite.Insert(tag); err != nil {
			t.Error(err)
		}
	}
}

type Dockerrun struct {
	Version    int `json:"AWSEBDockerrunVersion"`
	Containers []struct {
		Name        string `json:"name"`
		Image       string `json:"image"`
		Essential   bool   `json:"essential"`
		Environment []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"environment"`
	} `json:"containerDefinitions"`
}
