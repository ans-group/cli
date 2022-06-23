package draas

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/spf13/cobra"
)

func draasSolutionReplicaRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replica",
		Short: "sub-commands relating to solution replicas",
	}

	// Child root commands
	cmd.AddCommand(draasSolutionReplicaIOPSRootCmd(f))

	return cmd
}
