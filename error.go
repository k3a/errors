// Package errors is a drop-in replacement to standard errors, extending it with additional functionality
package errors

import "fmt"

// baseError is a trivial implementation of error.
type baseError struct {
	text string
}

func (e *baseError) Error() string {
	return e.text
}

// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func New(text string) error {
	return &baseError{text}
}

type wrapError struct {
	internalErr
	msg string
}

func (e *wrapError) Error() string {
	return e.msg
}

func (e *wrapError) Unwrap() error {
	return e.error
}

// Annotate is used to add extra context to an existing error (inspired by juju/errors)
func Annotate(err error, msg string) error {
	if err == nil {
		return nil
	}

	//TODO: tracing not implemented yet

	return &wrapError{
		internalErr: internalErr{err},
		msg:         msg + ": " + err.Error(),
	}
}

// Annotatef is used to add extra context to an existing error (inspired by juju/errors)
func Annotatef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	//TODO: tracing not implemented yet

	return &wrapError{
		internalErr: internalErr{err},
		msg:         fmt.Sprintf(format, args...) + ": " + err.Error(),
	}
}

// Trace adds the location of the Trace call to the stack (inspired by juju/errors)
func Trace(err error) error {
	//TODO: not implemented yet
	return err
}

// Errorf is an alias of fmt.Errorf
var Errorf = fmt.Errorf

// Cause returns the real underlying cause (by repeatly calling Unwrap). Returns nil for nil error.
func Cause(err error) error {
	for err != nil {
		unwrErr := Unwrap(err)
		if unwrErr == nil {
			return err
		}
		err = unwrErr
	}

	return nil
}
