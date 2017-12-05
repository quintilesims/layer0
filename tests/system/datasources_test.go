package system

import (
	"log"
	"testing"
)

// Test Resources:
// This test creates the following:
// - environment named dsrctest
// - load balancer named dsrctest
// - service named dsrctest
// - deploy named dsrctest
func TestDataSources(t *testing.T) {
	t.Parallel()

	s := NewSystemTest(t, "cases/datasources", nil)
	s.Terraform.Apply()
	defer s.Terraform.Destroy()

	// Compare outputs of data and resource values (resource
	// values have the '_expected' suffix)
	checkOutput := func(key string) {
		log.Printf("[DEBUG] Checking data source vs resource output for key: %s", key)

		if dVal, rVal := s.Terraform.Output(key), s.Terraform.Output(key+"_expected"); dVal != rVal {
			t.Fatalf(
				"Data value '%s' and Resource value '%s' do not match for key: %s",
				dVal,
				rVal,
				key)
		}
	}

	//check environment outputs
	checkOutput("environment_id")
	checkOutput("environment_name")
	checkOutput("environment_size")
	checkOutput("environment_min_count")
	checkOutput("environment_os")
	checkOutput("environment_ami")

	//check deploy output
	checkOutput("deploy_id")
	checkOutput("deploy_name")
	checkOutput("deploy_version")

	//check load balancer outputs
	checkOutput("load_balancer_id")
	checkOutput("load_balancer_name")
	checkOutput("load_balancer_environment_name")
	checkOutput("load_balancer_private")
	checkOutput("load_balancer_url")

	//check service outputs
	checkOutput("service_id")
	checkOutput("service_name")
	checkOutput("service_environment_id")
	checkOutput("service_environment_name")
	checkOutput("service_scale")

	log.Printf("L0 Terraform Provider Data sources Tests completed.")
}
