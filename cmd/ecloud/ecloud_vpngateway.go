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

func ecloudVPNGatewayRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vpngateway",
		Short: "sub-commands relating to VPN gateways",
	}

	// Child commands
	cmd.AddCommand(ecloudVPNGatewayUserRootCmd(f))
	cmd.AddCommand(ecloudVPNGatewaySpecificationRootCmd(f))
	cmd.AddCommand(ecloudVPNGatewayListCmd(f))
	cmd.AddCommand(ecloudVPNGatewayShowCmd(f))
	cmd.AddCommand(ecloudVPNGatewayCreateCmd(f))
	cmd.AddCommand(ecloudVPNGatewayUpdateCmd(f))
	cmd.AddCommand(ecloudVPNGatewayDeleteCmd(f))

	return cmd
}

func ecloudVPNGatewayListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists VPN gateways",
		Example: "ans ecloud vpngateway list",
		RunE:    ecloudCobraRunEFunc(f, ecloudVPNGatewayList),
	}

	cmd.Flags().String("name", "", "VPN gateway name for filtering")

	return cmd
}

func ecloudVPNGatewayList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	gateways, err := service.GetVPNGateways(params)
	if err != nil {
		return fmt.Errorf("Error retrieving VPN gateways: %s", err)
	}

	return output.CommandOutput(cmd, VPNGatewayCollection(gateways))
}

func ecloudVPNGatewayShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <gateway: id>...",
		Short:   "Show details of a VPN gateway",
		Example: "ans ecloud vpngateway show vpng-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VPN gateway")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPNGatewayShow),
	}
}

func ecloudVPNGatewayShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var vpnGateways []ecloud.VPNGateway
	for _, arg := range args {
		vpnGateway, err := service.GetVPNGateway(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving VPN gateway [%s]: %s", arg, err)
			continue
		}

		vpnGateways = append(vpnGateways, vpnGateway)
	}

	return output.CommandOutput(cmd, VPNGatewayCollection(vpnGateways))
}

func ecloudVPNGatewayCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a VPN gateway",
		Example: "ans ecloud vpngateway create --router rtr-abcdef12 --specification vpngs-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudVPNGatewayCreate),
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of gateway")
	cmd.Flags().String("router", "", "ID of router")
	cmd.MarkFlagRequired("router")
	cmd.Flags().String("specification", "", "ID of VPN gateway specification")
	cmd.MarkFlagRequired("specification")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the VPN gateway has been completely created")

	return cmd
}

func ecloudVPNGatewayCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateVPNGatewayRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.RouterID, _ = cmd.Flags().GetString("router")
	createRequest.SpecificationID, _ = cmd.Flags().GetString("specification")

	taskRef, err := service.CreateVPNGateway(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating VPN gateway: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for VPN gateway task to complete: %s", err)
		}
	}

	vpnGateway, err := service.GetVPNGateway(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new VPN gateway: %s", err)
	}

	return output.CommandOutput(cmd, VPNGatewayCollection([]ecloud.VPNGateway{vpnGateway}))
}

func ecloudVPNGatewayUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <gateway: id>...",
		Short:   "Updates a VPN gateway",
		Long:    "Update the name of a VPN gateway",
		Example: "ans ecloud vpngateway update vpng-abcdef12 --name \"my gateway\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VPN gateway")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPNGatewayUpdate),
	}

	cmd.Flags().String("name", "", "Name of gateway")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the VPN gateway has been completely updated")

	return cmd
}

func ecloudVPNGatewayUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchVPNGatewayRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var vpnGateways []ecloud.VPNGateway
	for _, arg := range args {
		task, err := service.PatchVPNGateway(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating VPN gateway [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for VPN gateway [%s]: %s", arg, err)
				continue
			}
		}

		vpnGateway, err := service.GetVPNGateway(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated VPN gateway [%s]: %s", arg, err)
			continue
		}

		vpnGateways = append(vpnGateways, vpnGateway)
	}

	return output.CommandOutput(cmd, VPNGatewayCollection(vpnGateways))
}

func ecloudVPNGatewayDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <gateway: id>...",
		Short:   "Removes a VPN gateway",
		Example: "ans ecloud vpngateway delete vpng-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VPN gateway")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPNGatewayDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the VPN gateway has been completely removed")

	return cmd
}

func ecloudVPNGatewayDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteVPNGateway(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing VPN gateway [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for VPN gateway [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
