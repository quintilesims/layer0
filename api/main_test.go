package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
	"github.com/quintilesims/layer0/api/backend/ecs"
	"github.com/quintilesims/layer0/api/logic"
	"github.com/quintilesims/layer0/common/config"
)

// Main test entrypoint
func TestMain(m *testing.M) {
	config.SetTestConfig()
	log.SetLevel(log.FatalLevel)
	retCode := m.Run()
	os.Exit(retCode)
}

func TestAPIDocs(t *testing.T) {
	logic := logic.NewLogic(nil, nil, &ecsbackend.ECSBackend{}, nil)
	setupRestful(*logic)

	httpRequest, _ := http.NewRequest("GET", "/apidocs.json", nil)
	httpWriter := httptest.NewRecorder()

	restful.DefaultContainer.ServeHTTP(httpWriter, httpRequest)

	if httpWriter.Code != 200 {
		t.Errorf("Expected Return Code 200 from apidocs.json")
	}

	body := httpWriter.Body.String()
	expected := []string{"environment", "deploy", "service"}
	for _, e := range expected {
		if !strings.Contains(body, e) {
			t.Errorf("Apidoc should list path %s", e)
		}
	}
}
