package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/quintilesims/layer0/plugins/terraform/layer0"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: layer0.Provider})
}
