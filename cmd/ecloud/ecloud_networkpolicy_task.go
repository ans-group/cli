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

func ecloudNetworkPolicyTaskRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "sub-commands relating to network policy tasks",
	}

	// Child commands
	cmd.AddCommand(ecloudNetworkPolicyTaskListCmd(f))

	return cmd
}

func ecloudNetworkPolicyTaskListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <policy: id>",
		Short:   "Lists network policy tasks",
		Long:    "This command lists network policy tasks",
		Example: "ans ecloud network policy task list i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing network policy")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNetworkPolicyTaskList),
	}

	cmd.Flags().String("id", "", "Task ID for filtering")
	cmd.Flags().String("name", "", "Task name for filtering")

	return cmd
}

func ecloudNetworkPolicyTaskList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("id", "id"),
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	tasks, err := service.GetNetworkPolicyTasks(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving network policy tasks: %s", err)
	}

	return output.CommandOutput(cmd, TaskCollection(tasks))
}
