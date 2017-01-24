package framework

import (
	"fmt"
	"github.com/quintilesims/layer0/cli/client"
	"github.com/quintilesims/layer0/cli/command"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	stsclient "github.com/quintilesims/layer0/tests/system/sts/client"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"testing"
)

type Config struct {
	T       *testing.T
	Dir     string
	DryRun  bool
	Verbose bool
	Vars    map[string]string
}

type SystemTestContext struct {
	T        *testing.T
	Dir      string
	DryRun   bool
	Verbose  bool
	Vars     map[string]string
	Client   *client.APIClient
	Resolver *command.TagResolver
}

func NewSystemTestContext(c Config) *SystemTestContext {
	apiClient := client.NewAPIClient(client.Config{
		Endpoint: config.APIEndpoint(),
		Token:    fmt.Sprintf("Basic %s", config.AuthToken()),
	})

	return &SystemTestContext{
		T:        c.T,
		Dir:      c.Dir,
		DryRun:   c.DryRun,
		Verbose:  c.Verbose,
		Vars:     c.Vars,
		Client:   apiClient,
		Resolver: command.NewTagResolver(apiClient),
	}
}

func (s *SystemTestContext) GetSystemTestService(environmentName, loadBalancerName string) *stsclient.SystemTestService {
	env := s.GetEnvironment(environmentName)
	lb := s.GetLoadBalancer(env.EnvironmentID, loadBalancerName)
	return stsclient.NewSystemTestService(s.T, lb.URL)
}

func (s *SystemTestContext) GetEnvironment(target string) *models.Environment {
	id := s.resolve("environment", target)
	environment, err := s.Client.GetEnvironment(id)
	if err != nil {
		s.T.Fatal(err)
	}

	return environment
}

func (s *SystemTestContext) GetLoadBalancer(environment, target string) *models.LoadBalancer {
	if environment != "" {
		target = fmt.Sprintf("%s:%s", environment, target)
	}

	id := s.resolve("load_balancer", target)
	loadBalancer, err := s.Client.GetLoadBalancer(id)
	if err != nil {
		s.T.Fatal(err)
	}

	return loadBalancer
}

func (s *SystemTestContext) GetService(environment, target string) *models.Service {
	if environment != "" {
		target = fmt.Sprintf("%s:%s", environment, target)
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
	if s.DryRun {
		s.runTerraform("plan")
		return
	}

	s.runTerraform("apply")
}

func (s *SystemTestContext) Destroy() {
	if s.DryRun {
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
	sigtermChan := make(chan os.Signal)
	signal.Notify(sigtermChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigtermChan
		cmd.Process.Kill()
	}()

	if s.Verbose {
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
