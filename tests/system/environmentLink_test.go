package system

import (
	"github.com/quintilesims/layer0/tests/system/clients"
	"testing"
	"time"
)

// Test Resources:
// This test creates two linked environments named 'el_alpha' and 'el_beta`
// The 'el_alpha' environment has a STS service running behind a public load balancer
// The 'el_beta' environment has a STS service running behind a private load balancer
func TestEnvironmentLink(t *testing.T) {
	t.Parallel()

	s := NewSystemTest(t, "cases/environment_link", nil)
	s.Terraform.Apply()
	defer s.Terraform.Destroy()

	alphaServiceURL := s.Terraform.Output("alpha_service_url")
	alphaEnvironmentID := s.Terraform.Output("alpha_environment_id")
	betaServiceURL := s.Terraform.Output("beta_service_url")
	betaEnvironmentID := s.Terraform.Output("beta_environment_id")

	alphaService := clients.NewSTSTestClient(t, alphaServiceURL)
	alphaService.WaitForHealthy(time.Minute * 3)

	// curl the private service in the beta environment from the public service in the alpha environment
	// the private service returns "Hello, World!" from its root path
	outputWithLink, err := alphaService.RunCommand("curl", "-s", betaServiceURL)
	if err != nil {
		t.Fatal(err)
	}

	if expected := "Hello, World!"; outputWithLink != expected {
		t.Fatalf("Output when link exists was '%s', expected '%s'", outputWithLink, expected)
	}
}
