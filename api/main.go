package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/quintilesims/layer0/api/controllers"
	"github.com/quintilesims/layer0/api/provider/aws"
	awsclient "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/db/tag_store"
	"github.com/quintilesims/layer0/common/logging"
	"github.com/urfave/cli"
	"github.com/zpatrick/fireball"
)

const (
	SWAGGER_URL     = "/api/"
	SWAGGER_UI_PATH = "static/swagger-ui/dist"
)

func serveSwaggerUI(w http.ResponseWriter, r *http.Request) {
	dir := http.Dir(SWAGGER_UI_PATH)
	fileServer := http.FileServer(dir)
	http.StripPrefix(SWAGGER_URL, fileServer).ServeHTTP(w, r)
}

var Version string

func main() {
	if Version == "" {
		Version = "unset/developer"
	}

	app := cli.NewApp()
	app.Name = "Layer0 API"
	app.Version = Version
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
		tagStore := tag_store.NewDynamoTagStore(awsConfig, cfg.DynamoTagTable())

		// todo: inject job_store.JobStore
		environmentProvider := aws.NewEnvironmentProvider(client, tagStore, cfg)
		serviceProvider := aws.NewServiceProvider(client, nil)
		deployProvider := aws.NewDeployProvider(client, nil)
		loadbalancerProvider := aws.NewLoadBalancerProvider(client, nil)
		taskProvider := aws.NewTaskProvider(client, nil)

		// todo: inject job scheduler
		routes := controllers.NewSwaggerController(Version).Routes()
		routes = append(routes, controllers.NewEnvironmentController(environmentProvider, nil).Routes()...)
		routes = append(routes, controllers.NewServiceController(serviceProvider, nil).Routes()...)
		routes = append(routes, controllers.NewDeployController(deployProvider).Routes()...)
		routes = append(routes, controllers.NewLoadBalancerController(loadbalancerProvider, nil).Routes()...)
		routes = append(routes, controllers.NewTaskController(taskProvider).Routes()...)

		// todo: add decorators to routes
		server := fireball.NewApp(routes)

		log.Printf("[INFO] Listening on port %d", cfg.Port())
		http.Handle("/", server)

		http.HandleFunc(SWAGGER_URL, serveSwaggerUI)
		return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port()), nil)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
