package output

import (
	"strconv"
	"strings"

	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ryanuber/go-glob"
)

// matchesFilters returns true if the row matches ALL of the given filters.
func matchesFilters(row *OrderedFields, filters []connection.APIRequestFiltering) bool {
	for _, f := range filters {
		if !matchesFilter(row, f) {
			return false
		}
	}
	return true
}

// matchesFilter returns true if the row matches a single filter.
func matchesFilter(row *OrderedFields, filter connection.APIRequestFiltering) bool {
	if !row.Exists(filter.Property) {
		return false
	}

	actual := row.Get(filter.Property)
	return compareValues(actual, filter.Operator, filter.Value)
}

// compareValues applies the operator comparison between the actual value and the expected values.
func compareValues(actual string, operator connection.APIRequestFilteringOperator, expected []string) bool {
	switch operator {
	case connection.EQOperator:
		return len(expected) > 0 && strings.EqualFold(actual, expected[0])
	case connection.NEQOperator:
		return len(expected) > 0 && !strings.EqualFold(actual, expected[0])
	case connection.LKOperator:
		return len(expected) > 0 && glob.Glob(strings.ToLower(expected[0]), strings.ToLower(actual))
	case connection.NLKOperator:
		return len(expected) > 0 && !glob.Glob(strings.ToLower(expected[0]), strings.ToLower(actual))
	case connection.GTOperator:
		return compareOrdered(actual, expected, func(cmp int) bool { return cmp > 0 })
	case connection.LTOperator:
		return compareOrdered(actual, expected, func(cmp int) bool { return cmp < 0 })
	case connection.INOperator:
		for _, v := range expected {
			if strings.EqualFold(actual, v) {
				return true
			}
		}
		return false
	case connection.NINOperator:
		for _, v := range expected {
			if strings.EqualFold(actual, v) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// compareOrdered attempts numeric comparison first, falling back to string comparison.
func compareOrdered(actual string, expected []string, predicate func(int) bool) bool {
	if len(expected) == 0 {
		return false
	}

	actualFloat, actualErr := strconv.ParseFloat(actual, 64)
	expectedFloat, expectedErr := strconv.ParseFloat(expected[0], 64)

	if actualErr == nil && expectedErr == nil {
		var cmp int
		switch {
		case actualFloat < expectedFloat:
			cmp = -1
		case actualFloat > expectedFloat:
			cmp = 1
		default:
			cmp = 0
		}
		return predicate(cmp)
	}

	cmp := strings.Compare(strings.ToLower(actual), strings.ToLower(expected[0]))
	return predicate(cmp)
}
