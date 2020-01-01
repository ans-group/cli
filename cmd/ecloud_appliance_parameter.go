package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudApplianceParameterRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parameter",
		Short: "sub-commands relating to appliance parameters",
	}

	// Child commands
	cmd.AddCommand(ecloudApplianceParameterListCmd())

	return cmd
}

func ecloudApplianceParameterListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists appliance parameters",
		Long:    "This command lists appliance parameters",
		Example: "ukfast ecloud appliance parameter list 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing appliance")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudApplianceParameterList(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudApplianceParameterList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	parameters, err := service.GetApplianceParameters(args[0], params)
	if err != nil {
		output.Fatalf("Error retrieving appliance parameters: %s", err)
		return
	}

	outputECloudApplianceParameters(parameters)
}
