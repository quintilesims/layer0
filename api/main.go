package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/common/config"
)

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

}
