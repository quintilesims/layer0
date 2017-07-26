package main

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
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
				Default:  "m3.medium",
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
				Default:  "linux",
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
	client := meta.(*Layer0Client)

	name := d.Get("name").(string)
	size := d.Get("size").(string)
	minCount := d.Get("min_count").(int)
	userData := d.Get("user_data").(string)
	os := d.Get("os").(string)
	ami := d.Get("ami").(string)

	environment, err := client.API.CreateEnvironment(name, size, minCount, []byte(userData), os, ami)
	if err != nil {
		return err
	}

	d.SetId(environment.EnvironmentID)
	return resourceLayer0EnvironmentRead(d, meta)
}

func resourceLayer0EnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Layer0Client)
	environmentID := d.Id()

	environment, err := client.API.GetEnvironment(environmentID)
	if err != nil {
		if strings.Contains(err.Error(), "No environment found") {
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
	client := meta.(*Layer0Client)
	environmentID := d.Id()

	if d.HasChange("min_count") {
		minCount := d.Get("min_count").(int)

		if _, err := client.API.UpdateEnvironment(environmentID, minCount); err != nil {
			return err
		}
	}

	return resourceLayer0EnvironmentRead(d, meta)
}

func resourceLayer0EnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Layer0Client)
	environmentID := d.Id()

	jobID, err := client.API.DeleteEnvironment(environmentID)
	if err != nil {
		if strings.Contains(err.Error(), "No environment found") {
			return nil
		}

		return err
	}

	if err := waitForJobWithContext(client, jobID); err != nil {
		return err
	}

	return nil
}
