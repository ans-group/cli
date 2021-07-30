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

func ecloudVolumeRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume",
		Short: "sub-commands relating to volumes",
	}

	// Child commands
	cmd.AddCommand(ecloudVolumeListCmd(f))
	cmd.AddCommand(ecloudVolumeShowCmd(f))
	cmd.AddCommand(ecloudVolumeCreateCmd(f))
	cmd.AddCommand(ecloudVolumeUpdateCmd(f))
	cmd.AddCommand(ecloudVolumeDeleteCmd(f))

	// Child root commands
	cmd.AddCommand(ecloudVolumeInstanceRootCmd(f))
	cmd.AddCommand(ecloudVolumeTaskRootCmd(f))

	return cmd
}

func ecloudVolumeListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists volumes",
		Long:    "This command lists volumes",
		Example: "ukfast ecloud volume list",
		RunE:    ecloudCobraRunEFunc(f, ecloudVolumeList),
	}

	cmd.Flags().String("name", "", "Volume name for filtering")
	cmd.Flags().String("vpc", "", "VPC ID for filtering")

	return cmd
}

func ecloudVolumeList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
		helper.NewStringFilterFlagOption("vpc", "vpc_id"),
	)
	if err != nil {
		return err
	}

	volumes, err := service.GetVolumes(params)
	if err != nil {
		return fmt.Errorf("Error retrieving volumes: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudVolumesProvider(volumes))
}

func ecloudVolumeShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <volume: id>...",
		Short:   "Shows a volume",
		Long:    "This command shows one or more volumes",
		Example: "ukfast ecloud volume show vol-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing volume")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVolumeShow),
	}
}

func ecloudVolumeShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var volumes []ecloud.Volume
	for _, arg := range args {
		volume, err := service.GetVolume(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving volume [%s]: %s", arg, err)
			continue
		}

		volumes = append(volumes, volume)
	}

	return output.CommandOutput(cmd, OutputECloudVolumesProvider(volumes))
}

func ecloudVolumeCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a volume",
		Long:    "This command creates a volume",
		Example: "ukfast ecloud volume create",
		RunE:    ecloudCobraRunEFunc(f, ecloudVolumeCreate),
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of volume")
	cmd.Flags().String("vpc", "", "ID of VPC")
	cmd.MarkFlagRequired("vpc")
	cmd.Flags().String("availability-zone", "", "ID of Availability Zone")
	cmd.MarkFlagRequired("availability-zone")
	cmd.Flags().Int("capacity", 0, "Capacity of volume in GiB")
	cmd.MarkFlagRequired("capacity")
	cmd.Flags().Int("iops", 0, "IOPS for volume")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the volume has been completely created")

	return cmd
}

func ecloudVolumeCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateVolumeRequest{}
	if cmd.Flags().Changed("name") {
		createRequest.Name, _ = cmd.Flags().GetString("name")
	}
	createRequest.VPCID, _ = cmd.Flags().GetString("vpc")
	createRequest.AvailabilityZoneID, _ = cmd.Flags().GetString("availability-zone")
	createRequest.Capacity, _ = cmd.Flags().GetInt("capacity")
	createRequest.IOPS, _ = cmd.Flags().GetInt("iops")

	taskRef, err := service.CreateVolume(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating volume: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for volume task to complete: %s", err)
		}
	}

	volume, err := service.GetVolume(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new volume: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudVolumesProvider([]ecloud.Volume{volume}))
}

func ecloudVolumeUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <volume: id>...",
		Short:   "Updates a volume",
		Long:    "This command updates one or more volumes",
		Example: "ukfast ecloud volume update vol-abcdef12 --name \"my volume\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing volume")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVolumeUpdate),
	}

	cmd.Flags().String("name", "", "Name of volume")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the volume has been updated")

	return cmd
}

func ecloudVolumeUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchVolumeRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var volumes []ecloud.Volume
	for _, arg := range args {
		task, err := service.PatchVolume(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating volume [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for volume [%s]: %s", arg, err)
				continue
			}
		}

		volume, err := service.GetVolume(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated volume [%s]: %s", arg, err)
			continue
		}

		volumes = append(volumes, volume)
	}

	return output.CommandOutput(cmd, OutputECloudVolumesProvider(volumes))
}

func ecloudVolumeDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <volume: id>...",
		Short:   "Removes a volume",
		Long:    "This command removes one or more volumes",
		Example: "ukfast ecloud volume delete vol-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing volume")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVolumeDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the volume has been completely removed")

	return cmd
}

func ecloudVolumeDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteVolume(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing volume [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for volume [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
