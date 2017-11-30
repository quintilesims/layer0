package layer0

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

func resourceLayer0Deploy() *schema.Resource {
	return &schema.Resource{
		Create: resourceLayer0DeployCreate,
		Read:   resourceLayer0DeployRead,
		Delete: resourceLayer0DeployDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLayer0DeployCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	req := models.CreateDeployRequest{
		DeployName: d.Get("name").(string),
		DeployFile: []byte(d.Get("content").(string)),
	}

	jobID, err := apiClient.CreateDeploy(req)
	if err != nil {
		return err
	}

	job, err := client.WaitForJob(apiClient, jobID, config.DefaultTimeout)
	if err != nil {
		return err
	}

	deployID := job.Result
	d.SetId(deployID)

	return resourceLayer0DeployRead(d, meta)
}

func resourceLayer0DeployRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	deployID := d.Id()

	deploy, err := apiClient.ReadDeploy(deployID)
	if err != nil {
		if err, ok := err.(*errors.ServerError); ok && err.Code == errors.DeployDoesNotExist {
			d.SetId("")
			log.Printf("[WARN] Error Reading Deploy (%s), deploy does not exist", deployID)
			return nil
		}

		return err
	}

	// do not set content as it fails to properly diff
	d.Set("name", deploy.DeployName)
	d.Set("version", deploy.Version)

	return nil
}

func resourceLayer0DeployDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	deployID := d.Id()

	jobID, err := apiClient.DeleteDeploy(deployID)
	if err != nil {
		return err
	}

	if _, err := client.WaitForJob(apiClient, jobID, config.DefaultTimeout); err != nil {
		if err, ok := err.(*errors.ServerError); ok && err.Code == errors.DeployDoesNotExist {
			return nil
		}

		return err
	}

	return nil
}
