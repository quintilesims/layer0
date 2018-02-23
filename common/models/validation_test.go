package models

import (
	"testing"
)

func TestRequestModelValidation(t *testing.T) {
	type Validator interface {
		Validate() error
	}

	containerOverride := func(fn func(*ContainerOverride)) *ContainerOverride {
		c := &ContainerOverride{
			ContainerName:        "c1",
			EnvironmentOverrides: map[string]string{"k1": "v1"},
		}

		fn(c)
		return c
	}

	createDeployRequest := func(fn func(*CreateDeployRequest)) *CreateDeployRequest {
		req := &CreateDeployRequest{
			DeployName: "dpl_name",
			DeployFile: []byte("content"),
		}

		fn(req)
		return req
	}

	createEnvironmentRequest := func(fn func(*CreateEnvironmentRequest)) *CreateEnvironmentRequest {
		req := &CreateEnvironmentRequest{
			EnvironmentName: "env",
			InstanceType:    "t2.small",
			Scale:           3,
			OperatingSystem: "linux",
			AMIID:           "ami123",
		}

		fn(req)
		return req
	}

	createLoadBalancerRequest := func(fn func(*CreateLoadBalancerRequest)) *CreateLoadBalancerRequest {
		req := &CreateLoadBalancerRequest{
			LoadBalancerName: "lb_name",
			EnvironmentID:    "env_id",
			IsPublic:         true,
			Ports: []Port{
				{HostPort: 443, ContainerPort: 80, Protocol: "https", CertificateARN: "cert"},
				{HostPort: 22, ContainerPort: 22, Protocol: "tcp"},
			},
			HealthCheck: HealthCheck{
				Target:             "tcp:80",
				Interval:           5,
				Timeout:            6,
				HealthyThreshold:   7,
				UnhealthyThreshold: 8,
			},
		}

		fn(req)
		return req
	}

	createServiceRequest := func(fn func(*CreateServiceRequest)) *CreateServiceRequest {
		req := &CreateServiceRequest{
			ServiceName:    "svc_name",
			EnvironmentID:  "env_id",
			DeployID:       "dpl_id",
			LoadBalancerID: "lb_id",
			Scale:          3,
		}

		fn(req)
		return req
	}

	createTaskRequest := func(fn func(*CreateTaskRequest)) *CreateTaskRequest {
		req := &CreateTaskRequest{
			TaskName:      "tsk_name",
			EnvironmentID: "env_id",
			DeployID:      "dpl_id",
			ContainerOverrides: []ContainerOverride{
				{ContainerName: "c1", EnvironmentOverrides: map[string]string{"k1": "v1"}},
				{ContainerName: "c2", EnvironmentOverrides: map[string]string{"k2": "v2"}},
			},
		}

		fn(req)
		return req
	}

	healthCheck := func(fn func(*HealthCheck)) *HealthCheck {
		h := &HealthCheck{
			Target:             "tcp:80",
			Interval:           5,
			Timeout:            6,
			HealthyThreshold:   7,
			UnhealthyThreshold: 8,
		}

		fn(h)
		return h
	}

	port := func(fn func(*Port)) *Port {
		p := &Port{
			HostPort:       443,
			ContainerPort:  80,
			Protocol:       "https",
			CertificateARN: "cert",
		}

		fn(p)
		return p
	}

	updateEnvironmentRequest := func(fn func(*UpdateEnvironmentRequest)) *UpdateEnvironmentRequest {
		req := &UpdateEnvironmentRequest{}
		fn(req)
		return req
	}

	updateServiceRequest := func(fn func(*UpdateServiceRequest)) *UpdateServiceRequest {
		req := &UpdateServiceRequest{}
		fn(req)
		return req
	}

	// todo: dynamic environment checks? may not be required depending on changes
	cases := map[string]Validator{
		"ContainerOverride: Missing ContainerName": containerOverride(func(c *ContainerOverride) {
			c.ContainerName = ""
		}),
		"ContainerOverride: Missing key": containerOverride(func(c *ContainerOverride) {
			c.EnvironmentOverrides[""] = "v1"
		}),
		"CreateDeployRequest: Missing DeployName": createDeployRequest(func(req *CreateDeployRequest) {
			req.DeployName = ""
		}),
		"CreateDeployRequest: Missing DeployContent": createDeployRequest(func(req *CreateDeployRequest) {
			req.DeployFile = nil
		}),
		"CreateEnvironmentRequest: Missing EnvironmentName": createEnvironmentRequest(func(req *CreateEnvironmentRequest) {
			req.EnvironmentName = ""
		}),
		"CreateEnvironmentRequest: Missing OperatingSystem": createEnvironmentRequest(func(req *CreateEnvironmentRequest) {
			req.OperatingSystem = ""
		}),
		"CreateEnvironmentRequest: Invalid OperatingSystem": createEnvironmentRequest(func(req *CreateEnvironmentRequest) {
			req.OperatingSystem = "darwin"
		}),
		"CreateEnvironmentRequest: Negative Scale": createEnvironmentRequest(func(req *CreateEnvironmentRequest) {
			req.Scale = -1
		}),
		"CreateLoadBalancerRequest: Missing LoadBalancerName": createLoadBalancerRequest(func(req *CreateLoadBalancerRequest) {
			req.LoadBalancerName = ""
		}),
		"CreateLoadBalancerRequest: Missing EnvironmentID": createLoadBalancerRequest(func(req *CreateLoadBalancerRequest) {
			req.EnvironmentID = ""
		}),
		"CreateServiceRequest: Missing ServiceName": createServiceRequest(func(req *CreateServiceRequest) {
			req.ServiceName = ""
		}),
		"CreateServiceRequest: Missing EnvironmentID": createServiceRequest(func(req *CreateServiceRequest) {
			req.EnvironmentID = ""
		}),
		"CreateServiceRequest: Missing DeployID": createServiceRequest(func(req *CreateServiceRequest) {
			req.DeployID = ""
		}),
		"CreateServiceRequest: Negative Scale": createServiceRequest(func(req *CreateServiceRequest) {
			req.Scale = -1
		}),
		"CreateTaskRequest: Missing TaskName": createTaskRequest(func(req *CreateTaskRequest) {
			req.TaskName = ""
		}),
		"CreateTaskRequest: Missing EnvironmentID": createTaskRequest(func(req *CreateTaskRequest) {
			req.EnvironmentID = ""
		}),
		"CreateTaskRequest: Missing DeployID": createTaskRequest(func(req *CreateTaskRequest) {
			req.DeployID = ""
		}),
		"HealthCheck: Missing Target": healthCheck(func(h *HealthCheck) {
			h.Target = ""
		}),
		"HealthCheck: Missing Interval": healthCheck(func(h *HealthCheck) {
			h.Interval = 0
		}),
		"HealthCheck: Missing Timeout": healthCheck(func(h *HealthCheck) {
			h.Timeout = 0
		}),
		"HealthCheck: Missing HealthyThreshold": healthCheck(func(h *HealthCheck) {
			h.HealthyThreshold = 0
		}),
		"HealthCheck: Missing UnhealthyThreshold": healthCheck(func(h *HealthCheck) {
			h.UnhealthyThreshold = 0
		}),
		"Port: Missing HostPort": port(func(p *Port) {
			p.HostPort = 0
		}),
		"Port: Missing ContainerPort": port(func(p *Port) {
			p.ContainerPort = 0
		}),
		"Port: Missing Protocol": port(func(p *Port) {
			p.Protocol = ""
		}),
		"UpdateEnvironmentRequest: Negative Scale": updateEnvironmentRequest(func(req *UpdateEnvironmentRequest) {
			scale := -1
			req.Scale = &scale
		}),
		"UpdateServiceRequest: Missing DeployID": updateServiceRequest(func(req *UpdateServiceRequest) {
			deployID := ""
			req.DeployID = &deployID
		}),
		"UpdateServiceRequest: Negative Scale": updateServiceRequest(func(req *UpdateServiceRequest) {
			scale := -1
			req.Scale = &scale
		}),
	}

	for name, v := range cases {
		t.Run(name, func(t *testing.T) {
			if err := v.Validate(); err == nil {
				t.Errorf("error was nil!")
			}
		})
	}
}
