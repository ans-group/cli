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

func ecloudTaskRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "sub-commands relating to tasks",
	}

	// Child commands
	cmd.AddCommand(ecloudTaskListCmd(f))
	cmd.AddCommand(ecloudTaskShowCmd(f))
	cmd.AddCommand(ecloudTaskWaitCmd(f))

	return cmd
}

func ecloudTaskListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists tasks",
		Long:    "This command lists tasks",
		Example: "ukfast ecloud task list",
		RunE:    ecloudCobraRunEFunc(f, ecloudTaskList),
	}

	cmd.Flags().String("name", "", "Task name for filtering")

	return cmd
}

func ecloudTaskList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	paginatedTasks, err := service.GetTasksPaginated(params)
	if err != nil {
		return fmt.Errorf("Error retrieving tasks: %s", err)
	}

	return output.CommandOutputPaginated(cmd, OutputECloudTasksProvider(paginatedTasks.Items), paginatedTasks)
}

func ecloudTaskShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <task: id>...",
		Short:   "Shows a task",
		Long:    "This command shows one or more tasks",
		Example: "ukfast ecloud task show vol-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing task")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudTaskShow),
	}
}

func ecloudTaskShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var tasks []ecloud.Task
	for _, arg := range args {
		task, err := service.GetTask(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving task [%s]: %s", arg, err)
			continue
		}

		tasks = append(tasks, task)
	}

	return output.CommandOutput(cmd, OutputECloudTasksProvider(tasks))
}

func ecloudTaskWaitCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "wait <task: id>...",
		Short:   "Waits for a task",
		Long:    "This command waits for one or more tasks to have expected status",
		Example: "ukfast ecloud task wait task-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing task")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudTaskWait),
	}

	cmd.Flags().String("status", "", fmt.Sprintf("Status to wait for. Defaults to '%s'", ecloud.TaskStatusComplete))

	return cmd
}

func ecloudTaskWait(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var expectedStatus ecloud.TaskStatus = ecloud.TaskStatusComplete
	if cmd.Flags().Changed("status") {
		status, _ := cmd.Flags().GetString("status")
		parsedStatus, err := ecloud.ParseTaskStatus(status)
		if err != nil {
			return fmt.Errorf("Failed to parse status: %s", err)
		}
		expectedStatus = parsedStatus
	}

	for _, arg := range args {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, arg, expectedStatus))
		if err != nil {
			output.OutputWithErrorLevelf("Error waiting for task [%s]: %s", arg, err)
		}
	}

	return nil
}
