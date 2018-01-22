package system

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/tests/clients"
	"github.com/quintilesims/tftest"
)

const (
	DEPLOY_SCALE_MIN        = 1
	DEPLOY_SCALE_MED        = 50
	DEPLOY_SCALE_MAX        = 100
	DEPLOY_FAMILY_SCALE_MIN = 1
	DEPLOY_FAMILY_SCALE_MED = 15
	DEPLOY_FAMILY_SCALE_MAX = 30
	ENVIRONMENT_SCALE_MIN   = 1
	ENVIRONMENT_SCALE_MED   = 15
	ENVIRONMENT_SCALE_MAX   = 30
	LOAD_BALANCER_SCALE_MIN = 1
	LOAD_BALANCER_SCALE_MED = 25
	LOAD_BALANCER_SCALE_MAX = 50
	SERVICE_SCALE_MIN       = 1
	SERVICE_SCALE_MED       = 25
	SERVICE_SCALE_MAX       = 50
	TASK_SCALE_MIN          = 1
	TASK_SCALE_MED          = 50
	TASK_SCALE_MAX          = 100
)

type StressTestCase struct {
	NumDeploys        int
	NumDeployFamilies int
	NumEnvironments   int
	NumLoadBalancers  int
	NumServices       int
	NumTasks          int
}

