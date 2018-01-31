package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/quintilesims/layer0/common/config"
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
	dhc := config.DefaultLoadBalancerHealthCheck()
	dp := config.DefaultLoadBalancerPort()
	defaultPortString := fmt.Sprintf("%d:%d/%s", dp.HostPort, dp.ContainerPort, dp.Protocol)
	defaultPortFlag := cli.StringSlice([]string{defaultPortString})

	return cli.Command{
		Name:  "loadbalancer",
		Usage: "Manage layer0 Load Balancers",
		Subcommands: []cli.Command{
			{
				Name:      "addport",
				Usage:     "Add a new listener port (HOST_PORT:CONTAINER_PORT/PROTOCOL) to load balancer LOAD_BALANCER_NAME",
				Action:    l.addPort,
				ArgsUsage: "LOAD_BALANCER_NAME HOST_PORT:CONTAINER_PORT/PROTOCOL",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "certificate",
						Usage: "Name of SSL certificate to use for port configuration (only required for https)",
					},
				},
			},
			{
				Name:      "create",
				Usage:     "Create a new load balancer LOAD_BALANCER_NAME in the environment specified in ENVIRONMENT_TARGET",
				Action:    l.create,
				ArgsUsage: "ENVIRONMENT_TARGET LOAD_BALANCER_NAME",
				Flags: []cli.Flag{
					cli.StringSliceFlag{
						Name:  "port",
						Value: &defaultPortFlag,
						Usage: "Port configuration in format 'HOST_PORT:CONTAINER_PORT/PROTOCOL'",
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
						Value: dhc.Target,
						Usage: "Health check target in format 'PROTOCOL:PORT' or 'PROTOCOL:PORT/WITH/PATH'",
					},
					cli.IntFlag{
						Name:  "healthcheck-interval",
						Value: dhc.Interval,
						Usage: "Health check interval in seconds",
					},
					cli.IntFlag{
						Name:  "healthcheck-timeout",
						Value: dhc.Timeout,
						Usage: "Health check timeout in seconds",
					},
					cli.IntFlag{
						Name:  "healthcheck-healthy-threshold",
						Value: dhc.HealthyThreshold,
						Usage: "Number of consecutive successes required to count as healthy",
					},
					cli.IntFlag{
						Name:  "healthcheck-unhealthy-threshold",
						Value: dhc.UnhealthyThreshold,
						Usage: "Number of consecutive failures required to count as unhealthy",
					},
				},
			},
			{
				Name:      "delete",
				Usage:     "Delete the load balancer LOAD_BALANCER_NAME",
				ArgsUsage: "LOAD_BALANCER_NAME",
				Action:    l.delete,
			},
			{
				Name:      "dropport",
				Usage:     "Drop the listener with host port HOST_PORT from load balancer LOAD_BALANCER_NAME",
				Action:    l.dropPort,
				ArgsUsage: "LOAD_BALANCER_NAME HOST_PORT",
			},
			{
				Name:      "get",
				Usage:     "Describe load balancer LOAD_BALANCER_NAME",
				Action:    l.read,
				ArgsUsage: "LOAD_BALANCER_NAME",
			},
			{
				Name:      "healthcheck",
				Usage:     "View or update the health check of load balancer LOAD_BALANCER_NAME",
				Action:    l.healthcheck,
				ArgsUsage: "LOAD_BALANCER_NAME",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "healthcheck-target",
						Usage: "Health check target in format 'PROTOCOL:PORT' or 'PROTOCOL:PORT/WITH/PATH'",
					},
					cli.IntFlag{
						Name:  "healthcheck-interval",
						Usage: "Health check interval in seconds",
					},
					cli.IntFlag{
						Name:  "healthcheck-timeout",
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
				},
			},
			{
				Name:      "list",
				Usage:     "List all load balancers",
				Action:    l.list,
				ArgsUsage: " ",
			},
		},
	}
}

