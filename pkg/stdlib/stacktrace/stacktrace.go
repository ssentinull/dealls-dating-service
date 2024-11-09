package stacktrace

import (
	"fmt"
	"math"
	"net/http"
	"runtime"

	openApiError "github.com/go-openapi/errors"
)

type ErrorCode uint32

const NoCode ErrorCode = math.MaxUint32

type stacktrace struct {
	message  string
	cause    error
	code     ErrorCode
	file     string
	function string
	line     int
}

func Wrap(cause error, msg string, vals ...interface{}) error {
	if cause == nil {
		return nil
	}
	return create(cause, NoCode, msg, vals...)
}

func New(msg string, vals ...interface{}) error {
	return create(nil, NoCode, msg, vals...)
}

func NewErrorWithCode(code ErrorCode, msg string, vals ...interface{}) error {
	return create(nil, code, msg, vals...)
}

func WrapWithCode(cause error, code ErrorCode, msg string, vals ...interface{}) error {
	if cause == nil {
		return nil
	}
	return create(cause, code, msg, vals...)
}

func create(cause error, code ErrorCode, msg string, vals ...interface{}) error {
	if code == 0 {
		code = http.StatusInternalServerError
	}

	if _, ok := cause.(*openApiError.CompositeError); ok {
		return cause
	}

	err := &stacktrace{
		message: fmt.Sprintf(msg, vals...),
		cause:   cause,
		code:    code,
	}

	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return err
	}

	file = RemoveGoPath(file)
	err.file, err.line = file, line

	f := runtime.FuncForPC(pc)
	if f == nil {
		return err
	}

	err.function = shortFuncName(f)

	return err
}

func GetCode(err error) ErrorCode {
	if err, ok := err.(*stacktrace); ok {
		return err.code
	}
	return NoCode
}

func GetCause(err error) error {
	if err, ok := err.(*stacktrace); ok {
		return err.cause
	}
	return err
}

func GetRootCode(err error) ErrorCode {
	rootCode := NoCode
	if err, ok := err.(*stacktrace); ok {
		cause := err
		for {
			if cause.code != NoCode {
				rootCode = cause.code
				break
			}
			if cause.code == NoCode {
				if errCause, ok := err.cause.(*stacktrace); ok {
					cause = errCause
					continue
				}
				break
			}
		}
	}
	return rootCode
}

func GetMessage(err error) string {
	if err, ok := err.(*stacktrace); ok {
		return err.message
	}

	return ""
}

func (st *stacktrace) Error() string {
	if st != nil {
		return fmt.Sprint(st)
	}

	return ""
}

func ParseStatusCodeToError(statusCode int, serviceName string) error {
	errorMaps := map[int]error{
		400: NewErrorWithCode(http.StatusBadRequest, fmt.Sprintf("got bad request from %s", serviceName)),
		401: NewErrorWithCode(http.StatusUnauthorized, "unauthorized"),
		404: NewErrorWithCode(http.StatusNotFound, "not found"),
		500: NewErrorWithCode(http.StatusInternalServerError, fmt.Sprintf("%s unavailable", serviceName)),
	}

	if statusCode >= 200 && statusCode <= 299 {
		return nil
	}

	if err, ok := errorMaps[statusCode]; ok {
		return err
	}

	if statusCode >= 500 {
		return errorMaps[500]
	}

	return NewErrorWithCode(http.StatusInternalServerError, fmt.Sprintf("error when call %s", serviceName))
}

func ParseStatusCodeToErrorWithCause(cause error, statusCode int, serviceName string) error {
	errorMaps := map[int]error{
		400: WrapWithCode(cause, http.StatusBadRequest, fmt.Sprintf("got bad request from %s", serviceName)),
		401: WrapWithCode(cause, http.StatusUnauthorized, "unauthorized"),
		404: WrapWithCode(cause, http.StatusNotFound, "not found"),
		500: WrapWithCode(cause, http.StatusInternalServerError, fmt.Sprintf("%s unavailable", serviceName)),
	}

	if statusCode >= 200 && statusCode <= 299 {
		return nil
	}

	if err, ok := errorMaps[statusCode]; ok {
		return err
	}

	if statusCode >= 500 {
		return errorMaps[500]
	}

	return WrapWithCode(cause, http.StatusInternalServerError, fmt.Sprintf("error when call %s", serviceName))
}

func IsStackTrace(err error) bool {
	if _, ok := err.(*stacktrace); ok {
		return true
	}
	return false
}

func WrapPassCode(cause error, msg string, vals ...interface{}) error {
	if cause == nil {
		return nil
	}

	return create(cause, GetCode(cause), msg, vals...)
}
