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

func ecloudInstanceTaskRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "sub-commands relating to instance tasks",
	}

	// Child commands
	cmd.AddCommand(ecloudInstanceTaskListCmd(f))
	cmd.AddCommand(ecloudInstanceTaskWaitCmd(f))

	return cmd
}

func ecloudInstanceTaskListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <instance: id>",
		Short:   "Lists instance tasks",
		Long:    "This command lists instance tasks",
		Example: "ukfast ecloud instance task list i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceTaskList),
	}

	cmd.Flags().String("id", "", "Task ID for filtering")
	cmd.Flags().String("name", "", "Task name for filtering")

	return cmd
}

func ecloudInstanceTaskList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("id", "id"),
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	tasks, err := service.GetInstanceTasks(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving instance tasks: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudTasksProvider(tasks))
}

func ecloudInstanceTaskWaitCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "wait <instance: id> <task: id>...",
		Short:   "Waits for a instance task",
		Long:    "This command waits for one or more instance tasks to have expected status",
		Example: "ukfast ecloud instance task wait i-abcdef12 task-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}
			if len(args) < 2 {
				return errors.New("Missing task")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceTaskWait),
	}

	cmd.Flags().String("status", "", fmt.Sprintf("Status to wait for. Defaults to '%s'", ecloud.TaskStatusComplete))

	return cmd
}

func ecloudInstanceTaskWait(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
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
		err := helper.WaitForCommand(InstanceTaskStatusWaitFunc(service, args[0], arg, expectedStatus))
		if err != nil {
			output.OutputWithErrorLevelf("Error waiting for instance task [%s]: %s", arg, err)
		}
	}

	return nil
}
