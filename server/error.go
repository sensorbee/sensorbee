package server

import (
	"fmt"
	"net/http"
)

const (
	requestURLNotFoundErrorCode = "E0001"
	internalServerErrorCode     = "E0002"
	requestBodyReadErrorCode    = "E0003"
	requestBodyParseErrorCode   = "E0004"

	// formValidationErrorCode means that validation of request body failed.
	// When this error happens, Error.Meta should have detailed error messages
	// for each field. Each field must have a slice of strings so that clients
	// can always write error handling codes assuming that they're arrays.
	formValidationErrorCode = "E0005"

	// bqlStmtParseErrorCode is returned when a statement cannot be parsed.
	// When this error happens, Error.Meta should have parse error messages
	// in Meta["parse_errors"] as an array of strings and the statement which
	// couldn't be parsed in Meta["statement"].
	bqlStmtParseErrorCode = "E0006"

	// bqlStmtProcessingErrorCode is returned when a statement cannot be
	// processed successfully. When this error happens, Error.Meta should have
	// an error message in Meta["error"] and statement in Meta["statement"].
	bqlStmtProcessingErrorCode = "E0007"
)

const (
	somethingWentWrong = "Something went wrong. Please try again later."
)

// NewInternalServerError creates an internal server error response having
// an error information.
func NewInternalServerError(err error) *Error {
	return NewError(internalServerErrorCode, somethingWentWrong,
		http.StatusInternalServerError, err)
}

// Error has the error information reported to clients of API.
type Error struct {
	// Code has an error code of this error. It always have
	// a single alphabet prefix and 4-digit error number.
	Code string `json:"code"`

	// Message contains a user-friendly error message for users
	// of WebUI or clients of API. It MUST NOT have internal
	// error information nor debug information which are
	// completely useless for them.
	Message string `json:"message"`

	// RequestID has an ID of the request which caused this error.
	// This value is convenient for users or clients to report their
	// problems to the development team. The team can look up incidents
	// easily by greping logs with the ID.
	//
	// The type of RequestID is string because it can be
	// very large and JavaScript might not be able to handle
	// it as an integer.
	RequestID string `json:"request_id"`

	// Status has an appropriate HTTP status code for this error.
	Status int `json:"-"`

	// Err is an internal error information which must not be shown to users
	// directly.
	Err error `json:"-"`

	// Meta contains arbitrary information of the error.
	// What the information means depends on each error.
	// For instance, a validation error might contain
	// error information of each field.
	Meta map[string]interface{} `json:"meta"`
}

// NewError creates a new Error instance.
func NewError(code string, msg string, status int, err error) *Error {
	return &Error{
		Code:    code,
		Message: msg,
		Status:  status,
		Err:     err,
		Meta:    map[string]interface{}{},
	}
}

// SetRequestID set the ID of the current request to Error.
func (e *Error) SetRequestID(id uint64) {
	e.RequestID = fmt.Sprint(id)
}
