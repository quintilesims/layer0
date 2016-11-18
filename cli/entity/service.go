package entity

import (
	"fmt"
	"gitlab.imshealth.com/xfra/layer0/cli/printer/table"
	"gitlab.imshealth.com/xfra/layer0/common/models"
	"strings"
)

type Service models.Service

func NewService(model *models.Service) *Service {
	service := Service(*model)
	return &service
}

func (this *Service) Table() table.Table {
	environment := this.EnvironmentName
	if environment == "" {
		environment = this.EnvironmentID
	}

	loadBalancer := this.LoadBalancerName
	if loadBalancer == "" {
		loadBalancer = this.LoadBalancerID
	}

	deploys := []string{}
	for _, deploy := range this.Deployments {
		var display string

		if deploy.DeployName != "" && deploy.DeployVersion != "" {
			display = fmt.Sprintf("%s:%s", deploy.DeployName, deploy.DeployVersion)
		} else {
			display = strings.Replace(deploy.DeployID, ".", ":", 1)
		}

		if deploy.RunningCount != deploy.DesiredCount {
			display = fmt.Sprintf("%s*", display)
		}

		deploys = append(deploys, display)
	}

	scale := fmt.Sprintf("%d/%d", this.RunningCount, this.DesiredCount)
	if this.PendingCount != 0 {
		scale = fmt.Sprintf("%s (%d)", scale, this.PendingCount)
	}

	table := []table.Column{
		table.NewSingleRowColumn("SERVICE ID", this.ServiceID),
		table.NewSingleRowColumn("SERVICE NAME", this.ServiceName),
		table.NewSingleRowColumn("ENVIRONMENT", environment),
		table.NewSingleRowColumn("LOADBALANCER", loadBalancer),
		table.NewMultiRowColumn("DEPLOYS", deploys),
		table.NewSingleRowColumn("SCALE", scale),
	}

	return table
}
