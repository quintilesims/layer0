package instance

import (
	"github.com/quintilesims/layer0/setup/terraform"
)

var MainModuleOutputs = map[string]terraform.Output{
	"name": {Value: "${module.main.name}"},
}
