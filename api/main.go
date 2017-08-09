package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	restful "github.com/emicklei/go-restful"
	"github.com/quintilesims/layer0/api/controllers"
	"github.com/quintilesims/layer0/api/providers/aws"
	awsclient "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/urfave/cli"
	"github.com/zpatrick/example/logging"
)

// todo: handle swagger
// todo: handle main.Version

func main() {
	app := cli.NewApp()
	app.Name = "Layer0 API"
	app.Flags = config.APIFlags()
	app.Action = func(c *cli.Context) error {
		cfg := config.NewContextAPIConfig(c)
		if err := cfg.Validate(); err != nil {
			return err
		}

		logger := logging.NewLogWriter(c.Bool("debug"))
		log.SetOutput(logger)

		awsConfig := defaults.Get().Config
		staticCreds := credentials.NewStaticCredentials(cfg.AccessKey(), cfg.SecretKey(), "")
		awsConfig.WithCredentials(staticCreds)
		awsConfig.WithRegion(cfg.Region())

		client := awsclient.NewClient(awsConfig)
		provider := aws.NewAWSProvider(client)

		// todo: inject job scheduler
		environmentController := controllers.NewEnvironmentController(provider, nil)
		restful.Add(environmentController.Routes())

		// todo: add restful.Filters
		// todo: add swagger service

		log.Printf("[INFO] Listening on port %d", cgf.Port())

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
