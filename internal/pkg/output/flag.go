package output

import "strings"

func ParseOutputFlag(flag string) (name, arg string) {
	if len(flag) < 1 {
		return
	}

	split := strings.SplitN(flag, "=", 2)

	name = split[0]

	if len(split) == 2 {
		arg = removeQuotes(split[1])
	}

	return
}

func removeQuotes(s string) string {
	if len(s) < 1 {
		return ""
	}

	if s[0] == '"' {
		return strings.Trim(s, "\"")
	}
	if s[0] == '\'' {
		return strings.Trim(s, "'")
	}

	return s
}
