package system

import (
	"testing"
)

func TestDataSources(t *testing.T) {
	t.Parallel()

	s := NewSystemTest(t, "cases/datasources", nil)
	s.Terraform.Apply()
	defer s.Terraform.Destroy()

	checkOutput := func(key string) {
		log.Debugf("Checking output: %s", key)
		s.Terraform.Output(key)
	}

	//check environment outputs
	checkOutput("environment_id")
	checkOutput("environment_size")
	checkOutput("environment_min_count")
	checkOutput("environment_os")
	checkOutput("environment_ami")

	//check deploy output
	checkOutput("deploy_id")

	//check load balancer outputs
	checkOutput("load_balancer_id")
	checkOutput("load_balancer_name")
	checkOutput("load_balancer_environment_name")
	checkOutput("load_balancer_private")
	checkOutput("load_balancer_url")
	checkOutput("load_balancer_service_id")
	checkOutput("load_balancer_service_name")

	//check service outputs
	checkOutput("service_id")
	checkOutput("service_environment_name")
	checkOutput("service_lb_name")
	checkOutput("service_lb_id")
	checkOutput("service_scale")

	log.Debugf("L0 Terraform Provider Data sources Tests completed.")
}
