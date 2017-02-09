package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/cli/client"
	"github.com/quintilesims/layer0/common/models"
)

func resourceLayer0LoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceLayer0LoadBalancerCreate,
		Read:   resourceLayer0LoadBalancerRead,
		Update: resourceLayer0LoadBalancerUpdate,
		Delete: resourceLayer0LoadBalancerDelete,

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
			"private": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeSet,
				Required: true,
				Set:      resourceLayer0PortHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"container_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"certificate": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"health_check": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "TCP:80",
						},
						"interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  30,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  5,
						},
						"healthy_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  2,
						},
						"unhealthy_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  2,
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
	healthCheck := expandHealthCheck(d.Get("health_check"))

	if healthCheck == nil {
		healthCheck = &models.HealthCheck{
			Target:             "TCP:80",
			Interval:           30,
			Timeout:            5,
			HealthyThreshold:   2,
			UnhealthyThreshold: 2,
		}
	}

	loadBalancer, err := client.CreateLoadBalancer(name, environmentID, *healthCheck, ports, !private)
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
	d.Set("health_check", flattenHealthCheck(loadBalancer.HealthCheck))
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

		if _, err := client.UpdateLoadBalancerPorts(loadBalancerID, ports); err != nil {
			return err
		}
	}

	if d.HasChange("health_check") {
		healthCheck := expandHealthCheck(d.Get("health_check"))

		if _, err := client.UpdateLoadBalancerHealthCheck(loadBalancerID, *healthCheck); err != nil {
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

func expandHealthCheck(flattened interface{}) *models.HealthCheck {
	hc := flattened.([]interface{})

	if len(hc) > 0 {
		check := hc[0].(map[string]interface{})

		return &models.HealthCheck{
			Target:             check["target"].(string),
			Interval:           check["interval"].(int),
			Timeout:            check["timeout"].(int),
			HealthyThreshold:   check["healthy_threshold"].(int),
			UnhealthyThreshold: check["unhealthy_threshold"].(int),
		}
	}

	return nil
}

func flattenHealthCheck(healthCheck models.HealthCheck) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)

	check := make(map[string]interface{})
	check["target"] = healthCheck.Target
	check["interval"] = healthCheck.Interval
	check["timeout"] = healthCheck.Timeout
	check["healthy_threshold"] = healthCheck.HealthyThreshold
	check["unhealthy_threshold"] = healthCheck.UnhealthyThreshold

	result = append(result, check)

	return result
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
			port.CertificateName = v.(string)
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

		if port.CertificateName != "" {
			data["certificate"] = port.CertificateName
		}

		flattened = append(flattened, data)
	}

	return flattened
}
