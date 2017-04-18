package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/cli/client"
)

func resourceLayer0EnvironmentLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceLayer0EnvironmentLinkCreate,
		Read:   resourceLayer0EnvironmentLinkRead,
		Delete: resourceLayer0EnvironmentLinkDelete,

		Schema: map[string]*schema.Schema{
			"source": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dest": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceLayer0EnvironmentLinkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(client.Client)

	if err := client.CreateEnvironmentLink(d.Get("source"), d.Get("dest")); err != nil {
		return err
	}

	return resourceLayer0EnvironmentLinkRead(d, meta)
}

func resourceLayer0EnvironmentLinkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(client.Client)
	sourceID := d.Get("source").(string)
	destID := d.Get("dest").(string)

	sourceEnvironment, err := client.GetEnvironment(sourceID)
	if err != nil {
		if strings.Contains(err.Error(), "No environment found") {
			d.SetId("")
			log.Printf("[WARN] Error Reading Environment (%s), environment does not exist", sourceID)
			return nil
		}

		return err
	}

	foundLink := false
	for _, v := range sourceEnvironment.Links {
		if v == destID {
			foundLink = true
		}
	}

	if foundLink == false {
		return fmt.Errorf("Link to %s not found in environment %s", destID, sourceID)
	}

	return nil
}

func resourceLayer0EnvironmentLinkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(client.Client)
	sourceID := d.Get("source")
	destID := d.Get("dest")

	jobID, err := client.DeleteLayer0EnvironmentLink(sourceID, destID)
	if err != nil {
		return err
	}

	if err := client.WaitForJob(jobID, defaultTimeout); err != nil {
		return err
	}

	return nil
}
