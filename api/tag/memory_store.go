package tag

import (
	"github.com/quintilesims/layer0/common/models"
)

type MemoryStore struct {
	tags models.Tags
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tags: models.Tags{},
	}
}

func (m *MemoryStore) Init() error {
	return nil
}

func (m *MemoryStore) Tags() models.Tags {
	return m.tags
}

func (m *MemoryStore) Delete(entityType, entityID, key string) error {
	for i := 0; i < len(m.tags); i++ {
		tag := m.tags[i]
		if tag.EntityType == entityType && tag.EntityID == entityID && tag.Key == key {
			m.tags = append(m.tags[:i], m.tags[i+1:]...)
			i--
		}
	}

	return nil
}

func (m *MemoryStore) Insert(tag models.Tag) error {
	m.tags = append(m.tags, tag)
	return nil
}

func (m *MemoryStore) SelectByType(entityType string) (models.Tags, error) {
	return m.tags.WithType(entityType), nil
}

func (m *MemoryStore) SelectByTypeAndID(entityType, entityID string) (models.Tags, error) {
	return m.tags.WithType(entityType).WithID(entityID), nil
}
