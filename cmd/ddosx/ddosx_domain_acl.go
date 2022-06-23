package ddosx

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/spf13/cobra"
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
