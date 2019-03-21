package cmd

import "github.com/spf13/cobra"

func ddosxDomainACLRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acl",
		Short: "sub-commands relating to domain ACLs",
	}

	// Child root commands
	cmd.AddCommand(ddosxDomainACLIPRuleRootCmd())
	cmd.AddCommand(ddosxDomainACLGeoIPRuleRootCmd())

	return cmd
}
