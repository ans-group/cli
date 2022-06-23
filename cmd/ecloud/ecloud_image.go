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

func ecloudImageRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "sub-commands relating to images",
	}

	// Child commands
	cmd.AddCommand(ecloudImageListCmd(f))
	cmd.AddCommand(ecloudImageShowCmd(f))

	// Child root commands
	cmd.AddCommand(ecloudImageParameterRootCmd(f))
	cmd.AddCommand(ecloudImageMetadataRootCmd(f))

	return cmd
}

func ecloudImageListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists images",
		Long:    "This command lists images",
		Example: "ukfast ecloud image list",
		RunE:    ecloudCobraRunEFunc(f, ecloudImageList),
	}
}

func ecloudImageList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	images, err := service.GetImages(params)
	if err != nil {
		return fmt.Errorf("Error retrieving images: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudImagesProvider(images))
}

func ecloudImageShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <image: id>...",
		Short:   "Shows a image",
		Long:    "This command shows one or more images",
		Example: "ukfast ecloud vm image img-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing image")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudImageShow),
	}
}

func ecloudImageShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var images []ecloud.Image
	for _, arg := range args {
		image, err := service.GetImage(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving image [%s]: %s", arg, err)
			continue
		}

		images = append(images, image)
	}

	return output.CommandOutput(cmd, OutputECloudImagesProvider(images))
}
