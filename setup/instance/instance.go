package instance

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/pkg/homedir"
	"github.com/quintilesims/layer0/setup/terraform"
	"io/ioutil"
	"os"
)

type Instance struct {
	Name      string
	Dir       string
	Terraform *terraform.Terraform
}

func NewInstance(name string) *Instance {
	return &Instance{
		Name: name,
		Dir:  fmt.Sprintf("%s/.layer0/%s", homedir.Get(), name),
	}
}

func (i *Instance) Apply() error {
	return i.Terraform.Apply(i.Dir)
}

func (i *Instance) Init() error {
	if err := os.MkdirAll(i.Dir, 0700); err != nil {
		return err
	}

	config := terraform.Config{
		Modules: map[string]terraform.Module{
			"layer0": {
				Source: "/home/ec2-user/go/src/github.com/quintilesims/layer0/setup/module",
				Inputs: map[string]string{
					"name": i.Name,
				},
			},
		},
	}

	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/main.tf.json", i.Dir)
	if err := ioutil.WriteFile(path, []byte(data), 0600); err != nil {
		return err
	}

	return i.Terraform.Get(i.Dir)
}
