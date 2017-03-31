package instance

import (
	"fmt"
	"github.com/docker/docker/pkg/homedir"
	"os"
)

type Instance interface {
	Name() string
	Init() error
	Exists() (bool, error)
	Apply() error
	Destroy() error
}

type Layer0Instance struct {
	name string
	dir  string
}

func NewLayer0Instance(name string) *Layer0Instance {
	return &Layer0Instance{
		name: name,
		dir:  fmt.Sprintf("%s/.layer0/%s", homedir.Get(), name),
	}
}

func (l *Layer0Instance) Name() string {
	return l.name
}

// todo: where to get credentials to check remote?
func (l *Layer0Instance) Exists() (bool, error) {
	return false, nil
}

func (l *Layer0Instance) Init() error {
	if err := os.MkdirAll(l.dir, 0700); err != nil {
		return err
	}

	// todo: get + convert dockercfg
	// todo: create l.dir/main.tf
	// todo: run terraform get

	return nil
}

func (l *Layer0Instance) Apply() error {
	return nil
}

func (l *Layer0Instance) Destroy() error {
	return nil
}
