package google

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/pathorcontents"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account_file": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("GOOGLE_ACCOUNT_FILE", nil),
				ValidateFunc: validateAccountFile,
				Deprecated:   "Use the credentials field instead",
			},

			"credentials": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_CREDENTIALS",
					"GOOGLE_CLOUD_KEYFILE_JSON",
					"GCLOUD_KEYFILE_JSON",
				}, nil),
				ValidateFunc: validateCredentials,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_PROJECT",
					"GCLOUD_PROJECT",
					"CLOUDSDK_CORE_PROJECT",
				}, nil),
			},

			"region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_REGION",
					"GCLOUD_REGION",
					"CLOUDSDK_COMPUTE_REGION",
				}, nil),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"google_iam_policy": dataSourceGoogleIamPolicy(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"google_compute_autoscaler":             resourceComputeAutoscaler(),
			"google_compute_address":                resourceComputeAddress(),
			"google_compute_backend_service":        resourceComputeBackendService(),
			"google_compute_disk":                   resourceComputeDisk(),
			"google_compute_firewall":               resourceComputeFirewall(),
			"google_compute_forwarding_rule":        resourceComputeForwardingRule(),
			"google_compute_global_address":         resourceComputeGlobalAddress(),
			"google_compute_global_forwarding_rule": resourceComputeGlobalForwardingRule(),
			"google_compute_health_check":           resourceComputeHealthCheck(),
			"google_compute_http_health_check":      resourceComputeHttpHealthCheck(),
			"google_compute_https_health_check":     resourceComputeHttpsHealthCheck(),
			"google_compute_image":                  resourceComputeImage(),
			"google_compute_instance":               resourceComputeInstance(),
			"google_compute_instance_group":         resourceComputeInstanceGroup(),
			"google_compute_instance_group_manager": resourceComputeInstanceGroupManager(),
			"google_compute_instance_template":      resourceComputeInstanceTemplate(),
			"google_compute_network":                resourceComputeNetwork(),
			"google_compute_project_metadata":       resourceComputeProjectMetadata(),
			"google_compute_region_backend_service": resourceComputeRegionBackendService(),
			"google_compute_route":                  resourceComputeRoute(),
			"google_compute_ssl_certificate":        resourceComputeSslCertificate(),
			"google_compute_subnetwork":             resourceComputeSubnetwork(),
			"google_compute_target_http_proxy":      resourceComputeTargetHttpProxy(),
			"google_compute_target_https_proxy":     resourceComputeTargetHttpsProxy(),
			"google_compute_target_pool":            resourceComputeTargetPool(),
			"google_compute_url_map":                resourceComputeUrlMap(),
			"google_compute_vpn_gateway":            resourceComputeVpnGateway(),
			"google_compute_vpn_tunnel":             resourceComputeVpnTunnel(),
			"google_container_cluster":              resourceContainerCluster(),
			"google_dns_managed_zone":               resourceDnsManagedZone(),
			"google_dns_record_set":                 resourceDnsRecordSet(),
			"google_sql_database":                   resourceSqlDatabase(),
			"google_sql_database_instance":          resourceSqlDatabaseInstance(),
			"google_sql_user":                       resourceSqlUser(),
			"google_project":                        resourceGoogleProject(),
			"google_pubsub_topic":                   resourcePubsubTopic(),
			"google_pubsub_subscription":            resourcePubsubSubscription(),
			"google_service_account":                resourceGoogleServiceAccount(),
			"google_storage_bucket":                 resourceStorageBucket(),
			"google_storage_bucket_acl":             resourceStorageBucketAcl(),
			"google_storage_bucket_object":          resourceStorageBucketObject(),
			"google_storage_object_acl":             resourceStorageObjectAcl(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	credentials := d.Get("credentials").(string)
	if credentials == "" {
		credentials = d.Get("account_file").(string)
	}
	config := Config{
		Credentials: credentials,
		Project:     d.Get("project").(string),
		Region:      d.Get("region").(string),
	}

	if err := config.loadAndValidate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateAccountFile(v interface{}, k string) (warnings []string, errors []error) {
	if v == nil {
		return
	}

	value := v.(string)

	if value == "" {
		return
	}

	contents, wasPath, err := pathorcontents.Read(value)
	if err != nil {
		errors = append(errors, fmt.Errorf("Error loading Account File: %s", err))
	}
	if wasPath {
		warnings = append(warnings, `account_file was provided as a path instead of
as file contents. This support will be removed in the future. Please update
your configuration to use ${file("filename.json")} instead.`)
	}

	var account accountFile
	if err := json.Unmarshal([]byte(contents), &account); err != nil {
		errors = append(errors,
			fmt.Errorf("account_file not valid JSON '%s': %s", contents, err))
	}

	return
}

func validateCredentials(v interface{}, k string) (warnings []string, errors []error) {
	if v == nil || v.(string) == "" {
		return
	}
	creds := v.(string)
	var account accountFile
	if err := json.Unmarshal([]byte(creds), &account); err != nil {
		errors = append(errors,
			fmt.Errorf("credentials are not valid JSON '%s': %s", creds, err))
	}

	return
}

// getRegionFromZone returns the region from a zone for Google cloud.
func getRegionFromZone(zone string) string {
	if zone != "" && len(zone) > 2 {
		region := zone[:len(zone)-2]
		return region
	}
	return ""
}

// getRegion reads the "region" field from the given resource data and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
func getRegion(d *schema.ResourceData, config *Config) (string, error) {
	res, ok := d.GetOk("region")
	if !ok {
		if config.Region != "" {
			return config.Region, nil
		}
		return "", fmt.Errorf("%q: required field is not set", "region")
	}
	return res.(string), nil
}

// getProject reads the "project" field from the given resource data and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
func getProject(d *schema.ResourceData, config *Config) (string, error) {
	res, ok := d.GetOk("project")
	if !ok {
		if config.Project != "" {
			return config.Project, nil
		}
		return "", fmt.Errorf("%q: required field is not set", "project")
	}
	return res.(string), nil
}

func getZonalResourceFromRegion(getResource func(string) (interface{}, error), region string, compute *compute.Service, project string) (interface{}, error) {
	zoneList, err := compute.Zones.List(project).Do()
	if err != nil {
		return nil, err
	}
	var resource interface{}
	for _, zone := range zoneList.Items {
		if strings.Contains(zone.Name, region) {
			resource, err = getResource(zone.Name)
			if err != nil {
				if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
					// Resource was not found in this zone
					continue
				}
				return nil, fmt.Errorf("Error reading Resource: %s", err)
			}
			// Resource was found
			return resource, nil
		}
	}
	// Resource does not exist in this region
	return nil, nil
}

// getNetworkLink reads the "network" field from the given resource data and if the value:
// - is a resource URL, returns the string unchanged
// - is the network name only, then looks up the resource URL using the google client
func getNetworkLink(d *schema.ResourceData, config *Config, field string) (string, error) {
	if v, ok := d.GetOk(field); ok {
		network := v.(string)

		project, err := getProject(d, config)
		if err != nil {
			return "", err
		}

		if !strings.HasPrefix(network, "https://www.googleapis.com/compute/") {
			// Network value provided is just the name, lookup the network SelfLink
			networkData, err := config.clientCompute.Networks.Get(
				project, network).Do()
			if err != nil {
				return "", fmt.Errorf("Error reading network: %s", err)
			}
			network = networkData.SelfLink
		}

		return network, nil

	} else {
		return "", nil
	}
}

// getNetworkName reads the "network" field from the given resource data and if the value:
// - is a resource URL, extracts the network name from the URL and returns it
// - is the network name only (i.e not prefixed with http://www.googleapis.com/compute/...), is returned unchanged
func getNetworkName(d *schema.ResourceData, field string) (string, error) {
	if v, ok := d.GetOk(field); ok {
		network := v.(string)

		if strings.HasPrefix(network, "https://www.googleapis.com/compute/") {
			// extract the network name from SelfLink URL
			networkName := network[strings.LastIndex(network, "/")+1:]
			if networkName == "" {
				return "", fmt.Errorf("network url not valid")
			}
			return networkName, nil
		}

		return network, nil
	}
	return "", nil
}
