package system

import (
	"strconv"
	"testing"
)

func BenchmarkStress1Environment0Deploys0Services(b *testing.B)   { benchmarkStress(1, 0, 0, b) }
func BenchmarkStress5Environments0Deploys0Services(b *testing.B)  { benchmarkStress(5, 0, 0, b) }
func BenchmarkStress10Environments0Deploys0Services(b *testing.B) { benchmarkStress(10, 0, 0, b) }

func BenchmarkStress1Environment1Deploy1Service(b *testing.B)   { benchmarkStress(1, 1, 1, b) }
func BenchmarkStress1Environment1Deploy5Services(b *testing.B)  { benchmarkStress(1, 1, 5, b) }
func BenchmarkStress1Environment1Deploy10Services(b *testing.B) { benchmarkStress(1, 1, 10, b) }

func benchmarkStress(env int, dep int, ser int, b *testing.B) {
	tfvars := map[string]string{
		"num_environments": strconv.Itoa(env),
		"num_deploys":      strconv.Itoa(dep),
		"num_services":     strconv.Itoa(ser),
	}

	log.Debugf("Testing with Environments: %v, Deploys: %v, Services: %v", env, dep, ser)

	s := NewStressTest(b, "cases/stress", tfvars)
	s.Terraform.Apply()
	defer s.Terraform.Destroy()

	log.Debug("Benchmarking list operations")

	for n := 0; n < b.N; n++ {
		s.Layer0.ListEnvironments()
	}
}
