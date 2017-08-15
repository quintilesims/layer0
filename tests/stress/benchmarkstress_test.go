package system

import (
	"strconv"
	"testing"
)

const (
	deployCommand = "while true ; do echo LONG RUNNING SERVICE ; sleep 5 ; done"
)

func BenchmarkStress5Environments0Deploys0Services(b *testing.B) {
	benchmarkStress(5, 0, 0, deployCommand, b)
}
func BenchmarkStress10Environments0Deploys0Services(b *testing.B) {
	benchmarkStress(10, 0, 0, deployCommand, b)
}
func BenchmarkStress20Environments0Deploys0Services(b *testing.B) {
	benchmarkStress(20, 0, 0, deployCommand, b)
}
func BenchmarkStress5Environments5Deploys0Services(b *testing.B) {
	benchmarkStress(5, 5, 0, deployCommand, b)
}
func BenchmarkStress10Environments10Deploys0Services(b *testing.B) {
	benchmarkStress(10, 10, 0, deployCommand, b)
}
func BenchmarkStress20Environments20Deploys0Services(b *testing.B) {
	benchmarkStress(20, 20, 0, deployCommand, b)
}
func BenchmarkStress10Environments10Deploys10Services(b *testing.B) {
	benchmarkStress(10, 10, 10, deployCommand, b)
}
func BenchmarkStress20Environments20Deploys20Services(b *testing.B) {
	benchmarkStress(20, 20, 20, deployCommand, b)
}
func BenchmarkStress20Environments20Deploys50Services(b *testing.B) {
	benchmarkStress(20, 20, 50, deployCommand, b)
}

func benchmarkStress(envs int, deps int, servs int, deploycomm string, b *testing.B) {
	tfvars := map[string]string{
		"num_environments": strconv.Itoa(envs),
		"num_deploys":      strconv.Itoa(deps),
		"num_services":     strconv.Itoa(servs),
		"deploy_command":   deploycomm,
	}

	log.Debugf("Testing with Environments: %v, Deploys: %v, Services: %v, Command: %v", envs, deps, servs, deploycomm)

	s := NewStressTest(b, "cases/modules", tfvars)
	s.Terraform.Apply()
	defer s.Terraform.Destroy()

	b.Run("ListEnvironments", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			s.Layer0.ListEnvironments()
		}
	})

	b.Run("ListLoadBalancers", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			s.Layer0.ListLoadBalancers()
		}
	})

	b.Run("ListDeploys", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			s.Layer0.ListDeploys()
		}
	})

	b.Run("ListServices", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			s.Layer0.ListServices()
		}
	})

	b.Run("ListTasks", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			s.Layer0.ListTasks()
		}
	})

	b.Run("ListJobs", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			s.Layer0.ListJobs()
		}
	})
}
