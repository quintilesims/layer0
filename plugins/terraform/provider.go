package main

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/quintilesims/layer0/cli/client"
)

var defaultTimeout = time.Minute * 15

func Provider() terraform.ResourceProvider {
	var p *schema.Provider
	p = &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Layer0 API endpoint.",
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Layer0 authentication token.",
			},
			"skip_ssl_verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Skip SSL Verification",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"layer0_deploy":           resourceLayer0Deploy(),
			"layer0_environment":      resourceLayer0Environment(),
			"layer0_environment_link": resourceLayer0EnvironmentLink(),
			"layer0_load_balancer":    resourceLayer0LoadBalancer(),
			"layer0_service":          resourceLayer0Service(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"layer0_api":           dataSourceLayer0API(),
			"layer0_environment":   dataSourceLayer0Environment(),
			"layer0_load_balancer": dataSourcelayer0LoadBalancer(),
			"layer0_deploy":        dataSourceLayer0Deploy(),
			"layer0_service":       dataSourcelayer0Service(),
		},
		ConfigureFunc: providerConfigure(p),
	}

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		layer0Client := &Layer0Client{
			API: client.NewAPIClient(client.Config{
				Endpoint:  d.Get("endpoint").(string),
				Token:     d.Get("token").(string),
				VerifySSL: !d.Get("skip_ssl_verify").(bool),
			}),
			StopContext: p.StopContext(),
		}

		return layer0Client, nil
	}
}
