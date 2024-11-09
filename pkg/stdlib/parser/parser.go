package parser

import (
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"

	goValidator "github.com/go-playground/validator/v10"
)

type Parser interface {
	JSONParser() JSONParser
	Validator() Validator
}

type parser struct {
	json      JSONParser
	validator Validator
	opt       Options
}

type Options struct {
	JSON      JSONOptions
	Validator ValidatorOptions
}

func Init(efLogger logger.Logger, bindingValidator *goValidator.Validate, opt Options) Parser {
	return &parser{
		json:      initJSONP(efLogger, opt.JSON),
		validator: initValidator(opt.Validator, bindingValidator),
		opt:       opt,
	}
}

func (p *parser) JSONParser() JSONParser {
	return p.json
}

func (p *parser) Validator() Validator {
	return p.validator
}
