package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/quintilesims/layer0/client"
	"github.com/quintilesims/layer0/common/config"
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
				Usage:     "create a new load balancer",
				Action:    l.create,
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
					cli.StringFlag{
						Name:  "healthcheck-target",
						Value: "TCP:80",
						Usage: "health check target in format 'PROTOCOL:PORT' or 'PROTOCOL:PORT/WITH/PATH'",
					},
					cli.IntFlag{
						Name:  "healthcheck-interval",
						Value: 30,
						Usage: "health check interval in seconds",
					},
					cli.IntFlag{
						Name:  "healthcheck-timeout",
						Value: 5,
						Usage: "health check timeout in seconds",
					},
					cli.IntFlag{
						Name:  "healthcheck-healthy-threshold",
						Value: 2,
						Usage: "number of consecutive successes required to count as healthy",
					},
					cli.IntFlag{
						Name:  "healthcheck-unhealthy-threshold",
						Value: 2,
						Usage: "number of consecutive failures required to count as unhealthy",
					},
					cli.BoolFlag{
						Name:  "nowait",
						Usage: "don't wait for the job to finish",
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
		},
	}
}

func (l *LoadBalancerCommand) create(c *cli.Context) error {
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

	healthCheck := models.HealthCheck{
		Target:             c.String("healthcheck-target"),
		Interval:           c.Int("healthcheck-interval"),
		Timeout:            c.Int("healthcheck-timeout"),
		HealthyThreshold:   c.Int("healthcheck-healthy-threshold"),
		UnhealthyThreshold: c.Int("healthcheck-unhealthy-threshold"),
	}

	req := models.CreateLoadBalancerRequest{
		LoadBalancerName: args["NAME"],
		// TODO: Retrieve the actual EnvironmentID
		EnvironmentID: "jparson4db98",
		IsPublic:      c.Bool("private"),
		Ports:         ports,
		HealthCheck:   healthCheck,
	}

	jobID, err := l.client.CreateLoadBalancer(req)
	if err != nil {
		return err
	}

	if c.GlobalBool("config.FLAG_NO_WAIT") {
		l.printer.Printf("Running as job '%s'", jobID)
		return nil
	}

	l.printer.StartSpinner("creating")
	defer l.printer.StopSpinner()

	job, err := client.WaitForJob(l.client, jobID, c.GlobalDuration(config.FLAG_TIMEOUT))
	if err != nil {
		return err
	}

	loadBalancerID := job.Result
	loadBalancer, err := l.client.ReadLoadBalancer(loadBalancerID)
	if err != nil {
		return err
	}

	return l.printer.PrintLoadBalancers(loadBalancer)
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
