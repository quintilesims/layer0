package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceLayer0Environment() *schema.Resource {
	return &schema.Resource{
		Read: datasourceLayer0EnvironmentRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"target_cap_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"os": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ami": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceLayer0EnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Layer0Client)

	environmentName := d.Get("name").(string)

	environmentID, err := resolveTags(client, environmentName, "environment", map[string]string{})
	if err != nil {
		return err
	}

	environment, err := client.API.GetEnvironment(environmentID)
	if err != nil {
		return err
	}

	d.SetId(environment.EnvironmentID)

	return setResourceData(d.Set, map[string]interface{}{
		"size":      environment.InstanceSize,
		"min_count": environment.ClusterCount,
		"os":        environment.OperatingSystem,
		"ami":       environment.AMIID,
	})
}
