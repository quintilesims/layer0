package system

import (
	"strconv"
	"testing"
)

const (
	deployCommand = "while true ; do echo LONG RUNNING SERVICE ; sleep 5 ; done"
)

func BenchmarkStress1Environment0Deploys0Services(b *testing.B) {
	benchmarkStress(1, 0, 0, deployCommand, b)
}
func BenchmarkStress5Environments0Deploys0Services(b *testing.B) {
	benchmarkStress(5, 0, 0, deployCommand, b)
}
func BenchmarkStress10Environments0Deploys0Services(b *testing.B) {
	benchmarkStress(10, 0, 0, deployCommand, b)
}
func BenchmarkStress1Environment1Deploy0Service(b *testing.B) {
	benchmarkStress(1, 1, 0, deployCommand, b)
}
func BenchmarkStress1Environment5Deploys0Services(b *testing.B) {
	benchmarkStress(1, 5, 0, deployCommand, b)
}
func BenchmarkStress1Environment10Deploys0Services(b *testing.B) {
	benchmarkStress(1, 10, 0, deployCommand, b)
}
func BenchmarkStress2Environments2Deploys1Service(b *testing.B) {
	benchmarkStress(2, 2, 1, deployCommand, b)
}
func BenchmarkStress2Environments2Deploy5Services(b *testing.B) {
	benchmarkStress(2, 2, 5, deployCommand, b)
}
func BenchmarkStress2Environments2Deploys10Services(b *testing.B) {
	benchmarkStress(2, 2, 10, deployCommand, b)
}

func benchmarkStress(envs int, deps int, servs int, deploycomm string, b *testing.B) {
	tfvars := map[string]string{
		"num_environments": strconv.Itoa(envs),
		"num_deploys":      strconv.Itoa(deps), "num_services": strconv.Itoa(servs),
		"deploy_command": deploycomm,
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

	b.Run("ListDeploys", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			s.Layer0.ListDeploys()
		}
	})

	b.Run("ListLoadBalancers", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			s.Layer0.ListLoadBalancers()
		}
	})

	b.Run("ListServices", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			s.Layer0.ListServices()
		}
	})
}
