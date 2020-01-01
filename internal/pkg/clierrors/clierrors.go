package clierrors

import "fmt"

type MalformedFlagError string

func (e *MalformedFlagError) Error() string {
	return string(*e)
}

// InvalidFlagValueString returns an error string for invalid flag values, with
// error err appended if not nil
func InvalidFlagValueString(name string, value string, err error) string {
	str := fmt.Sprintf("Invalid value '%s' provided for '%s'", value, name)
	if err != nil {
		str = fmt.Sprintf("%s: %s", str, err)
	}

	return str
}
