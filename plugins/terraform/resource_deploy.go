package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/cli/client"
	"log"
	"strings"
)

func resourceLayer0Deploy() *schema.Resource {
	return &schema.Resource{
		Create: resourceLayer0DeployCreate,
		Read:   resourceLayer0DeployRead,
		Delete: resourceLayer0DeployDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"content": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceLayer0DeployCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.APIClient)

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
	client := meta.(*client.APIClient)
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
	d.Set("version", deploy.Version)

	return nil
}

func resourceLayer0DeployDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.APIClient)
	deployID := d.Id()

	if err := client.DeleteDeploy(deployID); err != nil {
		if strings.Contains(err.Error(), "No deploy found") {
			return nil
		}

		return err
	}

	return nil
}
