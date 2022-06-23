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

func ecloudImageMetadataRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metadata",
		Short: "sub-commands relating to image metadata",
	}

	// Child commands
	cmd.AddCommand(ecloudImageMetadataListCmd(f))

	return cmd
}

func ecloudImageMetadataListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists image metadata",
		Long:    "This command lists image metadata",
		Example: "ukfast ecloud image metadata list img-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing image")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudImageMetadataList),
	}
}

func ecloudImageMetadataList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	metadata, err := service.GetImageMetadata(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving image metadata: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudImageMetadataProvider(metadata))
}
