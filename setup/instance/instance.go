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

func (i *Instance) Init(c *cli.Context) error {
	if err := os.MkdirAll(i.Dir, 0700); err != nil {
		return err
	}

	config := terraform.NewConfig()

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
	for _, input := range MainModuleInputs {
		v, err := input.Load(c, module[input.Name])
		if err != nil {
			return err
		}

		module[input.Name] = v
	}

	if err := terraform.WriteConfig(path, config); err != nil {
		return err
	}

	// todo: write outputs.tf
	output := &terraform.Config{
		Outputs: MainModuleOutputs,
	}

	outPath := fmt.Sprintf("%s/outputs.tf.json", i.Dir)
	if err := terraform.WriteConfig(outPath, output); err != nil {
		return err
	}

	// run 'terraform get' and `terraform fmt`
	if err := i.Terraform.Get(i.Dir); err != nil {
		return err
	}

	if err := i.Terraform.FMT(i.Dir); err != nil {
		return err
	}

	return nil
}
