package instance

import (
	"fmt"

	"github.com/quintilesims/layer0/setup/terraform"
)

func (l *LocalInstance) Unset(key string) error {
	if err := l.assertExists(); err != nil {
		return err
	}

	// load terraform config from ~/.layer0/<instance>/main.tf.json
	config, err := l.loadLayer0Config()
	if err != nil {
		return err
	}

	// create the 'layer0' module if it doesn't already exist
	if _, ok := config.Modules["layer0"]; !ok {
		config.Modules["layer0"] = terraform.Module{}
	}

	module := config.Modules["layer0"]
	if _, exists := module[key]; !exists {
		return fmt.Errorf("Key '%s' does not exist", key)
	}

	delete(module, key)

	// save the terraform config as ~/.layer0/<instance>/main.tf.json
	path := fmt.Sprintf("%s/main.tf.json", l.Dir)
	if err := terraform.WriteConfig(path, config); err != nil {
		return err
	}

	return nil
}
