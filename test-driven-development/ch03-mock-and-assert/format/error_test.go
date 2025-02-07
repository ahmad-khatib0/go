package format_test

import (
	"errors"
	"testing"

	"github.com/ahmad-khatib0/go/test-driven-development/ch03-mock-and-assert/format"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {

	// Arrange
	initialErr := errors.New("error msg")
	expr := "2%3"

	// Act
	wrappedErr := format.Error(expr, initialErr)

	// Assert
	assert.Contains(t, wrappedErr.Error(), initialErr.Error())
	assert.Contains(t, wrappedErr.Error(), expr)

}
