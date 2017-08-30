package system

import (
	"strconv"
	"strings"
	"testing"

	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/tests/clients"
	"github.com/quintilesims/tftest"
)

const (
	serviceDeployCommand = "while true ; do echo LONG RUNNING SERVICE ; sleep 5 ; done"
	taskDeployCommand    = "echo SHORT RUNNING TASK ; sleep 10"
)

type StressTestCase struct {
	DeployCommand     string
	NumDeploys        int
	NumDeployFamilies int
	NumEnvironments   int
	NumLoadBalancers  int
	NumServices       int
}

func runTest(b *testing.B, c StressTestCase) {
	vars := map[string]string{
		"endpoint":            config.APIEndpoint(),
		"token":               config.AuthToken(),
		"deploy_command":      c.DeployCommand,
		"num_deploys":         strconv.Itoa(c.NumDeploys),
		"num_deploy_families": strconv.Itoa(c.NumDeployFamilies),
		"num_environments":    strconv.Itoa(c.NumEnvironments),
		"num_load_balancers":  strconv.Itoa(c.NumLoadBalancers),
		"num_services":        strconv.Itoa(c.NumServices),
	}

	terraform := tftest.NewTestContext(
		b,
		tftest.Dir("module"),
		tftest.Vars(vars),
		tftest.DryRun(*dry),
		tftest.Log(b),
	)

	layer0 := clients.NewLayer0TestClient(b, vars["endpoint"], vars["token"])

	terraform.Apply()
	defer terraform.Destroy()

	methodsToBenchmark := map[string]func(){
		"ListDeploys":       func() { layer0.ListDeploys() },
		"ListEnvironments":  func() { layer0.ListEnvironments() },
		"ListLoadBalancers": func() { layer0.ListLoadBalancers() },
		"ListServices":      func() { layer0.ListServices() },
		"ListTasks":         func() { layer0.ListTasks() },
	}

	if c.NumDeploys > 0 {
		deployIDs := strings.Split(terraform.Output("deploy_ids"), ",\n")
		methodsToBenchmark["GetDeploy"] = func() { layer0.GetDeploy(deployIDs[0]) }
	}

	if c.NumEnvironments > 0 {
		environmentIDs := strings.Split(terraform.Output("environment_ids"), ",\n")
		methodsToBenchmark["GetEnvironment"] = func() { layer0.GetEnvironment(environmentIDs[0]) }
	}

	if c.NumLoadBalancers > 0 {
		loadBalancerIDs := strings.Split(terraform.Output("load_balancer_ids"), ",\n")
		methodsToBenchmark["GetLoadBalancer"] = func() { layer0.GetLoadBalancer(loadBalancerIDs[0]) }
	}

	if c.NumServices > 0 {
		serviceIDs := strings.Split(terraform.Output("service_ids"), ",\n")
		methodsToBenchmark["GetService"] = func() { layer0.GetService(serviceIDs[0]) }
	}

	benchmark(b, methodsToBenchmark)
}

func benchmark(b *testing.B, methods map[string]func()) {
	for name, fn := range methods {
		b.Run(name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				fn()
			}
		})
	}
}

func Benchmark5Services(b *testing.B) {
	runTest(b, StressTestCase{
		DeployCommand:   serviceDeployCommand,
		NumDeploys:      1,
		NumEnvironments: 2,
		NumServices:     5,
	})
}
