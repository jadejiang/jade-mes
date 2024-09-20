package ecode

import (
	"fmt"
	"strings"
)

// ECode ...
type ECode struct {
	code    int64
	message string
	errors  []error
}

// Error ...
func (e ECode) Error() string {
	errorMessage := e.message

	var subErrorMessages []string
	for _, err := range e.errors {
		subErrorMessages = append(subErrorMessages, err.Error())
	}

	if len(subErrorMessages) > 0 {
		errorMessage = fmt.Sprintf("%s: [%s]", errorMessage, strings.Join(subErrorMessages, " | "))
	}

	return errorMessage
}

// Errors ...
func (e ECode) Errors() []interface{} {
	errors := []interface{}{}

	for _, err := range e.errors {
		errors = append(errors, err.Error())
	}

	return errors
}

// Code ...
func (e ECode) Code() int64 {
	return e.code
}

// Message ...
func (e ECode) Message() string {
	return e.message
}

// Err ...
func (e ECode) Err(err error) error {
	e.errors = []error{err}
	return e
}

func newError(code int64, message string, errors ...error) ECode {
	err := ECode{
		code:    code,
		message: message,
		errors:  errors,
	}

	return err
}
