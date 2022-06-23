package sharedexchange

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/spf13/cobra"
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
