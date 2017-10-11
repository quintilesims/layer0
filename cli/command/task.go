package command

import (
	"github.com/quintilesims/layer0/client"
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

	return t.waitOnJobHelper(c, jobID, func(taskID string) error {
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

	taskID, err := t.resolveSingleEntityIDHelper("task", args["NAME"])
	if err != nil {
		return err
	}

	task, err := t.client.ReadTask(taskID)
	if err != nil {
		return err
	}

	return t.printer.PrintTasks(task)
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

	query := buildLogQueryHelper(id, c.String(client.LogQueryParamStart), c.String(client.LogQueryParamEnd), c.Int(client.LogQueryParamTail))
	logs, err := t.client.ReadTaskLogs(id, query)
	if err != nil {
		return err
	}

	return t.printer.PrintLogs(logs...)
}
