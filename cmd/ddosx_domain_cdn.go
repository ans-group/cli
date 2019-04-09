package cmd

import "github.com/spf13/cobra"

func ddosxDomainCDNRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cdn",
		Short: "sub-commands relating to domain CDN",
	}

	// Child root commands
	cmd.AddCommand(ddosxDomainCDNRuleRootCmd())

	return cmd
}
