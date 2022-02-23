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

func ecloudVIPRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vip",
		Short: "sub-commands relating to VIPs",
	}

	// Child commands
	cmd.AddCommand(ecloudVIPListCmd(f))
	cmd.AddCommand(ecloudVIPShowCmd(f))
	cmd.AddCommand(ecloudVIPCreateCmd(f))
	cmd.AddCommand(ecloudVIPUpdateCmd(f))
	cmd.AddCommand(ecloudVIPDeleteCmd(f))

	return cmd
}

func ecloudVIPListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists VIPs",
		Long:    "This command lists VIPs",
		Example: "ukfast ecloud vip list",
		RunE:    ecloudCobraRunEFunc(f, ecloudVIPList),
	}

	cmd.Flags().String("name", "", "Name for filtering")
	cmd.Flags().String("vpc", "", "VPC ID for filtering")

	return cmd
}

func ecloudVIPList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
		helper.NewStringFilterFlagOption("load-balancer", "load_balancer_id"),
	)
	if err != nil {
		return err
	}

	vips, err := service.GetVIPs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving VIPs: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudVIPsProvider(vips))
}

func ecloudVIPShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <vip: id>...",
		Short:   "Shows an VIP",
		Long:    "This command shows one or more VIPs",
		Example: "ukfast ecloud vip show vip-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VIP")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVIPShow),
	}
}

func ecloudVIPShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var vips []ecloud.VIP
	for _, arg := range args {
		vip, err := service.GetVIP(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving VIP [%s]: %s", arg, err)
			continue
		}

		vips = append(vips, vip)
	}

	return output.CommandOutput(cmd, OutputECloudVIPsProvider(vips))
}

func ecloudVIPCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a VIP",
		Long:    "This command creates a VIP",
		Example: "ukfast ecloud vip create --name testvip --load-balancer lb-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudVIPCreate),
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of VIP")
	cmd.Flags().String("load-balancer", "", "ID of load balancer")
	cmd.MarkFlagRequired("load-balancer")
	cmd.Flags().Bool("allocate-floating-ip", false, "Specifies a floating IP should be assigned to the VIP")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the VIP has been completely created")

	return cmd
}

func ecloudVIPCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateVIPRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.LoadBalancerID, _ = cmd.Flags().GetString("load-balancer")
	createRequest.AllocateFloatingIP, _ = cmd.Flags().GetBool("allocate-floating-ip")

	taskRef, err := service.CreateVIP(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating VIP: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for VIP task to complete: %s", err)
		}
	}

	vip, err := service.GetVIP(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new VIP: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudVIPsProvider([]ecloud.VIP{vip}))
}

func ecloudVIPUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <vip: id>...",
		Short:   "Updates a VIP",
		Long:    "This command updates one or more VIPs",
		Example: "ukfast ecloud vip update vip-abcdef12 --name \"my vip\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VIP")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVIPUpdate),
	}

	cmd.Flags().String("name", "", "Name of vip")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the VIP has been completely updated")

	return cmd
}

func ecloudVIPUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchVIPRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var vips []ecloud.VIP
	for _, arg := range args {
		task, err := service.PatchVIP(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating VIP [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for VIP [%s]: %s", arg, err)
				continue
			}
		}

		vip, err := service.GetVIP(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated VIP [%s]: %s", arg, err)
			continue
		}

		vips = append(vips, vip)
	}

	return output.CommandOutput(cmd, OutputECloudVIPsProvider(vips))
}

func ecloudVIPDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <vip: id>...",
		Short:   "Removes a VIP",
		Long:    "This command removes one or more VIPs",
		Example: "ukfast ecloud vip delete vip-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VIP")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVIPDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the VIP has been completely removed")

	return cmd
}

func ecloudVIPDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteVIP(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing VIP [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for VIP [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
