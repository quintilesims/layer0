package command

import (
	"fmt"
	"github.com/quintilesims/layer0/cli/entity"
	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
	"strconv"
	"strings"
)

type LoadBalancerCommand struct {
	*Command
}

func NewLoadBalancerCommand(command *Command) *LoadBalancerCommand {
	return &LoadBalancerCommand{command}
}

func (l *LoadBalancerCommand) GetCommand() cli.Command {
	return cli.Command{
		Name:  "loadbalancer",
		Usage: "manage layer0 load balancers",
		Subcommands: []cli.Command{
			{
				Name:      "addport",
				Usage:     "add a new listener port on a load balancer",
				Action:    wrapAction(l.Command, l.AddPort),
				ArgsUsage: "NAME HOST_PORT:CONTAINER_PORT/PROTOCOL",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "certificate",
						Usage: "name of certificate to use for port configuration (only required for https)",
					},
				},
			},
			{
				Name:      "create",
				Usage:     "create a new load balancer",
				Action:    wrapAction(l.Command, l.Create),
				ArgsUsage: "ENVIRONMENT NAME",
				Flags: []cli.Flag{
					cli.StringSliceFlag{
						Name:  "port",
						Usage: "port configuration in format 'HOST_PORT:CONTAINER_PORT/PROTOCOL' (default 80:80/tcp)",
					},
					cli.StringFlag{
						Name:  "certificate",
						Usage: "name of certificate to use for port configuration (only required for https)",
					},
					cli.BoolFlag{
						Name:  "private",
						Usage: "if specified, creates a private load balancer (default is public)",
					},
				},
			},
			{
				Name:      "delete",
				Usage:     "delete a load balancer",
				ArgsUsage: "NAME",
				Action:    wrapAction(l.Command, l.Delete),
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "wait",
						Usage: "wait for the job to complete before returning",
					},
				},
			},
			{
				Name:      "dropport",
				Usage:     "drop a listener port from a load balancer",
				Action:    wrapAction(l.Command, l.DropPort),
				ArgsUsage: "NAME HOST_PORT",
			},
			{
				Name:      "get",
				Usage:     "describe a load balancer",
				Action:    wrapAction(l.Command, l.Get),
				ArgsUsage: "NAME",
			},
			{
				Name:      "list",
				Usage:     "list all load balancers",
				Action:    wrapAction(l.Command, l.List),
				ArgsUsage: " ",
			},
		},
	}
}

func (l *LoadBalancerCommand) AddPort(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "NAME", "PORT")
	if err != nil {
		return err
	}

	port, err := parsePort(args["PORT"], c.String("certificate"))
	if err != nil {
		return err
	}

	id, err := l.resolveSingleID("load_balancer", args["NAME"])
	if err != nil {
		return err
	}

	loadBalancer, err := l.Client.GetLoadBalancer(id)
	if err != nil {
		return err
	}

	loadBalancer.Ports = append(loadBalancer.Ports, *port)
	loadBalancer, err = l.Client.UpdateLoadBalancer(id, loadBalancer.Ports)
	if err != nil {
		return err
	}

	return l.printLoadBalancer(loadBalancer)
}

func (l *LoadBalancerCommand) Create(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "ENVIRONMENT", "NAME")
	if err != nil {
		return err
	}

	ports := []models.Port{}
	for _, p := range c.StringSlice("port") {
		port, err := parsePort(p, c.String("certificate"))
		if err != nil {
			return err
		}

		ports = append(ports, *port)
	}

	if len(ports) == 0 {
		port := models.Port{
			HostPort:      80,
			ContainerPort: 80,
			Protocol:      "tcp",
		}

		ports = append(ports, port)
	}

	environmentID, err := l.resolveSingleID("environment", args["ENVIRONMENT"])
	if err != nil {
		return err
	}

	loadBalancer, err := l.Client.CreateLoadBalancer(args["NAME"], environmentID, ports, !c.Bool("private"))
	if err != nil {
		return err
	}

	return l.printLoadBalancer(loadBalancer)
}

func (l *LoadBalancerCommand) Delete(c *cli.Context) error {
	return l.deleteWithJob(c, "load_balancer", l.Client.DeleteLoadBalancer)
}

func (l *LoadBalancerCommand) DropPort(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "NAME", "HOST_PORT")
	if err != nil {
		return err
	}

	port, err := strconv.ParseInt(args["HOST_PORT"], 10, 64)
	if err != nil {
		return NewUsageError("'%s' is not a valid integer", args["PORT"])
	}

	id, err := l.resolveSingleID("load_balancer", args["NAME"])
	if err != nil {
		return err
	}

	loadBalancer, err := l.Client.GetLoadBalancer(id)
	if err != nil {
		return err
	}

	var exists bool
	for i, p := range loadBalancer.Ports {
		if p.HostPort == port {
			loadBalancer.Ports = append(loadBalancer.Ports[:i], loadBalancer.Ports[i+1:]...)
			exists = true
		}
	}

	if !exists {
		return fmt.Errorf("Host port '%v' doesn't exist on this Load Balancer", port)
	}

	loadBalancer, err = l.Client.UpdateLoadBalancer(id, loadBalancer.Ports)
	if err != nil {
		return err
	}

	return l.printLoadBalancer(loadBalancer)

}

func (l *LoadBalancerCommand) Get(c *cli.Context) error {
	return l.get(c, "load_balancer", func(id string) (entity.Entity, error) {
		loadBalancer, err := l.Client.GetLoadBalancer(id)
		if err != nil {
			return nil, err
		}

		return entity.NewLoadBalancer(loadBalancer), nil
	})
}

func (l *LoadBalancerCommand) List(c *cli.Context) error {
	loadBalancers, err := l.Client.ListLoadBalancers()
	if err != nil {
		return err
	}

	return l.printLoadBalancers(loadBalancers)
}

func parsePort(port, certificateName string) (*models.Port, error) {
	split := strings.FieldsFunc(port, func(r rune) bool {
		return r == ':' || r == '/'
	})

	if len(split) != 3 {
		return nil, NewUsageError("Port format is: HOST_PORT:CONTAINER_PORT/PROTOCOL")
	}

	hostPort, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil {
		return nil, NewUsageError("'%s' is not a valid integer", split[0])
	}

	containerPort, err := strconv.ParseInt(split[1], 10, 64)
	if err != nil {
		return nil, NewUsageError("'%s' is not a valid integer", split[1])
	}

	protocol := split[2]
	if strings.ToLower(protocol) != "https" {
		certificateName = ""
	}

	model := &models.Port{
		HostPort:        hostPort,
		ContainerPort:   containerPort,
		Protocol:        protocol,
		CertificateName: certificateName,
	}

	return model, nil
}

func (l *LoadBalancerCommand) printLoadBalancer(loadBalancer *models.LoadBalancer) error {
	entity := entity.NewLoadBalancer(loadBalancer)
	return l.Printer.PrintEntity(entity)
}

func (l *LoadBalancerCommand) printLoadBalancers(loadBalancers []*models.LoadBalancer) error {
	entities := []entity.Entity{}
	for _, loadBalancer := range loadBalancers {
		entities = append(entities, entity.NewLoadBalancer(loadBalancer))
	}

	return l.Printer.PrintEntities(entities)
}
