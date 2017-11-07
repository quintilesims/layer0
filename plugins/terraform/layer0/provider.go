package layer0

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/quintilesims/layer0/client"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The endpoint of the Layer0 API",
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The endpoint of the Layer0 API",
			},
			"skip_ssl_verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set, will skip SSL verification",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"layer0_deploy":        resourceLayer0Deploy(),
			"layer0_environment":   resourceLayer0Environment(),
			"layer0_load_balancer": resourceLayer0LoadBalancer(),
			"layer0_service":       resourceLayer0Service(),
			// todo: environment link
		},
		DataSourcesMap: map[string]*schema.Resource{
			"layer0_api":         dataSourceLayer0API(),
			"layer0_deploy":      dataSourceLayer0Deploy(),
			"layer0_environment": dataSourceLayer0Environment(),
			//			"layer0_load_balancer": dataSourcelayer0LoadBalancer(),
			"layer0_service": dataSourcelayer0Service(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiClient := client.NewAPIClient(client.Config{
		Endpoint:  d.Get("endpoint").(string),
		Token:     d.Get("token").(string),
		VerifySSL: !d.Get("skip_ssl_verify").(bool),
	})

	return apiClient, nil
}
