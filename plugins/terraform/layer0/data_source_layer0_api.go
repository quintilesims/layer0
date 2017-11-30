package layer0

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/client"
)

func dataSourceLayer0API() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLayer0APIRead,
		Schema: map[string]*schema.Schema{
			"instance": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_subnets": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"private_subnets": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceLayer0APIRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	config, err := apiClient.ReadConfig()
	if err != nil {
		return err
	}

	d.SetId(config.Instance)
	d.Set("instance", config.Instance)
	d.Set("vpc_id", config.VPCID)
	d.Set("version", config.Version)
	d.Set("public_subnets", config.PublicSubnets)
	d.Set("private_subnets", config.PrivateSubnets)

	return nil
}
