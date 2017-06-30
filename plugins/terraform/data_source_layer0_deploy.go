package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/cli/client"
)

func dataSourceLayer0Deploy() *schema.Resource {
	return &schema.Resource{
		Read: datasourceLayer0DeployRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceLayer0DeployRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(client.Client)

	deployName := d.Get("name").(string)
	version := d.Get("version").(string)
	params := map[string]string{
		"version": version,
	}

	deployID, err := resolveTags(client, deployName, "deploy", params)
	if err != nil {
		return err
	}

	deploy, err := client.GetDeploy(deployID)
	if err != nil {
		return err
	}

	d.SetId(deploy.DeployID)

	return setResourceData(d.Set, map[string]interface{}{
		"name":    deploy.DeployName,
		"version": deploy.Version,
	})
}
