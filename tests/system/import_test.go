package system

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/common/models"
)

// Test Resources:
// This test creates an environment named 'import' that has a
// SystemTestService named 'sts'
func TestImport(t *testing.T) {
	t.Parallel()

	s := NewSystemTest(t, "cases/import", nil)

	// Don't actually run this test if dryrun is specified
	// as it will first create resources then test imports
	if s.Terraform.DryRun() {
		t.Skipf("Test cannot execute during a dry run")
	}

	defer s.Terraform.Destroy()

	log.Printf("[DEBUG] Creating test resources")
	createEnvironmentReq := models.CreateEnvironmentRequest{
		EnvironmentName: "import",
		InstanceSize:    "t2.micro",
		MinClusterCount: 0,
		OperatingSystem: "linux",
	}

	environmentID := s.Layer0.CreateEnvironment(createEnvironmentReq)

	createLoadBalancerReq := models.CreateLoadBalancerRequest{
		LoadBalancerName: "sts",
		EnvironmentID:    environmentID,
		IsPublic:         true,
		Ports:            []models.Port{aws.DefaultLoadBalancerPort},
		HealthCheck:      aws.DefaultHealthCheck,
	}

	loadBalancerID := s.Layer0.CreateLoadBalancer(createLoadBalancerReq)

	data, err := ioutil.ReadFile("cases/modules/sts/Dockerrun.aws.json")
	if err != nil {
		t.Fatalf("Failed to read dockerrun: %v", err)
	}

	createDeployReq := models.CreateDeployRequest{
		DeployName: "sts",
		DeployFile: data,
	}

	deployID := s.Layer0.CreateDeploy(createDeployReq)

	createServiceReq := models.CreateServiceRequest{
		DeployID:       deployID,
		EnvironmentID:  environmentID,
		LoadBalancerID: loadBalancerID,
		ServiceName:    "sts",
	}

	serviceID := s.Layer0.CreateService(createServiceReq)

	s.Terraform.Import("layer0_environment.import", environmentID)
	s.Terraform.Import("module.sts.layer0_load_balancer.sts", loadBalancerID)
	s.Terraform.Import("module.sts.layer0_deploy.sts", deployID)
	s.Terraform.Import("module.sts.layer0_service.sts", serviceID)

	s.Terraform.Init()
	s.Terraform.Apply()
}
