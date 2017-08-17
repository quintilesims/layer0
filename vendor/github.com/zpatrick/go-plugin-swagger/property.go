package swagger

import (
	"fmt"
)

type Property struct {
	Type    string      `json:"type,omitempty"`
	Format  string      `json:"format,omitempty"`
	Items   *Items      `json:"items,omitempty"`
	Ref     string      `json:"$ref,omitempty"`
	Default interface{} `json:"default,omitempty"`
}

func NewIntProperty() Property {
	return Property{
		Type: "integer",
	}
}

func NewIntSliceProperty() Property {
	return Property{
		Type: "array",
		Items: &Items{
			Type: "integer",
		},
	}
}

func NewStringProperty() Property {
	return Property{
		Type: "string",
	}
}

func NewStringSliceProperty() Property {
	return Property{
		Type: "array",
		Items: &Items{
			Type: "string",
		},
	}
}

func NewBoolProperty() Property {
	return Property{
		Type: "boolean",
	}
}

func NewObjectProperty(name string) Property {
	return Property{
		Ref: fmt.Sprintf("#/definitions/%s", name),
	}
}

func NewObjectSliceProperty(name string) Property {
	return Property{
		Type: "array",
		Items: &Items{
			Ref: fmt.Sprintf("#/definitions/%s", name),
		},
	}
}
