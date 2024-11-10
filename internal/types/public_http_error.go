// Code generated by go-swagger; DO NOT EDIT.

package types

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PublicHTTPError public Http error
//
// swagger:model publicHttpError
type PublicHTTPError struct {

	// success
	// Example: true
	// Required: true
	Success bool `json:"success"`

	// Short, human-readable description of the error
	// Example: User is lacking permission to access this resource
	Message string `json:"message,omitempty"`

	// data
	Data *PublicHTTPErrorData `json:"data,omitempty"`

	// reference id to trace
	// Example: Ube4ab48e26e7d21c13dcbf07f8cebc0a
	TraceID string `json:"trace_id,omitempty"`
}

// Validate validates this public Http error
func (m *PublicHTTPError) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSuccess(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PublicHTTPError) validateSuccess(formats strfmt.Registry) error {

	if err := validate.Required("success", "body", bool(m.Success)); err != nil {
		return err
	}

	return nil
}

func (m *PublicHTTPError) validateData(formats strfmt.Registry) error {
	if swag.IsZero(m.Data) { // not required
		return nil
	}

	if m.Data != nil {
		if err := m.Data.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("data")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("data")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this public Http error based on the context it is used
func (m *PublicHTTPError) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateData(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PublicHTTPError) contextValidateData(ctx context.Context, formats strfmt.Registry) error {

	if m.Data != nil {

		if swag.IsZero(m.Data) { // not required
			return nil
		}

		if err := m.Data.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("data")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("data")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PublicHTTPError) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PublicHTTPError) UnmarshalBinary(b []byte) error {
	var res PublicHTTPError
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// PublicHTTPErrorData public HTTP error data
//
// swagger:model PublicHTTPErrorData
type PublicHTTPErrorData struct {

	// HTTP status code returned for the error
	// Example: 403
	// Maximum: 599
	// Minimum: 100
	Code int64 `json:"status,omitempty"`

	// More detailed, human-readable, optional explanation of the error
	// Example: Forbidden
	Detail string `json:"detail,omitempty"`

	// Type of error returned, should be used for client-side error handling
	// Example: generic
	Type string `json:"type,omitempty"`
}

// Validate validates this public HTTP error data
func (m *PublicHTTPErrorData) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCode(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PublicHTTPErrorData) validateCode(formats strfmt.Registry) error {
	if swag.IsZero(m.Code) { // not required
		return nil
	}

	if err := validate.MinimumInt("data"+"."+"status", "body", m.Code, 100, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("data"+"."+"status", "body", m.Code, 599, false); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this public HTTP error data based on context it is used
func (m *PublicHTTPErrorData) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *PublicHTTPErrorData) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PublicHTTPErrorData) UnmarshalBinary(b []byte) error {
	var res PublicHTTPErrorData
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}