package errors

import (
	"database/sql"
	"fmt"
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

	// 2. a parent function also testing for an error, assigning more generic
	// internal server error with no custom message
	if err != nil {
		err = HTTPInternalServerError(err, "")
	}

	// UnwrapHTTPError can get HTTPError-typed cause
	herr = UnwrapHTTPError(err)
	if herr == nil {
		t.Fatal("it must be possible to UnwrapHTTPError a HTTPError type")
	}

	if herr.Message != "Could not find the item in the database" {
		t.Fatal("UnwrapHTTPError must return the deepest *HTTPError type")
	}
}
