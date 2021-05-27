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

func ecloudFirewallPolicyTaskRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "sub-commands relating to firewall policy tasks",
	}

	// Child commands
	cmd.AddCommand(ecloudFirewallPolicyTaskListCmd(f))
	cmd.AddCommand(ecloudFirewallPolicyTaskWaitCmd(f))

	return cmd
}

func ecloudFirewallPolicyTaskListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <policy: id>",
		Short:   "Lists firewall policy tasks",
		Long:    "This command lists firewall policy tasks",
		Example: "ukfast ecloud firewall policy task list i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing firewall policy")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFirewallPolicyTaskList),
	}

	cmd.Flags().String("id", "", "Task ID for filtering")
	cmd.Flags().String("name", "", "Task name for filtering")

	return cmd
}

func ecloudFirewallPolicyTaskList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("id", "id"),
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	tasks, err := service.GetFirewallPolicyTasks(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving firewall policy tasks: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudTasksProvider(tasks))
}

func ecloudFirewallPolicyTaskWaitCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "wait <policy: id> <task: id>...",
		Short:   "Waits for a firewall policy task",
		Long:    "This command waits for one or more firewall policy tasks to have expected status",
		Example: "ukfast ecloud firewallpolicy task wait i-abcdef12 task-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing firewall policy")
			}
			if len(args) < 2 {
				return errors.New("Missing task")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFirewallPolicyTaskWait),
	}

	cmd.Flags().String("status", "", fmt.Sprintf("Status to wait for. Defaults to '%s'", ecloud.TaskStatusComplete))

	return cmd
}

func ecloudFirewallPolicyTaskWait(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
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
		err := helper.WaitForCommand(FirewallPolicyTaskStatusWaitFunc(service, args[0], arg, expectedStatus))
		if err != nil {
			output.OutputWithErrorLevelf("Error waiting for firewall policy task [%s]: %s", arg, err)
		}
	}

	return nil
}
