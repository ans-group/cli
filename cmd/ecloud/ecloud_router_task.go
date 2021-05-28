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
