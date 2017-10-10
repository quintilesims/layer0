package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

type LoadBalancerCommand struct {
	*CommandMediator
}

func NewLoadBalancerCommand(m *CommandMediator) *LoadBalancerCommand {
	return &LoadBalancerCommand{
		CommandMediator: m,
	}
}

func (l *LoadBalancerCommand) Command() cli.Command {
	return cli.Command{
		Name:  "loadbalancer",
		Usage: "manage layer0 load balancers",
		Subcommands: []cli.Command{
			{
				Name:      "create",
				Usage:     "Create a new Elastic Load Balancer",
				Action:    l.create,
				ArgsUsage: "ENVIRONMENT_NAME LOADBALANCER_NAME",
				Flags: []cli.Flag{
					cli.StringSliceFlag{
						Name:  "port",
						Usage: "Port configuration in format 'HOST_PORT:CONTAINER_PORT/PROTOCOL' (default 80:80/tcp)",
					},
					cli.StringFlag{
						Name:  "certificate",
						Usage: "Name of SSL certificate to use for port configuration (only required for https)",
					},
					cli.BoolFlag{
						Name:  "private",
						Usage: "If specified, creates a private load balancer (default is public)",
					},
					cli.StringFlag{
						Name:  "healthcheck-target",
						Value: "TCP:80",
						Usage: "Health check target in format 'PROTOCOL:PORT' or 'PROTOCOL:PORT/WITH/PATH' (default is TCP:80)",
					},
					cli.IntFlag{
						Name:  "healthcheck-interval",
						Value: 30,
						Usage: "Health check interval in seconds (default is 30)",
					},
					cli.IntFlag{
						Name:  "healthcheck-timeout",
						Value: 5,
						Usage: "Health check timeout in seconds (default is 5)",
					},
					cli.IntFlag{
						Name:  "healthcheck-healthy-threshold",
						Value: 2,
						Usage: "Number of consecutive successes required to count as healthy (default is 2)",
					},
					cli.IntFlag{
						Name:  "healthcheck-unhealthy-threshold",
						Value: 2,
						Usage: "Number of consecutive failures required to count as unhealthy (default is 2)",
					},
					cli.BoolFlag{
						Name:  "nowait",
						Usage: "Don't wait for the job to finish",
					},
				},
			},
			{
				Name:      "delete",
				Usage:     "delete a load balancer",
				ArgsUsage: "NAME",
				Action:    l.delete,
				Flags: []cli.Flag{
					cli.BoolTFlag{
						Name:  "wait",
						Usage: "wait for the job to complete before returning",
					},
				},
			},
			{
				Name:      "list",
				Usage:     "list all load balancers",
				Action:    l.list,
				ArgsUsage: " ",
			},
			{
				Name:      "read",
				Usage:     "describe a load balancer",
				Action:    l.read,
				ArgsUsage: "NAME",
			},
		},
	}
}

func (l *LoadBalancerCommand) create(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "ENVIRONMENT_NAME", "LOADBALANCER_NAME")
	if err != nil {
		return err
	}

	environmentID, err := resolveSingleEntityID(l.resolver, "environment", args["ENVIRONMENT_NAME"])
	if err != nil {
		return err
	}

	fmt.Printf("environmentID: %v\n", environmentID)

	ports := []models.Port{}
	for _, p := range c.StringSlice("port") {
		port, err := parsePort(p, c.String("certificate"))
		if err != nil {
			return err
		}

		ports = append(ports, *port)
	}

	// TODO: Should we be defaulting load balancers to this configuration or
	// require users to specify the listener ports?
	if len(ports) == 0 {
		port := models.Port{
			HostPort:      80,
			ContainerPort: 80,
			Protocol:      "tcp",
		}

		ports = append(ports, port)
	}

	healthCheck := models.HealthCheck{
		Target:             c.String("healthcheck-target"),
		Interval:           c.Int("healthcheck-interval"),
		Timeout:            c.Int("healthcheck-timeout"),
		HealthyThreshold:   c.Int("healthcheck-healthy-threshold"),
		UnhealthyThreshold: c.Int("healthcheck-unhealthy-threshold"),
	}

	req := models.CreateLoadBalancerRequest{
		LoadBalancerName: args["LOADBALANCER_NAME"],
		EnvironmentID:    environmentID,
		IsPublic:         c.Bool("private"),
		Ports:            ports,
		HealthCheck:      healthCheck,
	}

	jobID, err := l.client.CreateLoadBalancer(req)
	if err != nil {
		return err
	}

	if c.GlobalBool("config.FLAG_NO_WAIT") {
		l.printer.Printf("Running as job '%s'", jobID)
		return nil
	}

	return l.waitOnJobHelper(c, jobID, "creating", func(loadBalancerID string) error {
		loadBalancer, err := l.client.ReadLoadBalancer(loadBalancerID)
		if err != nil {
			return err
		}

		l.printer.StopSpinner()
		return l.printer.PrintLoadBalancers(loadBalancer)
	})
}

func (l *LoadBalancerCommand) delete(c *cli.Context) error {
	return l.deleteHelper(c, "loadbalancer", func(loadBalancerID string) (string, error) {
		return l.client.DeleteLoadBalancer(loadBalancerID)
	})
}

func (l *LoadBalancerCommand) list(c *cli.Context) error {
	loadBalancerSummaries, err := l.client.ListLoadBalancers()
	if err != nil {
		return err
	}

	return l.printer.PrintLoadBalancerSummaries(loadBalancerSummaries...)
}

func (l *LoadBalancerCommand) read(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "NAME")
	if err != nil {
		return err
	}

	loadBalancer, err := l.client.ReadLoadBalancer(args["NAME"])
	if err != nil {
		return err
	}

	return l.printer.PrintLoadBalancers(loadBalancer)
}

func parsePort(port, certificateName string) (*models.Port, error) {
	split := strings.FieldsFunc(port, func(r rune) bool {
		return r == ':' || r == '/'
	})

	if len(split) != 3 {
		return nil, fmt.Errorf("Port format is: HOST_PORT:CONTAINER_PORT/PROTOCOL")
	}

	hostPort, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("'%s' is not a valid integer", split[0])
	}

	containerPort, err := strconv.ParseInt(split[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("'%s' is not a valid integer", split[1])
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
