package layer0

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/models"
)

func dataSourcelayer0LoadBalancer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLayer0LoadBalancerRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environment_name": {
				Type:     schema.TypeString,
				Computed: true,
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
		},
	}
}

func dataSourceLayer0LoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	query := url.Values{}
	query.Set(models.TagQueryParamType, "load_balancer")
	query.Set(models.TagQueryParamName, d.Get("name").(string))
	query.Set(models.TagQueryParamEnvironmentID, d.Get("environment_id").(string))

	tags, err := apiClient.ListTags(query)
	if err != nil {
		return err
	}

	if len(tags) == 0 {
		return fmt.Errorf("No load balancer found matching %#v", query)
	}

	loadBalancerID := tags[0].EntityID
	loadBalancer, err := apiClient.ReadLoadBalancer(loadBalancerID)
	if err != nil {
		return err
	}

	d.SetId(loadBalancer.LoadBalancerID)
	d.Set("name", loadBalancer.LoadBalancerName)
	d.Set("type", loadBalancer.LoadBalancerType)
	d.Set("environment_id", loadBalancer.EnvironmentID)
	d.Set("environment_name", loadBalancer.EnvironmentName)
	d.Set("private", !loadBalancer.IsPublic)
	d.Set("url", loadBalancer.URL)

	return nil
}
