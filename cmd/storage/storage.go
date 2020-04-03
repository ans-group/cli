package storage

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
)

func StorageRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "storage",
		Short: "Commands relating to Storage service",
	}

	// Child root commands
	cmd.AddCommand(storageSolutionRootCmd(f))
	cmd.AddCommand(storageHostRootCmd(f))
	cmd.AddCommand(storageVolumeRootCmd(f))

	return cmd
}
