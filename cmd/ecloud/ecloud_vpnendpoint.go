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

func ecloudVPNEndpointRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vpnendpoint",
		Short: "sub-commands relating to VPN endpoints",
	}

	// Child commands
	cmd.AddCommand(ecloudVPNEndpointListCmd(f))
	cmd.AddCommand(ecloudVPNEndpointShowCmd(f))
	cmd.AddCommand(ecloudVPNEndpointCreateCmd(f))
	cmd.AddCommand(ecloudVPNEndpointUpdateCmd(f))
	cmd.AddCommand(ecloudVPNEndpointDeleteCmd(f))

	return cmd
}

func ecloudVPNEndpointListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists VPN endpoints",
		Long:    "This command lists VPN endpoints",
		Example: "ukfast ecloud vpnendpoint list",
		RunE:    ecloudCobraRunEFunc(f, ecloudVPNEndpointList),
	}

	cmd.Flags().String("name", "", "VPN endpoint name for filtering")

	return cmd
}

func ecloudVPNEndpointList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	endpoints, err := service.GetVPNEndpoints(params)
	if err != nil {
		return fmt.Errorf("Error retrieving VPN endpoints: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudVPNEndpointsProvider(endpoints))
}

func ecloudVPNEndpointShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <endpoint: id>...",
		Short:   "Shows a VPN endpoint",
		Long:    "This command shows one or more VPN endpoints",
		Example: "ukfast ecloud vpnendpoint show vpne-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VPN endpoint")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPNEndpointShow),
	}
}

func ecloudVPNEndpointShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var vpnEndpoints []ecloud.VPNEndpoint
	for _, arg := range args {
		vpnEndpoint, err := service.GetVPNEndpoint(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving VPN endpoint [%s]: %s", arg, err)
			continue
		}

		vpnEndpoints = append(vpnEndpoints, vpnEndpoint)
	}

	return output.CommandOutput(cmd, OutputECloudVPNEndpointsProvider(vpnEndpoints))
}

func ecloudVPNEndpointCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a VPN endpoint",
		Long:    "This command creates a VPN endpoint",
		Example: "ukfast ecloud vpnendpoint create --router rtr-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudVPNEndpointCreate),
	}

	// Setup flags
	cmd.Flags().String("vpn-service", "", "ID of VPN service")
	cmd.MarkFlagRequired("vpn-service")
	cmd.Flags().String("name", "", "Name of endpoint")
	cmd.Flags().String("floating-ip", "", "Floating IP ID")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the VPN endpoint has been completely created")

	return cmd
}

func ecloudVPNEndpointCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateVPNEndpointRequest{}
	createRequest.VPNServiceID, _ = cmd.Flags().GetString("vpn-service")
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.FloatingIPID, _ = cmd.Flags().GetString("floating-ip")

	taskRef, err := service.CreateVPNEndpoint(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating VPN endpoint: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for VPN endpoint task to complete: %s", err)
		}
	}

	vpnEndpoint, err := service.GetVPNEndpoint(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new VPN endpoint: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudVPNEndpointsProvider([]ecloud.VPNEndpoint{vpnEndpoint}))
}

func ecloudVPNEndpointUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <endpoint: id>...",
		Short:   "Updates a VPN endpoint",
		Long:    "This command updates one or more VPN endpoints",
		Example: "ukfast ecloud vpnendpoint update vpne-abcdef12 --name \"my endpoint\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VPN endpoint")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPNEndpointUpdate),
	}

	cmd.Flags().String("name", "", "Name of endpoint")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the VPN endpoint has been completely updated")

	return cmd
}

func ecloudVPNEndpointUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchVPNEndpointRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var vpnEndpoints []ecloud.VPNEndpoint
	for _, arg := range args {
		task, err := service.PatchVPNEndpoint(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating VPN endpoint [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for VPN endpoint [%s]: %s", arg, err)
				continue
			}
		}

		vpnEndpoint, err := service.GetVPNEndpoint(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated VPN endpoint [%s]: %s", arg, err)
			continue
		}

		vpnEndpoints = append(vpnEndpoints, vpnEndpoint)
	}

	return output.CommandOutput(cmd, OutputECloudVPNEndpointsProvider(vpnEndpoints))
}

func ecloudVPNEndpointDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <endpoint: id>...",
		Short:   "Removes a VPN endpoint",
		Long:    "This command removes one or more VPN endpoints",
		Example: "ukfast ecloud vpnendpoint delete vpne-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VPN endpoint")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPNEndpointDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the VPN endpoint has been completely removed")

	return cmd
}

func ecloudVPNEndpointDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteVPNEndpoint(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing VPN endpoint [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for VPN endpoint [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