func runTest(b *testing.B, c StressTestCase) {
	if c.NumTasks > 0 || c.NumServices > 0 {
		if c.NumEnvironments == 0 {
			c.NumEnvironments = ENVIRONMENT_SCALE_MIN
		}

		if c.NumDeploys == 0 {
			c.NumDeploys = DEPLOY_SCALE_MIN
		}
	}

	if c.NumLoadBalancers > 0 && c.NumEnvironments == 0 {
		c.NumEnvironments = ENVIRONMENT_SCALE_MIN
	}

	vars := map[string]string{
		"endpoint":            config.FLAG_ENDPOINT,
		"token":               config.FLAG_TOKEN,
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

	methodsToBenchmark := map[string]func(){}

	if c.NumDeploys > 0 {
		deployIDs := strings.Split(terraform.Output("deploy_ids"), ",\n")
		methodsToBenchmark["ReadDeploy"] = func() { layer0.ReadDeploy(deployIDs[0]) }
		methodsToBenchmark["ListDeploys"] = func() { layer0.ListDeploys() }
	}

	if c.NumEnvironments > 0 {
		environmentIDs := strings.Split(terraform.Output("environment_ids"), ",\n")
		methodsToBenchmark["ReadEnvironment"] = func() { layer0.ReadEnvironment(environmentIDs[0]) }
		methodsToBenchmark["ListEnvironments"] = func() { layer0.ListEnvironments() }
	}

	if c.NumLoadBalancers > 0 {
		loadBalancerIDs := strings.Split(terraform.Output("load_balancer_ids"), ",\n")
		methodsToBenchmark["ReadLoadBalancer"] = func() { layer0.ReadLoadBalancer(loadBalancerIDs[0]) }
		methodsToBenchmark["ListLoadBalancers"] = func() { layer0.ListLoadBalancers() }
	}

	if c.NumServices > 0 {
		serviceIDs := strings.Split(terraform.Output("service_ids"), ",\n")
		methodsToBenchmark["ReadService"] = func() { layer0.ReadService(serviceIDs[0]) }
		methodsToBenchmark["ListServices"] = func() { layer0.ListServices() }
	}

	if c.NumTasks > 0 {
		methodsToBenchmark["ListTasks"] = func() { layer0.ListTasks() }

		deployIDs := strings.Split(terraform.Output("deploy_ids"), ",\n")
		environmentIDs := strings.Split(terraform.Output("environment_ids"), ",\n")

		for i := 0; i < c.NumTasks; i++ {
			environmentID := environmentIDs[i%len(environmentIDs)]
			deployID := deployIDs[i%len(deployIDs)]
			taskName := fmt.Sprintf("tt%v", i)
			layer0.CreateTask(taskName, environmentID, deployID, 1, nil)
		}
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

func BenchmarkMinFamiliesMinDeploys(b *testing.B) {
	runTest(b, StressTestCase{
		NumDeployFamilies: DEPLOY_FAMILY_SCALE_MIN,
		NumDeploys:        DEPLOY_SCALE_MIN,
	})
}

func BenchmarkMedFamiliesMedDeploys(b *testing.B) {
	runTest(b, StressTestCase{
		NumDeployFamilies: DEPLOY_FAMILY_SCALE_MED,
		NumDeploys:        DEPLOY_SCALE_MED,
	})
}

func BenchmarkMaxFamiliesMaxDeploys(b *testing.B) {
	runTest(b, StressTestCase{
		NumDeployFamilies: DEPLOY_FAMILY_SCALE_MAX,
		NumDeploys:        DEPLOY_SCALE_MAX,
	})
}

func BenchmarkMinEnvironments(b *testing.B) {
	runTest(b, StressTestCase{
		NumEnvironments: ENVIRONMENT_SCALE_MIN,
	})
}

func BenchmarkMedEnvironments(b *testing.B) {
	runTest(b, StressTestCase{
		NumEnvironments: ENVIRONMENT_SCALE_MED,
	})
}

func BenchmarkMaxEnvironments(b *testing.B) {
	runTest(b, StressTestCase{
		NumEnvironments: ENVIRONMENT_SCALE_MAX,
	})
}

func BenchmarkMinServices(b *testing.B) {
	runTest(b, StressTestCase{
		NumServices: SERVICE_SCALE_MIN,
	})
}

func BenchmarkMedServices(b *testing.B) {
	runTest(b, StressTestCase{
		NumServices: SERVICE_SCALE_MED,
	})
}

func BenchmarkMaxServices(b *testing.B) {
	runTest(b, StressTestCase{
		NumServices: SERVICE_SCALE_MAX,
	})
}

func BenchmarkMinLoadBalancers(b *testing.B) {
	runTest(b, StressTestCase{
		NumLoadBalancers: LOAD_BALANCER_SCALE_MIN,
	})
}

func BenchmarkMedLoadBalancers(b *testing.B) {
	runTest(b, StressTestCase{
		NumLoadBalancers: LOAD_BALANCER_SCALE_MED,
	})
}

func BenchmarkMaxLoadBalancers(b *testing.B) {
	runTest(b, StressTestCase{
		NumLoadBalancers: LOAD_BALANCER_SCALE_MAX,
	})
}

func BenchmarkMinTasks(b *testing.B) {
	runTest(b, StressTestCase{
		NumTasks: TASK_SCALE_MIN,
	})
}

func BenchmarkMedTasks(b *testing.B) {
	runTest(b, StressTestCase{
		NumTasks: TASK_SCALE_MED,
	})
}

func BenchmarkMaxTasks(b *testing.B) {
	runTest(b, StressTestCase{
		NumTasks: TASK_SCALE_MAX,
	})
}

func BenchmarkAggregateMin(b *testing.B) {
	runTest(b, StressTestCase{
		NumDeploys:        DEPLOY_SCALE_MIN,
		NumDeployFamilies: DEPLOY_FAMILY_SCALE_MIN,
		NumEnvironments:   ENVIRONMENT_SCALE_MIN,
		NumLoadBalancers:  LOAD_BALANCER_SCALE_MIN,
		NumServices:       SERVICE_SCALE_MIN,
	})
}

func BenchmarkAggregateMed(b *testing.B) {
	runTest(b, StressTestCase{
		NumDeploys:        DEPLOY_SCALE_MED,
		NumDeployFamilies: DEPLOY_FAMILY_SCALE_MED,
		NumEnvironments:   ENVIRONMENT_SCALE_MED,
		NumLoadBalancers:  LOAD_BALANCER_SCALE_MED,
		NumServices:       SERVICE_SCALE_MED,
	})
}

func BenchmarkAggregateMax(b *testing.B) {
	runTest(b, StressTestCase{
		NumDeploys:        DEPLOY_SCALE_MAX,
		NumDeployFamilies: DEPLOY_FAMILY_SCALE_MAX,
		NumEnvironments:   ENVIRONMENT_SCALE_MAX,
		NumLoadBalancers:  LOAD_BALANCER_SCALE_MAX,
		NumServices:       SERVICE_SCALE_MAX,
	})
}
