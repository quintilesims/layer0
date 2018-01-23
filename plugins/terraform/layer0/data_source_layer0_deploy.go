package layer0

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/models"
)

func dataSourceLayer0Deploy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLayer0DeployRead,
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

func dataSourceLayer0DeployRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	query := url.Values{}
	query.Set(models.TagQueryParamType, "deploy")
	query.Set(models.TagQueryParamName, d.Get("name").(string))
	query.Set(models.TagQueryParamVersion, d.Get("version").(string))

	tags, err := apiClient.ListTags(query)
	if err != nil {
		return err
	}

	if len(tags) == 0 {
		return fmt.Errorf("No deploy found matching %#v", query)
	}

	deployID := tags[0].EntityID
	deploy, err := apiClient.ReadDeploy(deployID)
	if err != nil {
		return err
	}

	d.SetId(deploy.DeployID)
	d.Set("name", deploy.DeployName)
	d.Set("version", deploy.Version)

	return nil
}
