package instance

import (
	"fmt"
	"github.com/docker/docker/pkg/homedir"
	"github.com/quintilesims/layer0/setup/terraform"
	"os"
)

type LocalInstance struct {
	Name      string
	Dir       string
	Terraform *terraform.Terraform
}

func NewLocalInstance(name string) *LocalInstance {
	dir := fmt.Sprintf("%s/.layer0/%s", homedir.Get(), name)

	return &LocalInstance{
		Name: name,
		Dir:  dir,
	}
}

func (l *LocalInstance) assertExists() error {
	path := fmt.Sprintf("%s/main.tf.json", l.Dir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		text := fmt.Sprintf("Layer0 instance '%s' does not exist locally.\n", l.Name)
		text += fmt.Sprintf("Have you tried running `l0-setup pull %s`?", l.Name)
		return fmt.Errorf(text)
	}

	return nil
}
