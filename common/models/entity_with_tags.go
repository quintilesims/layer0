package models

type EntityWithTags struct {
	EntityID   string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	Tags       []Tag  `json:"tags"`
}
