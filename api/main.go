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
		if err := config.ValidateAPIContext(c); err != nil {
			return err
		}

		logger := logging.NewLogWriter(c.Bool("debug"))
		log.SetOutput(logger)

		awsConfig := defaults.Get().Config
		staticCreds := credentials.NewStaticCredentials(
			c.String(config.FlagAWSAccessKey.GetName()),
			c.String(config.FlagAWSSecretKey.GetName()),
			"")

		awsConfig.WithCredentials(staticCreds)
		awsConfig.WithRegion(c.String(config.FlagAWSRegion.GetName()))
		awsConfig.WithMaxRetries(config.DefaultMaxRetries)
		session := session.New(awsConfig)

		delay := c.Duration(config.FlagAWSRequestDelay.GetName())
		ticker := time.Tick(delay)
		session.Handlers.Send.PushBack(func(r *request.Request) {
			<-ticker
		})

		client := awsclient.NewClient(session)
		tagStore := tag.NewDynamoStore(session, c.String(config.FlagAWSTagTable.GetName()))
		jobStore := job.NewDynamoStore(session, c.String(config.FlagAWSJobTable.GetName()))

		adminProvider := aws.NewAdminProvider(client, tagStore, c)
		deployProvider := aws.NewDeployProvider(client, tagStore, c)
		environmentProvider := aws.NewEnvironmentProvider(client, tagStore, c)
		loadBalancerProvider := aws.NewLoadBalancerProvider(client, tagStore, c)
		serviceProvider := aws.NewServiceProvider(client, tagStore, c)
		taskProvider := aws.NewTaskProvider(client, tagStore, c)
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
		routes = append(routes, controllers.NewAdminController(c, Version).Routes()...)
		routes = append(routes, controllers.NewDeployController(deployProvider, jobStore, tagStore).Routes()...)
		routes = append(routes, controllers.NewEnvironmentController(environmentProvider, jobStore, tagStore).Routes()...)
		routes = append(routes, controllers.NewJobController(jobStore, tagStore).Routes()...)
		routes = append(routes, controllers.NewLoadBalancerController(loadBalancerProvider, jobStore, tagStore).Routes()...)
		routes = append(routes, controllers.NewServiceController(serviceProvider, jobStore, tagStore).Routes()...)
		routes = append(routes, controllers.NewTagController(tagStore).Routes()...)
		routes = append(routes, controllers.NewTaskController(taskProvider, jobStore, tagStore).Routes()...)

		user, pass, err := config.ParseAuthToken(c)
		if err != nil {
			return err
		}

		routes = fireball.Decorate(routes,
			fireball.LogDecorator(),
			fireball.BasicAuthDecorator(user, pass))

		server := fireball.NewApp(routes)
		server.ErrorHandler = controllers.ErrorHandler

		lockTable := c.String(config.FlagAWSLockTable.GetName())
		lockExpiry := c.Duration(config.FlagJobExpiry.GetName())
		jobLock := lock.NewDynamoLock(session, lockTable, lockExpiry)
		daemonLock := lock.NewDynamoLock(session, lockTable, time.Minute*5)

		// todo: get num workers from config
		jobTicker, stopWorkers := job.RunWorkersAndDispatcher(2, jobStore, jobRunner, jobLock)
		defer jobTicker.Stop()
		defer stopWorkers()

		sdFN := scaler.NewDaemonFN(jobStore, environmentProvider)
		scalerDaemon := daemon.NewDaemon("Scaler", "ScalerDaemon", daemonLock, sdFN)
		scalerDaemonTicker := scalerDaemon.RunEvery(time.Hour)
		defer scalerDaemonTicker.Stop()

		jobExpiry := c.Duration(config.FlagJobExpiry.GetName())
		jdFN := job.NewDaemonFN(jobStore, jobExpiry)
		jobDaemon := daemon.NewDaemon("Job", "JobDaemon", daemonLock, jdFN)
		jobDaemonTicker := jobDaemon.RunEvery(time.Hour)
		defer jobDaemonTicker.Stop()

		tdFN := tag.NewDaemonFN(tagStore, taskProvider)
		tagDaemon := daemon.NewDaemon("Tag", "TagDaemon", daemonLock, tdFN)
		tagDaemonTicker := tagDaemon.RunEvery(time.Hour)
		defer tagDaemonTicker.Stop()

		ldFN := lock.NewDaemonFN(daemonLock, lockExpiry)
		lockDaemon := daemon.NewDaemon("Lock", "LockDaemon", daemonLock, ldFN)
		lockDaemonTicker := lockDaemon.RunEvery(time.Hour)
		defer lockDaemonTicker.Stop()

		port := c.Int(config.FlagPort.GetName())
		log.Printf("[INFO] Listening on port %d", port)
		http.Handle("/", server)

		http.HandleFunc(SWAGGER_URL, serveSwaggerUI)
		return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
