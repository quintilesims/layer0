package system

import (
	"fmt"
	"github.com/quintilesims/layer0/cli/client"
	"github.com/quintilesims/layer0/cli/command"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"testing"
)

type SystemTestContext struct {
	T        *testing.T
	Dir      string
	Client   *client.APIClient
	Resolver *command.TagResolver
}

// todo: tfvars - Should we just use env vars - TF_VAR_token/endpoint?
func NewSystemTestContext(t *testing.T, dir string) *SystemTestContext {
	apiClient := client.NewAPIClient(client.Config{
		Endpoint:      config.APIEndpoint(),
		Token:         fmt.Sprintf("Basic %s", config.AuthToken()),
	})

	return &SystemTestContext{
		T:        t,
		Dir:      dir,
		Client:   apiClient,
		Resolver: command.NewTagResolver(apiClient),
	}
}

func (s *SystemTestContext) GetEnvironment(name string) *models.Environment {
	id := s.resolve("environment", name)
	environment, err := s.Client.GetEnvironment(id)
	if err != nil {
		s.T.Fatal(err)
	}

	return environment
}

func (s *SystemTestContext) GetLoadBalancer(name string) *models.LoadBalancer {
	id := s.resolve("load_balancer", name)
	loadBalancer, err := s.Client.GetLoadBalancer(id)
	if err != nil {
		s.T.Fatal(err)
	}

	return loadBalancer
}

func (s *SystemTestContext) GetService(name string) *models.Service {
	id := s.resolve("service", name)
	service, err := s.Client.GetService(id)
	if err != nil {
		s.T.Fatal(err)
	}

	return service
}

func (s *SystemTestContext) resolve(entityType, name string) string {
	ids, err := s.Resolver.Resolve(entityType, name)
	if err != nil {
		s.T.Fatal(err)
	}

	if len(ids) == 0 {
		s.T.Fatalf("Failed to resolve %s '%s' - no ids found", entityType, name)
	}

	if len(ids) > 1 {
		s.T.Fatalf("Failed to resolve %s '%s' - multiple ids found (%v)", entityType, name, ids)
	}

	return ids[0]
}

func (s *SystemTestContext) Apply() {
	print("!!! WARNING: USING TERRAFORM PLAN INSTEAD OF APPLY !!!\n")
	s.run("terraform", "apply")
}

func (s *SystemTestContext) Destroy() {
	s.run("terraform", "destroy", "-force")
}

func (s *SystemTestContext) run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = s.Dir

	// todo: send stdout to t.Log()  so it only shows up with 'go test -v'
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	s.cleanupCommandOnSIGTERM(cmd)

	if err := cmd.Start(); err != nil {
		s.T.Fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		s.T.Fatal(err)
	}
}

func (s *SystemTestContext) cleanupCommandOnSIGTERM(cmd *exec.Cmd) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cmd.Process.Kill()
	}()
}
