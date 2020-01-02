package cmd

import "github.com/spf13/cobra"

func loadtestDomainVerificationRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verification",
		Short: "sub-commands relating to domain verification",
	}

	// Child root commands
	cmd.AddCommand(loadtestDomainVerificationFileRootCmd())
	cmd.AddCommand(loadtestDomainVerificationDNSRootCmd())

	return cmd
}
