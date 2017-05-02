package instance

import (
	"fmt"
	"github.com/quintilesims/layer0/setup/terraform"
	//"github.com/Sirupsen/logrus"
)

func (l *LocalInstance) Upgrade(version string, force bool) error {
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

	// set new input values for 'source' and 'version'
	inputValues := map[string]string{
		INPUT_SOURCE:  fmt.Sprintf("%s?ref=%v", LAYER0_MODULE_SOURCE, version),
		INPUT_VERSION: version,
	}

	module := config.Modules["layer0"]
	for input, value := range inputValues {
		if current, ok := module[input]; ok && current != value && !force {
			fmt.Printf("This will update the '%s' input \n\tFrom: [%s] \n\tTo:   [%s]\n\n", input, current, value)
			fmt.Printf("\tPress 'enter' to accept this change: ")

			var input string
			fmt.Scanln(&input)
		}

		module[input] = value
	}

	// save the terraform config as ~/.layer0/<instance>/main.tf.json
	path := fmt.Sprintf("%s/main.tf.json", l.Dir)
	if err := terraform.WriteConfig(path, config); err != nil {
		return err
	}

	// run `terraform get` to download terraform modules
	if err := l.Terraform.Get(l.Dir); err != nil {
		return err
	}

	// run `terraform fmt` to validate the terraform syntax
	if err := l.Terraform.FMT(l.Dir); err != nil {
		return err
	}

	return nil
}
