package ddosx

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
)

func ddosxDomainACLRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acl",
		Short: "sub-commands relating to domain ACLs",
	}

	// Child root commands
	cmd.AddCommand(ddosxDomainACLIPRuleRootCmd(f))
	cmd.AddCommand(ddosxDomainACLGeoIPRuleRootCmd(f))

	return cmd
}
