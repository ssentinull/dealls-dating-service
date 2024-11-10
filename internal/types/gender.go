// Code generated by go-swagger; DO NOT EDIT.

package types

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// Gender User Gender
//
// swagger:model gender
type Gender string

func NewGender(value Gender) *Gender {
	return &value
}

// Pointer returns a pointer to a freshly-allocated Gender.
func (m Gender) Pointer() *Gender {
	return &m
}

const (

	// GenderMALE captures enum value "MALE"
	GenderMALE Gender = "MALE"

	// GenderFEMALE captures enum value "FEMALE"
	GenderFEMALE Gender = "FEMALE"
)

// for schema
var genderEnum []interface{}

func init() {
	var res []Gender
	if err := json.Unmarshal([]byte(`["MALE","FEMALE"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		genderEnum = append(genderEnum, v)
	}
}

func (m Gender) validateGenderEnum(path, location string, value Gender) error {
	if err := validate.EnumCase(path, location, value, genderEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this gender
func (m Gender) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateGenderEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this gender based on context it is used
func (m Gender) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}