package cmd

import (
	"errors"

	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"

	"github.com/spf13/cobra"
)

func ecloudTemplateDiskRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disk",
		Short: "sub-commands relating to template disks",
	}

	// Child commands
	cmd.AddCommand(ecloudTemplateDiskListCmd())

	return cmd
}

func ecloudTemplateDiskListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <template: id>",
		Short:   "lists template disks",
		Long:    "This command lists template disks",
		Example: "ukfast ecloud template disk list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing template")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudTemplateDiskList(getClient().ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Disk name for filtering")

	return cmd
}

func ecloudTemplateDiskList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	template, err := service.GetTemplate(args[0])
	if err != nil {
		output.Fatalf("Error retrieving template [%s]: %s", args[0], err)
		return
	}

	outputECloudVirtualMachineDisks(template.Disks)
}
