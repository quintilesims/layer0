package models

type Tag struct {
	TagID      int64  `json:"tag_id"`
	EntityID   string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	Key        string `json:"key"`
	Value      string `json:"value"`
}
