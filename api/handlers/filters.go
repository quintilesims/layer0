package handlers

import (
	"github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
	"github.com/quintilesims/layer0/common/config"
	"time"
)

func LogRequest(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	start := time.Now()
	chain.ProcessFilter(req, resp)
	duration := time.Since(start)
	
	if req.Request.URL != "/admin/health" {
		logrus.Infof("request %s %s (%v) %v", req.Request.Method, req.Request.URL, resp.StatusCode(), duration)
	}
}

func AddVersionHeader(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	resp.AddHeader("Version", config.APIVersion())
	chain.ProcessFilter(req, resp)
}

func EnableCORS(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	if origin := req.Request.Header.Get("Origin"); origin != "" {
		resp.AddHeader("Access-Control-Allow-Origin", origin)
		resp.AddHeader("Access-Control-Allow-Credentials", "true")
		resp.AddHeader("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token, Authorization")
	}

	chain.ProcessFilter(req, resp)
}
