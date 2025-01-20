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

func ecloudInstanceImageRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "sub-commands relating to instance images",
	}

	// Child commands
	cmd.AddCommand(ecloudInstanceImageCreateCmd(f))

	return cmd
}

func ecloudInstanceImageCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <instance: id>",
		Short:   "Creates an instance image",
		Long:    "This command creates an instance image",
		Example: "ans ecloud instance image create i-abcdef12 --name \"someimage\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceImageCreate),
	}

	cmd.Flags().String("name", "", "Name of image")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the instance restart task has been completed")

	return cmd
}

func ecloudInstanceImageCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateInstanceImageRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")

	taskRef, err := service.CreateInstanceImage(args[0], createRequest)
	if err != nil {
		return fmt.Errorf("Error creating instance image: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for task to complete: %s", err)
		}
	}

	image, err := service.GetImage(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new instance image: %s", err)
	}

	return output.CommandOutput(cmd, ImageCollection([]ecloud.Image{image}))
}
