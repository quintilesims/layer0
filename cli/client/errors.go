package client

import (
	"fmt"

	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
)

func sslError(err error) error {
	text := fmt.Sprintf("Currently configured endpoint (%v) expects HTTPS. ", config.APIEndpoint())
	text += "You should register a proper domain with GSD. "
	text += "You can set LAYER0_SKIP_SSL_VERIFY=1 to ignore this in dev scenarios. "
	text += fmt.Sprintf("(err: %v)", err)
	return fmt.Errorf(text)
}

type ServerError models.ServerError

func (s *ServerError) Error() string {
	return s.Message
}

func (s *ServerError) ToCommonError() *errors.ServerError {
	return errors.Newf(errors.ErrorCode(s.ErrorCode), s.Message)
}
