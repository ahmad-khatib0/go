package table_test

import (
	"errors"
	"testing"

	"github.com/ahmad-khatib0/go/test-driven-development/ch04_test_suites/table"
	"github.com/stretchr/testify/assert"
)

func TestDivide(t *testing.T) {
	tests := map[string]struct {
		x, y    int
		wantErr error
		want    string
	}{
		"pos x, pos y":   {x: 8, y: 4, want: "2.00"},
		"neg x, neg y":   {x: -4, y: -8, want: "0.50"},
		"equal x, y":     {x: 4, y: 4, want: "1.00"},
		"max x, pos y":   {x: 127, y: 2, want: "63.50"},
		"min x, pos y":   {x: -128, y: 2, want: "-64.00"},
		"zero x, pos y":  {x: 0, y: 2, want: "0.00"},
		"pos x, zero y":  {x: 10, y: 0, wantErr: errors.New("cannot divide by 0")},
		"zero x, zero y": {x: 0, y: 0, wantErr: errors.New("cannot divide by 0")},
		"max x, max y":   {x: 127, y: 127, want: "1.00"},
		"min x, min y":   {x: -128, y: -128, want: "1.00"},
	}

	for name, rtc := range tests {
		// We assign the current test case to a local tc variable to capture the test case range variable.
		// This is required as the subtest will now run in a goroutine under the hood. We need to create
		// a copy of the current value of the test case to the subtest closure, as opposed to the changing
		// range return value.
		tc := rtc
		t.Run(name, func(t *testing.T) {
			// The *testing.T type provides the t.Parallel() method, which allows us to specify which tests
			// can be run in parallel with other parallel marked tests from the same package. As the subtests
			// of our table-driven test run independently, we need to mark each as parallel and not just the
			// top-level test. The ability to mark certain tests for parallelization is particularly useful
			// together with table-driven tests, which contain independently running test cases
			t.Parallel()
			x, y := int8(tc.x), int8(tc.y)

			r, err := table.Divide(x, y)

			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, tc.want, *r)
		})
	}
}
