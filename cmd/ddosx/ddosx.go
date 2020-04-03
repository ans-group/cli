package ddosx

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
)

func DDoSXRootCmd(f factory.ClientFactory, appFilesystem afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ddosx",
		Short: "Commands relating to DDoSX service",
	}

	// Child root commands
	cmd.AddCommand(ddosxDomainRootCmd(f, appFilesystem))
	cmd.AddCommand(ddosxRecordRootCmd(f))
	cmd.AddCommand(ddosxSSLRootCmd(f))

	return cmd
}
