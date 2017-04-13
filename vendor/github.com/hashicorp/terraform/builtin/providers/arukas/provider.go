package arukas

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(JSONTokenParamName, nil),
				Description: "your Arukas APIKey(token)",
			},
			"secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(JSONSecretParamName, nil),
				Description: "your Arukas APIKey(secret)",
			},
			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(JSONUrlParamName, "https://app.arukas.io/api/"),
				Description: "default Arukas API url",
			},
			"trace": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(JSONDebugParamName, ""),
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(JSONTimeoutParamName, "600"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"arukas_container": resourceArukasContainer(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	config := Config{
		Token:   d.Get("token").(string),
		Secret:  d.Get("secret").(string),
		URL:     d.Get("api_url").(string),
		Trace:   d.Get("trace").(string),
		Timeout: d.Get("timeout").(int),
	}

	return config.NewClient()
}
