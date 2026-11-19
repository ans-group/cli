package cmd

import (
	"github.com/ans-group/cli/internal/pkg/output"
)

// DoctorResult represents the outcome of a single diagnostic check
type DoctorResult struct {
	Check   string `json:"check"`
	Status  string `json:"status"`
	Detail  string `json:"detail"`
	Context string `json:"context,omitempty"`
}

// DoctorResultCollection is a collection of diagnostic results, implementing the
// output package's field-providing interfaces
type DoctorResultCollection []DoctorResult

// DefaultColumns returns the columns rendered when no properties are specified. The
// context column is only included when at least one result is attributed to a named
// context, i.e. when running with --all-contexts
func (c DoctorResultCollection) DefaultColumns() []string {
	for _, result := range c {
		if len(result.Context) > 0 {
			return []string{"context", "check", "status", "detail"}
		}
	}

	return []string{"check", "status", "detail"}
}

// Fields returns the ordered fields for each diagnostic result
func (c DoctorResultCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, result := range c {
		fields := output.NewOrderedFields()
		fields.Set("context", result.Context)
		fields.Set("check", result.Check)
		fields.Set("status", result.Status)
		fields.Set("detail", result.Detail)

		data = append(data, fields)
	}

	return data
}
