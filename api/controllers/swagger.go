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
			Path: "/",
			Handlers: fireball.Handlers{
				"GET": s.Root,
			},
		},
		{
			Path: "/swagger.json",
			Handlers: fireball.Handlers{
				"GET": s.ServeSwaggerSpec,
			},
		},
	}

	return routes
}

func (s *SwaggerController) Root(c *fireball.Context) (fireball.Response, error) {
	html := "<a href='/api/?url=/swagger.json'>Swagger</a>"
	return fireball.NewResponse(200, []byte(html), nil), nil
}

func (s *SwaggerController) ServeSwaggerSpec(c *fireball.Context) (fireball.Response, error) {
	spec := swagger.Spec{
		SwaggerVersion: "2.0",
		Host:           c.Request.Host,
		Schemes:        []string{"http", "https"},
		Info: &swagger.Info{
			Title:   "Layer0",
			Version: s.version,
		},
		Definitions: map[string]swagger.Definition{
			"Admin":                     models.APIConfig{}.Definition(),
			"Container":                 models.Container{}.Definition(),
			"ContainerOverride":         models.ContainerOverride{}.Definition(),
			"CreateEntityResponse":      models.CreateEntityResponse{}.Definition(),
			"CreateEnvironmentRequest":  models.CreateEnvironmentRequest{}.Definition(),
			"CreateLoadBalancerRequest": models.CreateLoadBalancerRequest{}.Definition(),
			"CreateServiceRequest":      models.CreateServiceRequest{}.Definition(),
			"CreateTaskRequest":         models.CreateTaskRequest{}.Definition(),
			"CreateDeployRequest":       models.CreateDeployRequest{}.Definition(),
			"Deployment":                models.Deployment{}.Definition(),
			"Environment":               models.Environment{}.Definition(),
			"HealthCheck":               models.HealthCheck{}.Definition(),
			"LoadBalancer":              models.LoadBalancer{}.Definition(),
			"LogFile":                   models.LogFile{}.Definition(),
			"Port":                      models.Port{}.Definition(),
			"Service":                   models.Service{}.Definition(),
			"Tag":                       models.Tag{}.Definition(),
			"Task":                      models.Task{}.Definition(),
			"Deploy":                    models.Deploy{}.Definition(),
			"UpdateLoadBalancerRequest": models.UpdateLoadBalancerRequest{}.Definition(),
			"UpdateServiceRequest":      models.UpdateServiceRequest{}.Definition(),
			"UpdateEnvironmentRequest":  models.UpdateEnvironmentRequest{}.Definition(),
		},
		Tags: []swagger.Tag{
			{
				Name:        "Admin",
				Description: "Methods related to admin",
			},
			{
				Name:        "Environment",
				Description: "Methods related to environments",
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
			{
				Name:        "Tag",
				Description: "Methods related to tags",
			},
			{
				Name:        "Deploy",
				Description: "Methods related to deploys",
			},
		},
		Paths: map[string]swagger.Path{
			"/admin/config": map[string]swagger.Method{
				"get": {
					Summary: "Get Config",
					Tags:    []string{"Admin"},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "Config of API",
							Schema:      swagger.NewObjectSliceSchema("Admin"),
						},
					},
				},
			},
			"/health": map[string]swagger.Method{
				"get": {
					Summary: "Get Health",
					Tags:    []string{"Admin"},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "Health of API",
						},
					},
				},
			},
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
					Summary: "Create a new Environment",
					Tags:    []string{"Environment"},
					Parameters: []swagger.Parameter{
						swagger.NewBodyParam("CreateEnvironmentRequest", "Environment to create", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "The ID of the new Environment",
							Schema:      swagger.NewObjectSchema("CreateEntityResponse"),
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
				"patch": {
					Summary: "Update Environment links",
					Tags:    []string{"Environment"},
					Parameters: []swagger.Parameter{
						swagger.NewStringPathParam("id", "ID of the environment to update", true),
						swagger.NewBodyParam("UpdateEnvironmentRequest", "Environment update request", true),
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
					Summary: "Create a new LoadBalancer",
					Tags:    []string{"LoadBalancer"},
					Parameters: []swagger.Parameter{
						swagger.NewBodyParam("CreateLoadBalancerRequest", "LoadBalancer to create", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "The ID of the new LoadBalancer",
							Schema:      swagger.NewObjectSchema("CreateEntityResponse"),
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
				"patch": {
					Summary: "Update a LoadBalancer",
					Tags:    []string{"LoadBalancer"},
					Parameters: []swagger.Parameter{
						swagger.NewBodyParam("UpdateLoadBalancerRequest", "LoadBalancer to update", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "Success",
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
				"post": {
					Summary: "Create a new Service",
					Tags:    []string{"Service"},
					Parameters: []swagger.Parameter{
						swagger.NewBodyParam("CreateServiceRequest", "Service to create", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "The ID of the new Service",
							Schema:      swagger.NewObjectSchema("CreateEntityResponse"),
						},
					},
				},
				"get": {
					Summary: "List all Services",
					Tags:    []string{"Service"},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "An array of services",
							Schema:      swagger.NewObjectSliceSchema("Service"),
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
				"patch": {
					Summary: "Update a Service",
					Tags:    []string{"Service"},
					Parameters: []swagger.Parameter{
						swagger.NewBodyParam("UpdateServiceRequest", "Service to update", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "Success",
						},
					},
				},
				"delete": {
					Summary: "Delete a Service",
					Tags:    []string{"Service"},
					Parameters: []swagger.Parameter{
						swagger.NewStringPathParam("id", "ID of the service to delete", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "Success",
						},
					},
				},
			},
			"/service/{id}/logs": map[string]swagger.Method{
				"get": {
					Summary: "Get service logs",
					Tags:    []string{"Service"},
					Parameters: []swagger.Parameter{
						swagger.NewStringPathParam("id", "ID of the service to describe", true),
						swagger.NewIntQueryParam("tail", "The number of lines from the end to return", false),
						swagger.NewStringQueryParam("start", "The start of the time range to fetch logs (format YYYY-MM-DD HH:MM)", false),
						swagger.NewStringQueryParam("end", "The end of the time range to fetch logs (format YYYY-MM-DD HH:MM)", false),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "The service's logs",
							Schema:      swagger.NewObjectSliceSchema("LogFile"),
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
					Summary: "Create a new Task",
					Tags:    []string{"Task"},
					Parameters: []swagger.Parameter{
						swagger.NewBodyParam("CreateTaskRequest", "Task to create", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "The ID of the new Task",
							Schema:      swagger.NewObjectSchema("CreateEntityResponse"),
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
			"/tag": map[string]swagger.Method{
				"delete": {
					Summary: "Delete a Tag",
					Tags:    []string{"Tag"},
					Parameters: []swagger.Parameter{
						swagger.NewBodyParam("Tag", "Tag to delete", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "Success",
						},
					},
				},
				"get": {
					Summary: "List Tags",
					Tags:    []string{"Tag"},
					Parameters: []swagger.Parameter{
						swagger.NewStringQueryParam(models.TagQueryParamEnvironmentID, "Filter entities that have a matching 'environment_id' tag", false),
						swagger.NewStringQueryParam(models.TagQueryParamFuzz, "Filter entities that have a matching entity id or name tag (glob patterns allowed)", false),
						swagger.NewStringQueryParam(models.TagQueryParamID, "Filter entities that have a matching entity id", false),
						swagger.NewStringQueryParam(models.TagQueryParamName, "Filter entities that have a matching name tag", false),
						swagger.NewStringQueryParam(models.TagQueryParamType, "Filter entities that have a matching type", false),
						swagger.NewStringQueryParam(models.TagQueryParamVersion, "Filter entities that have a version tag (version='latest' will return only the latest version)", false),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "An array of tags",
							Schema:      swagger.NewObjectSliceSchema("Tag"),
						},
					},
				},
				"post": {
					Summary: "Create a new Tag",
					Tags:    []string{"Tag"},
					Parameters: []swagger.Parameter{
						swagger.NewBodyParam("Tag", "Tag to create", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "The new Tag",
							Schema:      swagger.NewObjectSchema("Tag"),
						},
					},
				},
			},
			"/deploy": map[string]swagger.Method{
				"get": {
					Summary: "List all Deploys",
					Tags:    []string{"Deploy"},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "An array of deploys",
							Schema:      swagger.NewObjectSliceSchema("Deploy"),
						},
					},
				},
				"post": {
					Summary: "Create a new Deploy",
					Tags:    []string{"Deploy"},
					Parameters: []swagger.Parameter{
						swagger.NewBodyParam("CreateDeployRequest", "Deploy to create (base64 encoded)", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "The ID of the new Deploy",
							Schema:      swagger.NewObjectSchema("CreateEntityResponse"),
						},
					},
				},
			},
			"/deploy/{id}": map[string]swagger.Method{
				"get": {
					Summary: "Describe a Deploy",
					Tags:    []string{"Deploy"},
					Parameters: []swagger.Parameter{
						swagger.NewStringPathParam("id", "ID of the deploy to describe", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "The desired deploy",
							Schema:      swagger.NewObjectSchema("Deploy"),
						},
					},
				},
				"delete": {
					Summary: "Delete a Deploy",
					Tags:    []string{"Deploy"},
					Parameters: []swagger.Parameter{
						swagger.NewStringPathParam("id", "ID of the deploy to delete", true),
					},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "Success",
						},
					},
				},
			},
		},
	}

	return fireball.NewJSONResponse(200, spec)
}
