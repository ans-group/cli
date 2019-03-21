package cmd

import "github.com/spf13/cobra"

func ddosxDomainVerificationRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verification",
		Short: "sub-commands relating to domain verification",
	}

	// Child root commands
	cmd.AddCommand(ddosxDomainVerificationFileUploadRootCmd())

	return cmd
}
