package config

import (
	"strconv"

	"github.com/ans-group/cli/internal/pkg/output"
)

type Context struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

func OutputConfigContextsProvider(contexts []Context) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(contexts),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, context := range contexts {
				fields := output.NewOrderedFields()
				fields.Set("name", output.NewFieldValue(context.Name, true))
				fields.Set("active", output.NewFieldValue(strconv.FormatBool(context.Active), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}
