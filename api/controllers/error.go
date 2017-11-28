package controllers

import (
	"log"
	"net/http"

	"github.com/quintilesims/layer0/common/errors"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	serverError, ok := err.(*errors.ServerError)
	if !ok {
		serverError = errors.New(errors.UnexpectedError, err)
	}

	log.Printf("[DEBUG] an api error occured: %v", serverError)
	serverError.Write(w, r)
}
