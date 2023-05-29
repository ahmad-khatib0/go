package adder_test

import (
	"testing"

	"github.com/ahmad-khatib0/go/idiomatic-approach-book/testing/adder"
)

func TestAddNumbers(t *testing.T) {
	result := adder.AddNumbers(2, 3)
	if result != 5 {
		t.Error("incorrect result: expected 5, got", result)
	}
}
