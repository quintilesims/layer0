package docker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config *Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal %s: %v", path, err)
	}

	if config.Auths == nil {
		config.Auths = map[string]Auth{}
	}

	return config, nil
}

func WriteConfig(path string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return fmt.Errorf("Failed to marshal docker config: %v", err)
	}

	if err := ioutil.WriteFile(path, data, 0600); err != nil {
		return err
	}

	return nil
}

func CopyCredentialHelper(inputPath, outputPath string) error {
	data, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(outputPath, data, 0600); err != nil {
		return err
	}

	return nil
}
