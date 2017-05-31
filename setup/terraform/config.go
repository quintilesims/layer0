package terraform

type Config struct {
	Variables map[string]Variable `json:"variable,omitempty"`
	Modules   map[string]Module   `json:"module,omitempty"`
	Outputs   map[string]Output   `json:"output,omitempty"`
}

func NewConfig() *Config {
	return &Config{
		Variables: map[string]Variable{},
		Modules:   map[string]Module{},
		Outputs:   map[string]Output{},
	}
}

type Variable struct {
	Type        string      `json:"type,omitempty"`
	Description string      `json:"description,omitempty"`
	Default     interface{} `json:"default,omitempty"`
}

type Module map[string]interface{}

type Output struct {
	Value interface{} `json:"value"`
}
