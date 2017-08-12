package system

import (
	"fmt"
	"strconv"
	"testing"
)

const (
	serviceCommand = "while true ; do echo LONG RUNNING SERVICE ; sleep 10 ; done"
	taskCommand    = "sleep 10"
)

type BenchmarkTestCase struct {
	Environments  int
	Deploys       int
	Services      int
	DeployCommand string
}

/*
func BenchmarkStress1Environment0Deploys0Services(b *testing.B)   { benchmarkStress(1, 0, 0, b) }
func BenchmarkStress5Environments0Deploys0Services(b *testing.B)  { benchmarkStress(5, 0, 0, b) }
func BenchmarkStress10Environments0Deploys0Services(b *testing.B) { benchmarkStress(10, 0, 0, b) }

func BenchmarkStress1Environment1Deploy0Service(b *testing.B)    { benchmarkStress(1, 1, 0, b) }
func BenchmarkStress1Environment5Deploys0Services(b *testing.B)  { benchmarkStress(1, 5, 0, b) }
func BenchmarkStress1Environment10Deploys0Services(b *testing.B) { benchmarkStress(1, 10, 0, b) }

func BenchmarkStress2Environments2Deploys1Service(b *testing.B)   { benchmarkStress(2, 2, 1, b) }
func BenchmarkStress2Environments2Deploy5Services(b *testing.B)   { benchmarkStress(2, 2, 5, b) }
func BenchmarkStress2Environments2Deploys10Services(b *testing.B) { benchmarkStress(2, 2, 10, b) }
*/

func BenchmarkStress(b *testing.B) {

	testCases := []BenchmarkTestCase{
		{1, 0, 0, ""},
		{5, 0, 0, ""},
		{10, 0, 0, ""},
		{1, 1, 0, ""},
		{1, 5, 0, ""},
		{1, 10, 0, ""},
		{2, 2, 1, serviceCommand},
		{2, 2, 5, serviceCommand},
		{2, 2, 10, serviceCommand},
	}

	for _, tc := range testCases {
		tfvars := map[string]string{
			"num_environments": strconv.Itoa(tc.Environments),
			"num_deploys":      strconv.Itoa(tc.Deploys),
			"num_services":     strconv.Itoa(tc.Services),
			"deploy_command":   tc.DeployCommand,
		}

		log.Debugf("Testing with Environments: %v, Deploys: %v, Services: %v, Deploy Command: '%v'", tc.Environments, tc.Deploys, tc.Services, tc.DeployCommand)

		title := fmt.Sprintf("Envs:%v_Deps:%v_Svcs:%v/", tc.Environments, tc.Deploys, tc.Services)

		s := NewStressTest(b, "cases/stress", tfvars)
		s.Terraform.Apply()
		defer s.Terraform.Destroy()

		b.Run(title+"ListEnvironments", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				s.Layer0.ListEnvironments()
			}
		})

		b.Run(title+"ListDeploys", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				s.Layer0.ListDeploys()
			}
		})

		b.Run(title+"ListLoadBalancers", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				s.Layer0.ListLoadBalancers()
			}
		})

		b.Run(title+"ListServices", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				s.Layer0.ListServices()
			}
		})
	}
}
