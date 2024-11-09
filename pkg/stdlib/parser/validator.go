package parser

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	goValidator "github.com/go-playground/validator/v10"
)

type Validator interface {
	BindAndValidateBody(c *gin.Context, val runtime.Validatable) error
	BindOnly(c *gin.Context, val any) error
	ValidateByBindingTag(c context.Context, val any) error
}

type Bindable interface {
	BindRequest(c *gin.Context) error
}

type validator struct {
	opt              ValidatorOptions
	bindingValidator *goValidator.Validate
}

type ValidatorOptions struct {
}

func initValidator(opt ValidatorOptions, bindingValidator *goValidator.Validate) Validator {
	return &validator{
		opt:              opt,
		bindingValidator: bindingValidator,
	}
}

func (v *validator) BindAndValidateBody(c *gin.Context, val runtime.Validatable) error {
	var unformatedErr error
	if vBind, ok := val.(Bindable); ok {
		unformatedErr = vBind.BindRequest(c)
	} else {
		err := c.Bind(&v)
		if err != nil {
			return err
		}

		unformatedErr = val.Validate(strfmt.Default)
	}

	return unformatedErr
}

func (v *validator) BindOnly(c *gin.Context, val any) error {
	binding := binding.Default(c.Request.Method, c.ContentType())
	err := c.ShouldBindWith(val, binding)

	switch err.(type) {
	case goValidator.ValidationErrors, goValidator.FieldError:
		return nil
	}

	return err
}

func (v *validator) ValidateByBindingTag(c context.Context, val any) error {

	err := v.bindingValidator.Struct(val)
	if err == nil {
		return nil
	}
	err = errors.NewParseError("body", "body", "", err)

	return errors.CompositeValidationError(err)
}
