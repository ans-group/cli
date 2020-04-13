package pss

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
)

func PSSRootCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pss",
		Short: "Commands relating to PSS service",
	}

	// Child root commands
	cmd.AddCommand(pssRequestRootCmd(f))
	cmd.AddCommand(pssReplyRootCmd(f, fs))

	return cmd
}
