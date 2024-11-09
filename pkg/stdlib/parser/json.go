package parser

import (
	"github.com/ssentinull/dealls-dating-service/pkg/stdlib/logger"

	jsoniter "github.com/json-iterator/go"
)

type jsonConfig string

const (
	JSONConfigDefault                  jsonConfig = `default`
	JSONConfigCompatibleWithStdLibrary jsonConfig = `standard`
	JSONConfigFastest                  jsonConfig = `fastest`
	JSONConfigCustom                   jsonConfig = `custom`
)

type JSONParser interface {
	Marshal(orig interface{}) ([]byte, error)
	Unmarshal(blob []byte, dest interface{}) error
}

type jsonParser struct {
	efLogger logger.Logger
	API      jsoniter.API
	opt      JSONOptions
}

type JSONOptions struct {
	Config                        jsonConfig
	IndentionStep                 int
	MarshalFloatWith6Digits       bool
	EscapeHTML                    bool
	SortMapKeys                   bool
	UseNumber                     bool
	DisallowUnknownFields         bool
	TagKey                        string
	OnlyTaggedField               bool
	ValidateJSONRawMessage        bool
	ObjectFieldMustBeSimpleString bool
	CaseSensitive                 bool
}

func initJSONP(efLogger logger.Logger, opt JSONOptions) JSONParser {
	var jsonAPI jsoniter.API
	switch opt.Config {

	case JSONConfigDefault:
		jsonAPI = jsoniter.ConfigDefault

	case JSONConfigFastest:
		jsonAPI = jsoniter.ConfigFastest

	case JSONConfigCompatibleWithStdLibrary:
		jsonAPI = jsoniter.ConfigCompatibleWithStandardLibrary

	case JSONConfigCustom:
		jsonAPI = jsoniter.Config{
			IndentionStep:                 opt.IndentionStep,
			MarshalFloatWith6Digits:       opt.MarshalFloatWith6Digits,
			EscapeHTML:                    opt.EscapeHTML,
			SortMapKeys:                   opt.SortMapKeys,
			UseNumber:                     opt.UseNumber,
			DisallowUnknownFields:         opt.DisallowUnknownFields,
			TagKey:                        opt.TagKey,
			OnlyTaggedField:               opt.OnlyTaggedField,
			ValidateJsonRawMessage:        opt.ValidateJSONRawMessage,
			ObjectFieldMustBeSimpleString: opt.ObjectFieldMustBeSimpleString,
			CaseSensitive:                 opt.CaseSensitive,
		}.Froze()

	default:
		jsonAPI = jsoniter.ConfigCompatibleWithStandardLibrary
	}
	p := &jsonParser{
		efLogger: efLogger,
		API:      jsonAPI,
		opt:      opt,
	}

	return p
}

func (p *jsonParser) Marshal(orig interface{}) ([]byte, error) {
	stream := p.API.BorrowStream(nil)
	defer p.API.ReturnStream(stream)
	stream.WriteVal(orig)
	result := make([]byte, stream.Buffered())
	if stream.Error != nil {
		return nil, stream.Error
	}
	copy(result, stream.Buffer())
	return result, nil
}

func (p *jsonParser) Unmarshal(blob []byte, dest interface{}) error {
	iter := p.API.BorrowIterator(blob)
	defer p.API.ReturnIterator(iter)
	iter.ReadVal(dest)
	if iter.Error != nil {
		return iter.Error
	}
	return nil
}
