package handlers

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"net/http"
)

func WriteJobResponse(response *restful.Response, jobID string) {
	response.AddHeader("Location", fmt.Sprintf("/job/%s", jobID))
	response.AddHeader("X-JobID", jobID)
	response.WriteHeader(http.StatusAccepted)
	response.WriteAsJson(``)
}
