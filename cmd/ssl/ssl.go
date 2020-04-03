package ssl

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
)

func SSLRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssl",
		Short: "Commands relating to SSL service",
	}

	// Child root commands
	cmd.AddCommand(sslCertificateRootCmd(f))

	return cmd
}
