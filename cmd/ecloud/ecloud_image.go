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
	cmd.AddCommand(ecloudImageUpdateCmd(f))
	cmd.AddCommand(ecloudImageDeleteCmd(f))

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
		Example: "ans ecloud image list",
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
		return fmt.Errorf("error retrieving images: %s", err)
	}

	return output.CommandOutput(cmd, ImageCollection(images))
}

func ecloudImageShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <image: id>...",
		Short:   "Shows a image",
		Long:    "This command shows one or more images",
		Example: "ans ecloud vm image img-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing image")
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

	return output.CommandOutput(cmd, ImageCollection(images))
}

func ecloudImageUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <image: id>...",
		Short:   "Updates a image",
		Long:    "This command updates one or more images",
		Example: "ans ecloud image update img-abcdef12 --name \"my image\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing image")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudImageUpdate),
	}

	cmd.Flags().String("name", "", "Name of image")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the image has been completely updated")

	return cmd
}

func ecloudImageUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.UpdateImageRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var images []ecloud.Image
	for _, arg := range args {
		task, err := service.UpdateImage(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating image [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for image [%s]: %s", arg, err)
				continue
			}
		}

		image, err := service.GetImage(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated image [%s]: %s", arg, err)
			continue
		}

		images = append(images, image)
	}

	return output.CommandOutput(cmd, ImageCollection(images))
}

func ecloudImageDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <image: id>...",
		Short:   "Removes a image",
		Long:    "This command removes one or more images",
		Example: "ans ecloud image delete img-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing image")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudImageDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the image has been completely removed")

	return cmd
}

func ecloudImageDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteImage(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing image [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for image [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
