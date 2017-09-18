package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/quintilesims/layer0/api/controllers"
	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/api/provider/aws"
	"github.com/quintilesims/layer0/api/scaler"
	"github.com/quintilesims/layer0/api/tag"
	awsclient "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
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
		session := session.New(awsConfig)

		client := awsclient.NewClient(awsConfig)
		tagStore := tag.NewDynamoStore(session, cfg.DynamoTagTable())
		jobStore := job.NewDynamoStore(session, cfg.DynamoJobTable())

		deployProvider := aws.NewDeployProvider(client, tagStore)
		environmentProvider := aws.NewEnvironmentProvider(client, tagStore, cfg)
		loadBalancerProvider := aws.NewLoadBalancerProvider(client, tagStore, cfg)
		serviceProvider := aws.NewServiceProvider(client, tagStore)
		taskProvider := aws.NewTaskProvider(client, tagStore)
		environmentScaler := aws.NewEnvironmentScaler()
		scalerDispatcher := scaler.NewDispatcher(environmentProvider, environmentScaler)

		jobRunner := aws.NewJobRunner(
			deployProvider,
			environmentProvider,
			loadBalancerProvider,
			serviceProvider,
			taskProvider,
			scalerDispatcher)

		routes := controllers.NewSwaggerController(Version).Routes()
		routes = append(routes, controllers.NewDeployController(deployProvider, jobStore).Routes()...)
		routes = append(routes, controllers.NewEnvironmentController(environmentProvider, jobStore).Routes()...)
		routes = append(routes, controllers.NewJobController(jobStore).Routes()...)
		routes = append(routes, controllers.NewLoadBalancerController(loadBalancerProvider, jobStore).Routes()...)
		routes = append(routes, controllers.NewServiceController(serviceProvider, jobStore).Routes()...)
		routes = append(routes, controllers.NewTaskController(taskProvider, jobStore).Routes()...)

		// todo: add decorators to routes
		server := fireball.NewApp(routes)

		// todo: get num workers from config
		jobTicker := job.RunWorkersAndDispatcher(2, jobStore, jobRunner)
		defer jobTicker.Stop()

		scalerTicker := scalerDispatcher.RunEvery(time.Minute * 5)
		defer scalerTicker.Stop()

		log.Printf("[INFO] Listening on port %d", cfg.Port())
		http.Handle("/", server)

		http.HandleFunc(SWAGGER_URL, serveSwaggerUI)
		return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port()), nil)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
