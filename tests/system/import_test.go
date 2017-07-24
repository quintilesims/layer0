package system

import (
	"io/ioutil"
	"testing"
)

// Test Resources:
// This test creates an environment named 'import' that has a
// SystemTestService named 'sts'
func TestImport(t *testing.T) {
	t.Parallel()

	s := NewSystemTest(t, "cases/import", nil)
	defer s.Terraform.Destroy()

	data, err := ioutil.ReadFile("cases/modules/sts/Dockerrun.aws.json")
	if err != nil {
		t.Fatalf("Failed to read dockerrun: %v", err)
	}

	log.Debugf("Creating test resources")
	environment := s.Layer0.CreateEnvironment("import")
	loadBalancer := s.Layer0.CreateLoadBalancer("sts", environment.EnvironmentID)
	deploy := s.Layer0.CreateDeploy("sts", data)
	service := s.Layer0.CreateService("sts", environment.EnvironmentID, deploy.DeployID, loadBalancer.LoadBalancerID)

	s.Terraform.Import("layer0_environment.import", environment.EnvironmentID)
	s.Terraform.Import("module.sts.layer0_load_balancer.sts", loadBalancer.LoadBalancerID)
	s.Terraform.Import("module.sts.layer0_deploy.sts", deploy.DeployID)
	s.Terraform.Import("module.sts.layer0_service.sts", service.ServiceID)

	s.Terraform.Apply()
}
