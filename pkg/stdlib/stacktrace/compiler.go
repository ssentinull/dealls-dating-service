package stacktrace

import (
	"net/http"
)

type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	DebugError *string   `json:"debug,omitempty"`
	sys        error
}

func (e *AppError) Error() string {
	return e.sys.Error()
}

func Compile(err error, debugMode bool) (int, AppError) {
	// Developer Debug Error
	var debugErr *string
	if debugMode {
		errStr := err.Error()
		if len(errStr) > 0 {
			debugErr = &errStr
		}
	}

	// Get Error Code
	code := GetCode(err)

	// Get Common Error
	if errMessage, ok := ErrorMessages[code]; ok {
		msg := errMessage.Message
		return errMessage.StatusCode, AppError{
			Code:       code,
			Message:    msg,
			sys:        err,
			DebugError: debugErr,
		}
	}

	// Set Default Error
	return http.StatusInternalServerError, AppError{
		Code:       code,
		Message:    "service error not defined!",
		sys:        err,
		DebugError: debugErr,
	}
}
