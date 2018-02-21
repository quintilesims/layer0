package layer0

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
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
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"environment_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  config.DefaultEnvironmentType,
				ForceNew: true,
			},
			"scale": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"os": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  config.DefaultEnvironmentOS,
				ForceNew: true,
			},
			"ami": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"current_scale": {
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
		InstanceType:     d.Get("instance_type").(string),
		EnvironmentType:  d.Get("environment_type").(string),
		UserDataTemplate: []byte(d.Get("user_data").(string)),
		Scale:            d.Get("scale").(int),
		OperatingSystem:  d.Get("os").(string),
		AMIID:            d.Get("ami").(string),
	}

	environmentID, err := apiClient.CreateEnvironment(req)
	if err != nil {
		return err
	}

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
	d.Set("instance_type", environment.InstanceType)
	d.Set("environment_type", environment.EnvironmentType)
	d.Set("scale", environment.DesiredScale)
	d.Set("security_group_id", environment.SecurityGroupID)
	d.Set("os", environment.OperatingSystem)
	d.Set("ami", environment.AMIID)

	return nil
}

func resourceLayer0EnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	environmentID := d.Id()

	req := models.UpdateEnvironmentRequest{}
	if d.HasChange("scale") {
		scale := d.Get("scale").(int)
		req.Scale = &scale
	}

	if req.Scale != nil {
		if err := apiClient.UpdateEnvironment(environmentID, req); err != nil {
			return err
		}
	}

	return resourceLayer0EnvironmentRead(d, meta)
}

func resourceLayer0EnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	environmentID := d.Id()

	if err := apiClient.DeleteEnvironment(environmentID); err != nil {
		if err, ok := err.(*errors.ServerError); ok && err.Code == errors.EnvironmentDoesNotExist {
			return nil
		}

		return err
	}

	return nil
}
