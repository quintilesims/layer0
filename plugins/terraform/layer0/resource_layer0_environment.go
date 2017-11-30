package layer0

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

func resourceLayer0Environment() *schema.Resource {
	return &schema.Resource{
		Create: resourceLayer0EnvironmentCreate,
		Read:   resourceLayer0EnvironmentRead,
		Update: resourceLayer0EnvironmentUpdate,
		Delete: resourceLayer0EnvironmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"size": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  aws.DefaultInstanceSize,
				ForceNew: true,
			},
			"min_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"os": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  aws.DefaultEnvironmentOS,
				ForceNew: true,
			},
			"ami": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"cluster_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLayer0EnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	req := models.CreateEnvironmentRequest{
		EnvironmentName:  d.Get("name").(string),
		InstanceSize:     d.Get("size").(string),
		UserDataTemplate: []byte(d.Get("user_data").(string)),
		MinClusterCount:  d.Get("min_count").(int),
		OperatingSystem:  d.Get("os").(string),
		AMIID:            d.Get("ami").(string),
	}

	jobID, err := apiClient.CreateEnvironment(req)
	if err != nil {
		return err
	}

	job, err := client.WaitForJob(apiClient, jobID, config.DefaultTimeout)
	if err != nil {
		return err
	}

	environmentID := job.Result
	d.SetId(environmentID)

	return resourceLayer0EnvironmentRead(d, meta)
}

func resourceLayer0EnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	environmentID := d.Id()

	environment, err := apiClient.ReadEnvironment(environmentID)
	if err != nil {
		if err, ok := err.(*errors.ServerError); ok && err.Code == errors.EnvironmentDoesNotExist {
			d.SetId("")
			log.Printf("[WARN] Error Reading Environment (%s), environment does not exist", environmentID)
			return nil
		}

		return err
	}

	d.Set("name", environment.EnvironmentName)
	d.Set("size", environment.InstanceSize)
	d.Set("cluster_count", environment.ClusterCount)
	d.Set("security_group_id", environment.SecurityGroupID)
	d.Set("os", environment.OperatingSystem)
	d.Set("ami", environment.AMIID)

	return nil
}

func resourceLayer0EnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	environmentID := d.Id()

	if d.HasChange("min_count") {
		minCount := d.Get("min_count").(int)

		req := models.UpdateEnvironmentRequest{
			MinClusterCount: &minCount,
		}

		jobID, err := apiClient.UpdateEnvironment(environmentID, req)
		if err != nil {
			return err
		}

		if _, err := client.WaitForJob(apiClient, jobID, config.DefaultTimeout); err != nil {
			return err
		}
	}

	return resourceLayer0EnvironmentRead(d, meta)
}

func resourceLayer0EnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	environmentID := d.Id()

	jobID, err := apiClient.DeleteEnvironment(environmentID)
	if err != nil {
		return err
	}

	if _, err := client.WaitForJob(apiClient, jobID, config.DefaultTimeout); err != nil {
		if err, ok := err.(*errors.ServerError); ok && err.Code == errors.EnvironmentDoesNotExist {
			return nil
		}

		return err
	}

	return nil
}
