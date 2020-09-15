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
