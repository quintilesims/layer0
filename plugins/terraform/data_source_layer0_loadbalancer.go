package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/cli/client"
)

func dataSourcelayer0LoadBalancer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcelayer0LoadBalancerRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environment_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"healthcheck_healthy_threshold": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"healthcheck_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"healthcheck_target": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"healthcheck_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"healthcheck_unhealthy_threshold": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourcelayer0LoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(client.Client)

	lbName := d.Get("name").(string)
	environmentID := d.Get("environment_id").(string)

	loadbalancerID, err := resolveTags(client, lbName, map[string]string{
		"type":           "load_balancer",
		"environment_id": environmentID,
	})
	if err != nil {
		return err
	}

	loadbalancer, err := client.GetLoadBalancer(loadbalancerID)
	if err != nil {
		return err
	}

	d.SetId(loadbalancer.LoadBalancerID)
	return setResourceData(d.Set, map[string]interface{}{
		"id":                              loadbalancer.LoadBalancerID,
		"name":                            loadbalancer.LoadBalancerName,
		"private":                         !loadbalancer.IsPublic,
		"url":                             loadbalancer.URL,
		"service_id":                      loadbalancer.ServiceID,
		"service_name":                    loadbalancer.ServiceName,
		"environment_id":                  loadbalancer.EnvironmentID,
		"environment_name":                loadbalancer.EnvironmentName,
		"healthcheck_healthy_threshold":   loadbalancer.HealthCheck.HealthyThreshold,
		"healthcheck_interval":            loadbalancer.HealthCheck.Interval,
		"healthcheck_target":              loadbalancer.HealthCheck.Target,
		"healthcheck_timeout":             loadbalancer.HealthCheck.Timeout,
		"healthcheck_unhealthy_threshold": loadbalancer.HealthCheck.UnhealthyThreshold,
	})
}
