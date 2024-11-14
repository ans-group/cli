package ecloud

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/spf13/cobra"
)

func ecloudLoadBalancerRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "loadbalancer",
		Short: "sub-commands relating to load balancers",
	}

	// Child commands
	cmd.AddCommand(ecloudLoadBalancerListCmd(f))
	cmd.AddCommand(ecloudLoadBalancerShowCmd(f))
	cmd.AddCommand(ecloudLoadBalancerCreateCmd(f))
	cmd.AddCommand(ecloudLoadBalancerUpdateCmd(f))
	cmd.AddCommand(ecloudLoadBalancerDeleteCmd(f))

	return cmd
}

func ecloudLoadBalancerListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists load balancers",
		Long:    "This command lists load balancers",
		Example: "ans ecloud loadbalancer list",
		RunE:    ecloudCobraRunEFunc(f, ecloudLoadBalancerList),
	}

	cmd.Flags().String("name", "", "Name for filtering")
	cmd.Flags().String("vpc", "", "VPC ID for filtering")

	return cmd
}

func ecloudLoadBalancerList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
		helper.NewStringFilterFlagOption("vpc", "vpc_id"),
	)
	if err != nil {
		return err
	}

	lbs, err := service.GetLoadBalancers(params)
	if err != nil {
		return fmt.Errorf("Error retrieving load balancers: %s", err)
	}

	return output.CommandOutput(cmd, LoadBalancerCollection(lbs))
}

func ecloudLoadBalancerShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <lb: id>...",
		Short:   "Shows an load balancer",
		Long:    "This command shows one or more load balancers",
		Example: "ans ecloud loadbalancer show lb-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing load balancer")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudLoadBalancerShow),
	}
}

func ecloudLoadBalancerShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var lbs []ecloud.LoadBalancer
	for _, arg := range args {
		lb, err := service.GetLoadBalancer(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving load balancer [%s]: %s", arg, err)
			continue
		}

		lbs = append(lbs, lb)
	}

	return output.CommandOutput(cmd, LoadBalancerCollection(lbs))
}

func ecloudLoadBalancerCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a load balancer",
		Long:    "This command creates a load balancer",
		Example: "ans ecloud loadbalancer create --vpc vpc-abcdef12 --availability-zone az-abcdef12 --spec lbs-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudLoadBalancerCreate),
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of load balancer")
	cmd.Flags().String("vpc", "", "ID of VPC")
	cmd.MarkFlagRequired("vpc")
	cmd.Flags().String("availability-zone", "", "ID of availability zone")
	cmd.MarkFlagRequired("availability-zone")
	cmd.Flags().String("spec", "", "ID of load balancer specification")
	cmd.MarkFlagRequired("spec")
	cmd.Flags().String("network", "", "Network ID for load balancer")
	cmd.MarkFlagRequired("network")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the load balancer has been completely created")

	return cmd
}

func ecloudLoadBalancerCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateLoadBalancerRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.VPCID, _ = cmd.Flags().GetString("vpc")
	createRequest.AvailabilityZoneID, _ = cmd.Flags().GetString("availability-zone")
	createRequest.LoadBalancerSpecID, _ = cmd.Flags().GetString("spec")
	createRequest.NetworkID, _ = cmd.Flags().GetString("network")

	taskRef, err := service.CreateLoadBalancer(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating load balancer: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for load balancer task to complete: %s", err)
		}
	}

	lb, err := service.GetLoadBalancer(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new load balancer: %s", err)
	}

	return output.CommandOutput(cmd, LoadBalancerCollection([]ecloud.LoadBalancer{lb}))
}

func ecloudLoadBalancerUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <lb: id>...",
		Short:   "Updates a load balancer",
		Long:    "This command updates one or more load balancers",
		Example: "ans ecloud loadbalancer update lb-abcdef12 --name \"my lb\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing load balancer")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudLoadBalancerUpdate),
	}

	cmd.Flags().String("name", "", "Name of lb")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the load balancer has been completely updated")

	return cmd
}

func ecloudLoadBalancerUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchLoadBalancerRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var lbs []ecloud.LoadBalancer
	for _, arg := range args {
		task, err := service.PatchLoadBalancer(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating load balancer [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for load balancer [%s]: %s", arg, err)
				continue
			}
		}

		lb, err := service.GetLoadBalancer(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated load balancer [%s]: %s", arg, err)
			continue
		}

		lbs = append(lbs, lb)
	}

	return output.CommandOutput(cmd, LoadBalancerCollection(lbs))
}

func ecloudLoadBalancerDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <lb: id>...",
		Short:   "Removes a load balancer",
		Long:    "This command removes one or more load balancers",
		Example: "ans ecloud loadbalancer delete lb-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing load balancer")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudLoadBalancerDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the load balancer has been completely removed")

	return cmd
}

func ecloudLoadBalancerDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteLoadBalancer(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing load balancer [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for load balancer [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
