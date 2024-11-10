// Code generated by go-swagger; DO NOT EDIT.

package feed

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NewGetFeedParams creates a new GetFeedParams object
// with the default values initialized.
func NewGetFeedParams() GetFeedParams {

	var (
		// initialize parameters with default values

		paginationPageDefault = int64(1)
		paginationSizeDefault = int64(10)
	)

	return GetFeedParams{
		PaginationPage: &paginationPageDefault,

		PaginationSize: &paginationSizeDefault,
	}
}

// GetFeedParams contains all the bound params for the get feed operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetFeed
type GetFeedParams struct {
	/*X Client Id
	  Required: true
	  In: header
	*/
	XClientID string `header:"X-Client-Id"`
	/*Pagination page
	  Minimum: 1
	  In: query
	  Default: 1
	*/
	PaginationPage *int64 `query:"page"`
	/*Pagination size
	  Minimum: 1
	  In: query
	  Default: 10
	*/
	PaginationSize *int64 `query:"size"`
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetFeedParams() beforehand.
func (o *GetFeedParams) BindRequest(c *gin.Context) error {
	var res []error

	var r *http.Request

	r = c.Request

	qs := runtime.Values(r.URL.Query())

	if err := o.bindXClientID(r.Header[http.CanonicalHeaderKey("X-Client-Id")], true, strfmt.Default); err != nil {
		res = append(res, err)
	}

	qPage, qhkPage, _ := qs.GetOK("page")
	if err := o.bindPaginationPage(qPage, qhkPage, strfmt.Default); err != nil {
		res = append(res, err)
	}

	qSize, qhkSize, _ := qs.GetOK("size")
	if err := o.bindPaginationSize(qSize, qhkSize, strfmt.Default); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetFeedParams) Validate(formats strfmt.Registry) error {
	var res []error

	// X-Client-Id
	// Required: true

	if err := validate.Required("X-Client-Id", "header", o.XClientID); err != nil {
		res = append(res, err)
	}

	// page
	// Required: false
	// AllowEmptyValue: false

	if err := o.validatePaginationPage(formats); err != nil {
		res = append(res, err)
	}

	// size
	// Required: false
	// AllowEmptyValue: false

	if err := o.validatePaginationSize(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindXClientID binds and validates parameter XClientID from header.
func (o *GetFeedParams) bindXClientID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("X-Client-Id", "header", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("X-Client-Id", "header", raw); err != nil {
		return err
	}

	o.XClientID = raw

	return nil
}

// bindPaginationPage binds and validates parameter PaginationPage from query.
func (o *GetFeedParams) bindPaginationPage(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewGetFeedParams()
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("page", "query", "int64", raw)
	}
	o.PaginationPage = &value

	if err := o.validatePaginationPage(formats); err != nil {
		return err
	}

	return nil
}

// validatePaginationPage carries on validations for parameter PaginationPage
func (o *GetFeedParams) validatePaginationPage(formats strfmt.Registry) error {

	// Required: false
	if o.PaginationPage == nil {
		return nil
	}

	if err := validate.MinimumInt("page", "query", *o.PaginationPage, 1, false); err != nil {
		return err
	}

	return nil
}

// bindPaginationSize binds and validates parameter PaginationSize from query.
func (o *GetFeedParams) bindPaginationSize(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewGetFeedParams()
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("size", "query", "int64", raw)
	}
	o.PaginationSize = &value

	if err := o.validatePaginationSize(formats); err != nil {
		return err
	}

	return nil
}

// validatePaginationSize carries on validations for parameter PaginationSize
func (o *GetFeedParams) validatePaginationSize(formats strfmt.Registry) error {

	// Required: false
	if o.PaginationSize == nil {
		return nil
	}

	if err := validate.MinimumInt("size", "query", *o.PaginationSize, 1, false); err != nil {
		return err
	}

	return nil
}
