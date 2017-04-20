package main

import (
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
	sourceID := d.Get("source").(string)
	destID := d.Get("dest").(string)

	if err := client.CreateLink(sourceID, destID); err != nil {
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

	for _, v := range sourceEnvironment.Links {
		if v == destID {
			return nil
		}
	}

	d.SetId("")
	log.Printf("[WARN] Error Reading Environment Link (%s), link does not exist", sourceID)
	return nil
}

func resourceLayer0EnvironmentLinkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(client.Client)
	sourceID := d.Get("source").(string)
	destID := d.Get("dest").(string)

	if err := client.DeleteLink(sourceID, destID); err != nil {
		return err
	}

	return nil
}
