package castx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToFloatSliceE(t *testing.T) {
	tests := []struct {
		input    any
		expected []float64
		iserr    bool
	}{
		{[]int{1, 3}, []float64{1, 3}, false},
		{[]any{1.2, 3.2}, []float64{1.2, 3.2}, false},
		{[]string{"2", "3"}, []float64{2, 3}, false},
		{[]string{"2.2", "3.2"}, []float64{2.2, 3.2}, false},
		{[2]string{"2", "3"}, []float64{2, 3}, false},
		{[2]string{"2.2", "3.2"}, []float64{2.2, 3.2}, false},
		// errors
		{nil, nil, true},
		{testing.T{}, nil, true},
		{[]string{"foo", "bar"}, nil, true},
	}

	for i, test := range tests {
		errMsg := fmt.Sprintf("i = %d", i)

		v, err := ToFloatSliceE(test.input)
		if test.iserr {
			assert.Error(t, err, errMsg)
			continue
		}

		assert.NoError(t, err, errMsg)
		assert.Equal(t, test.expected, v, errMsg)

		// Non-E test
		v = ToFloatSlice(test.input)
		assert.Equal(t, test.expected, v, errMsg)
	}
}

func TestToStringSlice(t *testing.T) {
	assert.Equal(t, []string{"foo", "bar"}, ToStringSlice("foo,bar"))
	assert.NotEqual(t, []string{"foo bar baz"}, ToStringSlice("foo bar baz,"))
	assert.Equal(t, []string{"foo bar baz", ""}, ToStringSlice("foo bar baz,"))
	assert.NotEqual(t, []string{"foo", "bar", "baz"}, ToStringSlice("foo bar baz"))
	assert.Equal(t, []string{"foo bar baz"}, ToStringSlice("foo bar baz"))
	assert.Equal(t, []string{"foo", "bar", "baz,", " asdf"}, ToStringSlice("foo,bar,\"baz,\", asdf"))
	assert.Equal(t, []string{"'foo'", "x\"bar", "baz"}, ToStringSlice("'foo',\"x\"\"bar\",baz"))
}
