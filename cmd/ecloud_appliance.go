package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudApplianceRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "appliance",
		Short: "sub-commands relating to appliances",
	}

	// Child commands
	cmd.AddCommand(ecloudApplianceListCmd())
	cmd.AddCommand(ecloudApplianceShowCmd())

	// Child root commands
	// cmd.AddCommand(ecloudApplianceTemplateRootCmd())

	return cmd
}

func ecloudApplianceListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists appliances",
		Long:    "This command lists appliances",
		Example: "ukfast ecloud appliance list",
		Run: func(cmd *cobra.Command, args []string) {
			ecloudApplianceList(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudApplianceList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	appliances, err := service.GetAppliances(params)
	if err != nil {
		output.Fatalf("Error retrieving appliances: %s", err)
		return
	}

	outputECloudAppliances(appliances)
}

func ecloudApplianceShowCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ecloudApplianceShow(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudApplianceShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	var appliances []ecloud.Appliance
	for _, arg := range args {
		appliance, err := service.GetAppliance(arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving appliance [%s]: %s", arg, err)
			continue
		}

		appliances = append(appliances, appliance)
	}

	outputECloudAppliances(appliances)
}
