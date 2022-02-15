package ecloud

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudLoadBalancerLoadBalancerNetworkRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "sub-commands relating to load balancer networks",
	}

	// Child commands
	cmd.AddCommand(ecloudLoadBalancerLoadBalancerNetworkListCmd(f))

	return cmd
}

func ecloudLoadBalancerLoadBalancerNetworkListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists networks for load balancer",
		Long:    "This command lists networks for load balancer",
		Example: "ukfast ecloud loadbalancer network list lb-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing load balancer")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudLoadBalancerLoadBalancerNetworkList),
	}
}

func ecloudLoadBalancerLoadBalancerNetworkList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	networks, err := service.GetLoadBalancerLoadBalancerNetworks(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving load balancer networks: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudLoadBalancerNetworksProvider(networks))
}
