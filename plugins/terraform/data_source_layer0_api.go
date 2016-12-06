package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/cli/client"
)

func dataSourceLayer0API() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLayer0APIRead,

		Schema: map[string]*schema.Schema{
			"prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
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
	client := meta.(*client.APIClient)

	apiConfig, err := client.GetConfig()
	if err != nil {
		return err
	}

	d.SetId(apiConfig.Prefix)
	d.Set("prefix", apiConfig.Prefix)
	d.Set("vpc_id", apiConfig.VPCID)
	d.Set("public_subnets", apiConfig.PublicSubnets)
	d.Set("private_subnets", apiConfig.PrivateSubnets)

	return nil
}
