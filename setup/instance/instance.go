package instance

import (
	"fmt"
	"github.com/docker/docker/pkg/homedir"
	"github.com/quintilesims/layer0/setup/terraform"
	"github.com/urfave/cli"
	"os"
)

type Instance struct {
	Name      string
	Dir       string
	Terraform *terraform.Terraform
}

func NewInstance(name string) *Instance {
	dir := fmt.Sprintf("%s/.layer0/%s", homedir.Get(), name)

	return &Instance{
		Name: name,
		Dir:  dir,
	}
}

func (i *Instance) Apply() error {
	return i.Terraform.Apply(i.Dir)
}

func (i *Instance) Destroy(force bool) error {
	if err := i.Terraform.Destroy(i.Dir, force); err != nil {
		return err
	}

	return os.RemoveAll(i.Dir)
}

func (i *Instance) Init(c *cli.Context, inputOverrides map[string]interface{}) error {
	if err := os.MkdirAll(i.Dir, 0700); err != nil {
		return err
	}

	config := terraform.NewConfig()

	// load main.tf.json if it exists
	path := fmt.Sprintf("%s/main.tf.json", i.Dir)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		c, err := terraform.LoadConfig(path)
		if err != nil {
			return err
		}

		config = c
	}

	if _, ok := config.Modules["main"]; !ok {
		config.Modules["main"] = terraform.Module{}
	}

	module := config.Modules["main"]
	module["name"] = i.Name

	// configure inputs for the main module
	for _, input := range MainModuleInputs {
		if v, ok := inputOverrides[input.Name]; ok {
			module[input.Name] = v
			continue
		}

		v, err := input.Prompt(module[input.Name])
		if err != nil {
			return err
		}

		module[input.Name] = v
	}

	// write main.tf.json and outputs.tf.json
	if err := terraform.WriteConfig(path, config); err != nil {
		return err
	}

	output := &terraform.Config{
		Outputs: MainModuleOutputs,
	}

	outPath := fmt.Sprintf("%s/outputs.tf.json", i.Dir)
	if err := terraform.WriteConfig(outPath, output); err != nil {
		return err
	}

	// run `terraform get` and `terraform fmt`
	if err := i.Terraform.Get(i.Dir); err != nil {
		return err
	}

	if err := i.Terraform.FMT(i.Dir); err != nil {
		return err
	}

	return nil
}

func (i *Instance) Plan() error {
	return i.Terraform.Plan(i.Dir)
}
