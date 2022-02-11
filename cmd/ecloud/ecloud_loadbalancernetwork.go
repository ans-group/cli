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

func ecloudLoadBalancerNetworkRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "loadbalancernetwork",
		Short: "sub-commands relating to load balancer networks",
	}

	// Child commands
	cmd.AddCommand(ecloudLoadBalancerNetworkListCmd(f))
	cmd.AddCommand(ecloudLoadBalancerNetworkShowCmd(f))
	cmd.AddCommand(ecloudLoadBalancerNetworkCreateCmd(f))
	cmd.AddCommand(ecloudLoadBalancerNetworkUpdateCmd(f))
	cmd.AddCommand(ecloudLoadBalancerNetworkDeleteCmd(f))

	return cmd
}

func ecloudLoadBalancerNetworkListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists load balancer networks",
		Long:    "This command lists load balancer networks",
		Example: "ukfast ecloud loadbalancer list",
		RunE:    ecloudCobraRunEFunc(f, ecloudLoadBalancerNetworkList),
	}

	cmd.Flags().String("name", "", "Name for filtering")

	return cmd
}

func ecloudLoadBalancerNetworkList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	lbs, err := service.GetLoadBalancerNetworks(params)
	if err != nil {
		return fmt.Errorf("Error retrieving load balancer networks: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudLoadBalancerNetworksProvider(lbs))
}

func ecloudLoadBalancerNetworkShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <network: id>...",
		Short:   "Shows an load balancer network",
		Long:    "This command shows one or more load balancer networks",
		Example: "ukfast ecloud loadbalancernetwork show lbn-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing load balancer network")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudLoadBalancerNetworkShow),
	}
}

func ecloudLoadBalancerNetworkShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var lbs []ecloud.LoadBalancerNetwork
	for _, arg := range args {
		lb, err := service.GetLoadBalancerNetwork(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving load balancer network [%s]: %s", arg, err)
			continue
		}

		lbs = append(lbs, lb)
	}

	return output.CommandOutput(cmd, OutputECloudLoadBalancerNetworksProvider(lbs))
}

func ecloudLoadBalancerNetworkCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a load balancer network",
		Long:    "This command creates a load balancer network",
		Example: "ukfast ecloud loadbalancernetwork create --vpc vpc-abcdef12 --availability-zone az-abcdef12 --spec lbs-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudLoadBalancerNetworkCreate),
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of load balancer network")
	cmd.Flags().String("load-balancer", "", "ID of load balancer")
	cmd.MarkFlagRequired("load-balancer")
	cmd.Flags().String("network", "", "ID of network")
	cmd.MarkFlagRequired("network")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the load balancer network has been completely created")

	return cmd
}

func ecloudLoadBalancerNetworkCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateLoadBalancerNetworkRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.LoadBalancerID, _ = cmd.Flags().GetString("load-balancer")
	createRequest.NetworkID, _ = cmd.Flags().GetString("network")

	taskRef, err := service.CreateLoadBalancerNetwork(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating load balancer network: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for load balancer network task to complete: %s", err)
		}
	}

	lb, err := service.GetLoadBalancerNetwork(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new load balancer network: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudLoadBalancerNetworksProvider([]ecloud.LoadBalancerNetwork{lb}))
}

func ecloudLoadBalancerNetworkUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <network: id>...",
		Short:   "Updates a load balancer network",
		Long:    "This command updates one or more load balancer networks",
		Example: "ukfast ecloud loadbalancernetwork update lbn-abcdef12 --name \"my lb\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing load balancer network")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudLoadBalancerNetworkUpdate),
	}

	cmd.Flags().String("name", "", "Name of lb")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the load balancer network has been completely updated")

	return cmd
}

func ecloudLoadBalancerNetworkUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchLoadBalancerNetworkRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var lbs []ecloud.LoadBalancerNetwork
	for _, arg := range args {
		task, err := service.PatchLoadBalancerNetwork(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating load balancer network [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for load balancer network [%s]: %s", arg, err)
				continue
			}
		}

		lb, err := service.GetLoadBalancerNetwork(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated load balancer network [%s]: %s", arg, err)
			continue
		}

		lbs = append(lbs, lb)
	}

	return output.CommandOutput(cmd, OutputECloudLoadBalancerNetworksProvider(lbs))
}

func ecloudLoadBalancerNetworkDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <network: id>...",
		Short:   "Removes a load balancer network",
		Long:    "This command removes one or more load balancer networks",
		Example: "ukfast ecloud loadbalancernetwork delete lbn-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing load balancer network")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudLoadBalancerNetworkDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the load balancer network has been completely removed")

	return cmd
}

func ecloudLoadBalancerNetworkDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteLoadBalancerNetwork(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing load balancer network [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for load balancer network [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
