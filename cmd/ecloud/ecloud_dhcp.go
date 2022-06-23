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

func ecloudDHCPRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dhcp",
		Short: "sub-commands relating to DHCPs",
	}

	// Child commands
	cmd.AddCommand(ecloudDHCPListCmd(f))
	cmd.AddCommand(ecloudDHCPShowCmd(f))

	return cmd
}

func ecloudDHCPListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists DHCPs",
		Long:    "This command lists DHCPs",
		Example: "ukfast ecloud dhcp list",
		RunE:    ecloudCobraRunEFunc(f, ecloudDHCPList),
	}

	cmd.Flags().String("name", "", "DHCP name for filtering")
	cmd.Flags().String("vpc", "", "VPC ID for filtering")

	return cmd
}

func ecloudDHCPList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
		helper.NewStringFilterFlagOption("vpc", "vpc_id"),
	)
	if err != nil {
		return err
	}

	dhcps, err := service.GetDHCPs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving DHCPs: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudDHCPsProvider(dhcps))
}

func ecloudDHCPShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <dhcp: id>...",
		Short:   "Shows a DHCP",
		Long:    "This command shows one or more DHCPs",
		Example: "ukfast ecloud dhcp show dhcp-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing dhcp")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudDHCPShow),
	}
}

func ecloudDHCPShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var dhcps []ecloud.DHCP
	for _, arg := range args {
		dhcp, err := service.GetDHCP(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving DHCP [%s]: %s", arg, err)
			continue
		}

		dhcps = append(dhcps, dhcp)
	}

	return output.CommandOutput(cmd, OutputECloudDHCPsProvider(dhcps))
}
