package system

import (
	"log"
	"testing"
	"time"

	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"github.com/quintilesims/layer0/tests/clients"
)

// Test Resources:
// This test creates two linked environments named 'el_public' and 'el_private`
// The 'el_public' environment has a STS service running behind a public load balancer
// The 'el_private' environment has a STS service running behind a private load balancer
func TestEnvironmentLink(t *testing.T) {
	t.Parallel()

	s := NewSystemTest(t, "cases/environment_link", nil)
	s.Terraform.Apply()
	defer s.Terraform.Destroy()

	publicEnvironmentID := s.Terraform.Output("public_environment_id")
	publicServiceURL := s.Terraform.Output("public_service_url")
	privateServiceURL := s.Terraform.Output("private_service_url")

	publicService := clients.NewSTSTestClient(t, publicServiceURL)
	publicService.WaitForHealthy(time.Minute * 3)

	// curl the private service in the private environment from the public service in the public environment
	// the private service returns "Hello, World!" from its root path
	testutils.WaitFor(t, time.Second*10, time.Minute*5, func() bool {
		log.Printf("[DEBUG] Running curl while link exists")
		output, err := publicService.RunCommand("curl", "-m", "10", "-s", privateServiceURL)
		if err != nil {
			log.Printf("[ERROR] Error running curl: %v", err)
			return false
		}

		if expected := "Hello, World!"; output != expected {
			log.Printf("[ERROR] Output from curl was '%s', expected '%s'", output, expected)
			return false
		}

		return true
	})

	log.Printf("[DEBUG] Removing environment link")
	req := models.UpdateEnvironmentRequest{
		Links: &[]string{},
	}

	s.Layer0.UpdateEnvironment(publicEnvironmentID, req)

	testutils.WaitFor(t, time.Second*10, time.Minute*2, func() bool {
		log.Printf("[DEBUG] Running curl without link")
		output, err := publicService.RunCommand("curl", "-m", "10", "-s", privateServiceURL)
		if err != nil {
			log.Printf("[ERROR] Error running curl: %v", err)
			return false
		}

		if output != "" {
			log.Printf("[ERROR] Output from curl was '%s', expected no output", output)
			return false
		}

		return true
	})
}
