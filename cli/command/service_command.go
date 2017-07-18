package command

import (
	"strconv"

	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

type ServiceCommand struct {
	*Command
}

func NewServiceCommand(command *Command) *ServiceCommand {
	return &ServiceCommand{command}
}

func (s *ServiceCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:  "service",
		Usage: "manage layer0 services",
		Subcommands: []cli.Command{
			{
				Name:      "create",
				Usage:     "create a new service",
				Action:    wrapAction(s.Command, s.Create),
				ArgsUsage: "ENVIRONMENT NAME DEPLOY",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "loadbalancer",
						Usage: "attach the service to the specified load balancer",
					},
					cli.BoolFlag{
						Name:  "wait",
						Usage: "wait until deployment completes before returning",
					},
				},
			},
			{
				Name:      "delete",
				Usage:     "delete a Service",
				ArgsUsage: "NAME",
				Action:    wrapAction(s.Command, s.Delete),
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "wait",
						Usage: "wait for the job to complete before returning",
					},
				},
			},
			{
				Name:      "update",
				Usage:     "run a new deploy on a service",
				Action:    wrapAction(s.Command, s.Update),
				ArgsUsage: "NAME DEPLOY",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "wait",
						Usage: "wait until the deployment completes before returning",
					},
				},
			},
			{
				Name:      "get",
				Usage:     "describe a service",
				Action:    wrapAction(s.Command, s.Get),
				ArgsUsage: "NAME",
			},
			{
				Name:      "list",
				Usage:     "list all services",
				Action:    wrapAction(s.Command, s.List),
				ArgsUsage: " ",
			},
			{
				Name:      "logs",
				Usage:     "get the logs for a service",
				Action:    wrapAction(s.Command, s.Logs),
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
			{
				Name:      "scale",
				Usage:     "scale a service",
				Action:    wrapAction(s.Command, s.Scale),
				ArgsUsage: "NAME COUNT",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "wait",
						Usage: "wait until the deployment completes before returning",
					},
				},
			},
		},
	}
}

func (s *ServiceCommand) Create(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "ENVIRONMENT", "NAME", "DEPLOY")
	if err != nil {
		return err
	}

	environmentID, err := s.resolveSingleID("environment", args["ENVIRONMENT"])
	if err != nil {
		return err
	}

	deployID, err := s.resolveSingleID("deploy", args["DEPLOY"])
	if err != nil {
		return err
	}

	var loadBalancerID string
	if loadBalancerName := c.String("loadbalancer"); loadBalancerName != "" {
		id, err := s.resolveSingleID("load_balancer", loadBalancerName)
		if err != nil {
			return err
		}

		loadBalancerID = id
	}

	service, err := s.Client.CreateService(args["NAME"], environmentID, deployID, loadBalancerID)
	if err != nil {
		return err
	}

	if !c.Bool("wait") {
		return s.Printer.PrintServices(service)
	}

	timeout, err := getTimeout(c)
	if err != nil {
		return err
	}

	s.Printer.StartSpinner("Waiting for Deployment")
	service, err = s.Client.WaitForDeployment(service.ServiceID, timeout)
	if err != nil {
		return err
	}

	return s.Printer.PrintServices(service)
}

func (s *ServiceCommand) Delete(c *cli.Context) error {
	return s.deleteWithJob(c, "service", s.Client.DeleteService)
}

func (s *ServiceCommand) Update(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "NAME", "DEPLOY")
	if err != nil {
		return err
	}

	serviceID, err := s.resolveSingleID("service", args["NAME"])
	if err != nil {
		return err
	}

	deployID, err := s.resolveSingleID("deploy", args["DEPLOY"])
	if err != nil {
		return err
	}

	service, err := s.Client.UpdateService(serviceID, deployID)
	if err != nil {
		return err
	}

	if !c.Bool("wait") {
		return s.Printer.PrintServices(service)
	}

	timeout, err := getTimeout(c)
	if err != nil {
		return err
	}

	s.Printer.StartSpinner("Waiting for Deployment")
	service, err = s.Client.WaitForDeployment(serviceID, timeout)
	if err != nil {
		return err
	}

	return s.Printer.PrintServices(service)
}

func (s *ServiceCommand) Get(c *cli.Context) error {
	services := []*models.Service{}
	getServicef := func(id string) error {
		service, err := s.Client.GetService(id)
		if err != nil {
			return err
		}

		services = append(services, service)
		return nil
	}

	if err := s.get(c, "service", getServicef); err != nil {
		return err
	}

	return s.Printer.PrintServices(services...)
}

func (s *ServiceCommand) List(c *cli.Context) error {
	serviceSummaries, err := s.Client.ListServices()
	if err != nil {
		return err
	}

	return s.Printer.PrintServiceSummaries(serviceSummaries...)
}

func (s *ServiceCommand) Logs(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "NAME")
	if err != nil {
		return err
	}

	id, err := s.resolveSingleID("service", args["NAME"])
	if err != nil {
		return err
	}

	logs, err := s.Client.GetServiceLogs(id, c.String("start"), c.String("end"), c.Int("tail"))
	if err != nil {
		return err
	}

	return s.Printer.PrintLogs(logs...)
}

func (s *ServiceCommand) Scale(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "NAME", "COUNT")
	if err != nil {
		return err
	}

	count, err := strconv.ParseInt(args["COUNT"], 10, 64)
	if err != nil {
		return NewUsageError("'%s' is not a valid integer", args["COUNT"])
	}

	id, err := s.resolveSingleID("service", args["NAME"])
	if err != nil {
		return err
	}

	service, err := s.Client.ScaleService(id, int(count))
	if err != nil {
		return err
	}

	if !c.Bool("wait") {
		return s.Printer.PrintServices(service)
	}

	timeout, err := getTimeout(c)
	if err != nil {
		return err
	}

	s.Printer.StartSpinner("Waiting for Deployment")
	service, err = s.Client.WaitForDeployment(id, timeout)
	if err != nil {
		return err
	}

	return s.Printer.PrintServices(service)
}
