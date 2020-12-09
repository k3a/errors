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

func (t temporaryErr) Temporary() bool {
	return true
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
