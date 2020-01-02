package clierrors

import (
	"fmt"
)

type ErrInvalidFlagValue struct {
	Name  string
	Value string
	Err   error
}

func (e *ErrInvalidFlagValue) Error() string {
	str := fmt.Sprintf("Invalid value '%s' provided for '%s'", e.Value, e.Name)
	if e.Err != nil {
		str = fmt.Sprintf("%s: %s", str, e.Err)
	}

	return str
}

func (e *ErrInvalidFlagValue) Unwrap() error { return e.Err }

func NewErrInvalidFlagValue(name string, value string, err error) *ErrInvalidFlagValue {
	return &ErrInvalidFlagValue{Name: name, Value: value, Err: err}
}
