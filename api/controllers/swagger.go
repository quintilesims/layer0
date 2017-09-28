package controllers

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/zpatrick/fireball"
	swagger "github.com/zpatrick/go-plugin-swagger"
)

type SwaggerController struct {
	version string
}

func NewSwaggerController(version string) *SwaggerController {
	return &SwaggerController{
		version: version,
	}
}

func (s *SwaggerController) Routes() []*fireball.Route {
	routes := []*fireball.Route{
		{
			Path: "/swagger.json",
			Handlers: fireball.Handlers{
				"GET": s.ServeSwaggerSpec,
			},
		},
	}

	return routes
}

func (s *SwaggerController) ServeSwaggerSpec(c *fireball.Context) (fireball.Response, error) {
	spec := swagger.Spec{
		SwaggerVersion: "2.0",
		Host:           c.Request.Host,
		Schemes:        []string{"https"},
		Info: &swagger.Info{
			Title:   "Layer0",
			Version: s.version,
		},
		Definitions: map[string]swagger.Definition{
			"Container":                 models.Container{}.Definition(),
			"CreateEnvironmentRequest":  models.CreateEnvironmentRequest{}.Definition(),
			"CreateLoadBalancerRequest": models.CreateLoadBalancerRequest{}.Definition(),
			"CreateTaskRequest":         models.CreateTaskRequest{}.Definition(),
			"Deployment":                models.Deployment{}.Definition(),
			"Environment":               models.Environment{}.Definition(),
			"HealthCheck":               models.HealthCheck{}.Definition(),
			"Job":                       models.Job{}.Definition(),
			"LoadBalancer":              models.LoadBalancer{}.Definition(),
			"LogFile":                   models.LogFile{}.Definition(),
			"Port":                      models.Port{}.Definition(),
			"Service":                   models.Service{}.Definition(),
			"Task":                      models.Task{}.Definition(),
			"UpdateLoadBalancerRequest": models.UpdateLoadBalancerRequest{}.Definition(),
			"UpdateServiceRequest":      models.UpdateServiceRequest{}.Definition(),
		},
		Tags: []swagger.Tag{
			{
				Name:        "Environment",
				Description: "Methods related to environments",
			},
			{
				Name:        "Job",
				Description: "Methods related to jobs",
			},
			{
				Name:        "LoadBalancer",
				Description: "Methods related to load balancers",
			},
			{
				Name:        "Service",
				Description: "Methods related to services",
			},
			{
				Name:        "Task",
				Description: "Methods related to tasks",
			},
		},
		Paths: map[string]swagger.Path{
			"/environment": map[string]swagger.Method{
				"get": {
					Summary: "List all Environments",
					Tags:    []string{"Environment"},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "An array of environments",
							Schema:      swagger.NewObjectSliceSchema("Environment"),
						},
					},
				},
				"post": {
					Summary: "Add an Environment",
					Tags:    []string{"Environment"},
					Parameters: []swagger.Parameter{
						swagger.NewBodyParam("CreateEnvironmentRequest", "Environment to add", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "The added environment",
							Schema:      swagger.NewObjectSchema("Environment"),
						},
					},
				},
			},
			"/environment/{id}": map[string]swagger.Method{
				"get": {
					Summary: "Describe an Environment",
					Tags:    []string{"Environment"},
					Parameters: []swagger.Parameter{
						swagger.NewStringPathParam("id", "ID of the environment to describe", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "The desired environment",
							Schema:      swagger.NewObjectSchema("Environment"),
						},
					},
				},
				"delete": {
					Summary: "Delete an Environment",
					Tags:    []string{"Environment"},
					Parameters: []swagger.Parameter{
						swagger.NewStringPathParam("id", "ID of the environment to delete", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "Success",
						},
					},
				},
			},
			"/job": map[string]swagger.Method{
				"get": {
					Summary: "List all Jobs",
					Tags:    []string{"Job"},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "An array of jobs",
							Schema:      swagger.NewObjectSliceSchema("Job"),
						},
					},
				},
			},
			"/job/{id}": map[string]swagger.Method{
				"get": {
					Summary: "Describe a Job",
					Tags:    []string{"Job"},
					Parameters: []swagger.Parameter{
						swagger.NewStringPathParam("id", "ID of the job to describe", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "The desired job",
							Schema:      swagger.NewObjectSchema("Job"),
						},
					},
				},
				"delete": {
					Summary: "Delete a Job",
					Tags:    []string{"Job"},
					Parameters: []swagger.Parameter{
						swagger.NewStringPathParam("id", "ID of the job to delete", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "Success",
						},
					},
				},
			},
			"/loadbalancer": map[string]swagger.Method{
				"get": {
					Summary: "List all LoadBalancers",
					Tags:    []string{"LoadBalancer"},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "An array of load balancers",
							Schema:      swagger.NewObjectSliceSchema("LoadBalancer"),
						},
					},
				},
				"post": {
					Summary: "Add a LoadBalancer",
					Tags:    []string{"LoadBalancer"},
					Parameters: []swagger.Parameter{
						swagger.NewBodyParam("CreateLoadBalancerRequest", "LoadBalancer to add", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "The added load balancer",
							// todo: this actually returns a CreateJobResponse
							Schema: swagger.NewObjectSchema("LoadBalancer"),
						},
					},
				},
				"put": {
					Summary: "Update a LoadBalancer",
					Tags:    []string{"LoadBalancer"},
					Parameters: []swagger.Parameter{
						swagger.NewBodyParam("UpdateLoadBalancerRequest", "LoadBalancer to update", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "The updated load balancer",
							Schema:      swagger.NewObjectSchema("LoadBalancer"),
						},
					},
				},
			},
			"/loadbalancer/{id}": map[string]swagger.Method{
				"get": {
					Summary: "Describe a LoadBalancer",
					Tags:    []string{"LoadBalancer"},
					Parameters: []swagger.Parameter{
						swagger.NewStringPathParam("id", "ID of the load balancer to describe", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "The desired load balancer",
							Schema:      swagger.NewObjectSchema("LoadBalancer"),
						},
					},
				},
				"delete": {
					Summary: "Delete a LoadBalancer",
					Tags:    []string{"LoadBalancer"},
					Parameters: []swagger.Parameter{
						swagger.NewStringPathParam("id", "ID of the load balancer to delete", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "Success",
						},
					},
				},
			},
			"/service": {
				"put": {
					Summary: "Update a Service",
					Tags:    []string{"Service"},
					Parameters: []swagger.Parameter{
						swagger.NewBodyParam("UpdateServiceRequest", "Service to update", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "The updated Service",
							Schema:      swagger.NewObjectSchema("Service"),
						},
					},
				},
			},
			"/service/{id}": {
				"get": {
					Summary: "Describe a Service",
					Tags:    []string{"Service"},
					Parameters: []swagger.Parameter{
						swagger.NewStringPathParam("id", "ID of the service to describe", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "The desired service",
							Schema:      swagger.NewObjectSchema("Service"),
						},
					},
				},
			},
			"/task": map[string]swagger.Method{
				"get": {
					Summary: "List all Tasks",
					Tags:    []string{"Task"},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "An array of tasks",
							Schema:      swagger.NewObjectSliceSchema("Task"),
						},
					},
				},
				"post": {
					Summary: "Add a Task",
					Tags:    []string{"Task"},
					Parameters: []swagger.Parameter{
						swagger.NewBodyParam("CreateTaskRequest", "Task to add", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "The added task",
							Schema:      swagger.NewObjectSchema("Task"),
						},
					},
				},
			},
			"/task/{id}": map[string]swagger.Method{
				"get": {
					Summary: "Describe a Task",
					Tags:    []string{"Task"},
					Parameters: []swagger.Parameter{
						swagger.NewStringPathParam("id", "ID of the task to describe", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "The desired task",
							Schema:      swagger.NewObjectSchema("Task"),
						},
					},
				},
				"delete": {
					Summary: "Delete a Task",
					Tags:    []string{"Task"},
					Parameters: []swagger.Parameter{
						swagger.NewStringPathParam("id", "ID of the task to delete", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "Success",
						},
					},
				},
			},
			"/task/{id}/logs": map[string]swagger.Method{
				"get": {
					Summary: "Get task logs",
					Tags:    []string{"Task"},
					Parameters: []swagger.Parameter{
						swagger.NewStringPathParam("id", "ID of the task to describe", true),
						swagger.NewIntQueryParam("tail", "The number of lines from the end to return", false),
						swagger.NewStringQueryParam("start", "The start of the time range to fetch logs (format YYYY-MM-DD HH:MM)", false),
						swagger.NewStringQueryParam("end", "The end of the time range to fetch logs (format YYYY-MM-DD HH:MM)", false),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "The task's logs",
							Schema:      swagger.NewObjectSliceSchema("LogFile"),
						},
					},
				},
			},
		},
	}

	return fireball.NewJSONResponse(200, spec)
}
