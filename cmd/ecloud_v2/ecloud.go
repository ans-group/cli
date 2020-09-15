package ecloud_v2

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
)

func ECloudV2RootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ecloud",
		Short: "Commands relating to eCloud service",
	}

	// Child root commands
	cmd.AddCommand(ecloudVPCRootCmd(f))
	cmd.AddCommand(ecloudAvailabilityZoneRootCmd(f))
	cmd.AddCommand(ecloudNetworkRootCmd(f))
	cmd.AddCommand(ecloudDHCPRootCmd(f))
	cmd.AddCommand(ecloudVPNRootCmd(f))
	cmd.AddCommand(ecloudInstanceRootCmd(f))
	cmd.AddCommand(ecloudFloatingIPRootCmd(f))
	cmd.AddCommand(ecloudFirewallRuleRootCmd(f))
	cmd.AddCommand(ecloudRouterRootCmd(f))
	cmd.AddCommand(ecloudRegionRootCmd(f))
	cmd.AddCommand(ecloudLoadBalancerClusterRootCmd(f))

	return cmd
}
