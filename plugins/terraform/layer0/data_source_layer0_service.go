package layer0

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/models"
)

func dataSourcelayer0Service() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLayer0ServiceRead,
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

func dataSourceLayer0ServiceRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	query := url.Values{}
	query.Set(models.TagQueryParamType, "service")
	query.Set(models.TagQueryParamName, d.Get("name").(string))
	query.Set(models.TagQueryParamEnvironmentID, d.Get("environment_id").(string))

	tags, err := apiClient.ListTags(query)
	if err != nil {
		return err
	}

	if len(tags) == 0 {
		return fmt.Errorf("No service found matching %#v", query)
	}

	serviceID := tags[0].EntityID
	service, err := apiClient.ReadService(serviceID)
	if err != nil {
		return err
	}

	d.SetId(service.ServiceID)
	d.Set("name", service.ServiceName)
	d.Set("environment_id", service.EnvironmentID)
	d.Set("environment_name", service.EnvironmentName)
	d.Set("scale", service.DesiredCount)

	return nil
}
