package layer0

import (
	"fmt"
	"os"
)

type Context interface {
	LoadProfile(name string) (Profile, error)
	LoadInstance(name string) (Instance, error)
	Init(InstanceConfig) (Instance, error)
}

type LocalContext struct {
}

func NewLocalContext() *LocalContext {
	return &LocalContext{}
}

func (c *LocalContext) LoadProfile(name string) (Profile, error) {
	return nil, nil
}

func (c *LocalContext) LoadInstance(name string) (Instance, error) {
	return nil, nil
}

func (c *LocalContext) Init(config InstanceConfig) (Instance, error) {
	// todo: error if already exists
	// "If you'd like to update an existing instance, please use `l0-setup config`"

	requiredVars := map[string]*string{
		"Name":           &config.Name,
		"AWS Access Key": &config.AccessKey,
		"AWS Secret Key": &config.SecretKey,
	}

	for name, val := range requiredVars {
		if *val == "" {
			if err := scanInput(name, val); err != nil {
				return nil, err
			}
		}
	}

	instance := NewLayer0Instance(config.Name)
	if err := os.MkdirAll(instance.Dir(), 0700); err != nil {
		return nil, err
	}

	// create profile
	// create main.tf
	// run terraform get
	// run terraform backend configure

	return instance, nil
}

func scanInput(display string, p *string) error {
	for i := 0; i < 5; i++ {
		fmt.Printf("%s: ", display)
		fmt.Scanln(p)

		if *p != "" {
			return nil
		}
	}

	return fmt.Errorf("Failed to get input for variable '%s'", display)
}
