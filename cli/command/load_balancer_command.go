package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/quintilesims/layer0/common/models"
	"github.com/urfave/cli"
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
						Usage: "name or arn of certificate to use for port configuration (only required for https)",
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
					cli.IntFlag{
						Name:  "idle-timeout",
						Value: 60,
						Usage: "idle timeout of the load balancer in seconds",
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
				Name:      "healthcheck",
				Usage:     "view or update the health check for a load balancer",
				Action:    wrapAction(l.Command, l.HealthCheck),
				ArgsUsage: "NAME",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "set-target",
						Usage: "health check target in format 'PROTOCOL:PORT' or 'PROTOCOL:PORT/WITH/PATH'",
					},
					cli.StringFlag{
						Name:  "set-interval",
						Usage: "health check interval in seconds",
					},
					cli.StringFlag{
						Name:  "set-timeout",
						Usage: "health check timeout in seconds",
					},
					cli.StringFlag{
						Name:  "set-healthy-threshold",
						Usage: "number of consecutive successes required to count as healthy",
					},
					cli.StringFlag{
						Name:  "set-unhealthy-threshold",
						Usage: "number of consecutive failures required to count as unhealthy",
					},
				},
			},
			{
				Name:      "idletimeout",
				Usage:     "update the idle timeout for a load balancer",
				Action:    wrapAction(l.Command, l.IdleTimeout),
				ArgsUsage: "NAME TIMEOUT",
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
	loadBalancer, err = l.Client.UpdateLoadBalancerPorts(id, loadBalancer.Ports)
	if err != nil {
		return err
	}

	return l.Printer.PrintLoadBalancers(loadBalancer)
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

	healthCheck := models.HealthCheck{
		Target:             c.String("healthcheck-target"),
		Interval:           c.Int("healthcheck-interval"),
		Timeout:            c.Int("healthcheck-timeout"),
		HealthyThreshold:   c.Int("healthcheck-healthy-threshold"),
		UnhealthyThreshold: c.Int("healthcheck-unhealthy-threshold"),
	}

	environmentID, err := l.resolveSingleID("environment", args["ENVIRONMENT"])
	if err != nil {
		return err
	}

	idleTimeout := c.Int("idle-timeout")
	loadBalancer, err := l.Client.CreateLoadBalancer(args["NAME"], environmentID, healthCheck, ports, !c.Bool("private"), idleTimeout)
	if err != nil {
		return err
	}

	return l.Printer.PrintLoadBalancers(loadBalancer)
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

	loadBalancer, err = l.Client.UpdateLoadBalancerPorts(id, loadBalancer.Ports)
	if err != nil {
		return err
	}

	return l.Printer.PrintLoadBalancers(loadBalancer)
}

func (l *LoadBalancerCommand) Get(c *cli.Context) error {
	loadBalancers := []*models.LoadBalancer{}
	getLoadBalancerf := func(id string) error {
		loadBalancer, err := l.Client.GetLoadBalancer(id)
		if err != nil {
			return err
		}

		loadBalancers = append(loadBalancers, loadBalancer)
		return nil
	}

	if err := l.get(c, "load_balancer", getLoadBalancerf); err != nil {
		return err
	}

	return l.Printer.PrintLoadBalancers(loadBalancers...)
}

func (l *LoadBalancerCommand) HealthCheck(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "NAME")
	if err != nil {
		return err
	}

	updateIsRequired := false
	healthCheck := models.HealthCheck{}

	if target := c.String("set-target"); target != "" {
		updateIsRequired = true
		healthCheck.Target = target
	}

	if interval := c.String("set-interval"); interval != "" {
		i, err := strconv.Atoi(interval)
		if err != nil {
			return err
		}

		updateIsRequired = true
		healthCheck.Interval = i
	}

	if timeout := c.String("set-timeout"); timeout != "" {
		t, err := strconv.Atoi(timeout)
		if err != nil {
			return err
		}

		updateIsRequired = true
		healthCheck.Timeout = t
	}

	if healthyThreshold := c.String("set-healthy-threshold"); healthyThreshold != "" {
		h, err := strconv.Atoi(healthyThreshold)
		if err != nil {
			return err
		}

		updateIsRequired = true
		healthCheck.HealthyThreshold = h
	}

	if unhealthyThreshold := c.String("set-unhealthy-threshold"); unhealthyThreshold != "" {
		u, err := strconv.Atoi(unhealthyThreshold)
		if err != nil {
			return err
		}

		updateIsRequired = true
		healthCheck.UnhealthyThreshold = u
	}

	id, err := l.resolveSingleID("load_balancer", args["NAME"])
	if err != nil {
		return err
	}

	loadBalancer, err := l.Client.GetLoadBalancer(id)
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

		loadBalancer, err = l.Client.UpdateLoadBalancerHealthCheck(id, loadBalancer.HealthCheck)
		if err != nil {
			return err
		}
	}

	return l.Printer.PrintLoadBalancerHealthCheck(loadBalancer)
}

func (l *LoadBalancerCommand) IdleTimeout(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "NAME", "TIMEOUT")
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

	idleTimeout, err := strconv.Atoi(args["TIMEOUT"])
	if err != nil {
		return err
	}

	loadBalancer, err = l.Client.UpdateLoadBalancerIdleTimeout(id, idleTimeout)
	if err != nil {
		return err
	}

	return l.Printer.PrintLoadBalancerIdleTimeout(loadBalancer)
}

func (l *LoadBalancerCommand) List(c *cli.Context) error {
	loadBalancerSummaries, err := l.Client.ListLoadBalancers()
	if err != nil {
		return err
	}

	return l.Printer.PrintLoadBalancerSummaries(loadBalancerSummaries...)
}

func parsePort(port, certificate string) (*models.Port, error) {
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
	var certificateName string
	var certificateARN string

	if strings.ToLower(protocol) == "https" {
		if strings.HasPrefix(strings.ToLower(certificate), "arn:") {
			certificateARN = certificate
		} else {
			certificateName = certificate
		}
	}

	model := &models.Port{
		HostPort:        hostPort,
		ContainerPort:   containerPort,
		Protocol:        protocol,
		CertificateName: certificateName,
		CertificateARN:  certificateARN,
	}

	return model, nil
}
