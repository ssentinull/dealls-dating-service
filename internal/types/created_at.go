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

// CreatedAt created at
//
// swagger:model createdAt
type CreatedAt strfmt.DateTime

// UnmarshalJSON sets a CreatedAt value from JSON input
func (m *CreatedAt) UnmarshalJSON(b []byte) error {
	return ((*strfmt.DateTime)(m)).UnmarshalJSON(b)
}

// MarshalJSON retrieves a CreatedAt value as JSON output
func (m CreatedAt) MarshalJSON() ([]byte, error) {
	return (strfmt.DateTime(m)).MarshalJSON()
}

// Validate validates this created at
func (m CreatedAt) Validate(formats strfmt.Registry) error {
	var res []error

	if err := validate.FormatOf("", "body", "date-time", strfmt.DateTime(m).String(), formats); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this created at based on context it is used
func (m CreatedAt) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *CreatedAt) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CreatedAt) UnmarshalBinary(b []byte) error {
	var res CreatedAt
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
