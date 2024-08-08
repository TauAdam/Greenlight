package data

import (
"errors"
)

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")
type Runtime int32