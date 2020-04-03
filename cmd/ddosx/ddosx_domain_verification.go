package ddosx

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
)

func ddosxDomainVerificationRootCmd(f factory.ClientFactory, appFilesystem afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verification",
		Short: "sub-commands relating to domain verification",
	}

	// Child root commands
	cmd.AddCommand(ddosxDomainVerificationFileUploadRootCmd(f, appFilesystem))
	cmd.AddCommand(ddosxDomainVerificationDNSRootCmd(f))

	return cmd
}
