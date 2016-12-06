package entity

import (
	"fmt"
	"github.com/quintilesims/layer0/cli/printer/table"
	"github.com/quintilesims/layer0/common/models"
	"strings"
)

type LoadBalancer models.LoadBalancer

func NewLoadBalancer(model *models.LoadBalancer) *LoadBalancer {
	loadBalancer := LoadBalancer(*model)
	return &loadBalancer
}

func (this *LoadBalancer) Table() table.Table {
	ports := []string{}
	for _, p := range this.Ports {
		port := fmt.Sprintf("%d:%d/%s", p.HostPort, p.ContainerPort, strings.ToUpper(p.Protocol))
		ports = append(ports, port)
	}

	environment := this.EnvironmentID
	if this.EnvironmentName != "" {
		environment = this.EnvironmentName
	}

	table := []table.Column{
		table.NewSingleRowColumn("LOADBALANCER ID", this.LoadBalancerID),
		table.NewSingleRowColumn("LOADBALANCER NAME", this.LoadBalancerName),
		table.NewSingleRowColumn("ENVIRONMENT", environment),
		table.NewMultiRowColumn("SERVICES", this.Services),
		table.NewMultiRowColumn("PORTS", ports),
		table.NewSingleRowColumn("PUBLIC", fmt.Sprintf("%v", this.IsPublic)),
		table.NewSingleRowColumn("URL", this.URL),
	}

	return table
}
