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

func ecloudVPNServiceRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vpnservice",
		Short: "sub-commands relating to VPN services",
	}

	// Child commands
	cmd.AddCommand(ecloudVPNServiceListCmd(f))
	cmd.AddCommand(ecloudVPNServiceShowCmd(f))
	cmd.AddCommand(ecloudVPNServiceCreateCmd(f))
	cmd.AddCommand(ecloudVPNServiceUpdateCmd(f))
	cmd.AddCommand(ecloudVPNServiceDeleteCmd(f))

	return cmd
}

func ecloudVPNServiceListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists VPN services",
		Long:    "This command lists VPN services",
		Example: "ans ecloud vpnservice list",
		RunE:    ecloudCobraRunEFunc(f, ecloudVPNServiceList),
	}

	cmd.Flags().String("name", "", "VPN service name for filtering")

	return cmd
}

func ecloudVPNServiceList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	services, err := service.GetVPNServices(params)
	if err != nil {
		return fmt.Errorf("Error retrieving VPN services: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudVPNServicesProvider(services))
}

func ecloudVPNServiceShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <service: id>...",
		Short:   "Shows a VPN service",
		Long:    "This command shows one or more VPN services",
		Example: "ans ecloud vpnservice show vpn-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VPN service")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPNServiceShow),
	}
}

func ecloudVPNServiceShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var vpnServices []ecloud.VPNService
	for _, arg := range args {
		vpnService, err := service.GetVPNService(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving VPN service [%s]: %s", arg, err)
			continue
		}

		vpnServices = append(vpnServices, vpnService)
	}

	return output.CommandOutput(cmd, OutputECloudVPNServicesProvider(vpnServices))
}

func ecloudVPNServiceCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a VPN service",
		Long:    "This command creates a VPN service",
		Example: "ans ecloud vpnservice create --router rtr-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudVPNServiceCreate),
	}

	// Setup flags
	cmd.Flags().String("router", "", "ID of router")
	cmd.MarkFlagRequired("router")
	cmd.Flags().String("name", "", "Name of service")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the VPN service has been completely created")

	return cmd
}

func ecloudVPNServiceCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateVPNServiceRequest{}
	createRequest.RouterID, _ = cmd.Flags().GetString("router")
	createRequest.Name, _ = cmd.Flags().GetString("name")

	taskRef, err := service.CreateVPNService(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating VPN service: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for VPN service task to complete: %s", err)
		}
	}

	vpnService, err := service.GetVPNService(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new VPN service: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudVPNServicesProvider([]ecloud.VPNService{vpnService}))
}

func ecloudVPNServiceUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <service: id>...",
		Short:   "Updates a VPN service",
		Long:    "This command updates one or more VPN services",
		Example: "ans ecloud vpnservice update vpn-abcdef12 --name \"my service\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VPN service")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPNServiceUpdate),
	}

	cmd.Flags().String("name", "", "Name of service")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the VPN service has been completely updated")

	return cmd
}

func ecloudVPNServiceUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchVPNServiceRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var vpnServices []ecloud.VPNService
	for _, arg := range args {
		task, err := service.PatchVPNService(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating VPN service [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for VPN service [%s]: %s", arg, err)
				continue
			}
		}

		vpnService, err := service.GetVPNService(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated VPN service [%s]: %s", arg, err)
			continue
		}

		vpnServices = append(vpnServices, vpnService)
	}

	return output.CommandOutput(cmd, OutputECloudVPNServicesProvider(vpnServices))
}

func ecloudVPNServiceDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <service: id>...",
		Short:   "Removes a VPN service",
		Long:    "This command removes one or more VPN services",
		Example: "ans ecloud vpnservice delete vpn-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VPN service")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPNServiceDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the VPN service has been completely removed")

	return cmd
}

func ecloudVPNServiceDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteVPNService(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing VPN service [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for VPN service [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
