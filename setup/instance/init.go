package instance

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/setup/docker"
	"github.com/quintilesims/layer0/setup/terraform"
	"os"
)

func (l *LocalInstance) Init(dockerInputPath string, inputOverrides map[string]interface{}) error {
	if err := os.MkdirAll(l.Dir, 0700); err != nil {
		return err
	}

	// load terraform config from ~/.layer0/<instance>/main.tf.json, or create a new one
	config, err := l.loadMainConfig()
	if err != nil {
		return err
	}

	// add/update the inputs of the terraform config
	if err := l.setMainModuleInputs(config, inputOverrides); err != nil {
		return err
	}

	// save the terraform config as ~/.layer0/<instance>/main.tf.json
	path := fmt.Sprintf("%s/main.tf.json", l.Dir)
	if err := terraform.WriteConfig(path, config); err != nil {
		return err
	}

	// create/write ~/.layer0/<instance>/outputs.tf.json
	output := &terraform.Config{
		Outputs: MainModuleOutputs,
	}

	outPath := fmt.Sprintf("%s/outputs.tf.json", l.Dir)
	if err := terraform.WriteConfig(outPath, output); err != nil {
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

	// create/write ~/.layer/<instance>/dockercfg.json
	if err := l.createOrWriteDockerCFG(dockerInputPath); err != nil {
		return err
	}

	return nil
}

func (l *LocalInstance) loadMainConfig() (*terraform.Config, error) {
	path := fmt.Sprintf("%s/main.tf.json", l.Dir)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return terraform.LoadConfig(path)
	}

	return terraform.NewConfig(), nil
}

func (l *LocalInstance) setMainModuleInputs(config *terraform.Config, inputOverrides map[string]interface{}) error {
	// create the 'main' module if it doesn't already exist
	if _, ok := config.Modules["main"]; !ok {
		config.Modules["main"] = terraform.Module{}
	}

	module := config.Modules["main"]
	for _, input := range MainModuleInputs {
		// if the input has a static value, it should always be set as the static value
		if input.StaticValue != nil {
			module[input.Name] = input.StaticValue
			continue
		}

		// if the user specified the input with a cli flag, use it
		if v, ok := inputOverrides[input.Name]; ok {
			module[input.Name] = v
			continue
		}

		// prompt the user for a new/updated input
		v, err := input.Prompt(module[input.Name])
		if err != nil {
			return err
		}

		module[input.Name] = v
	}

	// the 'name' input is always the name of the layer0 instance
	module["name"] = l.Name
	return nil
}

func (l *LocalInstance) createOrWriteDockerCFG(dockerInputPath string) error {
	dockerOutputPath := fmt.Sprintf("%s/dockercfg.json", l.Dir)

	// if user didn't specify a dockercfg, create an empty one if it doesn't already exist
	if dockerInputPath == "" {
		if _, err := os.Stat(dockerOutputPath); os.IsNotExist(err) {
			text := "No docker config specified. Please run "
			text += fmt.Sprintf("`l0-setup init --docker-path=<path/to/config.json> %s` ", l.Name)
			text += "if you would like to add private registry authentication."
			logrus.Warningf(text)

			return docker.WriteConfig(dockerOutputPath, docker.NewConfig())
		}

		return nil
	}

	config, err := docker.LoadConfig(dockerInputPath)
	if err != nil {
		return err
	}

	return docker.WriteConfig(dockerOutputPath, config)
}
