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
	config, err := l.loadLayer0Config()
	if err != nil {
		return err
	}

	// add/update the inputs of the terraform config
	if err := l.setLayer0ModuleInputs(config, inputOverrides); err != nil {
		return err
	}

	// save the terraform config as ~/.layer0/<instance>/main.tf.json
	path := fmt.Sprintf("%s/main.tf.json", l.Dir)
	if err := terraform.WriteConfig(path, config); err != nil {
		return err
	}

	// create/write ~/.layer0/<instance>/outputs.tf.json
	output := &terraform.Config{
		Outputs: Layer0ModuleOutputs,
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

func (l *LocalInstance) loadLayer0Config() (*terraform.Config, error) {
	path := fmt.Sprintf("%s/main.tf.json", l.Dir)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return terraform.LoadConfig(path)
	}

	return terraform.NewConfig(), nil
}

func (l *LocalInstance) setLayer0ModuleInputs(config *terraform.Config, inputOverrides map[string]interface{}) error {
	// create the 'layer0' module if it doesn't already exist
	if _, ok := config.Modules["layer0"]; !ok {
		config.Modules["layer0"] = terraform.Module{}
	}

	module := config.Modules["layer0"]
	for _, input := range Layer0ModuleInputs {
		// if the input has a static value, it should always be set as the static value
		if input.StaticValue != nil {
			module[input.Name] = input.StaticValue
			continue
		}

		// if the user specified the input with a cli flag, use it
		if v, ok := inputOverrides[input.Name]; ok {
			logrus.Infof("Using cli flag/environment variable for %s", input.Name)
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
