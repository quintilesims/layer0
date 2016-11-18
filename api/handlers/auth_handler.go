package handlers

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"gitlab.imshealth.com/xfra/layer0/common/config"
)

func basicAuthenticate(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	encoded := req.Request.Header.Get("Authorization")
	expected := "Basic " + config.APIAuthToken()

	// a better implementation would connect to an external service
	if len(encoded) == 0 || expected != encoded {
		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		resp.WriteErrorString(401, "401: Not Authorized")
		return
	}

	chain.ProcessFilter(req, resp)
}

func HttpsRedirect(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	proto := req.Request.Header.Get("X-Forwarded-Proto")
	if proto == "http" {
		url := fmt.Sprintf("https://%v%v", req.Request.Host, req.Request.URL)
		resp.AddHeader("Location", url)
		resp.WriteHeader(301)
		resp.WriteAsJson(``)
		return
	}

	chain.ProcessFilter(req, resp)
}
