package swagger

import (
	"fmt"
)

type Parameter struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	In          string      `json:"in"`
	Required    bool        `json:"required"`
	Type        string      `json:"type,omitempty"`
	Format      string      `json:"format,omitempty"`
	Schema      interface{} `json:"schema,omitempty"`
}

func NewIntPathParam(name, description string, required bool) Parameter {
	return Parameter{
		Name:        name,
		Description: description,
		Required:    required,
		In:          "path",
		Type:        "integer",
		Format:      "int",
	}
}

func NewStringPathParam(name, description string, required bool) Parameter {
	return Parameter{
		Name:        name,
		Description: description,
		Required:    required,
		In:          "path",
		Type:        "string",
	}
}

func NewBodyParam(name, description string, required bool) Parameter {
	return Parameter{
		Name:        name,
		Description: description,
		Required:    required,
		In:          "body",
		Schema: map[string]string{
			"$ref": fmt.Sprintf("#/definitions/%s", name),
		},
	}
}

func NewBodySliceParam(name, description string, required bool) Parameter {
	return Parameter{
		Name:        name,
		Description: description,
		Required:    required,
		In:          "body",
		Schema: map[string]interface{}{
			"type": "array",
			"items": map[string]interface{}{
				"$ref": fmt.Sprintf("#/definitions/%s", name),
			},
		},
	}
}
