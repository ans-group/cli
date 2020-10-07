package ecloud_v2

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudVPNRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vpn",
		Short: "sub-commands relating to VPNs",
	}

	// Child commands
	cmd.AddCommand(ecloudVPNListCmd(f))
	cmd.AddCommand(ecloudVPNShowCmd(f))
	cmd.AddCommand(ecloudVPNCreateCmd(f))
	cmd.AddCommand(ecloudVPNDeleteCmd(f))

	return cmd
}

func ecloudVPNListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists VPNs",
		Long:    "This command lists VPNs",
		Example: "ukfast ecloud vpn list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudVPNList(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudVPNList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	vpns, err := service.GetVPNs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving VPNs: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudVPNsProvider(vpns))
}

func ecloudVPNShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <vpn: id>...",
		Short:   "Shows a VPN",
		Long:    "This command shows one or more VPNs",
		Example: "ukfast ecloud vpn show vpn-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VPN")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudVPNShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudVPNShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var vpns []ecloud.VPN
	for _, arg := range args {
		vpn, err := service.GetVPN(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving VPN [%s]: %s", arg, err)
			continue
		}

		vpns = append(vpns, vpn)
	}

	return output.CommandOutput(cmd, OutputECloudVPNsProvider(vpns))
}

func ecloudVPNCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a VPN",
		Long:    "This command creates a VPN",
		Example: "ukfast ecloud vpn create",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudVPNCreate(c.ECloudService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("router", "", "ID of router")
	cmd.MarkFlagRequired("router")

	return cmd
}

func ecloudVPNCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {

	createRequest := ecloud.CreateVPNRequest{}
	createRequest.RouterID, _ = cmd.Flags().GetString("router")

	vpnID, err := service.CreateVPN(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating VPN: %s", err)
	}

	vpn, err := service.GetVPN(vpnID)
	if err != nil {
		return fmt.Errorf("Error retrieving new VPN: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudVPNsProvider([]ecloud.VPN{vpn}))
}

func ecloudVPNDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <vpn: id...>",
		Short:   "Removes a VPN",
		Long:    "This command removes one or more VPNs",
		Example: "ukfast ecloud vpn delete vpn-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing vpn")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			ecloudVPNDelete(c.ECloudService(), cmd, args)
			return nil
		},
	}
}

func ecloudVPNDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	for _, arg := range args {
		err := service.DeleteVPN(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing VPN [%s]: %s", arg, err)
		}
	}
}
