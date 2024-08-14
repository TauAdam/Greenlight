package data

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int32

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	// Unquote the JSON value
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// Split the string to isolate the number
	parts := strings.Split(unquotedJSONValue, " ")

	// Sanity check the format
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	// Parse the number
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// Set the receiver value
	*r = Runtime(i)
	return nil
}
