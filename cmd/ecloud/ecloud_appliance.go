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

func ecloudApplianceRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "appliance",
		Short: "sub-commands relating to appliances",
	}

	// Child commands
	cmd.AddCommand(ecloudApplianceListCmd(f))
	cmd.AddCommand(ecloudApplianceShowCmd(f))

	// Child root commands
	cmd.AddCommand(ecloudApplianceParameterRootCmd(f))

	return cmd
}

func ecloudApplianceListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists appliances",
		Long:    "This command lists appliances",
		Example: "ukfast ecloud appliance list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudApplianceList(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudApplianceList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	appliances, err := service.GetAppliances(params)
	if err != nil {
		return fmt.Errorf("Error retrieving appliances: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudAppliancesProvider(appliances))
}

func ecloudApplianceShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <appliance: id>...",
		Short:   "Shows a appliance",
		Long:    "This command shows one or more appliances",
		Example: "ukfast ecloud vm appliance 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing appliance")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudApplianceShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudApplianceShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var appliances []ecloud.Appliance
	for _, arg := range args {
		appliance, err := service.GetAppliance(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving appliance [%s]: %s", arg, err)
			continue
		}

		appliances = append(appliances, appliance)
	}

	return output.CommandOutput(cmd, OutputECloudAppliancesProvider(appliances))
}
