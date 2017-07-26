package main

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLayer0Service() *schema.Resource {
	return &schema.Resource{
		Create: resourceLayer0ServiceCreate,
		Read:   resourceLayer0ServiceRead,
		Update: resourceLayer0ServiceUpdate,
		Delete: resourceLayer0ServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"environment": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"deploy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"load_balancer": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"scale": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
		},
	}
}

func resourceLayer0ServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Layer0Client)

	environmentID := d.Get("environment").(string)
	name := d.Get("name").(string)
	deployID := d.Get("deploy").(string)
	loadBalancerID := d.Get("load_balancer").(string)
	scale := d.Get("scale").(int)

	service, err := client.API.CreateService(name, environmentID, deployID, loadBalancerID)
	if err != nil {
		return err
	}

	// set id first to tell terraform resource has been created
	d.SetId(service.ServiceID)

	if scale != 1 {
		if _, err := client.API.ScaleService(service.ServiceID, scale); err != nil {
			return err
		}
	}

	if err := waitForDeploymentWithContext(client, service.ServiceID); err != nil {
		return err
	}

	return resourceLayer0ServiceRead(d, meta)
}

func resourceLayer0ServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Layer0Client)
	serviceID := d.Id()

	service, err := client.API.GetService(serviceID)
	if err != nil {
		if strings.Contains(err.Error(), "No service found") {
			d.SetId("")
			log.Printf("[WARN] Error Reading Service (%s), service does not exist", serviceID)
			return nil
		}

		return err
	}

	d.Set("environment", service.EnvironmentID)
	d.Set("name", service.ServiceName)
	d.Set("load_balancer", service.LoadBalancerID)
	d.Set("scale", service.DesiredCount)

	for _, deployment := range service.Deployments {
		if deployment.Status == "PRIMARY" {
			d.Set("deploy", deployment.DeployID)
		}
	}

	return nil
}

func resourceLayer0ServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Layer0Client)
	serviceID := d.Id()

	if d.HasChange("deploy") {
		deployID := d.Get("deploy").(string)

		if _, err := client.API.UpdateService(serviceID, deployID); err != nil {
			return err
		}
	}

	if d.HasChange("scale") {
		scale := d.Get("scale").(int)

		if _, err := client.API.ScaleService(serviceID, scale); err != nil {
			return err
		}
	}

	 if err := waitForDeploymentWithContext(client, serviceID); err != nil {
                return err
        }

	return resourceLayer0ServiceRead(d, meta)
}

func resourceLayer0ServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Layer0Client)
	serviceID := d.Id()

	jobID, err := client.API.DeleteService(serviceID)
	if err != nil {
		if strings.Contains(err.Error(), "No service found") {
			return nil
		}

		return err
	}

	if err := waitForJobWithContext(client, jobID); err != nil {
		return err
	}

	return nil
}
