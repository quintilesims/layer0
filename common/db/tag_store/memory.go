package tag_store

import (
	"github.com/quintilesims/layer0/common/models"
)

type MemoryTagStore struct {
	tags models.Tags
}

func NewMemoryTagStore() *MemoryTagStore {
	return &MemoryTagStore{
		tags: models.Tags{},
	}
}

func (m *MemoryTagStore) Init() error {
	return nil
}

func (m *MemoryTagStore) Tags() models.Tags {
	return m.tags
}

func (m *MemoryTagStore) Delete(entityType, entityID, key string) error {
	for i := 0; i < len(m.tags); i++ {
		tag := m.tags[i]
		if tag.EntityType == entityType && tag.EntityID == entityID && tag.Key == key {
			m.tags = append(m.tags[:i], m.tags[i+1:]...)
			i--
		}
	}

	return nil
}

func (m *MemoryTagStore) Insert(tag models.Tag) error {
	m.tags = append(m.tags, tag)
	return nil
}

func (m *MemoryTagStore) SelectByType(entityType string) (models.Tags, error) {
	return m.tags.WithType(entityType), nil
}

func (m *MemoryTagStore) SelectByTypeAndID(entityType, entityID string) (models.Tags, error) {
	return m.tags.WithType(entityType).WithID(entityID), nil
}
