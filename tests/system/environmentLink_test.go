package system

import (
	"log"
	"testing"
	"time"

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
	log.Printf("[DEBUG] Running curl while link exists")
	for start := time.Now(); time.Since(start) < time.Minute*5; time.Sleep(time.Second * 10) {
		output, err := publicService.RunCommand("curl", "-m", "10", "-s", privateServiceURL)
		if err != nil {
			t.Fatalf("[ERROR] Error running curl: %v", err)
		}

		if expected := "Hello, World!"; output != expected {
			t.Fatalf("[ERROR] Output from curl was '%s', expected '%s'", output, expected)
		}
	}

	log.Printf("[DEBUG] Removing environment link")
	links := []string{}
	s.Layer0.UpdateEnvironmentLink(publicEnvironmentID, links)

	log.Printf("[DEBUG] Running curl without link")
	for start := time.Now(); time.Since(start) < time.Minute*2; time.Sleep(time.Second * 10) {
		output, err := publicService.RunCommand("curl", "-m", "10", "-s", privateServiceURL)
		if err != nil {
			t.Fatalf("[ERROR] Error running curl: %v", err)
		}

		if output != "" {
			t.Fatalf("[ERROR] Output from curl was '%s', expected no output", output)
		}
	}
}
