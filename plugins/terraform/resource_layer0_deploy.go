package main

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/common/errors"
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
	client := meta.(*Layer0Client)

	name := d.Get("name").(string)
	content := d.Get("content").(string)

	deploy, err := client.API.CreateDeploy(name, []byte(content))
	if err != nil {
		return err
	}

	d.SetId(deploy.DeployID)
	return resourceLayer0DeployRead(d, meta)
}

func resourceLayer0DeployRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Layer0Client)
	deployID := d.Id()

	deploy, err := client.API.GetDeploy(deployID)
	if err != nil {
		if err, ok := err.(*errors.ServerError); ok && err.Code == errors.DeployDoesNotExist {
			d.SetId("")
			log.Printf("[WARN] Error Reading Deploy (%s), deploy does not exist", deployID)
			return nil
		}

		return err
	}

	d.Set("name", deploy.DeployName)
	d.Set("version", deploy.Version)

	// do not set content as it fails to properly diff against what's
	// returned by the Layer0 API
	// TODO: improve suppressEquivalentDockerrunDiffs to ignore non-critical
	// differences between dockerruns

	return nil
}

func resourceLayer0DeployDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Layer0Client)
	deployID := d.Id()

	if err := client.API.DeleteDeploy(deployID); err != nil {
		if err, ok := err.(*errors.ServerError); ok && err.Code == errors.DeployDoesNotExist {
			return nil
		}

		return err
	}

	return nil
}
