package swagger

import (
	"fmt"
)

type Schemable interface {
	Schema() *Schema
}

type Schema struct {
	Type               string              `json:"type,omitempty"`
	RequiredProperties []string            `json:"required,omitempty"`
	Properties         map[string]Property `json:"property,omitempty"`
	Items              map[string]string   `json:"items,omitempty"`
	Ref                string              `json:"$ref,omitempty"`
}

func NewObjectSchema(name string) *Schema {
	return &Schema{
		Ref: fmt.Sprintf("#/definitions/%s", name),
	}
}

func NewObjectSliceSchema(name string) *Schema {
	return &Schema{
		Type: "array",
		Items: map[string]string{
			"$ref": fmt.Sprintf("#/definitions/%s", name),
		},
	}
}

func NewIntSchema() *Schema {
	return &Schema{
		Type: "integer",
	}
}

func NewIntSliceSchema() *Schema {
	return &Schema{
		Type: "array",
		Items: map[string]string{
			"type": "integer",
		},
	}
}
