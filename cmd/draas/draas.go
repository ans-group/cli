package draas

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/spf13/cobra"
)

func DRaaSRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "draas",
		Short: "Commands relating to DRaaS service",
	}

	// Child root commands
	cmd.AddCommand(draasSolutionRootCmd(f))
	cmd.AddCommand(draasIOPSTierRootCmd(f))
	cmd.AddCommand(draasBillingTypeRootCmd(f))

	return cmd
}
