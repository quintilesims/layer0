package terraform

import (
	"encoding/json"
)

type Config struct {
	Modules map[string]Module `json:"module"`
}

type Module struct {
	Source string            `json:"source"`
	Inputs map[string]string `json:"-"`
}

func (m Module) MarshalJSON() ([]byte, error) {
	v := map[string]string{
		"source": m.Source,
	}

	for key, val := range m.Inputs {
		v[key] = val
	}

	return json.Marshal(v)
}
