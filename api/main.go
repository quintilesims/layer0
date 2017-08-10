package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/quintilesims/layer0/api/controllers"
	"github.com/quintilesims/layer0/api/providers/aws"
	"github.com/quintilesims/layer0/api/scheduler"
	awsclient "github.com/quintilesims/layer0/common/aws"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/lock"
	"github.com/quintilesims/layer0/common/logging"
	"github.com/urfave/cli"
	"github.com/zpatrick/fireball"
)

// todo: handle main.Version

const (
	SWAGGER_URL     = "/api/"
	SWAGGER_UI_PATH = "static/swagger-ui/dist"
)

func serveSwaggerUI(w http.ResponseWriter, r *http.Request) {
	dir := http.Dir(SWAGGER_UI_PATH)
	fileServer := http.FileServer(dir)
	http.StripPrefix(SWAGGER_URL, fileServer).ServeHTTP(w, r)
}

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

		// todo: inject job_store.JobStore
		provider := aws.NewAWSProvider(client, nil)

		taskSchedulerLock := lock.NewDynamoDBExpiringLock(awsConfig, cfg.LockTable(), "TaskScheduler", time.Minute*5)
		taskScheduler := scheduler.NewECSTaskScheduler(taskSchedulerLock)
		defer taskScheduler.RunEvery(time.Minute).Stop()

		// todo: inject job scheduler
		routes := controllers.NewEnvironmentController(provider, nil).Routes()
		routes = append(routes, controllers.NewJobController(provider).Routes()...)

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
