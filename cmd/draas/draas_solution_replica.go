package draas

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
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
