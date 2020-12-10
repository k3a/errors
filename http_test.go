package errors

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestHTTP(t *testing.T) {
	var err error

	err = &HTTPError{Code: http.StatusNotFound}
	if strings.TrimSpace(err.Error()) == "" {
		t.Fatal("HTTP error with just code should still produce an error string")
	}

	err = HTTPBadRequest(Temporary(errors.New("internal err")), "Unable to process your request")
	if IsTemporary(err) == false {
		t.Fatal("HTTP error with a temporary internal error should satisfy IsTemporary()")
	}

	errStr := err.Error()
	if strings.Contains(errStr, "internal err") {
		t.Fatal("HTTP error's Error() result must not contain internal error details")
	}
	if strings.Contains(errStr, "Unable to process your request") == false {
		t.Fatal("HTTP error's Error() result must contain the public-facing message")
	}
}
