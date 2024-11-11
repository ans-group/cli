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

func ecloudVPNGatewaySpecificationRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "spec",
		Short: "sub-commands relating to VPN gateway specifications",
	}

	// Child commands
	cmd.AddCommand(ecloudVPNGatewaySpecificationListCmd(f))
	cmd.AddCommand(ecloudVPNGatewaySpecificationShowCmd(f))

	return cmd
}

func ecloudVPNGatewaySpecificationListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists VPN gateway specifications",
		Example: "ans ecloud vpngateway spec list",
		RunE:    ecloudCobraRunEFunc(f, ecloudVPNGatewaySpecificationList),
	}

	cmd.Flags().String("name", "", "VPN gateway specification name for filtering")

	return cmd
}

func ecloudVPNGatewaySpecificationList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	specs, err := service.GetVPNGatewaySpecifications(params)
	if err != nil {
		return fmt.Errorf("Error retrieving VPN gateway specifications: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudVPNGatewaySpecificationsProvider(specs))
}

func ecloudVPNGatewaySpecificationShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <specification: id>...",
		Short:   "Show details of a VPN gateway specification",
		Example: "ans ecloud vpngateway spec show vpngs-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VPN gateway specification")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPNGatewaySpecificationShow),
	}
}

func ecloudVPNGatewaySpecificationShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var specs []ecloud.VPNGatewaySpecification
	for _, arg := range args {
		spec, err := service.GetVPNGatewaySpecification(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving VPN gateway specification [%s]: %s", arg, err)
			continue
		}

		specs = append(specs, spec)
	}

	return output.CommandOutput(cmd, OutputECloudVPNGatewaySpecificationsProvider(specs))
}
