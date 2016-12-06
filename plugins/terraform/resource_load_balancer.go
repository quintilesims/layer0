package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.imshealth.com/xfra/layer0/cli/client"
	"gitlab.imshealth.com/xfra/layer0/common/models"
)

func resourceLayer0LoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceLayer0LoadBalancerCreate,
		Read:   resourceLayer0LoadBalancerRead,
		Update: resourceLayer0LoadBalancerUpdate,
		Delete: resourceLayer0LoadBalancerDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"environment": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Set:      resourceLayer0PortHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"container_port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"protocol": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"certificate": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceLayer0LoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.APIClient)

	name := d.Get("name").(string)
	environmentID := d.Get("environment").(string)
	private := d.Get("private").(bool)
	ports := expandPorts(d.Get("port").(*schema.Set).List())

	loadBalancer, err := client.CreateLoadBalancer(name, environmentID, ports, !private)
	if err != nil {
		return err
	}

	d.SetId(loadBalancer.LoadBalancerID)
	return resourceLayer0LoadBalancerRead(d, meta)
}

func resourceLayer0LoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.APIClient)
	loadBalancerID := d.Id()

	loadBalancer, err := client.GetLoadBalancer(loadBalancerID)
	if err != nil {
		if strings.Contains(err.Error(), "No load_balancer found") {
			d.SetId("")
			log.Printf("[WARN] Error Reading Load Balancer (%s), load balancer does not exist", loadBalancerID)
			return nil
		}

		return err
	}

	d.Set("name", loadBalancer.LoadBalancerName)
	d.Set("environment", loadBalancer.EnvironmentID)
	d.Set("private", !loadBalancer.IsPublic)
	d.Set("port", flattenPorts(loadBalancer.Ports))
	d.Set("url", loadBalancer.URL)

	return nil
}

func resourceLayer0LoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.APIClient)
	loadBalancerID := d.Id()

	if d.HasChange("port") {
		ports := expandPorts(d.Get("port").(*schema.Set).List())

		if _, err := client.UpdateLoadBalancer(loadBalancerID, ports); err != nil {
			return err
		}
	}

	return resourceLayer0LoadBalancerRead(d, meta)
}

func resourceLayer0LoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.APIClient)
	loadBalancerID := d.Id()

	jobID, err := client.DeleteLoadBalancer(loadBalancerID)
	if err != nil {
		if strings.Contains(err.Error(), "No load_balancer found") {
			return nil
		}

		return err
	}

	if err := client.WaitForJob(jobID, defaultTimeout); err != nil {
		return err
	}

	return nil
}

func resourceLayer0PortHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	hostPort := m["host_port"].(int)
	containerPort := m["container_port"].(int)
	protocol := strings.ToLower(m["protocol"].(string))

	buf.WriteString(fmt.Sprintf("%d-", hostPort))
	buf.WriteString(fmt.Sprintf("%d-", containerPort))
	buf.WriteString(fmt.Sprintf("%s-", protocol))

	if v, ok := m["certificate"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return hashcode.String(buf.String())
}

func expandPorts(flattened []interface{}) []models.Port {
	ports := []models.Port{}

	for _, flat := range flattened {
		data := flat.(map[string]interface{})

		port := models.Port{
			HostPort:      int64(data["host_port"].(int)),
			ContainerPort: int64(data["container_port"].(int)),
			Protocol:      data["protocol"].(string),
		}

		if v, ok := data["certificate"]; ok {
			port.CertificateID = v.(string)
		}

		ports = append(ports, port)
	}

	return ports
}

func flattenPorts(ports []models.Port) []map[string]interface{} {
	flattened := []map[string]interface{}{}

	for _, port := range ports {
		data := map[string]interface{}{
			"host_port":      port.HostPort,
			"container_port": port.ContainerPort,
			"protocol":       port.Protocol,
		}

		if port.CertificateID != "" {
			data["certificate"] = port.CertificateID
		}

		flattened = append(flattened, data)
	}

	return flattened
}
