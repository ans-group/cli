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

func ecloudRouterTaskRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "sub-commands relating to router tasks",
	}

	// Child commands
	cmd.AddCommand(ecloudRouterTaskListCmd(f))
	cmd.AddCommand(ecloudRouterTaskWaitCmd(f))

	return cmd
}

func ecloudRouterTaskListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <router: id>",
		Short:   "Lists router tasks",
		Long:    "This command lists router tasks",
		Example: "ukfast ecloud router task list rtr-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing router")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudRouterTaskList),
	}

	cmd.Flags().String("id", "", "Task ID for filtering")
	cmd.Flags().String("name", "", "Task name for filtering")

	return cmd
}

func ecloudRouterTaskList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("id", "id"),
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	tasks, err := service.GetRouterTasks(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving router tasks: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudTasksProvider(tasks))
}

func ecloudRouterTaskWaitCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "wait <router: id> <task: id>...",
		Short:   "Waits for a router task",
		Long:    "This command waits for one or more router tasks to have expected status",
		Example: "ukfast ecloud router task wait rtr-abcdef12 task-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing router")
			}
			if len(args) < 2 {
				return errors.New("Missing task")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudRouterTaskWait),
	}

	cmd.Flags().String("status", "", fmt.Sprintf("Status to wait for. Defaults to '%s'", ecloud.TaskStatusComplete))

	return cmd
}

func ecloudRouterTaskWait(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
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
		err := helper.WaitForCommand(RouterTaskStatusWaitFunc(service, args[0], arg, expectedStatus))
		if err != nil {
			output.OutputWithErrorLevelf("Error waiting for router task [%s]: %s", arg, err)
		}
	}

	return nil
}
