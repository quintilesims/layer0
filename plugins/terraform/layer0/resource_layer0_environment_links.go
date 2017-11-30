package layer0

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

func resourceLayer0EnvironmentLinks() *schema.Resource {
	return &schema.Resource{
		Create: resourceLayer0EnvironmentLinksCreate,
		Read:   resourceLayer0EnvironmentLinksRead,
		Update: resourceLayer0EnvironmentLinksCreate,
		Delete: resourceLayer0EnvironmentLinksDelete,
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"links": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func resourceLayer0EnvironmentLinksCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	environmentID := d.Get("environment_id").(string)
	links := expandStringList(d.Get("links").([]interface{}))

	req := models.UpdateEnvironmentRequest{
		Links: &links,
	}

	jobID, err := apiClient.UpdateEnvironment(environmentID, req)
	if err != nil {
		return err
	}

	if _, err := client.WaitForJob(apiClient, jobID, config.DefaultTimeout); err != nil {
		return err
	}

	return resourceLayer0EnvironmentLinksRead(d, meta)
}

func resourceLayer0EnvironmentLinksRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	environmentID := d.Get("environment_id").(string)

	environment, err := apiClient.ReadEnvironment(environmentID)
	if err != nil {
		if err, ok := err.(*errors.ServerError); ok && err.Code == errors.EnvironmentDoesNotExist {
			d.SetId("")
			log.Printf("[WARN] Error Reading Environment Link (%s): %v", environmentID, err)
			return nil
		}

		return err
	}

	d.SetId(environmentID)
	d.Set("links", flattenStringList(environment.Links))
	return nil
}

func resourceLayer0EnvironmentLinksDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	links := []string{}
	req := models.UpdateEnvironmentRequest{
		Links: &links,
	}

	jobID, err := apiClient.UpdateEnvironment(d.Id(), req)
	if err != nil {
		return err
	}

	if _, err := client.WaitForJob(apiClient, jobID, config.DefaultTimeout); err != nil {
		return err
	}

	return nil
}

func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, val)
		}
	}

	return vs
}

func flattenStringList(list []string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, v)
	}

	return vs
}
