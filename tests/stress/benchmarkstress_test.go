package system

import (
	"strconv"
	"testing"
)

const (
	deployCommand = "while true ; do echo LONG RUNNING SERVICE ; sleep 5 ; done"
)

func BenchmarkStress25Environments0Deploys0Services0LoadBalancers(b *testing.B) {
	benchmarkStress(25, 0, 0, 0, deployCommand, b)
}
func BenchmarkStress10Environments10Deploys0Services0LoadBalancers(b *testing.B) {
	benchmarkStress(10, 10, 0, 0, deployCommand, b)
}
func BenchmarkStress20Environments20Deploys0Services0LoadBalancers(b *testing.B) {
	benchmarkStress(20, 20, 0, 0, deployCommand, b)
}
func BenchmarkStress5Environments50Deploys0Services0LoadBalancers(b *testing.B) {
	benchmarkStress(5, 50, 0, 0, deployCommand, b)
}
func BenchmarkStress5Environments100Deploys0Services0LoadBalancers(b *testing.B) {
	benchmarkStress(5, 100, 0, 0, deployCommand, b)
}
func BenchmarkStress10Environments10Deploys10Services0LoadBalancers(b *testing.B) {
	benchmarkStress(10, 10, 10, 0, deployCommand, b)
}
func BenchmarkStress5Environments5Deploys50Services0LoadBalancers(b *testing.B) {
	benchmarkStress(5, 5, 50, 0, deployCommand, b)
}
func BenchmarkStress15Environments15Deploys15Services15LoadBalancers(b *testing.B) {
	benchmarkStress(15, 15, 15, 15, deployCommand, b)
}
func BenchmarkStress25Environments25Deploys25Services25LoadBalancers(b *testing.B) {
	benchmarkStress(25, 25, 25, 25, deployCommand, b)
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

	methodsToBenchmark := map[string]func(){
		"ListEnvironments":  func() { s.Layer0.ListEnvironments() },
		"ListLoadBalancers": func() { s.Layer0.ListLoadBalancers() },
		"ListDeploys":       func() { s.Layer0.ListDeploys() },
		"ListServices":      func() { s.Layer0.ListServices() },
		"ListTasks":         func() { s.Layer0.ListTasks() },
		"ListJobs":          func() { s.Layer0.ListJobs() },
	}

	benchmark(b, methodsToBenchmark)
}
