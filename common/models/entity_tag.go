package models

type EntityTag struct {
	EntityID   string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	Key        string `json:"key"`
	Value      string `json:"value"`
}
