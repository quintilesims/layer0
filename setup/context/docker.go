package context

import (
	"encoding/json"
	"io/ioutil"
)

type Auth struct {
        Auth  string `json:"auth,omitempty"`
        Email string `json:"email,ommitempty"`
}

type DockerConfigFile struct {
	Auths map[string]Auth `json:"auths,omitempty"`
}

func (d *DockerConfigFile) Write(path string) error {
	data, err := json.MarshalIndent(d.Auths, "", "\t")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, data, 0600); err != nil {
		return err
	}

	return nil
}
