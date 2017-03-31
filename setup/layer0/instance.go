package layer0

import (
	"fmt"
	"github.com/docker/docker/pkg/homedir"
)

type Instance interface {
	Name() string
	Dir() string
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

func (l *Layer0Instance) Dir() string {
	return fmt.Sprintf("%s/.layer0/%s", homedir.Get(), l.name)
}
