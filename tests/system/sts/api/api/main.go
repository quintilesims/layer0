package main

import (
	"fmt"
	"github.com/quintilesims/layer0/tests/system/sts/controllers"
	"github.com/urfave/cli"
	"github.com/zpatrick/fireball"
	"log"
	"net/http"
	"os"
)

func index(c *fireball.Context) (fireball.Response, error) {
	return fireball.NewResponse(200, []byte("Hello, World!"), nil), nil
}

func main() {
	app := cli.NewApp()
	app.Action = run
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "p, port",
			Value: 80,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	routes := []*fireball.Route{
		{
			Path: "/",
			Handlers: fireball.Handlers{
				"GET": index,
			},
		},
	}

	routes = append(routes, controllers.NewHealthController().Routes()...)
	routes = append(routes, controllers.NewCommandController().Routes()...)
	routes = fireball.Decorate(routes, fireball.LogDecorator())
	app := fireball.NewApp(routes)

	address := fmt.Sprintf("0.0.0.0:%d", c.Int("port"))
	log.Printf("Listening on %s\n", address)
	return http.ListenAndServe(address, app)
}
