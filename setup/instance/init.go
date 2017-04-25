package instance

import (
	"fmt"
	"github.com/quintilesims/layer0/setup/terraform"
	"github.com/urfave/cli"
	"os"
)

func (i *Instance) Init(c *cli.Context, inputOverrides map[string]interface{}) error {
	if err := os.MkdirAll(i.Dir, 0700); err != nil {
		return err
	}

	// load terraform config from ~/.layer0/<instance>/main.tf.json, or create a new one
	config, err := i.loadMainConfig()
	if err != nil {
		return err
	}

	// add/update the inputs of the terraform config
	if err := i.setMainModuleInputs(config, inputOverrides); err != nil {
		return err
	}

	// save the terraform config as ~/.layer0/<instance>/main.tf.json
	path := fmt.Sprintf("%s/main.tf.json", i.Dir)
	if err := terraform.WriteConfig(path, config); err != nil {
		return err
	}

	// create/write ~/.layer0/<instance>/outputs.tf.json
	output := &terraform.Config{
		Outputs: MainModuleOutputs,
	}

	outPath := fmt.Sprintf("%s/outputs.tf.json", i.Dir)
	if err := terraform.WriteConfig(outPath, output); err != nil {
		return err
	}

	// run `terraform get` to download terraform modules
	if err := i.Terraform.Get(i.Dir); err != nil {
		return err
	}

	// run `terraform fmt` to validate the terraform syntax
	if err := i.Terraform.FMT(i.Dir); err != nil {
		return err
	}

	return nil
}

func (i *Instance) loadMainConfig() (*terraform.Config, error) {
	path := fmt.Sprintf("%s/main.tf.json", i.Dir)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return terraform.LoadConfig(path)
	}

	return terraform.NewConfig(), nil
}

func (i *Instance) setMainModuleInputs(config *terraform.Config, inputOverrides map[string]interface{}) error {
	// create the 'main' module if it doesn't already exist
	if _, ok := config.Modules["main"]; !ok {
		config.Modules["main"] = terraform.Module{}
	}

	module := config.Modules["main"]
	for _, input := range MainModuleInputs {
		// if the user specified a cli flag or env var, use that for the input
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
	module["name"] = i.Name
	return nil
}
