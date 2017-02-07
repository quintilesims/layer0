package context

// Utilities to load and save docker config files (config.json or dockercfg).

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	// CONFIGFILE supplies current docker config file name.
	CONFIGFILE = "config.json"
	// OLD_CONFIGFILE supplies old docker config file name.
	OLD_CONFIGFILE = "dockercfg"
)

// AuthConfig contains authorization information for connecting to a Registry.
type AuthConfig struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Auth          string `json:"auth,omitempty"`
	Email         string `json:"email,omitempty"`
	ServerAddress string `json:"serveraddress,omitempty"`
	IdentityToken string `json:"identitytoken,omitempty"`
	RegistryToken string `json:"registrytoken,omitempty"`
}

// DockerConfigFile contains `config.json` file info.
// `dockercfg` files use just the `AuthConfigs` field.
type DockerConfigFile struct {
	AuthConfigs map[string]AuthConfig `json:"auths"`
	HTTPHeaders map[string]string     `json:"HttpHeaders,omitempty"`
	filename    string
	empty       bool
}

// LoadDockerConfig loads docker config files of config.json or dockercfg formats.
func LoadDockerConfig(dockercfg string) (*DockerConfigFile, error) {
	configFile := DockerConfigFile{
		AuthConfigs: make(map[string]AuthConfig),
		filename:    filepath.Clean(dockercfg),
	}
	configFileType := filepath.Base(configFile.filename)

	if _, err := os.Stat(configFile.filename); err != nil {
		return nil, err
	}

	if configFileType == CONFIGFILE {
		file, err := os.Open(configFile.filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		if err := json.NewDecoder(file).Decode(&configFile); err != nil {
			return nil, err
		}

		return &configFile, nil
	}

	if configFileType == OLD_CONFIGFILE {
		b, err := ioutil.ReadFile(configFile.filename)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(b, &configFile.AuthConfigs); err != nil {
			return nil, err
		}
		return &configFile, nil
	}

	return nil, fmt.Errorf("Could load docker config file.")
}

// Save a docker config file as `dockercfg` in JSON format.
func (configFile *DockerConfigFile) Save(saveDir string) error {
	data, err := json.MarshalIndent(configFile.AuthConfigs, "", "\t")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(saveDir), 0600); err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(saveDir, OLD_CONFIGFILE), data, 0600); err != nil {
		return err
	}

	return nil
}
