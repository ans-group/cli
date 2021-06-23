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

func ecloudHostRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "host",
		Short: "sub-commands relating to hosts",
	}

	// Child commands
	cmd.AddCommand(ecloudHostListCmd(f))
	cmd.AddCommand(ecloudHostShowCmd(f))
	cmd.AddCommand(ecloudHostCreateCmd(f))
	cmd.AddCommand(ecloudHostUpdateCmd(f))
	cmd.AddCommand(ecloudHostDeleteCmd(f))

	return cmd
}

func ecloudHostListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists hosts",
		Long:    "This command lists hosts",
		Example: "ukfast ecloud host list",
		RunE:    ecloudCobraRunEFunc(f, ecloudHostList),
	}

	cmd.Flags().String("name", "", "Host name for filtering")

	return cmd
}

func ecloudHostList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	groups, err := service.GetHosts(params)
	if err != nil {
		return fmt.Errorf("Error retrieving hosts: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudHostsProvider(groups))
}

func ecloudHostShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <group: id>...",
		Short:   "Shows an host",
		Long:    "This command shows one or more hosts",
		Example: "ukfast ecloud host show h-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing host")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudHostShow),
	}
}

func ecloudHostShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var groups []ecloud.Host
	for _, arg := range args {
		group, err := service.GetHost(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving host [%s]: %s", arg, err)
			continue
		}

		groups = append(groups, group)
	}

	return output.CommandOutput(cmd, OutputECloudHostsProvider(groups))
}

func ecloudHostCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a host",
		Long:    "This command creates a host",
		Example: "ukfast ecloud host create --policy np-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudHostCreate),
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of host")
	cmd.Flags().String("host-group", "", "ID of host group")
	cmd.MarkFlagRequired("host-group")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the host has been completely created")

	return cmd
}

func ecloudHostCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateHostRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.HostGroupID, _ = cmd.Flags().GetString("host-group")

	taskRef, err := service.CreateHost(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating host: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for host task to complete: %s", err)
		}
	}

	group, err := service.GetHost(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new host: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudHostsProvider([]ecloud.Host{group}))
}

func ecloudHostUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <group: id>...",
		Short:   "Updates a host",
		Long:    "This command updates one or more hosts",
		Example: "ukfast ecloud host update np-abcdef12 --name \"my group\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing host")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudHostUpdate),
	}

	cmd.Flags().String("name", "", "Name of host")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the host has been completely updated")

	return cmd
}

func ecloudHostUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchHostRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var groups []ecloud.Host
	for _, arg := range args {
		task, err := service.PatchHost(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating host [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for host [%s]: %s", arg, err)
				continue
			}
		}

		group, err := service.GetHost(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated host [%s]: %s", arg, err)
			continue
		}

		groups = append(groups, group)
	}

	return output.CommandOutput(cmd, OutputECloudHostsProvider(groups))
}

func ecloudHostDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <group: id>...",
		Short:   "Removes a host",
		Long:    "This command removes one or more hosts",
		Example: "ukfast ecloud host delete h-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing host")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudHostDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the host has been completely removed")

	return cmd
}

func ecloudHostDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteHost(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing host [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for host [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
