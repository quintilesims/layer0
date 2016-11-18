package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
	"gitlab.imshealth.com/xfra/layer0/api/backend/ecs"
	"gitlab.imshealth.com/xfra/layer0/api/logic"
	"gitlab.imshealth.com/xfra/layer0/common/config"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Main test entrypoint
func TestMain(m *testing.M) {
	config.SetTestConfig()
	log.SetLevel(log.FatalLevel)
	retCode := m.Run()
	os.Exit(retCode)
}

func TestAPIDocs(t *testing.T) {
	logic := logic.NewLogic(nil, nil, nil, &ecsbackend.ECSBackend{})
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
