package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
	swagger "github.com/emicklei/go-restful-swagger12"
	"github.com/quintilesims/layer0/api/handlers"
	"github.com/quintilesims/layer0/api/logic"
	"github.com/quintilesims/layer0/common/aws/provider"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/logutils"
	"github.com/quintilesims/layer0/common/startup"
)

const (
	SCALER_SLEEP_DURATION = time.Hour
)

func setupRestful(lgc logic.Logic) {
	adminLogic := logic.NewL0AdminLogic(lgc)
	deployLogic := logic.NewL0DeployLogic(lgc)
	environmentLogic := logic.NewL0EnvironmentLogic(lgc)
	healthLogic := logic.NewL0HealthLogic(lgc)
	loadBalancerLogic := logic.NewL0LoadBalancerLogic(lgc)
	serviceLogic := logic.NewL0ServiceLogic(lgc)
	taskLogic := logic.NewL0TaskLogic(lgc)
	jobLogic := logic.NewL0JobLogic(lgc, taskLogic, deployLogic)

	adminHandler := handlers.NewAdminHandler(adminLogic)
	deployHandler := handlers.NewDeployHandler(deployLogic)
	environmentHandler := handlers.NewEnvironmentHandler(environmentLogic, jobLogic)
	healthHandler := handlers.NewHealthHandler(healthLogic)
	jobHandler := handlers.NewJobHandler(jobLogic)
	loadBalancerHandler := handlers.NewLoadBalancerHandler(loadBalancerLogic, jobLogic)
	serviceHandler := handlers.NewServiceHandler(serviceLogic, jobLogic)
	tagHandler := handlers.NewTagHandler(lgc.TagStore)
	taskHandler := handlers.NewTaskHandler(taskLogic, jobLogic)

	restful.SetLogger(logutils.SilentLogger{})
	restful.Add(deployHandler.Routes())
	restful.Add(serviceHandler.Routes())
	restful.Add(environmentHandler.Routes())
	restful.Add(healthHandler.Routes())
	restful.Add(tagHandler.Routes())
	restful.Add(adminHandler.Routes())
	restful.Add(loadBalancerHandler.Routes())
	restful.Add(taskHandler.Routes())
	restful.Add(jobHandler.Routes())

	restful.Filter(handlers.LogRequest)
	restful.Filter(handlers.AddVersionHeader)
	restful.Filter(handlers.EnableCORS)
	restful.Filter(restful.OPTIONSFilter())
	restful.DefaultContainer.Filter(handlers.HttpsRedirect)

	config := swagger.Config{
		WebServices:     restful.RegisteredWebServices(),
		ApiPath:         "/apidocs.json",
		SwaggerPath:     swaggerPath,
		SwaggerFilePath: swaggerFilePath,
		StaticHandler:   new(SwaggerRedirectHandler),
	}

	swagger.InstallSwaggerService(config)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "custom 404")
			return
		}

		http.Redirect(w, r, config.SwaggerPath, 302)
	})
}

type SwaggerRedirectHandler struct{}

var swaggerPath = "/apidocs/"
var swaggerFilePath = "api/external/swagger-ui/dist"

func (*SwaggerRedirectHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	proto := req.Header.Get("X-Forwarded-Proto")
	if proto == "http" {
		url := fmt.Sprintf("https://%v%v", req.Host, req.URL)
		http.Redirect(writer, req, url, 301)
	} else {
		http.StripPrefix(swaggerPath, http.FileServer(http.Dir(swaggerFilePath))).ServeHTTP(writer, req)
	}
}

var Version string

func main() {
	if err := config.Validate(config.RequiredAPIVariables); err != nil {
		logrus.Fatal(err)
	}

	switch strings.ToLower(config.APILogLevel()) {
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

	logger := logutils.NewStackTraceLogger("Main")
	logutils.SetGlobalLogger(logger)

	if Version == "" {
		Version = "unset/developer"
	}

	config.SetAPIVersion(Version)
	logrus.Printf("l0-api %v", Version)

	port := ":" + config.APIPort()
	region := config.AWSRegion()
	credProvider := config.NewConfigCredProvider()

	backend, err := startup.GetBackend(credProvider, region)
	if err != nil {
		logrus.Fatal(err)
	}

	rateLimit, err := time.ParseDuration(config.AWSTimeBetweenRequests())
	if err != nil {
		logrus.Fatal(err)
	}

	provider.ResetRateLimiter(rateLimit)
	defer provider.StopRateLimiter()

	lgc, err := startup.GetLogic(backend)
	if err != nil {
		logrus.Fatal(err)
	}

	setupRestful(*lgc)

	taskLogic := logic.NewL0TaskLogic(*lgc)
	deployLogic := logic.NewL0DeployLogic(*lgc)
	jobLogic := logic.NewL0JobLogic(*lgc, taskLogic, deployLogic)
	environmentLogic := logic.NewL0EnvironmentLogic(*lgc)
	adminLogic := logic.NewL0AdminLogic(*lgc)

	if err := adminLogic.UpdateSQL(); err != nil {
		logrus.Errorf("Failed to update sql: %v", err)
	}

	Janitor := logic.NewJanitor(jobLogic, taskLogic, lgc.JobStore, lgc.TagStore)
	go runEnvironmentScaler(environmentLogic)

	logrus.Infof("Starting  Janitor")
	Janitor.Run()

	logrus.Print("Service on localhost" + port)
	logrus.Fatal(http.ListenAndServe(port, nil))
}

func runEnvironmentScaler(environmentLogic *logic.L0EnvironmentLogic) {
	logger := logutils.NewStandardLogger("AUTO Environment Scaler")

	for {
		environments, err := environmentLogic.ListEnvironments()
		if err != nil {
			logger.Errorf("Failed to list environments: %v", err)
			continue
		}

		for _, environment := range environments {
			logger.Infof("Scaling Environment %s", environment.EnvironmentID)

			if _, err := environmentLogic.Scaler.Scale(environment.EnvironmentID); err != nil {
				logger.Errorf("Failed to scale environment %s: %v", environment.EnvironmentID, err)
				continue
			}

			logger.Infof("Finished scaling environment %s", environment.EnvironmentID)
		}

		time.Sleep(SCALER_SLEEP_DURATION)
	}
}
