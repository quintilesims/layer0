package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/quintilesims/sts/models"
	"github.com/zpatrick/fireball"
	"log"
	"os/exec"
)

type CommandController struct {
	commands []*models.Command
}

func NewCommandController() *CommandController {
	return &CommandController{
		commands: []*models.Command{},
	}
}

func (cc *CommandController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/command",
			Handlers: fireball.Handlers{
				"GET":  cc.getCommands,
				"POST": cc.addCommand,
			},
		},
		{
			Path: "/command/:name",
			Handlers: fireball.Handlers{
				"GET": cc.getCommand,
			},
		},
	}
}

func (cc *CommandController) getCommands(c *fireball.Context) (fireball.Response, error) {
	return fireball.NewJSONResponse(200, cc.commands)
}

func (cc *CommandController) getCommand(c *fireball.Context) (fireball.Response, error) {
	name := c.PathVariables["name"]

	for _, cmd := range cc.commands {
		if cmd.Name == name {
			return fireball.NewJSONResponse(200, cmd)
		}
	}

	return nil, fmt.Errorf("Command with name '%s' not found", name)
}

func (cc *CommandController) addCommand(c *fireball.Context) (fireball.Response, error) {
	var req models.CreateCommandRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, err
	}

	if len(req.Args) == 0 {
		return nil, fmt.Errorf("must have at least one argument")
	}

	command := &models.Command{
		Name: req.Name,
		Args: req.Args,
	}

	cc.commands = append(cc.commands, command)
	cc.runCommand(command)
	return fireball.NewJSONResponse(200, command)
}

func (cc *CommandController) runCommand(cmd *models.Command) {
	c := exec.Command(cmd.Args[0], cmd.Args[1:]...)

	out, err := c.CombinedOutput()
	if err != nil {
		log.Println(err)
	}

	cmd.Output = string(out)
	log.Printf("Output from command '%s': %s\n", cmd.Name, cmd.Output)
}
