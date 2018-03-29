package layer0

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

func resourceLayer0LoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceLayer0LoadBalancerCreate,
		Read:   resourceLayer0LoadBalancerRead,
		Update: resourceLayer0LoadBalancerUpdate,
		Delete: resourceLayer0LoadBalancerDelete,
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
			"load_balancer_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  config.DefaultLoadBalancerType,
			},
			"private": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
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
						"certificate_arn": {
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
							Default:  config.DefaultLoadBalancerHealthCheck().Target,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  config.DefaultLoadBalancerHealthCheck().Path,
						},
						"interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  config.DefaultLoadBalancerHealthCheck().Interval,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  config.DefaultLoadBalancerHealthCheck().Timeout,
						},
						"healthy_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  config.DefaultLoadBalancerHealthCheck().HealthyThreshold,
						},
						"unhealthy_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  config.DefaultLoadBalancerHealthCheck().UnhealthyThreshold,
						},
					},
				},
			},
		},
	}
}

func resourceLayer0LoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)

	ports := expandPorts(d.Get("port").(*schema.Set).List())
	if len(ports) == 0 {
		ports = []models.Port{config.DefaultLoadBalancerPort()}
	}

	req := models.CreateLoadBalancerRequest{
		LoadBalancerName: d.Get("name").(string),
		LoadBalancerType: strings.ToLower(d.Get("load_balancer_type").(string)),
		EnvironmentID:    d.Get("environment").(string),
		IsPublic:         !d.Get("private").(bool),
		Ports:            ports,
		HealthCheck:      expandHealthCheck(d.Get("health_check")),
	}

	loadBalancerID, err := apiClient.CreateLoadBalancer(req)
	if err != nil {
		return err
	}

	d.SetId(loadBalancerID)
	return resourceLayer0LoadBalancerRead(d, meta)
}

func resourceLayer0LoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	loadBalancerID := d.Id()

	loadBalancer, err := apiClient.ReadLoadBalancer(loadBalancerID)
	if err != nil {
		if err, ok := err.(*errors.ServerError); ok && err.Code == errors.LoadBalancerDoesNotExist {
			d.SetId("")
			log.Printf("[WARN] Error Reading LoadBalancer (%s), loadBalancer does not exist", loadBalancerID)
			return nil
		}

		return err
	}

	d.Set("name", loadBalancer.LoadBalancerName)
	d.Set("environment", loadBalancer.EnvironmentID)
	d.Set("load_balancer_type", loadBalancer.LoadBalancerType)
	d.Set("private", !loadBalancer.IsPublic)
	d.Set("health_check", flattenHealthCheck(loadBalancer.HealthCheck))
	d.Set("port", flattenPorts(loadBalancer.Ports))
	d.Set("url", loadBalancer.URL)

	return nil
}

func resourceLayer0LoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	loadBalancerID := d.Id()

	req := models.UpdateLoadBalancerRequest{}

	if d.HasChange("port") {
		ports := expandPorts(d.Get("port").(*schema.Set).List())
		req.Ports = &ports
	}

	if d.HasChange("health_check") {
		healthCheck := expandHealthCheck(d.Get("health_check"))
		req.HealthCheck = &healthCheck
	}

	if req.Ports != nil || req.HealthCheck != nil {
		if err := apiClient.UpdateLoadBalancer(loadBalancerID, req); err != nil {
			return err
		}
	}

	return resourceLayer0LoadBalancerRead(d, meta)
}

func resourceLayer0LoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(client.Client)
	loadBalancerID := d.Id()

	if err := apiClient.DeleteLoadBalancer(loadBalancerID); err != nil {
		if err, ok := err.(*errors.ServerError); ok && err.Code == errors.LoadBalancerDoesNotExist {
			return nil
		}

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

	if v, ok := m["certificate_arn"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return hashcode.String(buf.String())
}

func expandHealthCheck(flattened interface{}) models.HealthCheck {
	hc := flattened.([]interface{})

	if len(hc) > 0 {
		check := hc[0].(map[string]interface{})

		return models.HealthCheck{
			Target:             check["target"].(string),
			Path:               check["path"].(string),
			Interval:           check["interval"].(int),
			Timeout:            check["timeout"].(int),
			HealthyThreshold:   check["healthy_threshold"].(int),
			UnhealthyThreshold: check["unhealthy_threshold"].(int),
		}
	}

	return config.DefaultLoadBalancerHealthCheck()
}

func flattenHealthCheck(healthCheck models.HealthCheck) []map[string]interface{} {
	result := make([]map[string]interface{}, 1)

	check := make(map[string]interface{})
	check["target"] = healthCheck.Target
	check["path"] = healthCheck.Path
	check["interval"] = healthCheck.Interval
	check["timeout"] = healthCheck.Timeout
	check["healthy_threshold"] = healthCheck.HealthyThreshold
	check["unhealthy_threshold"] = healthCheck.UnhealthyThreshold

	result[0] = check
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

		if v, ok := data["certificate_arn"]; ok {
			port.CertificateARN = v.(string)
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

		if port.CertificateARN != "" {
			data["certificate_arn"] = port.CertificateARN
		}

		flattened = append(flattened, data)
	}

	return flattened
}
