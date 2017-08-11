package system

import (
	"strconv"
	"testing"
)

const (
	serviceCommand = "while true ; do echo LONG RUNNING SERVICE ; sleep 10 ; done"
	taskCommand    = "sleep 10"
)

func BenchmarkStress1Environment0Deploys0Services(b *testing.B) {
	benchmarkStress(1, 0, 0, "", b)
}
func BenchmarkStress10Environments0Deploys0Services(b *testing.B) {
	benchmarkStress(10, 0, 0, "", b)
}
func BenchmarkStress50Environments0Deploys0Services(b *testing.B) {
	benchmarkStress(50, 0, 0, "", b)
}

func BenchmarkStress1Environment1Deploy1Service(b *testing.B) {
	benchmarkStress(1, 1, 1, serviceCommand, b)
}
func BenchmarkStress1Environment1Deploy5Services(b *testing.B) {
	benchmarkStress(1, 1, 5, serviceCommand, b)
}
func BenchmarkStress1Environment1Deploy10Services(b *testing.B) {
	benchmarkStress(1, 1, 10, serviceCommand, b)
}

func benchmarkStress(env int, dep int, ser int, cmd string, b *testing.B) {
	tfvars := map[string]string{
		"num_environments": strconv.Itoa(env),
		"num_deploys":      strconv.Itoa(dep),
		"num_services":     strconv.Itoa(ser),
		"deploy_command":   cmd,
	}

	log.Debugf("Testing with Environments: %v, Deploys: %v, Services: %v, Deploy Command: '%v'", env, dep, ser, cmd)

	s := NewStressTest(b, "cases/stress", tfvars)
	s.Terraform.Apply()
	defer s.Terraform.Destroy()

	log.Debug("Benchmarking list operations")

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

	b.Run("ListServices", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			s.Layer0.ListServices()
		}
	})
}
