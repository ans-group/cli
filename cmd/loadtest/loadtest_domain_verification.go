package loadtest

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
)

func loadtestDomainVerificationRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verification",
		Short: "sub-commands relating to domain verification",
	}

	// Child root commands
	cmd.AddCommand(loadtestDomainVerificationFileRootCmd(f))
	cmd.AddCommand(loadtestDomainVerificationDNSRootCmd(f))

	return cmd
}
