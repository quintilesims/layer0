package job

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

var DeleteLoadBalancerSteps = []Step{
	{
		Name:    "Delete Load Balancer",
		Timeout: time.Minute * 10,
		Action:  DeleteLoadBalancer,
	},
}

func DeleteLoadBalancer(quit chan bool, context *JobContext) error {
	loadBalancerID := context.Request()

	return runAndRetry(quit, time.Second*10, func() error {
		log.Infof("Running Action: DeleteLoadBalancer on '%s'", loadBalancerID)
		return context.LoadBalancerLogic.DeleteLoadBalancer(loadBalancerID)
	})
}
