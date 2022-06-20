package config

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func ConfigRootCmd(fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "sub-commands relating to CLI config",
	}

	// Child root commands
	cmd.AddCommand(configContextRootCmd(fs))

	return cmd
}
