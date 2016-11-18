package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-swagger/go-swagger/strfmt"

	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit/validate"
)

/*ServiceLog service log

swagger:model ServiceLog
*/
type ServiceLog struct {

	/* message
	 */
	Message *string `json:"Message,omitempty"`

	/* log time

	Required: true
	*/
	LogTime strfmt.DateTime `json:"log_time"`
}

// Validate validates this service log
func (m *ServiceLog) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLogTime(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ServiceLog) validateLogTime(formats strfmt.Registry) error {

	if err := validate.Required("log_time", "body", strfmt.DateTime(m.LogTime)); err != nil {
		return err
	}

	return nil
}
