package docker

import (
	"encoding/json"
	"io/ioutil"
)

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config *Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if config.Auths == nil {
		config.Auths = map[string]Auth{}
	}

	return config, nil
}

func WriteConfig(path string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return err
	}

	return nil
}
