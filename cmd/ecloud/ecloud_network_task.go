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

func ecloudNetworkTaskRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "sub-commands relating to network tasks",
	}

	// Child commands
	cmd.AddCommand(ecloudNetworkTaskListCmd(f))
	cmd.AddCommand(ecloudNetworkTaskWaitCmd(f))

	return cmd
}

func ecloudNetworkTaskListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <network: id>",
		Short:   "Lists network tasks",
		Long:    "This command lists network tasks",
		Example: "ukfast ecloud network task list net-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing network")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNetworkTaskList),
	}

	cmd.Flags().String("id", "", "Task ID for filtering")
	cmd.Flags().String("name", "", "Task name for filtering")

	return cmd
}

func ecloudNetworkTaskList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("id", "id"),
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	tasks, err := service.GetNetworkTasks(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving network tasks: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudTasksProvider(tasks))
}

func ecloudNetworkTaskWaitCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "wait <network: id> <task: id>...",
		Short:   "Waits for a network task",
		Long:    "This command waits for one or more network tasks to have expected status",
		Example: "ukfast ecloud network task wait net-abcdef12 task-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing network")
			}
			if len(args) < 2 {
				return errors.New("Missing task")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNetworkTaskWait),
	}

	cmd.Flags().String("status", "", fmt.Sprintf("Status to wait for. Defaults to '%s'", ecloud.TaskStatusComplete))

	return cmd
}

func ecloudNetworkTaskWait(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var expectedStatus ecloud.TaskStatus = ecloud.TaskStatusComplete
	if cmd.Flags().Changed("status") {
		status, _ := cmd.Flags().GetString("status")
		parsedStatus, err := ecloud.ParseTaskStatus(status)
		if err != nil {
			return fmt.Errorf("Failed to parse status: %s", err)
		}
		expectedStatus = parsedStatus
	}

	for _, arg := range args[1:] {
		err := helper.WaitForCommand(NetworkTaskStatusWaitFunc(service, args[0], arg, expectedStatus))
		if err != nil {
			output.OutputWithErrorLevelf("Error waiting for network task [%s]: %s", arg, err)
		}
	}

	return nil
}
