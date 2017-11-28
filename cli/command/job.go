package command

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

type JobCommand struct {
	*CommandBase
}

func NewJobCommand(b *CommandBase) *JobCommand {
	return &JobCommand{b}
}

func (j *JobCommand) Command() cli.Command {
	return cli.Command{
		Name:  "job",
		Usage: "manage layer0 jobs",
		Subcommands: []cli.Command{
			{
				Name:      "delete",
				Usage:     "delete a job",
				Action:    j.delete,
				ArgsUsage: "NAME",
			},
			{
				Name:      "read",
				Usage:     "describe a job",
				Action:    j.read,
				ArgsUsage: "NAME",
			},
			{
				Name:      "list",
				Usage:     "list all jobs",
				Action:    j.list,
				ArgsUsage: " ",
			},
		},
	}
}

func (j *JobCommand) delete(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "NAME")
	if err != nil {
		return err
	}

	jobID, err := j.resolveSingleEntityIDHelper("job", args["NAME"])
	if err != nil {
		return err
	}

	return j.client.DeleteJob(jobID)
}

func (j *JobCommand) read(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "NAME")
	if err != nil {
		return err
	}

	jobIDs, err := j.resolver.Resolve("job", args["NAME"])
	if err != nil {
		return err
	}

	jobs := make([]*models.Job, len(jobIDs))
	for i, jobID := range jobIDs {
		job, err := j.client.ReadJob(jobID)
		if err != nil {
			return err
		}

		jobs[i] = job
	}

	return j.printer.PrintJobs(jobs...)
}

func (j *JobCommand) list(c *cli.Context) error {
	jobs, err := j.client.ListJobs()
	if err != nil {
		return err
	}

	return j.printer.PrintJobs(jobs...)
}
