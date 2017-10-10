package command

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

type ServiceCommand struct {
	*CommandBase
}

func NewServiceCommand(b *CommandBase) *ServiceCommand {
	return &ServiceCommand{b}
}

func (s *ServiceCommand) Command() cli.Command {
	return cli.Command{
		Name:  "service",
		Usage: "manage Layer0 services",
		Subcommands: []cli.Command{
			{
				Name:      "create",
				Usage:     "create a new service",
				Action:    s.create,
				ArgsUsage: "ENVIRONMENT NAME DEPLOY",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "loadbalancer",
						Usage: "attach the service to the specified load balancer",
					},
					cli.BoolFlag{
						Name:  "nowait",
						Usage: "don't wait for deployment to complete before returning",
					},
				},
			},
			{
				Name:      "delete",
				Usage:     "delete a service",
				Action:    s.delete,
				ArgsUsage: "NAME",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "nowait",
						Usage: "don't wait for the job to complete before returning",
					},
				},
			},
			{
				Name:   "list",
				Usage:  "list all services",
				Action: s.list,
			},
			{
				Name:      "logs",
				Usage:     "get the logs for a service",
				Action:    s.logs,
				ArgsUsage: "NAME",
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "tail",
						Usage: "number of lines from the end to return",
						Value: 0,
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
			{
				Name:      "read",
				Usage:     "describe a service",
				Action:    s.read,
				ArgsUsage: "NAME",
			},
			{
				Name:      "scale",
				Usage:     "scale a service",
				Action:    s.scale,
				ArgsUsage: "NAME COUNT",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "nowait",
						Usage: "don't wait for the job to complete before returning",
					},
				},
			},
			{
				Name:      "update",
				Usage:     "run a new dploy on a service",
				Action:    s.update,
				ArgsUsage: "NAME DEPLOY",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "nowait",
						Usage: "don't wait for the job to complete before returning",
					},
				},
			},
		},
	}
}

func (s *ServiceCommand) create(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "ENVIRONMENT", "NAME", "DEPLOY")
	if err != nil {
		return err
	}

	deployID, err := s.resolveSingleEntityIDHelper("deploy", args["DEPLOY"])
	if err != nil {
		return err
	}

	environmentID, err := s.resolveSingleEntityIDHelper("environment", args["ENVIRONMENT"])
	if err != nil {
		return err
	}

	var loadBalancerID string
	if loadBalancerName := c.String("loadbalancer"); loadBalancerName != "" {
		id, err := s.resolveSingleEntityIDHelper("load_balancer", loadBalancerName)
		if err != nil {
			return err
		}

		loadBalancerID = id
	}

	req := models.CreateServiceRequest{
		DeployID:       deployID,
		EnvironmentID:  environmentID,
		LoadBalancerID: loadBalancerID,
		ServiceName:    args["NAME"],
	}

	jobID, err := s.client.CreateService(req)
	if err != nil {
		return err
	}

	onCompleteFn := func(serviceID string) error {
		service, err := s.client.ReadService(serviceID)
		if err != nil {
			return err
		}

		return s.printer.PrintServices(service)
	}

	return s.waitOnJobHelper(c, jobID, "creating", onCompleteFn)
}

func (s *ServiceCommand) delete(c *cli.Context) error {
	deleteFn := func(serviceID string) (string, error) {
		return s.client.DeleteService(serviceID)
	}

	return s.deleteHelper(c, "service", deleteFn)
}

func (s *ServiceCommand) list(c *cli.Context) error {
	serviceSummaries, err := s.client.ListServices()
	if err != nil {
		return err
	}

	return s.printer.PrintServiceSummaries(serviceSummaries...)
}

func (s *ServiceCommand) logs(c *cli.Context) error {
	return nil
}

func (s *ServiceCommand) read(c *cli.Context) error {
	return nil
}

func (s *ServiceCommand) scale(c *cli.Context) error {
	return nil
}

func (s *ServiceCommand) update(c *cli.Context) error {
	return nil
}
