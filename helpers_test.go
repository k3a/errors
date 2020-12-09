package errors

import (
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

	err := fmt.Errorf("temp err: %w", tempErr)
	if !IsTemporary(err) {
		t.Fatal("temporary error wrapped by fmt.Errorf() must still satisfy IsTemporary()")
	}

	err = Annotate(tempErr, "temp err")
	if !IsTemporary(err) {
		t.Fatal("temporary error wrapped by Annotate() must still satisfy IsTemporary()")
	}
}
