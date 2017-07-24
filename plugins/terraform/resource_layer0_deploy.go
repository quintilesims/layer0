package main

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/cli/client"
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
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppressEquivalentDockerrunDiffs,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLayer0DeployCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(client.Client)

	name := d.Get("name").(string)
	content := d.Get("content").(string)

	deploy, err := client.CreateDeploy(name, []byte(content))
	if err != nil {
		return err
	}

	d.SetId(deploy.DeployID)
	return resourceLayer0DeployRead(d, meta)
}

func resourceLayer0DeployRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(client.Client)
	deployID := d.Id()

	deploy, err := client.GetDeploy(deployID)
	if err != nil {
		if strings.Contains(err.Error(), "No deploy found") {
			d.SetId("")
			log.Printf("[WARN] Error Reading Deploy (%s), deploy does not exist", deployID)
			return nil
		}

		return err
	}

	d.Set("name", deploy.DeployName)
	d.Set("content", string(deploy.Dockerrun))
	d.Set("version", deploy.Version)

	return nil
}

func resourceLayer0DeployDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(client.Client)
	deployID := d.Id()

	if err := client.DeleteDeploy(deployID); err != nil {
		if strings.Contains(err.Error(), "No deploy found") {
			return nil
		}

		return err
	}

	return nil
}
