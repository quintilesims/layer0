package models

import swagger "github.com/zpatrick/go-plugin-swagger"

type CreateEntityResponse struct {
	EntityID string `json:"entity_id"`
}

func (c CreateEntityResponse) Definition() swagger.Definition {
	return swagger.Definition{
		Type: "object",
		Properties: map[string]swagger.Property{
			"entity_id": swagger.NewStringProperty(),
		},
	}
}
