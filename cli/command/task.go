package command

import (
	"fmt"
	"log"
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
				ArgsUsage: "ENVIRONMENT TASK_NAME DEPLOY",
				Flags: []cli.Flag{
					cli.StringSliceFlag{
						Name:  "env",
						Usage: "environment variable override in format 'CONTAINER:VAR=VAL' (can be specified multiple times)",
					},
				},
			},
			{
				Name:      "delete",
				Usage:     "delete a Task",
				ArgsUsage: "TASK_NAME",
				Action:    t.delete,
			},
			{
				Name:      "list",
				Usage:     "list all Tasks",
				Action:    t.list,
				ArgsUsage: " ",
			},
			{
				Name:      "get",
				Usage:     "describe a Task",
				Action:    t.read,
				ArgsUsage: "TASK_NAME",
			},
			{
				Name:      "logs",
				Usage:     "get the logs for a task",
				Action:    t.logs,
				ArgsUsage: "TASK_NAME",
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
	args, err := extractArgs(c.Args(), "ENVIRONMENT", "TASK_NAME", "DEPLOY")
	if err != nil {
		return err
	}

	overrides, err := parseOverrides(c.StringSlice("env"))
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
		TaskName:           args["TASK_NAME"],
		EnvironmentID:      environmentID,
		DeployID:           deployID,
		ContainerOverrides: overrides,
	}

	if err := req.Validate(); err != nil {
		return err
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
	args, err := extractArgs(c.Args(), "TASK_NAME")
	if err != nil {
		return err
	}

	taskIDs, err := t.resolver.Resolve("task", args["TASK_NAME"])
	if err != nil {
		return err
	}

	taskSummaries, err := t.client.ListTasks()
	if err != nil {
		return err
	}

	taskExists := map[string]bool{}
	for _, taskSummary := range taskSummaries {
		taskExists[taskSummary.TaskID] = true
	}

	tasks := make([]*models.Task, 0, len(taskIDs))
	for _, taskID := range taskIDs {
		if !taskExists[taskID] {
			log.Printf("[DEBUG] Resolver returned an expired task '%s'", taskID)
			continue
		}

		task, err := t.client.ReadTask(taskID)
		if err != nil {
			return err
		}

		tasks = append(tasks, task)
	}

	return t.printer.PrintTasks(tasks...)
}

func (t *TaskCommand) logs(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "TASK_NAME")
	if err != nil {
		return err
	}

	taskID, err := t.resolveSingleEntityIDHelper("task", args["TASK_NAME"])
	if err != nil {
		return err
	}

	query := buildLogQueryHelper(c.String("start"), c.String("end"), c.Int("tail"))
	logs, err := t.client.ReadTaskLogs(taskID, query)
	if err != nil {
		return err
	}

	return t.printer.PrintLogs(logs...)
}

func parseOverrides(overrides []string) ([]models.ContainerOverride, error) {
	catalog := map[string]models.ContainerOverride{}

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
			catalog[container] = models.ContainerOverride{
				ContainerName:        container,
				EnvironmentOverrides: map[string]string{},
			}
		}

		catalog[container].EnvironmentOverrides[key] = val
	}

	models := []models.ContainerOverride{}
	for _, override := range catalog {
		models = append(models, override)
	}

	return models, nil

}
