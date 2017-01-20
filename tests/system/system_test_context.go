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
	"strings"
	"syscall"
	"testing"
	"time"
)

func startSystemTest(t *testing.T, dir string, vars map[string]string) *SystemTestContext {
	t.Parallel()

	if vars == nil {
		vars = map[string]string{}
	}

	// add default terraform variables
	vars["endpoint"] = config.APIEndpoint()
	vars["token"] = config.AuthToken()

	c := NewSystemTestContext(t, dir, vars)
	c.Apply()

	return c
}

type SystemTestContext struct {
	T        *testing.T
	Dir      string
	Client   *client.APIClient
	Resolver *command.TagResolver
	Vars     map[string]string
}

func NewSystemTestContext(t *testing.T, dir string, vars map[string]string) *SystemTestContext {
	apiClient := client.NewAPIClient(client.Config{
		Endpoint: config.APIEndpoint(),
		Token:    fmt.Sprintf("Basic %s", config.AuthToken()),
	})

	return &SystemTestContext{
		T:        t,
		Dir:      dir,
		Client:   apiClient,
		Resolver: command.NewTagResolver(apiClient),
		Vars:     vars,
	}
}

func (s *SystemTestContext) WaitForAllDeployments(timeout time.Duration) {
	services, err := s.Client.ListServices()
	if err != nil {
		s.T.Fatal(err)
	}

	for _, service := range services {
		print("Waiting for service deployment ", service.ServiceID)
		if _, err := s.Client.WaitForDeployment(service.ServiceID, timeout); err != nil {
			s.T.Fatal(err)
		}
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

func (s *SystemTestContext) GetLoadBalancer(environmentID, target string) *models.LoadBalancer {
	if environmentID != "" {
		target = fmt.Sprintf("%s:%s", environmentID, target)
	}

	id := s.resolve("load_balancer", target)
	loadBalancer, err := s.Client.GetLoadBalancer(id)
	if err != nil {
		s.T.Fatal(err)
	}

	return loadBalancer
}

func (s *SystemTestContext) GetService(environmentID, target string) *models.Service {
	if environmentID != "" {
		target = fmt.Sprintf("%s:%s", environmentID, target)
	}

	id := s.resolve("service", target)
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
	if *dry {
		s.runTerraform("plan")
		return
	}

	s.runTerraform("apply")
}

func (s *SystemTestContext) Destroy() {
	if *dry {
		s.runTerraform("plan", "-destroy")
		return
	}

	s.runTerraform("destroy", "-force")
}

func (s *SystemTestContext) runTerraform(args ...string) {
	// set terraform environment variables using TF_VAR_<var>
	env := []string{}
	for k, v := range s.Vars {
		env = append(env, fmt.Sprintf("TF_VAR_%s=%s", k, v))
	}

	cmd := exec.Command("terraform", args...)
	cmd.Dir = s.Dir
	cmd.Env = env

	// kill the process if a SIGTERM signal is sent
	sigtermChan := make(chan os.Signal, 2)
	signal.Notify(sigtermChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigtermChan
		cmd.Process.Kill()
	}()

	if verbose {
		fmt.Printf("Running terraform %s from %s \n", args[0], cmd.Dir)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		text := fmt.Sprintf("Error running terraform %s from %s: %v\n", args[0], cmd.Dir, err)
		for _, line := range strings.Split(string(output), "\n") {
			text += line + "\n"
		}

		s.T.Fatal(text)
	}
}
