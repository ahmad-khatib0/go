package format_test

import (
	"fmt"
	"testing"

	"github.com/ahmad-khatib0/go/test-driven-development/ch03-mock-and-assert/format"
	"github.com/stretchr/testify/assert"
)

func TestResult(t *testing.T) {
	// Arrange
	result := 5.55
	expr := "2+3"

	// Act
	wrappedResult := format.Result(expr, result)

	// Assert
	assert.Contains(t, wrappedResult, expr)
	assert.Contains(t, wrappedResult, fmt.Sprint(result))
}
