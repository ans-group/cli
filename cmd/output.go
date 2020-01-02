package cmd

import "github.com/ukfast/cli/internal/pkg/output"

// Output calls the relevant OutputProvider data retrieval methods for given value
// in global variable 'flagFormat'
func Output(out output.OutputHandlerProvider) error {
	handler := output.NewOutputHandler(out, flagFormat)
	handler.Properties = flagProperty
	handler.Template = flagOutputTemplate

	return handler.Handle()
}
