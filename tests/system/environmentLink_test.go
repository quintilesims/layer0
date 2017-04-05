package system

import (
	"github.com/quintilesims/layer0/tests/system/clients"
	"testing"
	"time"
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

	publicServiceURL := s.Terraform.Output("public_service_url")
	privateServiceURL := s.Terraform.Output("private_service_url")

	publicService := clients.NewSTSTestClient(t, publicServiceURL)
	publicService.WaitForHealthy(time.Minute * 3)

	// curl the private service in the private environment from the public service in the public environment
	// the private service returns "Hello, World!" from its root path
	output, err := publicService.RunCommand("curl", "-s", privateServiceURL)
	if err != nil {
		t.Fatal(err)
	}

	if expected := "Hello, World!"; output != expected {
		t.Fatalf("Output was '%s', expected '%s'", output, expected)
	}
}
