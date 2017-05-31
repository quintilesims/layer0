package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/common/aws/provider"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/logutils"
	"github.com/quintilesims/layer0/common/startup"
	"github.com/quintilesims/layer0/common/types"
	"github.com/quintilesims/layer0/runner/job"
	"github.com/urfave/cli"
	"os"
	"strings"
)

var Version string

func main() {
	app := cli.NewApp()
	app.Name = "Layer0 Runner"
	app.Usage = "Run a Layer0 Job"
	app.Version = getVersion()
	app.Action = Run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "job, j",
			Usage:  "job id to run",
			EnvVar: config.JOB_ID,
		},
		cli.StringFlag{
			Name:   "access_key, a",
			Usage:  "aws access key id",
			EnvVar: config.AWS_ACCESS_KEY_ID,
		},
		cli.StringFlag{
			Name:   "secret_key, s",
			Usage:  "aws secret access key",
			EnvVar: config.AWS_SECRET_ACCESS_KEY,
		},
		cli.StringFlag{
			Name:   "prefix, p",
			Usage:  "layer0 prefix",
			Value:  "l0",
			EnvVar: config.PREFIX,
		},
		cli.StringFlag{
			Name:   "region, r",
			Usage:  "aws region",
			Value:  "us-west-2",
			EnvVar: config.AWS_REGION,
		},
	}

	app.Run(os.Args)
}

func getVersion() string {
	if Version == "" {
		Version = "0.0.1x-unset-develop"
	}

	return Version
}

func Run(c *cli.Context) {
	if err := config.Validate(config.RequiredRunnerVariables); err != nil {
		logrus.Fatal(err)
	}

	switch strings.ToLower(config.RunnerLogLevel()) {
	case "0", "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "1", "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "2", "warning":
		logrus.SetLevel(logrus.WarnLevel)
	case "3", "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "4", "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	log := logutils.NewStackTraceLogger("Main")
	logutils.SetGlobalLogger(log)

	if c.String("job") == "" {
		log.Fatal("JOB_ID not specified")
	}

	credProvider := provider.NewExplicitCredProvider(c.String("access_key"), c.String("secret_key"))
	backend, err := startup.GetBackend(credProvider, c.String("region"))
	if err != nil {
		log.Fatal(err)
	}

	logic, err := startup.GetLogic(backend)
	if err != nil {
		log.Fatal(err)
	}

	if err := logic.JobStore.Init(); err != nil {
		log.Fatal(err)
	}

	if err := logic.TagStore.Init(); err != nil {
		log.Fatal(err)
	}

	runner := job.NewJobRunner(logic, c.String("job"))

	if err := runner.Load(); err != nil {
		runner.MarkStatus(types.Error)
		log.Fatal(err)
	}

	if err := runner.Run(); err != nil {
		runner.MarkStatus(types.Error)
		log.Fatal(err)
	}

	log.Info("Done")
}
