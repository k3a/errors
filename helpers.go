package errors

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
	error
}

func (e *temporaryErr) Temporary() bool {
	return true
}

func (e *temporaryErr) Unwrap() error {
	return e.error
}

// Temporary creates a new error satisfying IsTemporary
func Temporary(err error) error {
	return &temporaryErr{err}
}

// IsTemporary returns true if a wrapped error implements Temporary() bool and it returns true
func IsTemporary(err error) bool {
	return unwrapIsFunc(err, func(err error) bool {
		t, ok := err.(interface{ Temporary() bool })
		return ok && t.Temporary()
	})
}

type timeoutErr struct {
	error
}

func (e *timeoutErr) Timeout() bool {
	return true
}

func (e *timeoutErr) Unwrap() error {
	return e.error
}

// Timeout creates a new error satisfying IsTimeout
func Timeout(err error) error {
	return &timeoutErr{err}
}

// IsTimeout returns true if a wrapped error implements Timeout() bool and it returns true
func IsTimeout(err error) bool {
	return unwrapIsFunc(err, func(err error) bool {
		t, ok := err.(interface{ Timeout() bool })
		return ok && t.Timeout()
	})
}
