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

func ecloudVolumeTaskRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "sub-commands relating to volume tasks",
	}

	// Child commands
	cmd.AddCommand(ecloudVolumeTaskListCmd(f))
	cmd.AddCommand(ecloudVolumeTaskWaitCmd(f))

	return cmd
}

func ecloudVolumeTaskListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists volume tasks",
		Long:    "This command lists volume tasks",
		Example: "ukfast ecloud volume task list vol-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing volume")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVolumeTaskList),
	}

	cmd.Flags().String("id", "", "Task ID for filtering")
	cmd.Flags().String("name", "", "Task name for filtering")

	return cmd
}

func ecloudVolumeTaskList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("id", "id"),
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	tasks, err := service.GetVolumeTasks(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving volume tasks: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudTasksProvider(tasks))
}

func ecloudVolumeTaskWaitCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "wait",
		Short:   "Waits for a volume task",
		Long:    "This command waits for a volume task to have expected status",
		Example: "ukfast ecloud volume task wait vol-abcdef12 task-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing volume")
			}
			if len(args) < 2 {
				return errors.New("Missing task")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVolumeTaskWait),
	}

	cmd.Flags().String("status", "", fmt.Sprintf("Status to wait for. Defaults to '%s'", ecloud.TaskStatusComplete))

	return cmd
}

func ecloudVolumeTaskWait(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var expectedStatus ecloud.TaskStatus = ecloud.TaskStatusComplete
	if cmd.Flags().Changed("status") {
		status, _ := cmd.Flags().GetString("status")
		parsedStatus, err := ecloud.ParseTaskStatus(status)
		if err != nil {
			return fmt.Errorf("Failed to parse status: %s", err)
		}
		expectedStatus = parsedStatus
	}

	err := helper.WaitForCommand(VolumeTaskStatusWaitFunc(service, args[0], args[1], expectedStatus))
	if err != nil {
		return fmt.Errorf("Error waiting for volume task: %s", err)
	}

	return nil
}
