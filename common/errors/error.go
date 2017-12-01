package errors

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/quintilesims/layer0/common/models"
)

type ServerError struct {
	Code ErrorCode
	Err  error
}

func Newf(code ErrorCode, format string, tokens ...interface{}) *ServerError {
	return New(code, fmt.Errorf(format, tokens...))
}

func New(code ErrorCode, err error) *ServerError {
	return &ServerError{
		Code: code,
		Err:  err,
	}
}

func FromModel(se models.ServerError) *ServerError {
	return New(ErrorCode(se.ErrorCode), fmt.Errorf(se.Message))
}

func (s *ServerError) Error() string {
	return fmt.Sprintf("ServerError (code='%s') %s", s.Code, s.Err.Error())
}

func (s *ServerError) Model() models.ServerError {
	return models.ServerError{
		ErrorCode: s.Code.String(),
		Message:   s.Err.Error(),
	}
}

func (s *ServerError) Write(w http.ResponseWriter, r *http.Request) {
	switch s.Code {
	case InvalidRequest,
		DeployDoesNotExist,
		EnvironmentDoesNotExist,
		JobDoesNotExist,
		LoadBalancerDoesNotExist,
		ServiceDoesNotExist,
		TaskDoesNotExist:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	body, err := json.Marshal(s.Model())
	if err != nil {
		log.Printf("[ERROR] Failed to marshal server error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(body)
}
