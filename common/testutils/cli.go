package testutils

import (
	"io/ioutil"
	"strings"

	"github.com/urfave/cli"
)

func NewApp(command cli.Command) *cli.App {
	app := cli.NewApp()
	app.Commands = []cli.Command{command}
	app.Writer = ioutil.Discard
	app.ErrWriter = ioutil.Discard
	return app
}

func RunApp(command cli.Command, input string) error {
	app := NewApp(command)
	return app.Run(strings.Split(input, " "))
}
