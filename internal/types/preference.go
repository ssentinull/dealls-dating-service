// Code generated by go-swagger; DO NOT EDIT.

package types

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Preference Preference
//
// swagger:model preference
type Preference struct {

	// the Id
	ID int64 `json:"id" binding:"required"`

	// the User Id
	UserID int64 `json:"user_id" binding:"required"`

	// the Gender
	Gender string `json:"gender" binding:"required"`

	// the Minimum Age
	MinAge int64 `json:"min_age" binding:"required"`

	// the Maximum Age
	MaxAge int64 `json:"max_age" binding:"required"`

	// the Location
	Location string `json:"location" binding:"required"`

	// created at
	// Format: date-time
	CreatedAt CreatedAt `json:"created_at,omitempty"`
}

// Validate validates this preference
func (m *Preference) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Preference) validateCreatedAt(formats strfmt.Registry) error {
	if swag.IsZero(m.CreatedAt) { // not required
		return nil
	}

	if err := m.CreatedAt.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("created_at")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("created_at")
		}
		return err
	}

	return nil
}

// ContextValidate validate this preference based on the context it is used
func (m *Preference) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateCreatedAt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Preference) contextValidateCreatedAt(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.CreatedAt) { // not required
		return nil
	}

	if err := m.CreatedAt.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("created_at")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("created_at")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Preference) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Preference) UnmarshalBinary(b []byte) error {
	var res Preference
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
