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
	} else if len(config.Auths) == 0 {
		fmt.Println("[WARNING] Even though you have specified a path to a docker config file, " +
			"it does not contain any auth information. If you need to add auth information " +
			"to the docker config file, you can do so and re-run the l0-setup init command to " +
			"include private registry authentication.\n")

		fmt.Println("Press 'enter' to continue without private registry authentication: ")

		var input string
		fmt.Scanln(&input)
	} else if config.CredsStore != "" {
		fmt.Printf("[WARNING] Even though you have specified a path to a docker config file, "+
			"the config file is using a credential store '%s' to cache the credentials. This "+
			"means the credentials for private registry authentication aren't in the docker "+
			"config file. You can either not use a credential store by removing the 'credsStore' "+
			"section or add a 'credHelpers' section and exclude your private docker repository "+
			"so that the private registry credentials, for just your repository are stored in "+
			"the file you have specified.\n\n",
			config.CredsStore)

		fmt.Println("Press 'enter' to continue without private registry authentication: ")

		var input string
		fmt.Scanln(&input)
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
