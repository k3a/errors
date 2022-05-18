package errors

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

func TestTemporary(t *testing.T) {
	normalErr := fmt.Errorf("some error")
	if IsTemporary(normalErr) {
		t.Fatal("a new normal error cannot be temporary")
	}

	tempErr := Temporary(normalErr)
	if !IsTemporary(tempErr) {
		t.Fatal("error wrapped in Temporary() must satisfy IsTemporary()")
	}

	if Cause(tempErr) != normalErr {
		t.Fatalf("unwrap on a Temporary error should still work")
	}

	err := fmt.Errorf("temp err: %w", tempErr)
	if !IsTemporary(err) {
		t.Fatal("temporary error wrapped by fmt.Errorf() must still satisfy IsTemporary()")
	}

	err = Annotate(tempErr, "temp err")
	if !IsTemporary(err) {
		t.Fatal("temporary error wrapped by Annotate() must still satisfy IsTemporary()")
	}
}

func TestUnwrapHTTPError(t *testing.T) {
	err := sql.ErrNoRows

	herr := UnwrapHTTPError(err)
	if herr != nil {
		t.Fatal("UnwrapHTTPError must return nil on non-HTTPError error type")
	}

	herr = UnwrapHTTPError(nil)
	if herr != nil {
		t.Fatal("UnwrapHTTPError must return nil for nil value")
	}

	// 1. initial error spotted and was assigned a user-facing error
	err = HTTPNotFound(err, "Could not find the item in the database")

	// 2. optionally some annotations are added which are ignored by UnwrapHTTPError
	err = Annotatef(err, "while user %d tried to access document %s", 123, "doc123")

	// 3. a parent function also testing for an error, assigning more generic
	// internal server error with no custom message
	err = HTTPInternalServerError(err, "")

	var topmostHTTPError *HTTPError
	if As(err, &topmostHTTPError) {
		errStr := topmostHTTPError.Internal.Error()
		if !strings.HasPrefix(errStr, "while user ") {
			t.Fatalf("Internal error of the newest/topmost HTTPError in the chain should include annotations; Returned error string was '%v'", errStr)
		}
	}

	// UnwrapHTTPError can get HTTPError-typed cause
	herr = UnwrapHTTPError(err)
	if herr == nil {
		t.Fatal("it must be possible to UnwrapHTTPError a HTTPError type")
	}

	if herr.Message != "Could not find the item in the database" {
		t.Fatalf("UnwrapHTTPError must return the deepest/oldest *HTTPError type; it returned %v", herr)
	}

	cause := Cause(err)
	if cause != sql.ErrNoRows {
		t.Fatalf("cause must always return real internal cause (deepest/oldest error in the chain); returned %v instead", cause)
	}

	// backtrace
	bt := Backtrace(err)
	numTraces := 0
	for _, s := range bt {
		if s == '[' {
			numTraces++
		}
	}

	if numTraces != 3 {
		t.Fatal("expected 3 traces in Backtrace")
	}
}
