package command

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
	"strings"
)

type TaskCommand struct {
	*Command
}

func NewTaskCommand(command *Command) *TaskCommand {
	return &TaskCommand{command}
}

func (t *TaskCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:  "task",
		Usage: "manage layer0 tasks",
		Subcommands: []cli.Command{
			{
				Name:      "create",
				Usage:     "create a new task",
				Action:    wrapAction(t.Command, t.Create),
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
					cli.BoolFlag{
						Name:  "wait",
						Usage: "wait for the job to complete before returning",
					},
				},
			},
			{
				Name:      "delete",
				Usage:     "delete a task",
				ArgsUsage: "NAME",
				Action:    wrapAction(t.Command, t.Delete),
			},
			{
				Name:      "get",
				Usage:     "describe a task",
				Action:    wrapAction(t.Command, t.Get),
				ArgsUsage: "NAME",
			},
			{
				Name:      "list",
				Usage:     "list all tasks",
				Action:    wrapAction(t.Command, t.List),
				ArgsUsage: " ",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "all",
						Usage: "included deleted tasks",
					},
				},
			},
			{
				Name:      "logs",
				Usage:     "get the logs for a task",
				Action:    wrapAction(t.Command, t.Logs),
				ArgsUsage: "NAME",
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "tail",
						Usage: "number of lines from the end to return",
					},
				},
			},
		},
	}
}

func (t *TaskCommand) Create(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "ENVIRONMENT", "NAME", "DEPLOY")
	if err != nil {
		return err
	}

	overrides, err := parseOverrides(c.StringSlice("env"))
	if err != nil {
		return err
	}

	environmentID, err := t.resolveSingleID("environment", args["ENVIRONMENT"])
	if err != nil {
		return err
	}

	deployID, err := t.resolveSingleID("deploy", args["DEPLOY"])
	if err != nil {
		return err
	}

	jobID, err := t.Client.CreateTask(args["NAME"], environmentID, deployID, c.Int("copies"), overrides)
	if err != nil {
		return err
	}

	if !c.Bool("wait") {
		t.Printer.Printf("This operation is running as a job. Run `l0 job get %s` to see progress\n", jobID)
		return nil
	}

	timeout, err := getTimeout(c)
	if err != nil {
		return err
	}

	t.Printer.StartSpinner("Creating")
	if err := t.Client.WaitForJob(jobID, timeout); err != nil {
		return err
	}

	job, err := t.Client.GetJob(jobID)
	if err != nil {
		return err
	}

	taskIDs := []string{}
	for key, val := range job.Meta {	
		if strings.HasPrefix(key, "task_"){
			taskIDs = append(taskIDs, val)
		}
	}	

	tasks := make([]*models.Task, len(taskIDs))
	for i, taskID := range taskIDs {
		task, err := t.Client.GetTask(taskID)
		if err != nil {
			return err
		}

		tasks[i] = task
	}

	return t.Printer.PrintTasks(tasks...)
}

func (t *TaskCommand) Delete(c *cli.Context) error {
	return t.delete(c, "task", t.Client.DeleteTask)
}

func (t *TaskCommand) Get(c *cli.Context) error {
	tasks := []*models.Task{}
	getTaskf := func(id string) error {
		task, err := t.Client.GetTask(id)
		if err != nil {
			return err
		}

		tasks = append(tasks, task)
		return nil
	}

	if err := t.get(c, "task", getTaskf); err != nil {
		return err
	}

	return t.Printer.PrintTasks(tasks...)
}

func (t *TaskCommand) List(c *cli.Context) error {
	taskSummaries, err := t.Client.ListTasks()
	if err != nil {
		return err
	}

	if !c.Bool("all") {
		taskSummaries = filterTaskSummaries(taskSummaries)
	}

	return t.Printer.PrintTaskSummaries(taskSummaries...)
}

func (t *TaskCommand) Logs(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "NAME")
	if err != nil {
		return err
	}

	id, err := t.resolveSingleID("task", args["NAME"])
	if err != nil {
		return err
	}

	logs, err := t.Client.GetTaskLogs(id, c.Int("tail"))
	if err != nil {
		return err
	}

	return t.Printer.PrintLogs(logs...)
}

func filterTaskSummaries(tasks []*models.TaskSummary) []*models.TaskSummary {
	filtered := []*models.TaskSummary{}

	for _, task := range tasks {
		if task.TaskName != "" {
			filtered = append(filtered, task)
		}
	}

	return filtered
}

func parseOverrides(overrides []string) ([]models.ContainerOverride, error) {
	catalog := map[string]models.ContainerOverride{}

	for _, o := range overrides {
		split := strings.FieldsFunc(o, func(r rune) bool {
			return r == ':' || r == '='
		})

		if len(split) != 3 {
			return nil, NewUsageError("Environment Variable Override format is: CONTAINER:VAR=VAL")
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
