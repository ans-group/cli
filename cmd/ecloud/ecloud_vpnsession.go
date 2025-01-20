package ecloud

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/spf13/cobra"
)

func ecloudVPNSessionRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vpnsession",
		Short: "sub-commands relating to VPN sessions",
	}

	// Child commands
	cmd.AddCommand(ecloudVPNSessionListCmd(f))
	cmd.AddCommand(ecloudVPNSessionShowCmd(f))
	cmd.AddCommand(ecloudVPNSessionCreateCmd(f))
	cmd.AddCommand(ecloudVPNSessionUpdateCmd(f))
	cmd.AddCommand(ecloudVPNSessionDeleteCmd(f))

	// Child root commands
	cmd.AddCommand(ecloudVPNSessionPreSharedKeyRootCmd(f))

	return cmd
}

func ecloudVPNSessionListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists VPN sessions",
		Long:    "This command lists VPN sessions",
		Example: "ans ecloud vpnsession list",
		RunE:    ecloudCobraRunEFunc(f, ecloudVPNSessionList),
	}

	cmd.Flags().String("name", "", "VPN session name for filtering")

	return cmd
}

func ecloudVPNSessionList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	sessions, err := service.GetVPNSessions(params)
	if err != nil {
		return fmt.Errorf("Error retrieving VPN sessions: %s", err)
	}

	return output.CommandOutput(cmd, VPNSessionCollection(sessions))
}

func ecloudVPNSessionShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <session: id>...",
		Short:   "Shows a VPN session",
		Long:    "This command shows one or more VPN sessions",
		Example: "ans ecloud vpnsession show vpns-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VPN session")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPNSessionShow),
	}
}

func ecloudVPNSessionShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var vpnSessions []ecloud.VPNSession
	for _, arg := range args {
		vpnSession, err := service.GetVPNSession(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving VPN session [%s]: %s", arg, err)
			continue
		}

		vpnSessions = append(vpnSessions, vpnSession)
	}

	return output.CommandOutput(cmd, VPNSessionCollection(vpnSessions))
}

func ecloudVPNSessionCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a VPN session",
		Long:    "This command creates a VPN session",
		Example: "ans ecloud vpnsession create --router rtr-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudVPNSessionCreate),
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of session")
	cmd.Flags().String("vpn-profile-group", "", "ID of VPN profile group")
	cmd.MarkFlagRequired("vpn-profile-group")
	cmd.Flags().String("vpn-service", "", "ID of VPN service")
	cmd.MarkFlagRequired("vpn-service")
	cmd.Flags().String("vpn-endpoint", "", "ID of VPN endpoint")
	cmd.MarkFlagRequired("vpn-endpoint")
	cmd.Flags().String("remote-ip", "", "IP address of remote")
	cmd.MarkFlagRequired("remote-ip")
	cmd.Flags().String("remote-networks", "", "Comma seperated list of remote networks")
	cmd.MarkFlagRequired("remote-networks")
	cmd.Flags().String("local-networks", "", "Comma seperated list of local networks")
	cmd.MarkFlagRequired("local-networks")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the VPN session has been completely created")

	return cmd
}

func ecloudVPNSessionCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateVPNSessionRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.VPNProfileGroupID, _ = cmd.Flags().GetString("vpn-profile-group")
	createRequest.VPNServiceID, _ = cmd.Flags().GetString("vpn-service")
	createRequest.VPNEndpointID, _ = cmd.Flags().GetString("vpn-endpoint")
	createRequest.RemoteNetworks, _ = cmd.Flags().GetString("remote-networks")
	createRequest.LocalNetworks, _ = cmd.Flags().GetString("local-networks")
	remoteIP, _ := cmd.Flags().GetString("remote-ip")
	createRequest.RemoteIP = connection.IPAddress(remoteIP)

	taskRef, err := service.CreateVPNSession(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating VPN session: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for VPN session task to complete: %s", err)
		}
	}

	vpnSession, err := service.GetVPNSession(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new VPN session: %s", err)
	}

	return output.CommandOutput(cmd, VPNSessionCollection([]ecloud.VPNSession{vpnSession}))
}

func ecloudVPNSessionUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <session: id>...",
		Short:   "Updates a VPN session",
		Long:    "This command updates one or more VPN sessions",
		Example: "ans ecloud vpnsession update vpns-abcdef12 --name \"my session\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VPN session")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPNSessionUpdate),
	}

	cmd.Flags().String("name", "", "Name of session")
	cmd.Flags().String("vpn-profile-group", "", "ID of VPN profile group")
	cmd.Flags().String("remote-ip", "", "IP address of remote")
	cmd.Flags().String("remote-networks", "", "Comma seperated list of remote networks")
	cmd.Flags().String("local-networks", "", "Comma seperated list of local networks")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the VPN session has been completely updated")

	return cmd
}

func ecloudVPNSessionUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchVPNSessionRequest{}

	patchRequest.Name, _ = cmd.Flags().GetString("name")
	patchRequest.VPNProfileGroupID, _ = cmd.Flags().GetString("vpn-profile-group")
	patchRequest.RemoteNetworks, _ = cmd.Flags().GetString("remote-networks")
	patchRequest.LocalNetworks, _ = cmd.Flags().GetString("local-networks")

	if cmd.Flags().Changed("remote-ip") {
		remoteIP, _ := cmd.Flags().GetString("remote-networks")
		patchRequest.RemoteIP = connection.IPAddress(remoteIP)
	}

	var vpnSessions []ecloud.VPNSession
	for _, arg := range args {
		task, err := service.PatchVPNSession(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating VPN session [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for VPN session [%s]: %s", arg, err)
				continue
			}
		}

		vpnSession, err := service.GetVPNSession(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated VPN session [%s]: %s", arg, err)
			continue
		}

		vpnSessions = append(vpnSessions, vpnSession)
	}

	return output.CommandOutput(cmd, VPNSessionCollection(vpnSessions))
}

func ecloudVPNSessionDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <session: id>...",
		Short:   "Removes a VPN session",
		Long:    "This command removes one or more VPN sessions",
		Example: "ans ecloud vpnsession delete vpns-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VPN session")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPNSessionDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the VPN session has been completely removed")

	return cmd
}

func ecloudVPNSessionDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteVPNSession(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing VPN session [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for VPN session [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
