package terraform

import (
	"encoding/json"
	"io/ioutil"
)

func LoadTFVars(path string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var tfvars map[string]interface{}
	if err := json.Unmarshal(data, &tfvars); err != nil {
		return nil, err
	}

	return tfvars, nil
}

func WriteTFVars(path string, tfvars map[string]interface{}) error {
	data, err := json.MarshalIndent(tfvars, "", "   ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}
