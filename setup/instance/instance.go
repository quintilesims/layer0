package instance

import (
	"fmt"
	"github.com/docker/docker/pkg/homedir"
	"github.com/quintilesims/layer0/setup/terraform"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
)

type Instance struct {
	Name      string
	Dir       string
	TFVarPath string
	Terraform *terraform.Terraform
}

func NewInstance(name string) *Instance {
	dir := fmt.Sprintf("%s/.layer0/%s", homedir.Get(), name)

	return &Instance{
		Name:      name,
		Dir:       dir,
		TFVarPath: fmt.Sprintf("%s/terraform.tfvars.json", dir),
	}
}

func (i *Instance) Apply() error {
	return i.Terraform.Apply(i.Dir)
}

func (i *Instance) Init(c *cli.Context) error {
	if err := os.MkdirAll(i.Dir, 0700); err != nil {
		return err
	}

	// load variables from the cli, user, and/or tfvars
	tfvars := map[string]interface{}{}
	for _, schema := range InstanceVariableSchemas {
		v, err := schema.Load(c, i.TFVarPath)
		if err != nil {
			return err
		}

		tfvars[schema.Name] = v
	}

	// write the updated tfvars
	if err := terraform.WriteTFVars(i.TFVarPath, tfvars); err != nil {
		return err
	}

	// write variables.tf, main.tf, and outputs.tf if they don't already exist
	for fileName, data := range TFFiles {
		path := fmt.Sprintf("%s/%s", i.Dir, fileName)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := ioutil.WriteFile(path, data, 0644); err != nil {
				return err
			}
		}
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
