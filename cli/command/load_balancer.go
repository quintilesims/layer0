package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
)

type LoadBalancerCommand struct {
	*CommandBase
}

func NewLoadBalancerCommand(b *CommandBase) *LoadBalancerCommand {
	return &LoadBalancerCommand{b}
}

func (l *LoadBalancerCommand) Command() cli.Command {
	return cli.Command{
		Name:  "loadbalancer",
		Usage: "manage layer0 load balancers",
		Subcommands: []cli.Command{
			{
				Name:      "addport",
				Usage:     "Add a new listener port (HOST_PORT:CONTAINER_PORT/PROTOCOL) to Load Balancer LOADBALANCER_NAME",
				Action:    l.addport,
				ArgsUsage: "LOADBALANCER_NAME HOST_PORT:CONTAINER_PORT/PROTOCOL",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "certificate",
						Usage: "Name of SSL certificate to use for port configuration (only required for https)",
					},
					cli.BoolFlag{
						Name:  "nowait",
						Usage: "Don't wait for the job to finish",
					},
				},
			},
			{
				Name:      "create",
				Usage:     "Create a new Load Balancer LOADBALANCER_NAME in Environment ENVIRONMENT_NAME",
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
				Usage:     "Delete the Load Balancer LOADBALANCER_NAME",
				ArgsUsage: "LOADBALANCER_NAME",
				Action:    l.delete,
			},
			{
				Name:      "dropport",
				Usage:     "Drop the listener with host port HOST_PORT from Load Balancer LOADBALANCER_NAME",
				Action:    l.dropport,
				ArgsUsage: "LOADBALANCER_NAME HOST_PORT",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "nowait",
						Usage: "Don't wait for the job to finish",
					},
				},
			},
			{
				Name:      "healthcheck",
				Usage:     "View or update the health check of Load Balancer LOADBALANCER_NAME",
				Action:    l.healthcheck,
				ArgsUsage: "LOADBALANCER_NAME",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "target",
						Value: "TCP:80",
						Usage: "Health check target in format 'PROTOCOL:PORT' or 'PROTOCOL:PORT/WITH/PATH'",
					},
					cli.IntFlag{
						Name:  "interval",
						Usage: "Health check interval in seconds",
					},
					cli.IntFlag{
						Name:  "timeout",
						Usage: "Health check timeout in seconds",
					},
					cli.IntFlag{
						Name:  "healthy-threshold",
						Usage: "Number of consecutive successes required to count as healthy",
					},
					cli.IntFlag{
						Name:  "unhealthy-threshold",
						Usage: "Number of consecutive failures required to count as unhealthy",
					},
					cli.BoolFlag{
						Name:  "nowait",
						Usage: "Don't wait for the job to finish",
					},
				},
			},
			{
				Name:      "list",
				Usage:     "List all Load Balancers",
				Action:    l.list,
				ArgsUsage: " ",
			},
			{
				Name:      "read",
				Usage:     "Describe Load Balancer LOADBALANCER_NAME",
				Action:    l.read,
				ArgsUsage: "LOADBALANCER_NAME",
			},
		},
	}
}

