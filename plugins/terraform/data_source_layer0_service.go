package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourcelayer0Service() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcelayer0ServiceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"environment_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scale": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourcelayer0ServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Layer0Client)

	serviceName := d.Get("name").(string)
	environmentID := d.Get("environment_id").(string)
	params := map[string]string{
		"environment_id": environmentID,
	}

	serviceID, err := resolveTags(client, serviceName, "service", params)
	if err != nil {
		return err
	}

	service, err := client.API.GetService(serviceID)
	if err != nil {
		return err
	}

	d.SetId(service.ServiceID)

	return setResourceData(d.Set, map[string]interface{}{
		"name":             service.ServiceName,
		"environment_id":   service.EnvironmentID,
		"environment_name": service.EnvironmentName,
		"scale":            service.DesiredCount,
	})
}
