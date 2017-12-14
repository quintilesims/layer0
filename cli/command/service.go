package command

import (
	"fmt"
	"strconv"

	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/config"
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
				ArgsUsage: "ENVIRONMENT SERVICE_NAME DEPLOY",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "loadbalancer",
						Usage: "attach the service to the specified load balancer",
					},
					cli.IntFlag{
						Name:  "scale",
						Value: config.DefaultServiceScale,
						Usage: "The desired scale of the service",
					},
				},
			},
			{
				Name:      "delete",
				Usage:     "delete a service",
				Action:    s.delete,
				ArgsUsage: "SERVICE_NAME",
			},
			{
				Name:      "get",
				Usage:     "describe a service",
				Action:    s.read,
				ArgsUsage: "SERVICE_NAME",
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
				ArgsUsage: "SERVICE_NAME",
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "tail",
						Usage: "number of lines from the end to return (default: 0)",
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
				Name:      "scale",
				Usage:     "scale a service",
				Action:    s.scale,
				ArgsUsage: "SERVICE_NAME COUNT",
			},
			{
				Name:      "update",
				Usage:     "run a new dploy on a service",
				Action:    s.update,
				ArgsUsage: "SERVICE_NAME DEPLOY",
			},
		},
	}
}

func (s *ServiceCommand) create(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "ENVIRONMENT", "SERVICE_NAME", "DEPLOY")
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
		ServiceName:    args["SERVICE_NAME"],
		Scale:          c.Int("scale"),
	}

	if err := req.Validate(); err != nil {
		return err
	}

	jobID, err := s.client.CreateService(req)
	if err != nil {
		return err
	}

	return s.waitOnJobHelper(c, jobID, "creating", func(serviceID string) error {
		service, err := client.WaitForDeployment(s.client, serviceID, c.GlobalDuration(config.FlagTimeout.GetName()))
		if err != nil {
			return err
		}

		return s.printer.PrintServices(service)
	})
}

func (s *ServiceCommand) delete(c *cli.Context) error {
	return s.deleteHelper(c, "service", func(serviceID string) (string, error) {
		return s.client.DeleteService(serviceID)
	})
}

func (s *ServiceCommand) list(c *cli.Context) error {
	serviceSummaries, err := s.client.ListServices()
	if err != nil {
		return err
	}

	return s.printer.PrintServiceSummaries(serviceSummaries...)
}

func (s *ServiceCommand) logs(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "SERVICE_NAME")
	if err != nil {
		return err
	}

	serviceID, err := s.resolveSingleEntityIDHelper("service", args["SERVICE_NAME"])
	if err != nil {
		return err
	}

	query := buildLogQueryHelper(c.String("start"), c.String("end"), c.Int("tail"))

	logs, err := s.client.ReadServiceLogs(serviceID, query)
	if err != nil {
		return err
	}

	return s.printer.PrintLogs(logs...)
}

func (s *ServiceCommand) read(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "SERVICE_NAME")
	if err != nil {
		return err
	}

	serviceIDs, err := s.resolver.Resolve("service", args["SERVICE_NAME"])
	if err != nil {
		return err
	}

	services := make([]*models.Service, len(serviceIDs))
	for i, serviceID := range serviceIDs {
		service, err := s.client.ReadService(serviceID)
		if err != nil {
			return err
		}

		services[i] = service
	}

	return s.printer.PrintServices(services...)
}

func (s *ServiceCommand) scale(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "SERVICE_NAME", "COUNT")
	if err != nil {
		return err
	}

	serviceID, err := s.resolveSingleEntityIDHelper("service", args["SERVICE_NAME"])
	if err != nil {
		return err
	}

	scale, err := strconv.Atoi(args["COUNT"])
	if err != nil {
		return fmt.Errorf("Failed to parse COUNT argument: %v", args["COUNT"])
	}

	req := models.UpdateServiceRequest{
		Scale: &scale,
	}

	jobID, err := s.client.UpdateService(serviceID, req)
	if err != nil {
		return err
	}

	return s.waitOnJobHelper(c, jobID, "scaling", func(serviceID string) error {
		service, err := client.WaitForDeployment(s.client, serviceID, c.GlobalDuration(config.FlagTimeout.GetName()))
		if err != nil {
			return err
		}

		return s.printer.PrintServices(service)
	})
}

func (s *ServiceCommand) update(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "SERVICE_NAME", "DEPLOY")
	if err != nil {
		return err
	}

	serviceID, err := s.resolveSingleEntityIDHelper("service", args["SERVICE_NAME"])
	if err != nil {
		return err
	}

	deployID, err := s.resolveSingleEntityIDHelper("deploy", args["DEPLOY"])
	if err != nil {
		return err
	}

	req := models.UpdateServiceRequest{
		DeployID: &deployID,
	}

	jobID, err := s.client.UpdateService(serviceID, req)
	if err != nil {
		return err
	}

	return s.waitOnJobHelper(c, jobID, "updating", func(serviceID string) error {
		service, err := client.WaitForDeployment(s.client, serviceID, c.GlobalDuration(config.FlagTimeout.GetName()))
		if err != nil {
			return err
		}

		return s.printer.PrintServices(service)
	})
}
