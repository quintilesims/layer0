package instance

import (
	"fmt"
	"github.com/docker/docker/pkg/homedir"
	"github.com/quintilesims/layer0/setup/terraform"
	"os"
	"regexp"
)

type LocalInstance struct {
	Name      string
	Dir       string
	Terraform *terraform.Terraform
}

func NewLocalInstance(name string) Instance {
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

func (l *LocalInstance) validateInstanceName() error {
	re := regexp.MustCompile("^[a-z][a-z0-9]{0,15}$")
	if !re.MatchString(l.Name) {
		text := "INSTANCE argument violated one or more of the following constraints: \n"
		text += "1. Start with a lowercase letter \n"
		text += "2. Only contain lowercase alphanumeric characters \n"
		text += "3. Be between 1 and 15 characters in length \n"
		return fmt.Errorf(text)
	}

	return nil
}
