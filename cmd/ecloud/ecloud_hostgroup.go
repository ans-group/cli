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

func ecloudHostGroupRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hostgroup",
		Short: "sub-commands relating to host groups",
	}

	// Child commands
	cmd.AddCommand(ecloudHostGroupListCmd(f))
	cmd.AddCommand(ecloudHostGroupShowCmd(f))
	cmd.AddCommand(ecloudHostGroupCreateCmd(f))
	cmd.AddCommand(ecloudHostGroupUpdateCmd(f))
	cmd.AddCommand(ecloudHostGroupDeleteCmd(f))

	return cmd
}

func ecloudHostGroupListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists host groups",
		Long:    "This command lists host groups",
		Example: "ukfast ecloud hostgroup list",
		RunE:    ecloudCobraRunEFunc(f, ecloudHostGroupList),
	}

	cmd.Flags().String("name", "", "Host group name for filtering")

	return cmd
}

func ecloudHostGroupList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	groups, err := service.GetHostGroups(params)
	if err != nil {
		return fmt.Errorf("Error retrieving host groups: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudHostGroupsProvider(groups))
}

func ecloudHostGroupShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <group: id>...",
		Short:   "Shows an host group",
		Long:    "This command shows one or more host groups",
		Example: "ukfast ecloud hostgroup show hg-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing host group")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudHostGroupShow),
	}
}

func ecloudHostGroupShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var groups []ecloud.HostGroup
	for _, arg := range args {
		group, err := service.GetHostGroup(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving host group [%s]: %s", arg, err)
			continue
		}

		groups = append(groups, group)
	}

	return output.CommandOutput(cmd, OutputECloudHostGroupsProvider(groups))
}

func ecloudHostGroupCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a host group",
		Long:    "This command creates a host group",
		Example: "ukfast ecloud hostgroup create --policy hg-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudHostGroupCreate),
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of host group")
	cmd.Flags().String("vpc", "", "ID of VPC")
	cmd.MarkFlagRequired("vpc")
	cmd.Flags().String("availability-zone", "", "ID of availability zone")
	cmd.Flags().String("host-spec", "", "ID of host specification")
	cmd.MarkFlagRequired("host-spec")
	cmd.Flags().Bool("windows-enabled", false, "Specifies Windows OS should be enabled for instances")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the host group has been completely created")

	return cmd
}

func ecloudHostGroupCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateHostGroupRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.VPCID, _ = cmd.Flags().GetString("vpc")
	createRequest.AvailabilityZoneID, _ = cmd.Flags().GetString("availability-zone")
	createRequest.HostSpecID, _ = cmd.Flags().GetString("host-spec")
	createRequest.WindowsEnabled, _ = cmd.Flags().GetBool("windows-enabled")

	taskRef, err := service.CreateHostGroup(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating host group: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for host group task to complete: %s", err)
		}
	}

	group, err := service.GetHostGroup(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new host group: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudHostGroupsProvider([]ecloud.HostGroup{group}))
}

func ecloudHostGroupUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <group: id>...",
		Short:   "Updates a host group",
		Long:    "This command updates one or more host groups",
		Example: "ukfast ecloud hostgroup update hg-abcdef12 --name \"my group\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing host group")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudHostGroupUpdate),
	}

	cmd.Flags().String("name", "", "Name of host group")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the host group has been completely updated")

	return cmd
}

func ecloudHostGroupUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchHostGroupRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var groups []ecloud.HostGroup
	for _, arg := range args {
		task, err := service.PatchHostGroup(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating host group [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for host group [%s]: %s", arg, err)
				continue
			}
		}

		group, err := service.GetHostGroup(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated host group [%s]: %s", arg, err)
			continue
		}

		groups = append(groups, group)
	}

	return output.CommandOutput(cmd, OutputECloudHostGroupsProvider(groups))
}

func ecloudHostGroupDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <group: id>...",
		Short:   "Removes a host group",
		Long:    "This command removes one or more host groups",
		Example: "ukfast ecloud hostgroup delete hg-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing host group")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudHostGroupDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the host group has been completely removed")

	return cmd
}

func ecloudHostGroupDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteHostGroup(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing host group [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for host group [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
