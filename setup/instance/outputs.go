package instance

import (
	"github.com/quintilesims/layer0/setup/terraform"
)

const (
	OUTPUT_NAME     = "name"
	OUTPUT_ENDPOINT = "endpoint"
	OUTPUT_TOKEN    = "token"
)

var MainModuleOutputs = map[string]terraform.Output{
	OUTPUT_NAME:     {Value: "${module.main.name}"},
	OUTPUT_ENDPOINT: {Value: "TODO!"},
	OUTPUT_TOKEN:    {Value: "TODO!"},
}
