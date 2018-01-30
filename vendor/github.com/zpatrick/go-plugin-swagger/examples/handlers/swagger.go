package handlers

import (
	"github.com/zpatrick/go-plugin-swagger"
	"github.com/zpatrick/go-plugin-swagger/example/movie"
)

func SwaggerSpec() swagger.Spec {
	return swagger.Spec{
		SwaggerVersion: "2.0",
		Schemes:        []string{"http"},
		Info: &swagger.Info{
			Title:   "Swagger Example",
			Version: "0.0.1",
		},
		Definitions: map[string]swagger.Definition{
			"Movie": movie.Movie{}.Definition(),
		},
		Tags: []swagger.Tag{
			{
				Name:        "Movies",
				Description: "Methods related to movies",
			},
		},
		Paths: map[string]swagger.Path{
			"/movies": map[string]swagger.Method{
				"get": {
					Summary: "List all Movies",
					Tags:    []string{"Movies"},
					Responses: map[string]swagger.Response{
						"200": {
							Description: "An array of movies",
							Schema:      swagger.NewObjectSliceSchema("Movie"),
						},
					},
				},
				"post": {
					Summary: "Add a Movie",
					Tags:    []string{"Movies"},
					Parameters: []swagger.Parameter{
						swagger.NewBodyParam("Movie", "Movie to add", true),
					},
					Responses: map[string]swagger.Response{
						"204": {
							Description: "The added movie",
							Schema:      swagger.NewObjectSchema("Movie"),
						},
					},
				},
			},
		},
	}
}
