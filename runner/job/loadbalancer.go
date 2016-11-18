package job

import (
	log "github.com/Sirupsen/logrus"
	"time"
)

var DeleteLoadBalancerSteps = []Step{
	Step{
		Name:    "Delete Load Balancer",
		Timeout: time.Minute * 5,
		Action:  DeleteLoadBalancer,
	},
}

func DeleteLoadBalancer(quit chan bool, context JobContext) error {
	loadBalancerID := context.Request()

	return runAndRetry(quit, time.Second*10, func() error {
		log.Infof("Running Action: DeleteLoadBalancer on '%s'", loadBalancerID)
		return context.LoadBalancerLogic().DeleteLoadBalancer(loadBalancerID)
	})
}
