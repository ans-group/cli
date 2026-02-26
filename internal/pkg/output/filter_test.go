package output

import (
	"testing"

	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/stretchr/testify/assert"
)

func TestCompareValues(t *testing.T) {
	tests := []struct {
		name     string
		actual   string
		operator connection.APIRequestFilteringOperator
		expected []string
		want     bool
	}{
		// EQ operator
		{"EQ_Match", "hello", connection.EQOperator, []string{"hello"}, true},
		{"EQ_CaseInsensitive", "Hello", connection.EQOperator, []string{"hello"}, true},
		{"EQ_NoMatch", "hello", connection.EQOperator, []string{"world"}, false},

		// NEQ operator
		{"NEQ_Match", "hello", connection.NEQOperator, []string{"world"}, true},
		{"NEQ_NoMatch", "hello", connection.NEQOperator, []string{"hello"}, false},
		{"NEQ_CaseInsensitive", "Hello", connection.NEQOperator, []string{"hello"}, false},

		// LK operator
		{"LK_GlobMatch", "my-instance", connection.LKOperator, []string{"my-*"}, true},
		{"LK_GlobNoMatch", "your-instance", connection.LKOperator, []string{"my-*"}, false},
		{"LK_CaseInsensitive", "My-Instance", connection.LKOperator, []string{"my-*"}, true},

		// NLK operator
		{"NLK_GlobMatch", "your-instance", connection.NLKOperator, []string{"my-*"}, true},
		{"NLK_GlobNoMatch", "my-instance", connection.NLKOperator, []string{"my-*"}, false},

		// GT operator - numeric
		{"GT_NumericTrue", "5", connection.GTOperator, []string{"3"}, true},
		{"GT_NumericFalse", "3", connection.GTOperator, []string{"5"}, false},
		{"GT_NumericEqual", "3", connection.GTOperator, []string{"3"}, false},
		{"GT_FloatTrue", "3.5", connection.GTOperator, []string{"3.0"}, true},

		// LT operator - numeric
		{"LT_NumericTrue", "3", connection.LTOperator, []string{"5"}, true},
		{"LT_NumericFalse", "5", connection.LTOperator, []string{"3"}, false},
		{"LT_NumericEqual", "3", connection.LTOperator, []string{"3"}, false},

		// GT/LT string fallback
		{"GT_StringTrue", "banana", connection.GTOperator, []string{"apple"}, true},
		{"GT_StringFalse", "apple", connection.GTOperator, []string{"banana"}, false},
		{"LT_StringTrue", "apple", connection.LTOperator, []string{"banana"}, true},

		// IN operator
		{"IN_Match", "running", connection.INOperator, []string{"running", "stopped"}, true},
		{"IN_NoMatch", "pending", connection.INOperator, []string{"running", "stopped"}, false},
		{"IN_CaseInsensitive", "Running", connection.INOperator, []string{"running", "stopped"}, true},

		// NIN operator
		{"NIN_Match", "pending", connection.NINOperator, []string{"running", "stopped"}, true},
		{"NIN_NoMatch", "running", connection.NINOperator, []string{"running", "stopped"}, false},

		// Empty expected
		{"EQ_EmptyExpected", "hello", connection.EQOperator, []string{}, false},
		{"GT_EmptyExpected", "5", connection.GTOperator, []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := compareValues(tt.actual, tt.operator, tt.expected)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMatchesFilter(t *testing.T) {
	row := NewOrderedFields()
	row.Set("name", "my-instance")
	row.Set("vcpu", "4")
	row.Set("status", "Running")

	t.Run("MatchingFilter", func(t *testing.T) {
		f := connection.APIRequestFiltering{Property: "name", Operator: connection.EQOperator, Value: []string{"my-instance"}}
		assert.True(t, matchesFilter(row, f))
	})

	t.Run("NonMatchingFilter", func(t *testing.T) {
		f := connection.APIRequestFiltering{Property: "name", Operator: connection.EQOperator, Value: []string{"other"}}
		assert.False(t, matchesFilter(row, f))
	})

	t.Run("MissingProperty", func(t *testing.T) {
		f := connection.APIRequestFiltering{Property: "missing", Operator: connection.EQOperator, Value: []string{"value"}}
		assert.False(t, matchesFilter(row, f))
	})
}

func TestMatchesFilters(t *testing.T) {
	row := NewOrderedFields()
	row.Set("name", "my-instance")
	row.Set("vcpu", "4")
	row.Set("status", "Running")

	t.Run("AllFiltersMatch", func(t *testing.T) {
		filters := []connection.APIRequestFiltering{
			{Property: "name", Operator: connection.EQOperator, Value: []string{"my-instance"}},
			{Property: "vcpu", Operator: connection.GTOperator, Value: []string{"2"}},
		}
		assert.True(t, matchesFilters(row, filters))
	})

	t.Run("OneFilterFails", func(t *testing.T) {
		filters := []connection.APIRequestFiltering{
			{Property: "name", Operator: connection.EQOperator, Value: []string{"my-instance"}},
			{Property: "vcpu", Operator: connection.GTOperator, Value: []string{"10"}},
		}
		assert.False(t, matchesFilters(row, filters))
	})

	t.Run("EmptyFilters", func(t *testing.T) {
		assert.True(t, matchesFilters(row, nil))
	})
}
