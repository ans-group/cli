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

func ecloudVolumeGroupRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volumegroup",
		Short: "sub-commands relating to volumegroups",
	}

	// Child commands
	cmd.AddCommand(ecloudVolumeGroupListCmd(f))
	cmd.AddCommand(ecloudVolumeGroupShowCmd(f))
	cmd.AddCommand(ecloudVolumeGroupCreateCmd(f))
	cmd.AddCommand(ecloudVolumeGroupUpdateCmd(f))
	cmd.AddCommand(ecloudVolumeGroupDeleteCmd(f))

	// Child root commands
	cmd.AddCommand(ecloudVolumeGroupVolumeRootCmd(f))

	return cmd
}

func ecloudVolumeGroupListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists volumegroups",
		Long:    "This command lists volumegroups",
		Example: "ans ecloud volumegroup list",
		RunE:    ecloudCobraRunEFunc(f, ecloudVolumeGroupList),
	}

	cmd.Flags().String("name", "", "Volume Group name for filtering")
	cmd.Flags().String("vpc", "", "VPC ID for filtering")

	return cmd
}

func ecloudVolumeGroupList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
		helper.NewStringFilterFlagOption("vpc", "vpc_id"),
	)
	if err != nil {
		return err
	}

	volumegroups, err := service.GetVolumeGroups(params)
	if err != nil {
		return fmt.Errorf("Error retrieving volume groups: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudVolumeGroupsProvider(volumegroups))
}

func ecloudVolumeGroupShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <volumegroup: id>...",
		Short:   "Shows a volumegroup",
		Long:    "This command shows one or more volumegroups",
		Example: "ans ecloud volumegroup show vol-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing volume-group")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVolumeGroupShow),
	}
}

func ecloudVolumeGroupShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var volumegroups []ecloud.VolumeGroup
	for _, arg := range args {
		volumegroup, err := service.GetVolumeGroup(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving volume group [%s]: %s", arg, err)
			continue
		}

		volumegroups = append(volumegroups, volumegroup)
	}

	return output.CommandOutput(cmd, OutputECloudVolumeGroupsProvider(volumegroups))
}

func ecloudVolumeGroupCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a volumegroup",
		Long:    "This command creates a volumegroup",
		Example: "ans ecloud volumegroup create",
		RunE:    ecloudCobraRunEFunc(f, ecloudVolumeGroupCreate),
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of volume-group")
	cmd.Flags().String("vpc", "", "ID of VPC")
	cmd.MarkFlagRequired("vpc")
	cmd.Flags().String("availability-zone", "", "ID of Availability Zone")
	cmd.MarkFlagRequired("availability-zone")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the volume-group has been completely created")

	return cmd
}

func ecloudVolumeGroupCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateVolumeGroupRequest{}
	if cmd.Flags().Changed("name") {
		createRequest.Name, _ = cmd.Flags().GetString("name")
	}
	createRequest.VPCID, _ = cmd.Flags().GetString("vpc")
	createRequest.AvailabilityZoneID, _ = cmd.Flags().GetString("availability-zone")

	taskRef, err := service.CreateVolumeGroup(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating volume group: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for volume group task to complete: %s", err)
		}
	}

	volumegroup, err := service.GetVolumeGroup(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new volume group: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudVolumeGroupsProvider([]ecloud.VolumeGroup{volumegroup}))
}

func ecloudVolumeGroupUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <volumegroup: id>...",
		Short:   "Updates a volumegroup",
		Long:    "This command updates one or more volumegroups",
		Example: "ans ecloud volumegroup update volgroup-abcdef12 --name \"my volumegroup\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing volume-group")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVolumeGroupUpdate),
	}

	cmd.Flags().String("name", "", "Name of volume-group")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the volume-group has been updated")

	return cmd
}

func ecloudVolumeGroupUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchVolumeGroupRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var volumegroups []ecloud.VolumeGroup
	for _, arg := range args {
		task, err := service.PatchVolumeGroup(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating volume group [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for volume group [%s]: %s", arg, err)
				continue
			}
		}

		volumegroup, err := service.GetVolumeGroup(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated volume group [%s]: %s", arg, err)
			continue
		}

		volumegroups = append(volumegroups, volumegroup)
	}

	return output.CommandOutput(cmd, OutputECloudVolumeGroupsProvider(volumegroups))
}

func ecloudVolumeGroupDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <volumegroup: id>...",
		Short:   "Removes a volumegroup",
		Long:    "This command removes one or more volumegroups",
		Example: "ans ecloud volumegroup delete vol-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing volume-group")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVolumeGroupDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the volume group has been completely removed")

	return cmd
}

func ecloudVolumeGroupDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteVolumeGroup(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing volume group [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for volume group [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
