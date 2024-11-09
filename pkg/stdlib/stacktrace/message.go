package stacktrace

import "net/http"

type Message struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type errorMessage map[ErrorCode]Message

var (
	ErrMsgISE = Message{
		StatusCode: http.StatusInternalServerError,
		Message:    `Internal Server Error. Please Call Administrator.`,
	}

	ErrMsgNotFound = Message{
		StatusCode: http.StatusNotFound,
		Message:    `Record Does Not Exist. Please Validate Your Input Or Contact Administrator.`,
	}

	ErrMsgBadRequest = Message{
		StatusCode: http.StatusBadRequest,
		Message:    `Invalid Input. Please Validate Your Input.`,
	}

	ErrMsgUnauthorized = Message{
		StatusCode: http.StatusUnauthorized,
		Message:    `Unauthorized Access. You are not authorized to access this resource.`,
	}

	ErrMsgUniqueConst = Message{
		StatusCode: http.StatusConflict,
		Message:    `Record Has Existed and Must Be Unique. Please Validate Your Input Or Contact Administrator.`,
	}

	errMsgLocked = Message{
		StatusCode: http.StatusLocked,
		Message:    `The requested resource cannot be accessed at this time. Please try again later.`,
	}

	ErrorMessages = errorMessage{
		http.StatusInternalServerError: ErrMsgISE,
		http.StatusNotFound:            ErrMsgNotFound,
		http.StatusBadRequest:          ErrMsgBadRequest,
		http.StatusUnauthorized:        ErrMsgUnauthorized,
		http.StatusConflict:            ErrMsgUniqueConst,
		http.StatusLocked:              errMsgLocked,
	}
)
