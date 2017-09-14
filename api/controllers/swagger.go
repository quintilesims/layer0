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
			"CreateEnvironmentRequest":  models.CreateEnvironmentRequest{}.Definition(),
			"CreateLoadBalancerRequest": models.CreateLoadBalancerRequest{}.Definition(),
			"Environment":               models.Environment{}.Definition(),
			"HealthCheck":               models.HealthCheck{}.Definition(),
			"LoadBalancer":              models.LoadBalancer{}.Definition(),
			"Port":                      models.Port{}.Definition(),
			"UpdateLoadBalancerRequest": models.UpdateLoadBalancerRequest{}.Definition(),
		},
		Tags: []swagger.Tag{
			{
				Name:        "Environment",
				Description: "Methods related to environments",
			},
			{
				Name:        "LoadBalancer",
				Description: "Methods related to load balancers",
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
							Schema:      swagger.NewObjectSchema("LoadBalancer"),
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
		},
	}

	return fireball.NewJSONResponse(200, spec)
}
