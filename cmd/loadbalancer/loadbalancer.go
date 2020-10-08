package loadbalancer

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
)

func LoadBalancerRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "loadbalancer",
		Short: "Commands relating to load balancer service",
	}

	// Child root commands
	cmd.AddCommand(loadbalancerGroupRootCmd(f))
	cmd.AddCommand(loadbalancerConfigurationRootCmd(f))

	return cmd
}
