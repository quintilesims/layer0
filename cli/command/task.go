package command

import (
	"fmt"
	"strings"

	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

type TaskCommand struct {
	*CommandBase
}

func NewTaskCommand(b *CommandBase) *TaskCommand {
	return &TaskCommand{b}
}

func (t *TaskCommand) Command() cli.Command {
	return cli.Command{
		Name:  "task",
		Usage: "manage layer0 tasks",
		Subcommands: []cli.Command{
			{
				Name:      "create",
				Usage:     "create a new task",
				Action:    t.create,
				ArgsUsage: "ENVIRONMENT NAME DEPLOY",
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "copies",
						Value: 1,
						Usage: "number of copies of deploy to run (default: 1)",
					},
					cli.StringSliceFlag{
						Name:  "env",
						Usage: "environment variable override in format 'CONTAINER:VAR=VAL' (can be specified multiple times)",
					},
				},
			},
			{
				Name:      "delete",
				Usage:     "delete a Task",
				ArgsUsage: "NAME",
				Action:    t.delete,
			},
			{
				Name:      "list",
				Usage:     "list all Tasks",
				Action:    t.list,
				ArgsUsage: " ",
			},
			{
				Name:      "read",
				Usage:     "describe a Task",
				Action:    t.read,
				ArgsUsage: "NAME",
			},
			{
				Name:      "logs",
				Usage:     "get the logs for a task",
				Action:    t.logs,
				ArgsUsage: "NAME",
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "tail",
						Usage: "number of lines from the end to return",
					},
					cli.StringFlag{
						Name:  "start",
						Usage: "the start of the time range to fetch logs (format: YYYY-MM-DD HH:MM)",
					},
					cli.StringFlag{
						Name:  "end",
						Usage: "the end of the time range to fetch logs (format: YYYY-MM-DD HH:MM)",
					},
				},
			},
		},
	}
}

func (t *TaskCommand) create(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "ENVIRONMENT", "DEPLOY", "NAME")
	if err != nil {
		return err
	}

	// TODO: Add Parse override.
	// overrides, err := parseOverrides(c.StringSlice("env"))
	// if err != nil {
	// 	return err
	// }

	environmentID, err := t.resolveSingleEntityIDHelper("environment", args["ENVIRONMENT"])
	if err != nil {
		return err
	}

	deployID, err := t.resolveSingleEntityIDHelper("deploy", args["DEPLOY"])
	if err != nil {
		return err
	}

	req := models.CreateTaskRequest{
		TaskName:      args["NAME"],
		EnvironmentID: environmentID,
		DeployID:      deployID,
	}

	jobID, err := t.client.CreateTask(req)
	if err != nil {
		return err
	}

	return t.waitOnJobHelper(c, jobID, "creating", func(taskID string) error {
		task, err := t.client.ReadTask(taskID)
		if err != nil {
			return err
		}

		return t.printer.PrintTasks(task)
	})
}

func (t *TaskCommand) delete(c *cli.Context) error {
	return t.deleteHelper(c, "task", func(taskID string) (string, error) {
		return t.client.DeleteTask(taskID)
	})
}

func (t *TaskCommand) list(c *cli.Context) error {
	taskSummaries, err := t.client.ListTasks()
	if err != nil {
		return err
	}

	return t.printer.PrintTaskSummaries(taskSummaries...)
}

func (t *TaskCommand) read(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "NAME")
	if err != nil {
		return err
	}

	taskIDs, err := t.resolver.Resolve("task", args["NAME"])
	if err != nil {
		return err
	}

	// MARK: Working Still
	ts := make([]*models.Task, len(taskIDs))
	tasks, err := t.readEntitiesHelper(taskIDs, func(taskID string) (interface{}, error) {
		task, err := t.client.ReadTask(taskID)
		if err != nil {
			return nil, err
		}
		return task, nil
	})

	for i, t := range tasks {
		if t, ok := t.(*models.Task); ok {
			ts[i] = t
		}
	}

	return t.printer.PrintTasks(ts...)
}

func (t *TaskCommand) logs(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "NAME")
	if err != nil {
		return err
	}

	id, err := t.resolveSingleEntityIDHelper("task", args["NAME"])
	if err != nil {
		return err
	}

	query := buildLogQueryHelper(id, c.String("start"), c.String("end"), c.Int("tail"))
	logs, err := t.client.ReadTaskLogs(id, query)
	if err != nil {
		return err
	}

	return t.printer.PrintLogs(logs...)
}

// MARK: Move to Models
type ContainerOverride struct {
	ContainerName        string            `json:"container_name"`
	EnvironmentOverrides map[string]string `json:"environment_overrides"`
}

func parseOverrides(overrides []string) ([]ContainerOverride, error) {
	catalog := map[string]ContainerOverride{}

	for _, o := range overrides {
		split := strings.FieldsFunc(o, func(r rune) bool {
			return r == ':' || r == '='
		})

		if len(split) != 3 {
			return nil, fmt.Errorf("Environment Variable Override format is: CONTAINER:VAR=VAL")
		}

		container := split[0]
		key := split[1]
		val := split[2]

		if _, ok := catalog[container]; !ok {
			catalog[container] = ContainerOverride{
				ContainerName:        container,
				EnvironmentOverrides: map[string]string{},
			}
		}

		catalog[container].EnvironmentOverrides[key] = val
	}

	models := []ContainerOverride{}
	for _, override := range catalog {
		models = append(models, override)
	}

	return models, nil

}
