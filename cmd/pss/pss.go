package pss

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func PSSRootCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pss",
		Short: "Commands relating to PSS service",
	}

	// Child root commands
	cmd.AddCommand(pssRequestRootCmd(f))
	cmd.AddCommand(pssReplyRootCmd(f, fs))
	cmd.AddCommand(pssCaseOptionRootCmd(f))

	return cmd
}
