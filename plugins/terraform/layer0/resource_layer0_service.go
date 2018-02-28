package layer0

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
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
			"stateful": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceLayer0ServiceCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	serviceType := models.DeployCompatibilityStateless
	if d.Get("stateful").(bool) == true {
		serviceType = models.DeployCompatibilityStateful
	}

	req := models.CreateServiceRequest{
		ServiceName:    d.Get("name").(string),
		EnvironmentID:  d.Get("environment").(string),
		DeployID:       d.Get("deploy").(string),
		LoadBalancerID: d.Get("load_balancer").(string),
		Scale:          d.Get("scale").(int),
		ServiceType:    serviceType,
	}

	serviceID, err := apiClient.CreateService(req)
	if err != nil {
		return err
	}

	d.SetId(serviceID)

	if _, err := client.WaitForDeployment(apiClient, serviceID, config.DefaultTimeout); err != nil {
		return err
	}

	return resourceLayer0ServiceRead(d, meta)
}

func resourceLayer0ServiceRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	serviceID := d.Id()

	service, err := apiClient.ReadService(serviceID)
	if err != nil {
		if err, ok := err.(*errors.ServerError); ok && err.Code == errors.ServiceDoesNotExist {
			d.SetId("")
			log.Printf("[WARN] Error Reading Service (%s), service does not exist", serviceID)
			return nil
		}

		return err
	}

	d.Set("name", service.ServiceName)
	d.Set("environment", service.EnvironmentID)
	d.Set("load_balancer", service.LoadBalancerID)
	d.Set("scale", service.DesiredCount)
	d.Set("stateful", bool(service.ServiceType == models.DeployCompatibilityStateful))

	for _, deployment := range service.Deployments {
		if deployment.Status == "PRIMARY" {
			d.Set("deploy", deployment.DeployID)
		}
	}

	return nil
}

func resourceLayer0ServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	serviceID := d.Id()

	req := models.UpdateServiceRequest{}

	if d.HasChange("deploy") {
		deployID := d.Get("deploy").(string)
		req.DeployID = &deployID
	}

	if d.HasChange("scale") {
		scale := d.Get("scale").(int)
		req.Scale = &scale
	}

	if req.DeployID != nil || req.Scale != nil {
		if err := apiClient.UpdateService(serviceID, req); err != nil {
			return err
		}

		if _, err := client.WaitForDeployment(apiClient, serviceID, config.DefaultTimeout); err != nil {
			return err
		}
	}

	return resourceLayer0ServiceRead(d, meta)
}

func resourceLayer0ServiceDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	serviceID := d.Id()

	if err := apiClient.DeleteService(serviceID); err != nil {
		if err, ok := err.(*errors.ServerError); ok && err.Code == errors.ServiceDoesNotExist {
			return nil
		}

		return err
	}

	return nil
}
