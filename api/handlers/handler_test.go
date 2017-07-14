package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/testutils"
)

// Main test entrypoint
func TestMain(m *testing.M) {
	config.SetTestConfig()
	log.SetLevel(log.FatalLevel)
	retCode := m.Run()
	os.Exit(retCode)
}

type TestRequest struct {
	Body       interface{}
	Path       string
	Parameters map[string]string
	Query      string
}

func (this *TestRequest) RestfulRequest() (*restful.Request, error) {
	jsonBytes, err := json.Marshal(this.Body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("", this.Path, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	restfulRequest := restful.NewRequest(req)
	for key, val := range this.Parameters {
		restfulRequest.PathParameters()[key] = val
	}

	restfulRequest.Request.URL.RawQuery = this.Query

	return restfulRequest, nil
}

type HandlerTestCase struct {
	Name    string
	Request *TestRequest
	Setup   func(*gomock.Controller) interface{}
	Run     func(*testutils.Reporter, interface{}, *restful.Request, *restful.Response, Readf)
}

type Readf func(interface{})

func RunHandlerTestCase(t *testing.T, testCase HandlerTestCase) {
	reporter := testutils.NewReporter(t, testCase.Name)
	ctrl := gomock.NewController(reporter)
	defer ctrl.Finish()

	recorder := httptest.NewRecorder()
	response := restful.NewResponse(recorder)
	request, err := testCase.Request.RestfulRequest()
	if err != nil {
		reporter.Fatal(err)
	}

	var target interface{}
	if testCase.Setup != nil {
		target = testCase.Setup(ctrl)
	}

	read := func(response interface{}) {
		if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
			reporter.Fatal(err)
		}
	}

	testCase.Run(reporter, target, request, response, read)
}

func RunHandlerTestCases(t *testing.T, testCases []HandlerTestCase) {
	for _, testCase := range testCases {
		RunHandlerTestCase(t, testCase)
	}
}
