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

func ecloudVPNProfileGroupRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vpnprofilegroup",
		Short: "sub-commands relating to VPN sessions",
	}

	// Child commands
	cmd.AddCommand(ecloudVPNProfileGroupListCmd(f))
	cmd.AddCommand(ecloudVPNProfileGroupShowCmd(f))

	return cmd
}

func ecloudVPNProfileGroupListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists VPN sessions",
		Long:    "This command lists VPN sessions",
		Example: "ans ecloud vpnprofilegroup list",
		RunE:    ecloudCobraRunEFunc(f, ecloudVPNProfileGroupList),
	}

	cmd.Flags().String("name", "", "VPN session name for filtering")

	return cmd
}

func ecloudVPNProfileGroupList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	sessions, err := service.GetVPNProfileGroups(params)
	if err != nil {
		return fmt.Errorf("Error retrieving VPN sessions: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudVPNProfileGroupsProvider(sessions))
}

func ecloudVPNProfileGroupShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <session: id>...",
		Short:   "Shows a VPN session",
		Long:    "This command shows one or more VPN sessions",
		Example: "ans ecloud vpnprofilegroup show vpns-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VPN session")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPNProfileGroupShow),
	}
}

func ecloudVPNProfileGroupShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var vpnProfileGroups []ecloud.VPNProfileGroup
	for _, arg := range args {
		vpnProfileGroup, err := service.GetVPNProfileGroup(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving VPN session [%s]: %s", arg, err)
			continue
		}

		vpnProfileGroups = append(vpnProfileGroups, vpnProfileGroup)
	}

	return output.CommandOutput(cmd, OutputECloudVPNProfileGroupsProvider(vpnProfileGroups))
}
