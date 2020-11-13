package ecloudv2

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/ptr"
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
		Use:     "show <cluster: id>...",
		Short:   "Shows a load balancer cluster",
		Long:    "This command shows one or more load balancer clusters",
		Example: "ukfast ecloud loadbalancercluster show lbc-abcdef12",
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

func ecloudLoadBalancerClusterCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a load balancer cluster",
		Long:    "This command creates a load balancer cluster",
		Example: "ukfast ecloud loadbalancercluster create --vpc vpc-abcdef12 --az az-abcdef12",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudLoadBalancerClusterCreate(c.ECloudService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of load balancer cluster")
	cmd.Flags().String("vpc", "", "ID of VPC")
	cmd.MarkFlagRequired("vpc")
	cmd.Flags().String("az", "", "ID of availability zone")
	cmd.MarkFlagRequired("az")

	return cmd
}

func ecloudLoadBalancerClusterCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateLoadBalancerClusterRequest{}
	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		createRequest.Name = ptr.String(name)
	}
	createRequest.VPCID, _ = cmd.Flags().GetString("vpc")
	createRequest.AvailabilityZoneID, _ = cmd.Flags().GetString("az")

	lbcID, err := service.CreateLoadBalancerCluster(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating load balancer cluster: %s", err)
	}

	lbc, err := service.GetLoadBalancerCluster(lbcID)
	if err != nil {
		return fmt.Errorf("Error retrieving new load balancer cluster: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudLoadBalancerClustersProvider([]ecloud.LoadBalancerCluster{lbc}))
}

func ecloudLoadBalancerClusterUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <cluster: id>...",
		Short:   "Updates a load balancer cluster",
		Long:    "This command updates one or more load balancer clusters",
		Example: "ukfast ecloud loadbalancercluster update rtr-abcdef12 --name \"my lbc\"",
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

			return ecloudLoadBalancerClusterUpdate(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Name of load balancer cluster")

	return cmd
}

func ecloudLoadBalancerClusterUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchLoadBalancerClusterRequest{}

	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		patchRequest.Name = ptr.String(name)
	}

	var lbcs []ecloud.LoadBalancerCluster
	for _, arg := range args {
		err := service.PatchLoadBalancerCluster(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating load balancer cluster [%s]: %s", arg, err)
			continue
		}

		lbc, err := service.GetLoadBalancerCluster(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated load balancer cluster [%s]: %s", arg, err)
			continue
		}

		lbcs = append(lbcs, lbc)
	}

	return output.CommandOutput(cmd, OutputECloudLoadBalancerClustersProvider(lbcs))
}

func ecloudLoadBalancerClusterDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <cluster: id...>",
		Short:   "Removes a load balancer cluster",
		Long:    "This command removes one or more load balancer clusters",
		Example: "ukfast ecloud loadbalancercluster delete rtr-abcdef12",
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

			ecloudLoadBalancerClusterDelete(c.ECloudService(), cmd, args)
			return nil
		},
	}
}

func ecloudLoadBalancerClusterDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	for _, arg := range args {
		err := service.DeleteLoadBalancerCluster(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing load balancer cluster [%s]: %s", arg, err)
		}
	}
}
