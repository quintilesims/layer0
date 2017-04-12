package tag_store

import (
	"github.com/quintilesims/layer0/common/models"
	"time"
)

type MemoryTagStore struct {
	tags []*models.Tag
}

func NewMemoryTagStore() *MemoryTagStore {
	return &MemoryTagStore{
		tags: []*models.Tag{},
	}
}

func (m *MemoryTagStore) Init() error {
	return nil
}

func (m *MemoryTagStore) Insert(tag *models.Tag) error {
	tag.TagID = time.Now().UnixNano()
	m.tags = append(m.tags, tag)
	return nil
}

func (m *MemoryTagStore) Delete(tagID int64) error {
	for i := 0; i < len(m.tags); i++ {
		if m.tags[i].TagID == tagID {
			m.tags = append(m.tags[:i], m.tags[i+1:]...)
			i--
		}
	}

	return nil
}

func (m *MemoryTagStore) SelectAll() (models.Tags, error) {
	return m.tags, nil
}

func (m *MemoryTagStore) SelectByQuery(entityType, entityID string) (models.Tags, error) {
	tags := []*models.Tag{}
	for _, tag := range m.tags {
		ok := true
		if entityType != "" && tag.EntityType != entityType {
			ok = false
		}

		if entityID != "" && tag.EntityID != entityID {
			ok = false
		}

		if ok {
			tags = append(tags, tag)
		}
	}

	return tags, nil
}
