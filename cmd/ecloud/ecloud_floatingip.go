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

func ecloudFloatingIPRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "floatingip",
		Short: "sub-commands relating to floating IPs",
	}

	// Child commands
	cmd.AddCommand(ecloudFloatingIPListCmd(f))
	cmd.AddCommand(ecloudFloatingIPShowCmd(f))
	cmd.AddCommand(ecloudFloatingIPCreateCmd(f))
	cmd.AddCommand(ecloudFloatingIPUpdateCmd(f))
	cmd.AddCommand(ecloudFloatingIPDeleteCmd(f))
	cmd.AddCommand(ecloudFloatingIPAssignCmd(f))
	cmd.AddCommand(ecloudFloatingIPUnassignCmd(f))

	return cmd
}

func ecloudFloatingIPListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists floating IPs",
		Long:    "This command lists floating IPs",
		Example: "ans ecloud floatingip list",
		RunE:    ecloudCobraRunEFunc(f, ecloudFloatingIPList),
	}
}

func ecloudFloatingIPList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	fips, err := service.GetFloatingIPs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving floating IPs: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudFloatingIPsProvider(fips))
}

func ecloudFloatingIPShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <floatingip: id>...",
		Short:   "Shows a floating IP",
		Long:    "This command shows one or more floating IPs",
		Example: "ans ecloud floatingip show fip-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing floating IP")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFloatingIPShow),
	}
}

func ecloudFloatingIPShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var fips []ecloud.FloatingIP
	for _, arg := range args {
		fip, err := service.GetFloatingIP(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving floating IP [%s]: %s", arg, err)
			continue
		}

		fips = append(fips, fip)
	}

	return output.CommandOutput(cmd, OutputECloudFloatingIPsProvider(fips))
}

func ecloudFloatingIPCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a floating IP",
		Long:    "This command creates a floating IP address",
		Example: "ans ecloud floatingip create --vpc vpc-abcdef12 --availability-zone az-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudFloatingIPCreate),
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of floating IP")
	cmd.Flags().String("vpc", "", "ID of VPC")
	cmd.MarkFlagRequired("vpc")
	cmd.Flags().String("availability-zone", "", "ID of availability zone")
	cmd.MarkFlagRequired("availability-zone")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the floating IP has been completely created")

	return cmd
}

func ecloudFloatingIPCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateFloatingIPRequest{}
	if cmd.Flags().Changed("name") {
		createRequest.Name, _ = cmd.Flags().GetString("name")
	}
	createRequest.VPCID, _ = cmd.Flags().GetString("vpc")
	createRequest.AvailabilityZoneID, _ = cmd.Flags().GetString("availability-zone")

	taskRef, err := service.CreateFloatingIP(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating floating IP: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for floating IP task to complete: %s", err)
		}
	}

	fip, err := service.GetFloatingIP(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new floating IP: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudFloatingIPsProvider([]ecloud.FloatingIP{fip}))
}

func ecloudFloatingIPUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <fip: id>...",
		Short:   "Updates a floating IP",
		Long:    "This command updates one or more floating IPs",
		Example: "ans ecloud floatingip update fip-abcdef12 --name \"my fip\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing floating IP")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFloatingIPUpdate),
	}

	cmd.Flags().String("name", "", "Name of floating IP")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the floating IP has been completely updated")

	return cmd
}

func ecloudFloatingIPUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchFloatingIPRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var fips []ecloud.FloatingIP
	for _, arg := range args {
		taskRef, err := service.PatchFloatingIP(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating floating IP [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for floating ip [%s]: %s", arg, err)
				continue
			}
		}

		fip, err := service.GetFloatingIP(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated floating IP [%s]: %s", arg, err)
			continue
		}

		fips = append(fips, fip)
	}

	return output.CommandOutput(cmd, OutputECloudFloatingIPsProvider(fips))
}

func ecloudFloatingIPDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <fip: id>...",
		Short:   "Removes a floating IP",
		Long:    "This command removes one or more floating IPs",
		Example: "ans ecloud floatingip delete fip-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing floating IP")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFloatingIPDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the floating IP has been completely removed")

	return cmd
}

func ecloudFloatingIPDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteFloatingIP(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing floating IP [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for removal of floating IP [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}

func ecloudFloatingIPAssignCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "assign <fip: id>",
		Short:   "Assigns a floating IP to a resource",
		Long:    "This command assigns a floating IP to a resource",
		Example: "ans ecloud floatingip assign fip-abcdef12 --resource i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing floating IP")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFloatingIPAssign),
	}

	cmd.Flags().String("resource", "", "ID of resource to assign")
	cmd.MarkFlagRequired("resource")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the floating IP has been completely assigned")

	return cmd
}

func ecloudFloatingIPAssign(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	fipID := args[0]
	resource, _ := cmd.Flags().GetString("resource")
	req := ecloud.AssignFloatingIPRequest{
		ResourceID: resource,
	}

	taskID, err := service.AssignFloatingIP(fipID, req)
	if err != nil {
		return fmt.Errorf("Error assigning floating IP to resource: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for floating IP [%s] to be assigned: %s", fipID, err)
		}
	}

	fip, err := service.GetFloatingIP(fipID)
	if err != nil {
		return fmt.Errorf("Error retrieving new floating IP: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudFloatingIPsProvider([]ecloud.FloatingIP{fip}))
}

func ecloudFloatingIPUnassignCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unassign <fip: id>...",
		Short:   "Unassigns a floating IP",
		Long:    "This command unassigns one or more floating IPs from connected resources",
		Example: "ans ecloud floatingip unassign fip-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing floating IP")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFloatingIPUnassign),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the floating IP has been completely unassigned")

	return cmd
}

func ecloudFloatingIPUnassign(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.UnassignFloatingIP(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error unassigning floating IP [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for floating IP [%s] to be unassigned: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
