package data

import (
	"fmt"
	"strconv"
)

// Runtime has the underlying type int32 (the same as our Movie struct field).
type Runtime int32

// MarshalJSON() This should return the JSON-encoded value for the movie
// runtime (in our case, it will return a string in the format "<runtime> mins").
// We’re deliberately using a value receiver for our MarshalJSON() method rather than a
// pointer receiver This gives us more flexibility because it means that our custom JSON
// encoding will work on both Runtime values and pointers to Runtime values
func (r Runtime) MarshalJSON() ([]byte, error) {
	// Generate a string containing the movie runtime in the required format.
	jsonValue := fmt.Sprintf("%d mins", r)

	// Use the strconv.Quote() function on the string to wrap it in double quotes. It
	// needs to be surrounded by double quotes in order to be a valid *JSON string*.
	quotedJSONValue := strconv.Quote(jsonValue)

	// Convert the quoted string value to a byte slice and return it.
	return []byte(quotedJSONValue), nil
}
