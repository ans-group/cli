package helper

import (
	"strconv"
	"strings"
)

// JoinInt joins an array of integers with given sep
func JoinInt(a []int, sep string) string {
	if len(a) == 0 {
		return ""
	}

	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}

	return strings.Join(b, sep)
}