func (l *LoadBalancerCommand) addPort(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "LOAD_BALANCER_NAME", "PORT")
	if err != nil {
		return err
	}

	port, err := parsePort(args["PORT"], c.String("certificate"))
	if err != nil {
		return err
	}

	loadBalancerID, err := l.resolveSingleEntityIDHelper("load_balancer", args["LOAD_BALANCER_NAME"])
	if err != nil {
		return err
	}

	loadBalancer, err := l.client.ReadLoadBalancer(loadBalancerID)
	if err != nil {
		return err
	}

	ports := append(loadBalancer.Ports, *port)
	req := models.UpdateLoadBalancerRequest{
		Ports: &ports,
	}

	if err := l.client.UpdateLoadBalancer(loadBalancerID, req); err != nil {
		return err
	}

	loadBalancer, err = l.client.ReadLoadBalancer(loadBalancerID)
	if err != nil {
		return err
	}

	return l.printer.PrintLoadBalancers(loadBalancer)
}

func (l *LoadBalancerCommand) create(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "ENVIRONMENT_TARGET", "LOAD_BALANCER_NAME")
	if err != nil {
		return err
	}

	environmentID, err := l.resolveSingleEntityIDHelper("environment", args["ENVIRONMENT_TARGET"])
	if err != nil {
		return err
	}

	// remove the default port flag if --port was specified
	portFlags := c.StringSlice("port")
	if c.IsSet("port") {
		portFlags = portFlags[1:]
	}

	ports := []models.Port{}
	for _, p := range portFlags {
		port, err := parsePort(p, c.String("certificate"))
		if err != nil {
			return err
		}

		ports = append(ports, *port)
	}

	if target := c.String("healthcheck-target"); target != "" {
		if err := validateTarget(target); err != nil {
			return err
		}
	}

	healthCheck := models.HealthCheck{
		Target:             c.String("healthcheck-target"),
		Interval:           c.Int("healthcheck-interval"),
		Timeout:            c.Int("healthcheck-timeout"),
		HealthyThreshold:   c.Int("healthcheck-healthy-threshold"),
		UnhealthyThreshold: c.Int("healthcheck-unhealthy-threshold"),
	}

	req := models.CreateLoadBalancerRequest{
		LoadBalancerName: args["LOAD_BALANCER_NAME"],
		EnvironmentID:    environmentID,
		IsPublic:         !c.Bool("private"),
		Ports:            ports,
		HealthCheck:      healthCheck,
	}

	if err := req.Validate(); err != nil {
		return err
	}

	loadBalancerID, err := l.client.CreateLoadBalancer(req)
	if err != nil {
		return err
	}

	loadBalancer, err := l.client.ReadLoadBalancer(loadBalancerID)
	if err != nil {
		return err
	}

	return l.printer.PrintLoadBalancers(loadBalancer)
}

func (l *LoadBalancerCommand) delete(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "LOAD_BALANCER_NAME")
	if err != nil {
		return err
	}

	loadBalancerID, err := l.resolveSingleEntityIDHelper("load_balancer", args["LOAD_BALANCER_NAME"])
	if err != nil {
		return err
	}

	return l.client.DeleteLoadBalancer(loadBalancerID)
}

func (l *LoadBalancerCommand) dropPort(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "LOAD_BALANCER_NAME", "HOST_PORT")
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(args["HOST_PORT"])
	if err != nil {
		return fmt.Errorf("'%s' is not a valid integer", args["HOST_PORT"])
	}

	loadBalancerID, err := l.resolveSingleEntityIDHelper("load_balancer", args["LOAD_BALANCER_NAME"])
	if err != nil {
		return err
	}

	loadBalancer, err := l.client.ReadLoadBalancer(loadBalancerID)
	if err != nil {
		return err
	}

	var exists bool
	for i, p := range loadBalancer.Ports {
		if p.HostPort == int64(port) {
			loadBalancer.Ports = append(loadBalancer.Ports[:i], loadBalancer.Ports[i+1:]...)
			exists = true
			break
		}
	}

	if !exists {
		return fmt.Errorf("Host port '%v' doesn't exist on this Load Balancer", port)
	}

	req := models.UpdateLoadBalancerRequest{
		Ports: &loadBalancer.Ports,
	}

	if err := l.client.UpdateLoadBalancer(loadBalancerID, req); err != nil {
		return err
	}

	loadBalancer, err = l.client.ReadLoadBalancer(loadBalancerID)
	if err != nil {
		return err
	}

	return l.printer.PrintLoadBalancers(loadBalancer)
}

