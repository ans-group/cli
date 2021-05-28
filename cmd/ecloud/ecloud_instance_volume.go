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

func ecloudInstanceVolumeRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume",
		Short: "sub-commands relating to instance volumes",
	}

	// Child commands
	cmd.AddCommand(ecloudInstanceVolumeListCmd(f))
	cmd.AddCommand(ecloudInstanceVolumeAttachCmd(f))
	cmd.AddCommand(ecloudInstanceVolumeDetachCmd(f))

	return cmd
}

func ecloudInstanceVolumeListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists instance volumes",
		Long:    "This command lists instance volumes",
		Example: "ukfast ecloud instance volume list i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceVolumeList),
	}

	cmd.Flags().String("name", "", "Volume name for filtering")

	return cmd
}

func ecloudInstanceVolumeList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	volumes, err := service.GetInstanceVolumes(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving instance volumes: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudVolumesProvider(volumes))
}

func ecloudInstanceVolumeAttachCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "attach",
		Short:   "Attaches a volume to an instances",
		Long:    "This command attaches a volume to an instance",
		Example: "ukfast ecloud instance volume attach i-abcdef12 --volume vol-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceVolumeAttach),
	}

	cmd.Flags().String("volume", "", "ID of volume to attach")
	cmd.MarkFlagRequired("volume")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until volume has been attached")

	return cmd
}

func ecloudInstanceVolumeAttach(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	req := ecloud.AttachDetachInstanceVolumeRequest{}
	req.VolumeID, _ = cmd.Flags().GetString("volume")

	taskID, err := service.AttachInstanceVolume(args[0], req)
	if err != nil {
		return fmt.Errorf("Error attaching instance volume: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for task: %s", err)
		}
	}

	return nil
}

func ecloudInstanceVolumeDetachCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "detach",
		Short:   "Detaches a volume from an instance",
		Long:    "This command detaches a volume from an instance",
		Example: "ukfast ecloud instance volume detach i-abcdef12 --volume vol-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceVolumeDetach),
	}

	cmd.Flags().String("volume", "", "ID of volume to detach")
	cmd.MarkFlagRequired("volume")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until volume has been detached")

	return cmd
}

func ecloudInstanceVolumeDetach(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	req := ecloud.AttachDetachInstanceVolumeRequest{}
	req.VolumeID, _ = cmd.Flags().GetString("volume")

	taskID, err := service.DetachInstanceVolume(args[0], req)
	if err != nil {
		return fmt.Errorf("Error detaching instance volume: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for task: %s", err)
		}
	}

	return nil
}
