package draas

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
)

func DRaaSRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "draas",
		Short: "Commands relating to DRaaS service",
	}

	// Child root commands
	cmd.AddCommand(draasSolutionRootCmd(f))
	cmd.AddCommand(draasIOPSTierRootCmd(f))

	return cmd
}
