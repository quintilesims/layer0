package instance

import (
	"fmt"
	"github.com/docker/docker/pkg/homedir"
	"github.com/quintilesims/layer0/setup/terraform"
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

func (i *Instance) assertExists() error {
	path := fmt.Sprintf("%s/main.tf.json", i.Dir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		text := fmt.Sprintf("Layer0 instance '%s' does not exist locally.\n", i.Name)
		text += fmt.Sprintf("Have you tried running `l0-setup pull %s`?", i.Name)
		return fmt.Errorf(text)
	}

	return nil
}
