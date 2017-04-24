package instance

import(
	 "github.com/docker/docker/pkg/homedir"
	"io/ioutil"
	"fmt"
)

func ListLocalInstances() ([]string, error){
	dir := fmt.Sprintf("%s/.layer0", homedir.Get())
	files, err := ioutil.ReadDir(dir)
	if err != nil{
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
