package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/quintilesims/layer0/api/controllers"
	"github.com/quintilesims/layer0/api/daemon"
	"github.com/quintilesims/layer0/api/job"
	"github.com/quintilesims/layer0/api/lock"
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
		awsConfig.WithMaxRetries(cfg.MaxRetries())
		session := session.New(awsConfig)

		delay := c.Duration(config.FLAG_AWS_TIME_BETWEEN_REQUESTS)
		ticker := time.Tick(delay)
		session.Handlers.Send.PushBack(func(r *request.Request) {
			<-ticker
		})

		client := awsclient.NewClient(session)
		tagStore := tag.NewDynamoStore(session, cfg.DynamoTagTable())
		jobStore := job.NewDynamoStore(session, cfg.DynamoJobTable())

		adminProvider := aws.NewAdminProvider(client, tagStore, cfg)
		deployProvider := aws.NewDeployProvider(client, tagStore, cfg)
		environmentProvider := aws.NewEnvironmentProvider(client, tagStore, cfg)
		loadBalancerProvider := aws.NewLoadBalancerProvider(client, tagStore, cfg)
		serviceProvider := aws.NewServiceProvider(client, tagStore, cfg)
		taskProvider := aws.NewTaskProvider(client, tagStore, cfg)
		environmentScaler := aws.NewEnvironmentScaler()
		scalerDispatcher := scaler.NewDispatcher(jobStore, time.Second*15)

		if err := adminProvider.Init(); err != nil {
			return err
		}

		jobRunner := aws.NewJobRunner(
			deployProvider,
			environmentProvider,
			loadBalancerProvider,
			serviceProvider,
			taskProvider,
			environmentScaler,
			scalerDispatcher)

		routes := controllers.NewSwaggerController(Version).Routes()
		routes = append(routes, controllers.NewEnvironmentController(environmentProvider, jobStore, tagStore).Routes()...)
		routes = append(routes, controllers.NewJobController(jobStore, tagStore).Routes()...)
		routes = append(routes, controllers.NewLoadBalancerController(loadBalancerProvider, jobStore, tagStore).Routes()...)
		routes = append(routes, controllers.NewServiceController(serviceProvider, jobStore, tagStore).Routes()...)
		routes = append(routes, controllers.NewTagController(tagStore).Routes()...)
		routes = append(routes, controllers.NewTaskController(taskProvider, jobStore, tagStore).Routes()...)

		adminRoutes := controllers.NewAdminController(cfg, Version).Routes()

		// Admin health check should be safe to not require basic auth
		for i, v := range adminRoutes {
			if v.Path != "/admin/health" {
				routes = append(routes, adminRoutes[i])
			}
		}

		user, pass, err := cfg.ParseAuthToken()
		if err != nil {
			return err
		}

		routes = fireball.Decorate(routes,
			fireball.LogDecorator(),
			fireball.BasicAuthDecorator(user, pass))

		for i, v := range adminRoutes {
			if v.Path == "/admin/health" {
				routes = append(routes, adminRoutes[i])
			}
		}

		server := fireball.NewApp(routes)
		server.ErrorHandler = controllers.ErrorHandler

		jobLock := lock.NewDynamoLock(session, cfg.DynamoLockTable(), cfg.JobExpiry())
		daemonLock := lock.NewDynamoLock(session, cfg.DynamoLockTable(), time.Minute*5)

		// todo: get num workers from config
		jobTicker, stopWorkers := job.RunWorkersAndDispatcher(2, jobStore, jobRunner, jobLock)
		defer jobTicker.Stop()
		defer stopWorkers()

		sdFN := scaler.NewDaemonFN(jobStore, environmentProvider)
		scalerDaemon := daemon.NewDaemon("Scaler", "ScalerDaemon", daemonLock, sdFN)
		scalerDaemonTicker := scalerDaemon.RunEvery(time.Hour)
		defer scalerDaemonTicker.Stop()

		jdFN := job.NewDaemonFN(jobStore, cfg.JobExpiry())
		jobDaemon := daemon.NewDaemon("Job", "JobDaemon", daemonLock, jdFN)
		jobDaemonTicker := jobDaemon.RunEvery(time.Hour)
		defer jobDaemonTicker.Stop()

		tdFN := tag.NewDaemonFN(tagStore, taskProvider)
		tagDaemon := daemon.NewDaemon("Tag", "TagDaemon", daemonLock, tdFN)
		tagDaemonTicker := tagDaemon.RunEvery(time.Hour)
		defer tagDaemonTicker.Stop()

		ldFN := lock.NewDaemonFN(daemonLock, cfg.LockExpiry())
		lockDaemon := daemon.NewDaemon("Lock", "LockDaemon", daemonLock, ldFN)
		lockDaemonTicker := lockDaemon.RunEvery(time.Hour)
		defer lockDaemonTicker.Stop()

		log.Printf("[INFO] Listening on port %d", cfg.Port())
		http.Handle("/", server)

		http.HandleFunc(SWAGGER_URL, serveSwaggerUI)
		return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port()), nil)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
