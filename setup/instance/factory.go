package instance

import (
	"fmt"
	"github.com/docker/docker/pkg/homedir"
	"io/ioutil"
	"os"
)

type Factory interface {
	NewInstance(name string) (Instance, error)
	ListInstances() ([]string, error)
}

type Layer0Factory struct{}

func NewLayer0Factory() *Layer0Factory {
	return &Layer0Factory{}
}

func (f *Layer0Factory) NewInstance(name string) (Instance, error) {
	return NewLayer0Instance(name), nil
}

func (f *Layer0Factory) ListInstances() ([]string, error) {
	dir := fmt.Sprintf("%s/.layer0", homedir.Get())
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	instances := []string{}
	for _, file := range files {
		if file.IsDir() {
			instances = append(instances, file.Name())
		}
	}

	return instances, nil
}