func (l *LoadBalancerCommand) healthcheck(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "LOAD_BALANCER_NAME")
	if err != nil {
		return err
	}

	loadBalancerID, err := l.resolveSingleEntityIDHelper("load_balancer", args["LOAD_BALANCER_NAME"])
	if err != nil {
		return err
	}

	loadBalancer, err := l.client.ReadLoadBalancer(loadBalancerID)
	if err != nil {
		return err
	}

	healthCheck := loadBalancer.HealthCheck
	requiresUpdate := false

	if target := c.String("healthcheck-target"); target != "" {
		if err := validateTarget(target); err != nil {
			return err
		}

		healthCheck.Target = target
		requiresUpdate = true
	}

	if interval := c.Int("healthcheck-interval"); interval != 0 {
		healthCheck.Interval = interval
		requiresUpdate = true
	}

	if timeout := c.Int("healthcheck-timeout"); timeout != 0 {
		healthCheck.Timeout = timeout
		requiresUpdate = true
	}

	if healthyThreshold := c.Int("healthy-threshold"); healthyThreshold != 0 {
		healthCheck.HealthyThreshold = healthyThreshold
		requiresUpdate = true
	}

	if unhealthyThreshold := c.Int("unhealthy-threshold"); unhealthyThreshold != 0 {
		healthCheck.UnhealthyThreshold = unhealthyThreshold
		requiresUpdate = true
	}

	if !requiresUpdate {
		return l.printer.PrintLoadBalancerHealthCheck(loadBalancer)
	}

	req := models.UpdateLoadBalancerRequest{
		HealthCheck: &healthCheck,
	}

	if err := l.client.UpdateLoadBalancer(loadBalancerID, req); err != nil {
		return err
	}

	loadBalancer, err = l.client.ReadLoadBalancer(loadBalancerID)
	if err != nil {
		return err
	}

	return l.printer.PrintLoadBalancerHealthCheck(loadBalancer)
}

func (l *LoadBalancerCommand) list(c *cli.Context) error {
	loadBalancerSummaries, err := l.client.ListLoadBalancers()
	if err != nil {
		return err
	}

	return l.printer.PrintLoadBalancerSummaries(loadBalancerSummaries...)
}

func (l *LoadBalancerCommand) read(c *cli.Context) error {
	args, err := extractArgs(c.Args(), "LOAD_BALANCER_NAME")
	if err != nil {
		return err
	}

	loadBalancerIDs, err := l.resolver.Resolve("load_balancer", args["LOAD_BALANCER_NAME"])
	if err != nil {
		return err
	}

	loadBalancers := make([]*models.LoadBalancer, len(loadBalancerIDs))
	for i, loadBalancerID := range loadBalancerIDs {
		loadBalancer, err := l.client.ReadLoadBalancer(loadBalancerID)
		if err != nil {
			return err
		}

		loadBalancers[i] = loadBalancer
	}

	return l.printer.PrintLoadBalancers(loadBalancers...)
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
	if strings.ToLower(protocol) == "https" && certificateName == "" {
		return nil, fmt.Errorf("HTTPS protocol specified in a port, but no certificate provided")
	}

	model := &models.Port{
		HostPort:        hostPort,
		ContainerPort:   containerPort,
		Protocol:        protocol,
		CertificateName: certificateName,
	}

	return model, nil
}

func validateTarget(target string) error {
	split := strings.FieldsFunc(target, func(r rune) bool {
		return r == ':' || r == '/'
	})

	protocol := strings.ToLower(split[0])
	if len(split) < 3 && (protocol == "https" || protocol == "http") {
		text := "HTTP & HTTPS targets must specify a port followed by a path.\n"
		text += "For example, HTTPS:443/health"
		return fmt.Errorf(text)
	}

	return nil
}
