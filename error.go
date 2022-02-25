// Package errors is a drop-in replacement to standard errors, extending it with additional functionality
package errors

import (
	"fmt"
	"runtime"
)

// stringError is a simple error containing just an error string and PC
type stringError struct {
	programCounter
	s string
}

func (e stringError) Error() string {
	return e.s
}

// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func New(text string) error {
	return &stringError{
		programCounter: callerPC(3),
		s:              text,
	}
}

// internalErr is included in any error internal to this package (except stringError)
type internalErr struct {
	error
}

func (h internalErr) SkipInternalErr() error { return h.error }

type internalErrProvider interface {
	SkipInternalErr() error
}

func callerPC(skip int) programCounter {
	pc := make([]uintptr, 1)
	if runtime.Callers(skip, pc) > 0 {
		return programCounter(pc[0])
	}
	return 0
}

type programCounter uintptr

func (pc programCounter) FuncName() string {
	f := runtime.FuncForPC(uintptr(pc))
	return f.Name()
}

func (pc programCounter) FileLine() (file string, line int) {
	f := runtime.FuncForPC(uintptr(pc))
	return f.FileLine(uintptr(pc))
}

// ErrorWithPC represents an error which also hold PC (program counter)
// with the location at which the error ocurred
type ErrorWithPC interface {
	FuncName() string
	FileLine() (file string, line int)
}

type wrapError struct {
	// the actual internal error
	internalErr
	programCounter
	// message of this wrapped error
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

	return &wrapError{
		internalErr:    internalErr{err},
		msg:            msg + ": " + err.Error(),
		programCounter: callerPC(3),
	}
}

// Annotatef is used to add extra context to an existing error (inspired by juju/errors)
func Annotatef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &wrapError{
		internalErr:    internalErr{err},
		msg:            fmt.Sprintf(format, args...) + ": " + err.Error(),
		programCounter: callerPC(3),
	}
}

// Trace adds the location of the Trace call to the stack (inspired by juju/errors)
func Trace(err error) error {
	return &wrapError{
		internalErr:    internalErr{err},
		programCounter: callerPC(3),
	}
}

// Errorf is like fmt.Errorf - formatting a new stringError
func Errorf(format string, a ...interface{}) error {
	return &stringError{
		programCounter: callerPC(3),
		s:              fmt.Sprintf(format, a...),
	}
}

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
