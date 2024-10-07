package format_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ahmad-khatib0/go/test-driven-development/ch02-basics/format"
)

func TestResult(t *testing.T) {
	// Arrange
	result := 5.55
	expr := "2+3"

	// Act
	wrappedResult := format.Result(expr, result)

	// Assert
	if !strings.Contains(wrappedResult, expr) {
		t.Errorf("error does not contain: got %s, want %s", wrappedResult, expr)
	}
	if !strings.Contains(wrappedResult, fmt.Sprint(result)) {
		t.Errorf("error does not contain: got %s, want %s", wrappedResult, fmt.Sprint(result))
	}

}
