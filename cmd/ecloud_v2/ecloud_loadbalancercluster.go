package ecloud_v2

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudLoadBalancerClusterRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "loadbalancercluster",
		Short: "sub-commands relating to load balancer clusters",
	}

	// Child commands
	cmd.AddCommand(ecloudLoadBalancerClusterListCmd(f))
	cmd.AddCommand(ecloudLoadBalancerClusterShowCmd(f))

	return cmd
}

func ecloudLoadBalancerClusterListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists load balancer clusters",
		Long:    "This command lists load balancer clusters",
		Example: "ukfast ecloud loadbalancercluster list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudLoadBalancerClusterList(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudLoadBalancerClusterList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	rules, err := service.GetLoadBalancerClusters(params)
	if err != nil {
		return fmt.Errorf("Error retrieving load balancer clusters: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudLoadBalancerClustersProvider(rules))
}

func ecloudLoadBalancerClusterShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <rule: id>...",
		Short:   "Shows a load balancer cluster",
		Long:    "This command shows one or more load balancer clusters",
		Example: "ukfast ecloud loadbalancercluster show fwr-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing load balancer cluster")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudLoadBalancerClusterShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudLoadBalancerClusterShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var rules []ecloud.LoadBalancerCluster
	for _, arg := range args {
		rule, err := service.GetLoadBalancerCluster(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving load balancer cluster [%s]: %s", arg, err)
			continue
		}

		rules = append(rules, rule)
	}

	return output.CommandOutput(cmd, OutputECloudLoadBalancerClustersProvider(rules))
}
