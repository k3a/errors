package errors

import (
	"fmt"
	"os"
	"testing"
)

func TestAnnotate(t *testing.T) {
	originalErr := fmt.Errorf("new error")

	an1 := Annotate(originalErr, "first annotation")
	if Cause(an1) != originalErr {
		t.Fatal("result of Cause() on an annotated error must return the original error")
	}
	if !Is(an1, originalErr) {
		t.Fatal("result of standard Is() on an annotated error must success for originalErr")
	}

	_, errNotFound := os.Open("/sfsdfdsf/gdfeyjw/jytrreture/sdfghs")
	if errNotFound == nil {
		t.Fatal("non-existent path expected")
	}

	var pathErr *os.PathError
	if !As(errNotFound, &pathErr) {
		t.Fatal("error returned from os.Open for non-existen file is expected to be type-castable-to os.PathError type")
	}

	annotatedErrNotFound := Annotate(errNotFound, "annotated os.PathError")
	if !As(annotatedErrNotFound, &pathErr) {
		t.Fatal("annotated error returned from os.Open for non-existen file is expected to be type-castable-to os.PathError type")
	}
}
