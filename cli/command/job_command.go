package command

import (
	"github.com/quintilesims/layer0/cli/entity"
	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

type JobCommand struct {
	*Command
}

func NewJobCommand(command *Command) *JobCommand {
	return &JobCommand{command}
}

func (j *JobCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:     "job",
		Usage:    "manage layer0 jobs",
		HideHelp: true,
		Subcommands: []cli.Command{
			{
				Name:      "delete",
				Usage:     "delete a job",
				Action:    wrapAction(j.Command, j.Delete),
				ArgsUsage: "NAME",
			},
			{
				Name:      "get",
				Usage:     "describe a job",
				Action:    wrapAction(j.Command, j.Get),
				ArgsUsage: "NAME",
			},
			{
				Name:      "list",
				Usage:     "list all jobs",
				Action:    wrapAction(j.Command, j.List),
				ArgsUsage: " ",
			},
			{
				Name:      "logs",
				Usage:     "get the logs for a job",
				Action:    wrapAction(j.Command, j.Logs),
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

func (j *JobCommand) Delete(c *cli.Context) error {
	return j.delete(c, "job", j.Client.DeleteJob)
}

func (j *JobCommand) Get(c *cli.Context) error {
	return j.get(c, "job", func(id string) (entity.Entity, error) {
		job, err := j.Client.GetJob(id)
		if err != nil {
			return nil, err
		}

		return entity.NewJob(job), nil
	})
}

func (j *JobCommand) List(c *cli.Context) error {
	jobs, err := j.Client.ListJobs()
	if err != nil {
		return err
	}

	return j.printJobs(jobs)
}

func (j *JobCommand) Logs(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "NAME")
	if err != nil {
		return err
	}

	id, err := j.resolveSingleID("job", args["NAME"])
	if err != nil {
		return err
	}

	job, err := j.Client.GetJob(id)
	if err != nil {
		return err
	}

	logs, err := j.Client.GetTaskLogs(job.TaskID, c.Int("tail"))
	if err != nil {
		return err
	}

	return j.Printer.PrintLogs(logs)
}

func (j *JobCommand) printJob(job *models.Job) error {
	entity := entity.NewJob(job)
	return j.Printer.PrintEntity(entity)
}

func (j *JobCommand) printJobs(jobs []*models.Job) error {
	entities := []entity.Entity{}
	for _, job := range jobs {
		entities = append(entities, entity.NewJob(job))
	}

	return j.Printer.PrintEntities(entities)
}
