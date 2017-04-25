package docker

import (
	"encoding/json"
)

type Config struct {
	Auths map[string]Auth `json:"auths,omitempty"`
}

type Auth map[string]interface{}

func NewConfig() *Config {
	return &Config{
		Auths: map[string]Auth{},
	}
}

// this overrides the default json.Marshal function
// we only marshal config.Auths to match the older 'dockercfg' format
// see: http://docs.aws.amazon.com/AmazonECS/latest/developerguide/private-auth.html
func (c *Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Auths)
}
