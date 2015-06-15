package api

import (
	"errors"
	"fmt"
)

const (
	MaxErrorCode = 9999
)

var (
	InternalErrorCode = GenerateErrorCode("context", 6)
)

var packageErrorPrefixes = map[string]rune{}

// Error has error information occurred in Hawk.
// It has error messages for clients(users).
type Error struct {
	originalCode *ErrorCode `json:"-"`

	// Code has a error code of this error. It always have
	// a single alphabet prefix and 4-digit error number.
	Code string `json:"code"`

	// Message contains a user-friendly error message for users
	// of WebUI or clients of API. It MUST NOT have internal
	// error information nor debug information which are
	// completely useless for them.
	Message string `json:"message"`

	// RequestID has an ID of the request which caused this error.
	// This value is convenient for users or clients of Hawk
	// to report their problems to the development team.
	// The team can look up incidents easily by greping logs
	// with the ID.
	//
	// The type of RequestID is string because it can be
	// very large and JavaScript might not be able to handle
	// it as an integer.
	RequestID string `json:"request_id"`

	// Status has an appropriate HTTP status code for this error.
	Status int `json:"-"`

	// Err contains internal error information returned from
	// a method or a function which caused the error.
	Err error `json:"-"`

	// Meta contains arbitrary information of the error.
	// What the information means depends on each error.
	// For instance, a validation error might contain
	// error information of each field.
	Meta map[string]interface{} `json:"meta"`
}

// ErrorCode has an error code of an error defined in each package.
type ErrorCode struct {
	// Package is the name of the module which caused this error.
	Package string `json:"package"`

	// Code has an error code of this error.
	Code int `json:"code"`
}

// String returns a string representation of the error code.
// The string contains one letter package code and 4-digit
// error code (e.g. "E0105").
func (e *ErrorCode) String() string {
	p, ok := packageErrorPrefixes[e.Package]
	if !ok {
		p = 'X'
	}
	return fmt.Sprintf("%c%04d", p, e.Code)
}

func (e *Error) SetRequestID(id uint64) {
	e.RequestID = fmt.Sprint(id)
}

func NewError(code *ErrorCode, msg string, status int, err error) *Error {
	if err == nil {
		err = errors.New(msg)
	}

	return &Error{
		originalCode: code,
		Code:         code.String(),
		Message:      msg,
		Status:       status,
		Err:          err,
		Meta:         make(map[string]interface{}),
	}
}

// GenerateErrorCode creates a new error code.
func GenerateErrorCode(pkg string, code int) *ErrorCode {
	if code <= 0 || MaxErrorCode < code {
		// Since this function is only executed on initialization,
		// it can safely panic.
		panic(fmt.Errorf("An error code must be in [1, %v]: %v", MaxErrorCode, code))
	}
	return &ErrorCode{
		Package: pkg,
		Code:    code,
	}
}
