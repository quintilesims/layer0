package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/cli/client"
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
			"load_balancer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"load_balancer_name": {
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
	client := meta.(client.Client)

	serviceName := d.Get("name").(string)
	environmentID := d.Get("environment_id").(string)
	params := map[string]string{
		"environment_id": environmentID,
	}

	serviceID, err := resolveTags(client, serviceName, "service", params)
	if err != nil {
		return err
	}

	service, err := client.GetService(serviceID)
	if err != nil {
		return err
	}

	d.SetId(service.ServiceID)

	return setResourceData(d.Set, map[string]interface{}{
		"name":               service.ServiceName,
		"environment_id":     service.EnvironmentID,
		"environment_name":   service.EnvironmentName,
		"load_balancer_name": service.LoadBalancerName,
		"load_balancer_id":   service.LoadBalancerID,
		"scale":              service.DesiredCount,
	})
}
