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

func ecloudVPNGatewayUserRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "sub-commands relating to VPN gateway users",
	}

	// Child commands
	cmd.AddCommand(ecloudVPNGatewayUserListCmd(f))
	cmd.AddCommand(ecloudVPNGatewayUserShowCmd(f))
	cmd.AddCommand(ecloudVPNGatewayUserCreateCmd(f))
	cmd.AddCommand(ecloudVPNGatewayUserUpdateCmd(f))
	cmd.AddCommand(ecloudVPNGatewayUserDeleteCmd(f))

	return cmd
}

func ecloudVPNGatewayUserListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists VPN gateway users",
		Long:    "This command lists VPN gateway users",
		Example: "ans ecloud vpngateway user list",
		RunE:    ecloudCobraRunEFunc(f, ecloudVPNGatewayUserList),
	}

	cmd.Flags().String("name", "", "VPN gateway user name for filtering")
	cmd.Flags().String("vpngateway", "", "VPN gateway ID for filtering")

	return cmd
}

func ecloudVPNGatewayUserList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
		helper.NewStringFilterFlagOption("vpngateway", "vpn_gateway_id"),
	)
	if err != nil {
		return err
	}

	users, err := service.GetVPNGatewayUsers(params)
	if err != nil {
		return fmt.Errorf("error retrieving VPN gateway users: %s", err)
	}

	return output.CommandOutput(cmd, VPNGatewayUserCollection(users))
}

func ecloudVPNGatewayUserShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <user: id>...",
		Short:   "Shows a VPN gateway user",
		Example: "ans ecloud vpngateway user show vpngu-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing VPN gateway user")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPNGatewayUserShow),
	}
}

func ecloudVPNGatewayUserShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var users []ecloud.VPNGatewayUser
	for _, arg := range args {
		user, err := service.GetVPNGatewayUser(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving VPN gateway user [%s]: %s", arg, err)
			continue
		}

		users = append(users, user)
	}

	return output.CommandOutput(cmd, VPNGatewayUserCollection(users))
}

func ecloudVPNGatewayUserCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a VPN gateway user",
		Example: "ans ecloud vpngateway user create --vpngateway vpng-abcdef12 --username testuser --password testpass",
		RunE:    ecloudCobraRunEFunc(f, ecloudVPNGatewayUserCreate),
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of user")
	cmd.Flags().String("vpngateway", "", "ID of VPN gateway")
	_ = cmd.MarkFlagRequired("vpngateway")
	cmd.Flags().String("username", "", "Username for the VPN user")
	_ = cmd.MarkFlagRequired("username")
	cmd.Flags().String("password", "", "Password for the VPN user")
	_ = cmd.MarkFlagRequired("password")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the user has been created")

	return cmd
}

func ecloudVPNGatewayUserCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateVPNGatewayUserRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.VPNGatewayID, _ = cmd.Flags().GetString("vpngateway")
	createRequest.Username, _ = cmd.Flags().GetString("username")
	createRequest.Password, _ = cmd.Flags().GetString("password")

	taskRef, err := service.CreateVPNGatewayUser(createRequest)
	if err != nil {
		return fmt.Errorf("error creating VPN gateway user: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("error waiting for VPN gateway user task to complete: %s", err)
		}
	}

	user, err := service.GetVPNGatewayUser(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("error retrieving new VPN gateway user: %s", err)
	}

	return output.CommandOutput(cmd, VPNGatewayUserCollection([]ecloud.VPNGatewayUser{user}))
}

func ecloudVPNGatewayUserUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <user: id>...",
		Short:   "Updates a VPN gateway user",
		Long:    "Updates the specified VPN user's friendly name and/or password",
		Example: "ans ecloud vpngateway user update vpngu-abcdef12 --name \"my user\" --password newpass123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing VPN gateway user")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPNGatewayUserUpdate),
	}

	cmd.Flags().String("name", "", "Name of user")
	cmd.Flags().String("password", "", "Password for the VPN user")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the user has been updated")

	return cmd
}

func ecloudVPNGatewayUserUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchVPNGatewayUserRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}
	if cmd.Flags().Changed("password") {
		patchRequest.Password, _ = cmd.Flags().GetString("password")
	}

	var users []ecloud.VPNGatewayUser
	for _, arg := range args {
		task, err := service.PatchVPNGatewayUser(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating VPN gateway user [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for VPN gateway user [%s]: %s", arg, err)
				continue
			}
		}

		user, err := service.GetVPNGatewayUser(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated VPN gateway user [%s]: %s", arg, err)
			continue
		}

		users = append(users, user)
	}

	return output.CommandOutput(cmd, VPNGatewayUserCollection(users))
}

func ecloudVPNGatewayUserDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <user: id>...",
		Short:   "Removes a VPN gateway user",
		Example: "ans ecloud vpngateway user delete vpngu-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing VPN gateway user")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPNGatewayUserDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the user has been removed")

	return cmd
}

func ecloudVPNGatewayUserDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteVPNGatewayUser(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing VPN gateway user [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for VPN gateway user [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
