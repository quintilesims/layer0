package main

import (
	"fmt"
	"github.com/quintilesims/sts/api/controllers"
	"github.com/urfave/cli"
	"github.com/zpatrick/fireball"
	"log"
	"net/http"
	"os"
)

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
	routes := controllers.NewIndexController().Routes()
	routes = append(routes, controllers.NewHealthController().Routes()...)
	routes = append(routes, controllers.NewCommandController().Routes()...)
	routes = fireball.Decorate(routes, fireball.LogDecorator())
	app := fireball.NewApp(routes)
	app.ErrorHandler = controllers.JSONErrorHandler

	address := fmt.Sprintf("0.0.0.0:%d", c.Int("port"))
	log.Printf("Listening on %s\n", address)
	return http.ListenAndServe(address, app)
}
