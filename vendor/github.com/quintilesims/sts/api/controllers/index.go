package controllers

import (
	"github.com/zpatrick/fireball"
)

type IndexController struct{}

func NewIndexController() *IndexController {
	return &IndexController{}
}

func (i *IndexController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/",
			Handlers: fireball.Handlers{
				"GET": i.getIndex,
			},
		},
	}
}

func (i *IndexController) getIndex(c *fireball.Context) (fireball.Response, error) {
	return fireball.NewResponse(200, []byte("Hello, World!"), nil), nil
}
