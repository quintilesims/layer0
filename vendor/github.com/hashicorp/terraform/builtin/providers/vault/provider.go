package vault

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/hashicorp/vault/api"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VAULT_ADDR", nil),
				Description: "URL of the root of the target Vault server.",
			},
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VAULT_TOKEN", nil),
				Description: "Token to use to authenticate to Vault.",
			},
			"ca_cert_file": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ca_cert_dir"},
				DefaultFunc:   schema.EnvDefaultFunc("VAULT_CACERT", nil),
				Description:   "Path to a CA certificate file to validate the server's certificate.",
			},
			"ca_cert_dir": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ca_cert_file"},
				DefaultFunc:   schema.EnvDefaultFunc("VAULT_CAPATH", nil),
				Description:   "Path to directory containing CA certificate files to validate the server's certificate.",
			},
			"client_auth": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Client authentication credentials.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert_file": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc("VAULT_CLIENT_CERT", nil),
							Description: "Path to a file containing the client certificate.",
						},
						"key_file": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc("VAULT_CLIENT_KEY", nil),
							Description: "Path to a file containing the private key that the certificate was issued for.",
						},
					},
				},
			},
			"skip_tls_verify": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VAULT_SKIP_VERIFY", nil),
				Description: "Set this to true only if the target Vault server is an insecure development instance.",
			},
			"max_lease_ttl_seconds": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,

				// Default is 20min, which is intended to be enough time for
				// a reasonable Terraform run can complete but not
				// significantly longer, so that any leases are revoked shortly
				// after Terraform has finished running.
				DefaultFunc: schema.EnvDefaultFunc("TERRAFORM_VAULT_MAX_TTL", 1200),

				Description: "Maximum TTL for secret leases requested by this provider",
			},
		},

		ConfigureFunc: providerConfigure,

		DataSourcesMap: map[string]*schema.Resource{
			"vault_generic_secret": genericSecretDataSource(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"vault_generic_secret": genericSecretResource(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &api.Config{
		Address: d.Get("address").(string),
	}

	clientAuthI := d.Get("client_auth").([]interface{})
	if len(clientAuthI) > 1 {
		return nil, fmt.Errorf("client_auth block may appear only once")
	}

	clientAuthCert := ""
	clientAuthKey := ""
	if len(clientAuthI) == 1 {
		clientAuth := clientAuthI[0].(map[string]interface{})
		clientAuthCert = clientAuth["cert_file"].(string)
		clientAuthKey = clientAuth["key_file"].(string)
	}

	config.ConfigureTLS(&api.TLSConfig{
		CACert:   d.Get("ca_cert_file").(string),
		CAPath:   d.Get("ca_cert_dir").(string),
		Insecure: d.Get("skip_tls_verify").(bool),

		ClientCert: clientAuthCert,
		ClientKey:  clientAuthKey,
	})

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to configure Vault API: %s", err)
	}

	// In order to enforce our relatively-short lease TTL, we derive a
	// temporary child token that inherits all of the policies of the
	// token we were given but expires after max_lease_ttl_seconds.
	//
	// The intent here is that Terraform will need to re-fetch any
	// secrets on each run and so we limit the exposure risk of secrets
	// that end up stored in the Terraform state, assuming that they are
	// credentials that Vault is able to revoke.
	//
	// Caution is still required with state files since not all secrets
	// can explicitly be revoked, and this limited scope won't apply to
	// any secrets that are *written* by Terraform to Vault.

	client.SetToken(d.Get("token").(string))
	renewable := false
	childTokenLease, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		DisplayName:    "terraform",
		TTL:            fmt.Sprintf("%ds", d.Get("max_lease_ttl_seconds").(int)),
		ExplicitMaxTTL: fmt.Sprintf("%ds", d.Get("max_lease_ttl_seconds").(int)),
		Renewable:      &renewable,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create limited child token: %s", err)
	}

	childToken := childTokenLease.Auth.ClientToken
	policies := childTokenLease.Auth.Policies

	log.Printf("[INFO] Using Vault token with the following policies: %s", strings.Join(policies, ", "))

	client.SetToken(childToken)

	return client, nil
}
