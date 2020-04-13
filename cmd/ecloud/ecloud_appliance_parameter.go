package ecloud

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudApplianceParameterRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parameter",
		Short: "sub-commands relating to appliance parameters",
	}

	// Child commands
	cmd.AddCommand(ecloudApplianceParameterListCmd(f))

	return cmd
}

func ecloudApplianceParameterListCmd(f factory.ClientFactory) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudApplianceParameterList(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudApplianceParameterList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	parameters, err := service.GetApplianceParameters(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving appliance parameters: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudApplianceParametersProvider(parameters))
}
