package ddosx

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func DDoSXRootCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ddosx",
		Short: "Commands relating to DDoSX service",
	}

	// Child root commands
	cmd.AddCommand(ddosxDomainRootCmd(f, fs))
	cmd.AddCommand(ddosxRecordRootCmd(f))
	cmd.AddCommand(ddosxSSLRootCmd(f, fs))
	cmd.AddCommand(ddosxWAFRootCmd(f))

	return cmd
}
