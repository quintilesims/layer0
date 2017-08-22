package consul

import (
	"context"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/backend"
	"github.com/hashicorp/terraform/helper/schema"
)

// New creates a new backend for Consul remote state.
func New() backend.Backend {
	s := &schema.Backend{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path to store state in Consul",
			},

			"access_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Access token for a Consul ACL",
				Default:     "", // To prevent input
			},

			"address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Address to the Consul Cluster",
				Default:     "", // To prevent input
			},

			"scheme": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Scheme to communicate to Consul with",
				Default:     "", // To prevent input
			},

			"datacenter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Datacenter to communicate with",
				Default:     "", // To prevent input
			},

			"http_auth": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "HTTP Auth in the format of 'username:password'",
				Default:     "", // To prevent input
			},

			"gzip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Compress the state data using gzip",
				Default:     false,
			},
		},
	}

	result := &Backend{Backend: s}
	result.Backend.ConfigureFunc = result.configure
	return result
}

type Backend struct {
	*schema.Backend

	configData *schema.ResourceData
}

func (b *Backend) configure(ctx context.Context) error {
	// Grab the resource data
	b.configData = schema.FromContextBackendConfig(ctx)

	// Initialize a client to test config
	_, err := b.clientRaw()
	return err
}

func (b *Backend) clientRaw() (*consulapi.Client, error) {
	data := b.configData

	// Configure the client
	config := consulapi.DefaultConfig()
	if v, ok := data.GetOk("access_token"); ok && v.(string) != "" {
		config.Token = v.(string)
	}
	if v, ok := data.GetOk("address"); ok && v.(string) != "" {
		config.Address = v.(string)
	}
	if v, ok := data.GetOk("scheme"); ok && v.(string) != "" {
		config.Scheme = v.(string)
	}
	if v, ok := data.GetOk("datacenter"); ok && v.(string) != "" {
		config.Datacenter = v.(string)
	}
	if v, ok := data.GetOk("http_auth"); ok && v.(string) != "" {
		auth := v.(string)

		var username, password string
		if strings.Contains(auth, ":") {
			split := strings.SplitN(auth, ":", 2)
			username = split[0]
			password = split[1]
		} else {
			username = auth
		}

		config.HttpAuth = &consulapi.HttpBasicAuth{
			Username: username,
			Password: password,
		}
	}

	return consulapi.NewClient(config)
}
