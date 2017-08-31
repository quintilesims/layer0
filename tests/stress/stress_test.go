package system

import (
	"strconv"
	"strings"
	"testing"

	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/tests/clients"
	"github.com/quintilesims/tftest"
)

type StressTestCase struct {
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
		NumDeploys:      1,
		NumEnvironments: 2,
		NumServices:     5,
	})
}

func Benchmark25Environments(b *testing.B) {
	runTest(b, StressTestCase{
		NumEnvironments: 25,
	})
}

func Benchmark10Environments10Deploys(b *testing.B) {
	runTest(b, StressTestCase{
		NumDeploys:      10,
		NumEnvironments: 10,
	})
}

func Benchmark20Environments20Deploys(b *testing.B) {
	runTest(b, StressTestCase{
		NumDeploys:      20,
		NumEnvironments: 20,
	})
}

func Benchmark5Environments50Deploys(b *testing.B) {
	runTest(b, StressTestCase{
		NumDeploys:      50,
		NumEnvironments: 5,
	})
}

func Benchmark5Environments100Deploys(b *testing.B) {
	runTest(b, StressTestCase{
		NumDeploys:      100,
		NumEnvironments: 5,
	})
}

func Benchmark10Environments10Deploys10Services(b *testing.B) {
	runTest(b, StressTestCase{
		NumDeploys:      10,
		NumEnvironments: 10,
		NumServices:     10,
	})
}

func Benchmark5Environments5Deploys50Services(b *testing.B) {
	runTest(b, StressTestCase{
		NumDeploys:      5,
		NumEnvironments: 5,
		NumServices:     50,
	})
}

func Benchmark15Environments15Deploys15Services15LoadBalancers(b *testing.B) {
	runTest(b, StressTestCase{
		NumDeploys:       15,
		NumEnvironments:  15,
		NumLoadBalancers: 15,
		NumServices:      15,
	})
}

func Benchmark25Environments25Deploys25Services25LoadBalancers(b *testing.B) {
	runTest(b, StressTestCase{
		NumDeploys:       25,
		NumEnvironments:  25,
		NumLoadBalancers: 25,
		NumServices:      25,
	})
}
