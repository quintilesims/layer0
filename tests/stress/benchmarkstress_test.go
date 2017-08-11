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

func BenchmarkStress(b *testing.B) {
	testCases := []BenchmarkTestCase{
		{1, 0, 0, ""},
		{10, 0, 0, ""},
		{50, 0, 0, ""},
		{1, 1, 1, serviceCommand},
		{1, 1, 10, serviceCommand},
		{1, 1, 50, serviceCommand},
	}

	for _, tc := range testCases {
		tfvars := map[string]string{
			"num_environments": strconv.Itoa(tc.Environments),
			"num_deploys":      strconv.Itoa(tc.Deploys),
			"num_services":     strconv.Itoa(tc.Services),
			"deploy_command":   tc.DeployCommand,
		}

		log.Debugf("Testing with Environments: %v, Deploys: %v, Services: %v, Deploy Command: '%v'", tc.Environments, tc.Deploys, tc.Services, tc.DeployCommand)

		s := NewStressTest(b, "cases/stress", tfvars)
		s.Terraform.Apply()
		defer s.Terraform.Destroy()

		log.Debug("Benchmarking list operations")

		title := fmt.Sprintf("Envs:%v_Deps:%v_Svcs:%v/", tc.Environments, tc.Deploys, tc.Services)

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

		b.Run(title+"ListServices", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				s.Layer0.ListServices()
			}
		})
	}
}
