package terraform

type Config struct {
	Modules   map[string]Module   `json:"module,omitempty"`
	Outputs   map[string]Output   `json:"output,omitempty"`
}

func NewConfig() *Config {
	return &Config{
		Modules:   map[string]Module{},
		Outputs:   map[string]Output{},
	}
}

type Module map[string]interface{}

type Output struct {
	Value string `json:"value"`
}
