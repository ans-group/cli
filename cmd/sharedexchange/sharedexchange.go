package sharedexchange

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
)

func SharedExchangeRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sharedexchange",
		Short: "Commands relating to eCloud Flex service",
	}

	// Child root commands
	cmd.AddCommand(sharedexchangeDomainRootCmd(f))

	return cmd
}
