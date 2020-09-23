package formam

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Error codes.
const (
	ErrCodeNotAPointer  uint8 = iota // Didn't pass a pointer to Decode().
	ErrCodeArrayIndex                // Error attempting to use an array index (e.g. foo[2]).
	ErrCodeConversion                // Error converting field to the type.
	ErrCodeUnknownType               // Unknown type.
	ErrCodeUnknownField              // No struct field for passed parameter (will never be used if IgnoreUnknownKeys is set).
	ErrCodeRange                     // Number is out of range (e.g. parsing 300 in uint8 would overflow).
	ErrCodeArraySize                 // Array longer than MaxSize.
)

// Error indicates a error produced
type Error struct {
	code        uint8
	field, path string
	err         error
}

func (s *Error) Error() string {
	var b strings.Builder
	b.WriteString("formam: ")
	if s.field != "" {
		b.WriteString("field=")
		b.WriteString(s.field)
		b.WriteString("; path=")
		b.WriteString(s.path)
		b.WriteString(": ")
	}

	b.WriteString(s.err.Error())
	return b.String()
}

// MarshalJSON implements the interface Marshaler
func (s Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Error())
}

// Code for this error. See the ErrCode* constants.
func (s Error) Code() uint8 {
	return s.code
}

// Path for this error.
func (s Error) Path() string {
	return s.path
}

// Cause implements the causer interface from github.com/pkg/errors.
func (s *Error) Cause() error {
	return s.err
}

func newError(code uint8, field, path, format string, a ...interface{}) error {
	return &Error{code: code, field: field, path: path, err: fmt.Errorf(format, a...)}
}
