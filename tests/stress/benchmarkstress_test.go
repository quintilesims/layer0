package system

import (
	"strconv"
	"testing"
)

const (
	deployCommand = "while true ; do echo LONG RUNNING SERVICE ; sleep 5 ; done"
)

func BenchmarkStress5Environments0Deploys0Services0LoadBalancers(b *testing.B) {
	benchmarkStress(5, 0, 0, 0, deployCommand, b)
}
func BenchmarkStress10Environments0Deploys0Services0LoadBalancers(b *testing.B) {
	benchmarkStress(10, 0, 0, 0, deployCommand, b)
}
func BenchmarkStress15Environments0Deploys0Services0LoadBalancers(b *testing.B) {
	benchmarkStress(15, 0, 0, 0, deployCommand, b)
}
func BenchmarkStress5Environments5Deploys0Services0LoadBalancers(b *testing.B) {
	benchmarkStress(5, 5, 0, 0, deployCommand, b)
}
func BenchmarkStress10Environments10Deploys0Services0LoadBalancers(b *testing.B) {
	benchmarkStress(10, 10, 0, 0, deployCommand, b)
}
func BenchmarkStress15Environments15Deploys0Services0LoadBalancers(b *testing.B) {
	benchmarkStress(15, 15, 0, 0, deployCommand, b)
}
func BenchmarkStress10Environments10Deploys10Services0LoadBalancers(b *testing.B) {
	benchmarkStress(10, 10, 10, 0, deployCommand, b)
}
func BenchmarkStress10Environments10Deploys20Services0LoadBalancers(b *testing.B) {
	benchmarkStress(10, 10, 20, 0, deployCommand, b)
}
func BenchmarkStress15Environments15Deploys15Services0LoadBalancers(b *testing.B) {
	benchmarkStress(15, 15, 15, 0, deployCommand, b)
}
func BenchmarkStress15Environments15Deploys30Services0LoadBalancers(b *testing.B) {
	benchmarkStress(15, 15, 30, 0, deployCommand, b)
}
func BenchmarkStress15Environments15Deploys15Services5LoadBalancers(b *testing.B) {
	benchmarkStress(15, 15, 15, 5, deployCommand, b)
}
func BenchmarkStress15Environments15Deploys15Services10LoadBalancers(b *testing.B) {
	benchmarkStress(15, 15, 15, 10, deployCommand, b)
}
func BenchmarkStress15Environments15Deploys15Services15LoadBalancers(b *testing.B) {
	benchmarkStress(15, 15, 15, 15, deployCommand, b)
}

func benchmarkStress(envs, deps, servs, lbs int, deploycomm string, b *testing.B) {
	tfvars := map[string]string{
		"num_environments":  strconv.Itoa(envs),
		"num_deploys":       strconv.Itoa(deps),
		"num_services":      strconv.Itoa(servs),
		"num_loadbalancers": strconv.Itoa(lbs),
		"deploy_command":    deploycomm,
	}

	log.Debugf("Testing with Environments: %v, Deploys: %v, Services: %v, Load Balancers: %v, Command: %v", envs, deps, servs, lbs, deploycomm)

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
