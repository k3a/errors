package errors

var (
	// A simple static error satisfying IsTemporary
	ErrTemporary = &temporaryErr{internalErr{New("temporary error ocurred")}}
	// A simple static error satisfying IsTimeout
	ErrTimeout = &timeoutErr{internalErr{New("operation timed out")}}
)

func unwrapIsFunc(err error, fn func(err error) bool) bool {
	for err != nil {
		if fn(err) {
			return true
		}
		err = Unwrap(err)
	}
	return false
}

type temporaryErr struct {
	internalErr
}

func (e *temporaryErr) Temporary() bool {
	return true
}

// Temporary creates a new error satisfying IsTemporary
func Temporary(err error) error {
	if err == nil {
		err = ErrTemporary
	}
	return &temporaryErr{internalErr{err}}
}

// IsTemporary returns true if a wrapped error implements Temporary() bool and it returns true
func IsTemporary(err error) bool {
	return unwrapIsFunc(err, func(err error) bool {
		t, ok := err.(interface{ Temporary() bool })
		return ok && t.Temporary()
	})
}

type timeoutErr struct {
	internalErr
}

func (e *timeoutErr) Timeout() bool {
	return true
}

// Timeout creates a new error satisfying IsTimeout
func Timeout(err error) error {
	if err == nil {
		err = ErrTimeout
	}
	return &timeoutErr{internalErr{err}}
}

// IsTimeout returns true if a wrapped error implements Timeout() bool and it returns true
func IsTimeout(err error) bool {
	return unwrapIsFunc(err, func(err error) bool {
		t, ok := err.(interface{ Timeout() bool })
		return ok && t.Timeout()
	})
}

// If calls errConstruct if err is not nil with an optional msg passed to it.
// Function pointer errConstruct can be nil, in which case an error is returned as-is.
// For nil err, the function returns nil as well.
// Example: return If(err, HTTPServiceUnavailable, "Service is currently unavalable, sorry")
func If(err error, errConstruct func(error, string) error, msg string) error {
	if err == nil {
		return nil
	}

	if errConstruct == nil {
		return err
	}

	return errConstruct(err, msg)
}

// Iff calls errConstruct if err is not nil with an optional template and template args passed to it.
// This function is "format-capable" version of If.
// Function pointer errConstruct can be nil, in which case an error is returned as-is.
// For nil err, the function returns nil as well.
// Example: return Iff(err, HTTPServiceUnavailable, "We are overloaded, please wait %d seconds", waitTime)
func Iff(err error, errConstruct func(error, string, ...interface{}) error, template string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	if errConstruct == nil {
		return err
	}

	return errConstruct(err, template, args...)
}

// UnwrapHTTPError unwraps err (if it is *HTTPError instance)
// to find the deepest error of the type *HTTPError type
// Returns nil if err is nil or not *HTTPError type
func UnwrapHTTPError(err error) *HTTPError {
	if herr, ok := err.(*HTTPError); ok {
		return herr.UnwrapHTTPError()
	}
	return nil
}
