package handlers

import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
)

func WriteJobResponse(response *restful.Response, jobID string) {
	response.AddHeader("Location", fmt.Sprintf("/job/%s", jobID))
	response.AddHeader("X-JobID", jobID)
	response.WriteHeader(http.StatusAccepted)
	response.WriteAsJson(``)
}
