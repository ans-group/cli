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

func ecloudVPCTaskRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "sub-commands relating to VPC tasks",
	}

	// Child commands
	cmd.AddCommand(ecloudVPCTaskListCmd(f))
	cmd.AddCommand(ecloudVPCTaskWaitCmd(f))

	return cmd
}

func ecloudVPCTaskListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <vpc: id>",
		Short:   "Lists VPC tasks",
		Long:    "This command lists VPC tasks",
		Example: "ukfast ecloud vpc task list vpc-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VPC")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPCTaskList),
	}

	cmd.Flags().String("id", "", "Task ID for filtering")
	cmd.Flags().String("name", "", "Task name for filtering")

	return cmd
}

func ecloudVPCTaskList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("id", "id"),
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	tasks, err := service.GetVPCTasks(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving VPC tasks: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudTasksProvider(tasks))
}

func ecloudVPCTaskWaitCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "wait <vpc: id> <task: id>...",
		Short:   "Waits for a VPC task",
		Long:    "This command waits for one or more VPC tasks to have expected status",
		Example: "ukfast ecloud vpc task wait vpc-abcdef12 task-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VPC")
			}
			if len(args) < 2 {
				return errors.New("Missing task")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPCTaskWait),
	}

	cmd.Flags().String("status", "", fmt.Sprintf("Status to wait for. Defaults to '%s'", ecloud.TaskStatusComplete))

	return cmd
}

func ecloudVPCTaskWait(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
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
		err := helper.WaitForCommand(VPCTaskStatusWaitFunc(service, args[0], arg, expectedStatus))
		if err != nil {
			output.OutputWithErrorLevelf("Error waiting for VPC task [%s]: %s", arg, err)
		}
	}

	return nil
}
