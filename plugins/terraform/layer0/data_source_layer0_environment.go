package layer0

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/client"
)

func dataSourceLayer0Environment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLayer0EnvironmentRead,

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
			"cluster_count": {
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
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLayer0EnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	query := url.Values{}
	query.Set(client.TagQueryParamType, "environment")
	query.Set(client.TagQueryParamName, d.Get("name").(string))

	tags, err := apiClient.ListTags(query)
	if err != nil {
		return err
	}

	if len(tags) == 0 {
		return fmt.Errorf("No environment found matching the specified name")
	}

	environmentID := tags[0].EntityID
	environment, err := apiClient.ReadEnvironment(environmentID)
	if err != nil {
		return err
	}

	d.SetId(environment.EnvironmentID)
	d.Set("name", environment.EnvironmentName)
	d.Set("size", environment.InstanceSize)
	d.Set("cluster_count", environment.ClusterCount)
	d.Set("security_group_id", environment.SecurityGroupID)
	d.Set("os", environment.OperatingSystem)
	d.Set("ami", environment.AMIID)

	return nil
}
