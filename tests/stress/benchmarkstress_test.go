package system

import (
	"strconv"
	"testing"
)

const (
	deployCommand = "while true ; do echo LONG RUNNING SERVICE ; sleep 5 ; done"
)

func Benchmark25Environments(b *testing.B) {
	benchmarkStress(b, 25, 0, 0, 0)
}
func Benchmark10Environments10Deploys(b *testing.B) {
	benchmarkStress(b, 10, 10, 0, 0)
}
func Benchmark20Environments20Deploys(b *testing.B) {
	benchmarkStress(b, 20, 20, 0, 0)
}
func Benchmark5Environments50Deploys(b *testing.B) {
	benchmarkStress(b, 5, 50, 0, 0)
}
func Benchmark5Environments100Deploys(b *testing.B) {
	benchmarkStress(b, 5, 100, 0, 0)
}
func Benchmark10Environments10Deploys10Services(b *testing.B) {
	benchmarkStress(b, 10, 10, 10, 0)
}
func Benchmark5Environments5Deploys50Services(b *testing.B) {
	benchmarkStress(b, 5, 5, 50, 0)
}
func Benchmark15Environments15Deploys15Services15LoadBalancers(b *testing.B) {
	benchmarkStress(b, 15, 15, 15, 15)
}
func Benchmark25Environments25Deploys25Services25LoadBalancers(b *testing.B) {
	benchmarkStress(b, 25, 25, 25, 25)
}

func benchmarkStress(b *testing.B, envs, deps, servs, lbs int) {
	tfvars := map[string]string{
		"num_environments":  strconv.Itoa(envs),
		"num_deploys":       strconv.Itoa(deps),
		"num_services":      strconv.Itoa(servs),
		"num_loadbalancers": strconv.Itoa(lbs),
		"deploy_command":    deployCommand,
	}

	log.Debugf("Testing with params: %v", tfvars)

	terraform, layer0 := NewStressTest(b, "cases/modules", tfvars)
	terraform.Apply()
	defer terraform.Destroy()

	b.Run("ListEnvironments", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			layer0.ListEnvironments()
		}
	})

	b.Run("ListLoadBalancers", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			layer0.ListLoadBalancers()
		}
	})

	b.Run("ListDeploys", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			layer0.ListDeploys()
		}
	})

	b.Run("ListServices", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			layer0.ListServices()
		}
	})

	b.Run("ListTasks", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			layer0.ListTasks()
		}
	})

	b.Run("ListJobs", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			layer0.ListJobs()
		}
	})
}