func (l *LoadBalancerCommand) addport(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "LOADBALANCER_NAME", "PORT")
	if err != nil {
		return err
	}

	port, err := parsePort(args["PORT"], c.String("certificate"))
	if err != nil {
		return err
	}

	loadBalancerID, err := resolveSingleEntityID(l.resolver, "loadbalancer", args["LOADBALANCER_NAME"])
	if err != nil {
		return err
	}

	loadBalancer, err := l.client.ReadLoadBalancer(loadBalancerID)
	if err != nil {
		return err
	}

	loadBalancer.Ports = append(loadBalancer.Ports, *port)
	req := models.UpdateLoadBalancerRequest{
		LoadBalancerID: loadBalancerID,
		Ports:          &loadBalancer.Ports,
		HealthCheck:    &loadBalancer.HealthCheck,
	}
	jobID, err := l.client.UpdateLoadBalancer(req)
	if err != nil {
		return err
	}

	if c.GlobalBool("config.FLAG_NO_WAIT") {
		l.printer.Printf("Running as job '%s'", jobID)
		return nil
	}

	return l.waitOnJobHelper(c, jobID, "Adding port", func(loadBalancerID string) error {
		loadBalancer, err := l.client.ReadLoadBalancer(loadBalancerID)
		if err != nil {
			return err
		}

		l.printer.StopSpinner()
		return l.printer.PrintLoadBalancers(loadBalancer)
	})
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

	ports := []models.Port{}
	for _, p := range c.StringSlice("port") {
		port, err := parsePort(p, c.String("certificate"))
		if err != nil {
			return err
		}

		ports = append(ports, *port)
	}

	// The default Port is filled here as opposed to the cli.Flag.Value
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

	return l.waitOnJobHelper(c, jobID, "Creating", func(loadBalancerID string) error {
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

func (l *LoadBalancerCommand) dropport(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "LOADBALANCER_NAME", "HOST_PORT")
	if err != nil {
		return err
	}

	port, err := strconv.ParseInt(args["HOST_PORT"], 10, 64)
	if err != nil {
		return fmt.Errorf("'%s' is not a valid integer", args["HOST_PORT"])
	}

	loadBalancerID, err := resolveSingleEntityID(l.resolver, "loadbalancer", args["LOADBALANCER_NAME"])
	if err != nil {
		return err
	}

	loadBalancer, err := l.client.ReadLoadBalancer(loadBalancerID)
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

	req := models.UpdateLoadBalancerRequest{
		LoadBalancerID: loadBalancerID,
		Ports:          &loadBalancer.Ports,
		HealthCheck:    &loadBalancer.HealthCheck,
	}
	jobID, err := l.client.UpdateLoadBalancer(req)
	if err != nil {
		return err
	}

	if c.GlobalBool("config.FLAG_NO_WAIT") {
		l.printer.Printf("Running as job '%s'", jobID)
		return nil
	}

	return l.waitOnJobHelper(c, jobID, "Dropping port", func(loadBalancerID string) error {
		loadBalancer, err := l.client.ReadLoadBalancer(loadBalancerID)
		if err != nil {
			return err
		}

		l.printer.StopSpinner()
		return l.printer.PrintLoadBalancers(loadBalancer)
	})
}

func (l *LoadBalancerCommand) healthcheck(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "LOADBALANCER_NAME")
	if err != nil {
		return err
	}

	updateIsRequired := false
	healthCheck := models.HealthCheck{}

	if target := c.String("target"); target != "" {
		updateIsRequired = true
		healthCheck.Target = target
	}

	if interval := c.String("interval"); interval != "" {
		i, err := strconv.Atoi(interval)
		if err != nil {
			return err
		}

		updateIsRequired = true
		healthCheck.Interval = i
	}

	if timeout := c.String("timeout"); timeout != "" {
		t, err := strconv.Atoi(timeout)
		if err != nil {
			return err
		}

		updateIsRequired = true
		healthCheck.Timeout = t
	}

	if healthyThreshold := c.String("healthy-threshold"); healthyThreshold != "" {
		h, err := strconv.Atoi(healthyThreshold)
		if err != nil {
			return err
		}

		updateIsRequired = true
		healthCheck.HealthyThreshold = h
	}

	if unhealthyThreshold := c.String("unhealthy-threshold"); unhealthyThreshold != "" {
		u, err := strconv.Atoi(unhealthyThreshold)
		if err != nil {
			return err
		}

		updateIsRequired = true
		healthCheck.UnhealthyThreshold = u
	}

	loadBalancerID, err := resolveSingleEntityID(l.resolver, "loadbalancer", args["LOADBALANCER_NAME"])
	if err != nil {
		return err
	}

	loadBalancer, err := l.client.ReadLoadBalancer(loadBalancerID)
	if err != nil {
		return err
	}

	if updateIsRequired {
		if healthCheck.Target != "" {
			loadBalancer.HealthCheck.Target = healthCheck.Target
		}

		if healthCheck.Interval != 0 {
			loadBalancer.HealthCheck.Interval = healthCheck.Interval
		}

		if healthCheck.Timeout != 0 {
			loadBalancer.HealthCheck.Timeout = healthCheck.Timeout
		}

		if healthCheck.HealthyThreshold != 0 {
			loadBalancer.HealthCheck.HealthyThreshold = healthCheck.HealthyThreshold
		}

		if healthCheck.UnhealthyThreshold != 0 {
			loadBalancer.HealthCheck.UnhealthyThreshold = healthCheck.UnhealthyThreshold
		}
	}

	req := models.UpdateLoadBalancerRequest{
		LoadBalancerID: loadBalancerID,
		Ports:          &loadBalancer.Ports,
		HealthCheck:    &loadBalancer.HealthCheck,
	}

	jobID, err := l.client.UpdateLoadBalancer(req)
	if err != nil {
		return err
	}

	if c.GlobalBool("config.FLAG_NO_WAIT") {
		l.printer.Printf("Running as job '%s'", jobID)
		return nil
	}

	return l.waitOnJobHelper(c, jobID, "Updating health check", func(loadBalancerID string) error {
		loadBalancer, err := l.client.ReadLoadBalancer(loadBalancerID)
		if err != nil {
			return err
		}

		l.printer.StopSpinner()
		return l.printer.PrintLoadBalancers(loadBalancer)
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
	args, err := extractArgs(c.Args(), "LOADBALANCER_NAME")
	if err != nil {
		return err
	}

	loadBalancer, err := l.client.ReadLoadBalancer(args["LOADBALANCER_NAME"])
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
