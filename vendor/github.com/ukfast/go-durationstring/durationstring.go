package durationstring

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Parse takes string s and returns the duration in years, months, days, hours, minutes, seconds, nanoseconds, with
// a non-nil error if there was an error with parsing. Expects s in format ([0-9]+(y|mo|d|h|m|s|ns))+, and accepts any order e.g.
// 5y4m, 5s4d
func Parse(s string) (years, months, days, hours, minutes, seconds, milliseconds, microseconds, nanoseconds int, err error) {
	var digitBuf bytes.Buffer
	var unitBuf bytes.Buffer

	flushBuffers := func() {
		digitBuf = bytes.Buffer{}
		unitBuf = bytes.Buffer{}
	}

	flush := func() error {
		digit := digitBuf.String()
		unit := unitBuf.String()
		flushBuffers()

		if len(digit) < 1 {
			return fmt.Errorf("Digit not supplied for unit '%s'", unit)
		}
		if len(unit) < 1 {
			return fmt.Errorf("Unit not supplied for digit '%s'", digit)
		}

		digitInt, err := strconv.Atoi(digit)
		if err != nil {
			return fmt.Errorf("Failed to parse digit '%s' as int: %s", digit, err.Error())
		}

		switch strings.ToUpper(unit) {
		case "Y":
			years = digitInt
		case "MO":
			months = digitInt
		case "D":
			days = digitInt
		case "H":
			hours = digitInt
		case "M":
			minutes = digitInt
		case "S":
			seconds = digitInt
		case "MS":
			milliseconds = digitInt
		case "US":
			fallthrough
		case "µS":
			microseconds = digitInt
		case "NS":
			nanoseconds = digitInt
		default:
			return fmt.Errorf("invalid unit '%s'", unit)
		}

		return nil
	}

	isUnit := false
	flushBuffers()
	for i, char := range s {
		if unicode.IsDigit(char) {
			digitBuf.WriteRune(char)
			isUnit = false
		} else {
			unitBuf.WriteRune(char)
			isUnit = true
		}

		// if we're at the last rune in iteration, or looking at a unit and next rune is a digit, flush
		if len(s)-1 == i || (isUnit && unicode.IsDigit(rune(s[i+1]))) {
			err := flush()
			if err != nil {
				return 0, 0, 0, 0, 0, 0, 0, 0, 0, err
			}
		}
	}

	return
}

// String takes years, months, days, hours, minutes, seconds, milliseconds, microseconds, nanoseconds and
// returns a formatted string. This function excludes units which are < 1
func String(years, months, days, hours, minutes, seconds, milliseconds, microseconds, nanoseconds int) string {
	buf := bytes.Buffer{}

	add := func(digit int, unit string) {
		if digit < 1 {
			return
		}

		buf.WriteString(fmt.Sprintf("%d%s", digit, unit))
	}

	add(years, "y")
	add(months, "mo")
	add(days, "d")
	add(hours, "h")
	add(minutes, "m")
	add(seconds, "s")
	add(milliseconds, "ms")
	add(microseconds, "µS")
	add(nanoseconds, "ns")

	return buf.String()
}
