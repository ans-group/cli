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

func ecloudImageParameterRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parameter",
		Short: "sub-commands relating to image parameters",
	}

	// Child commands
	cmd.AddCommand(ecloudImageParameterListCmd(f))

	return cmd
}

func ecloudImageParameterListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists image parameters",
		Long:    "This command lists image parameters",
		Example: "ans ecloud image parameter list img-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing image")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudImageParameterList),
	}
}

func ecloudImageParameterList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	parameters, err := service.GetImageParameters(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving image parameters: %s", err)
	}

	return output.CommandOutput(cmd, ImageParameterCollection(parameters))
}
