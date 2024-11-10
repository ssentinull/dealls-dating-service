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

// CreateFeedPreferenceRequest Create Feed Preference Request
//
// swagger:model createFeedPreferenceRequest
type CreateFeedPreferenceRequest struct {

	// gender
	Gender Gender `json:"gender,omitempty"`

	// location
	Location Location `json:"location,omitempty"`

	// the Maximum Age
	MaxAge int64 `json:"max_age" binding:"required"`

	// the Minimum Age
	MinAge int64 `json:"min_age" binding:"required"`
}

// Validate validates this create feed preference request
func (m *CreateFeedPreferenceRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateGender(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLocation(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CreateFeedPreferenceRequest) validateGender(formats strfmt.Registry) error {
	if swag.IsZero(m.Gender) { // not required
		return nil
	}

	if err := m.Gender.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("gender")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("gender")
		}
		return err
	}

	return nil
}

func (m *CreateFeedPreferenceRequest) validateLocation(formats strfmt.Registry) error {
	if swag.IsZero(m.Location) { // not required
		return nil
	}

	if err := m.Location.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("location")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("location")
		}
		return err
	}

	return nil
}

// ContextValidate validate this create feed preference request based on the context it is used
func (m *CreateFeedPreferenceRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateGender(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateLocation(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CreateFeedPreferenceRequest) contextValidateGender(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.Gender) { // not required
		return nil
	}

	if err := m.Gender.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("gender")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("gender")
		}
		return err
	}

	return nil
}

func (m *CreateFeedPreferenceRequest) contextValidateLocation(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.Location) { // not required
		return nil
	}

	if err := m.Location.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("location")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("location")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *CreateFeedPreferenceRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CreateFeedPreferenceRequest) UnmarshalBinary(b []byte) error {
	var res CreateFeedPreferenceRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
