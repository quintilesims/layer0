package models

import swagger "github.com/zpatrick/go-plugin-swagger"

type Tag struct {
	EntityID   string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	Key        string `json:"key"`
	Value      string `json:"value"`
}

func (t Tag) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"entity_id":   swagger.NewStringProperty(),
			"entity_type": swagger.NewStringProperty(),
			"key":         swagger.NewStringProperty(),
			"value":       swagger.NewStringProperty(),
		},
	}
}
