package ecloud

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/pkg/browser"

	"github.com/spf13/cobra"
)

func ecloudVirtualMachineConsoleRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "consolesession",
		Short: "sub-commands relating to virtual machine Consoles",
	}

	// Child commands
	cmd.AddCommand(ecloudVirtualMachineConsoleSessionCreateCmd(f))

	return cmd
}

func ecloudVirtualMachineConsoleSessionCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <virtualmachine: id>",
		Short:   "Creates a virtual machine console session",
		Long:    "This command creates a virtual machine console session",
		Example: "ans ecloud vm consolesession create 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing virtual machine")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudVirtualMachineConsoleSessionCreate(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().Bool("browser", false, "Indicates session should be opened in default browser")

	return cmd
}

func ecloudVirtualMachineConsoleSessionCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	vmID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid virtual machine ID [%s]", args[0])
	}

	console, err := service.CreateVirtualMachineConsoleSession(vmID)
	if err != nil {
		return fmt.Errorf("Error creating virtual machine console session: %s", err)
	}

	openBrowser, _ := cmd.Flags().GetBool("browser")
	if openBrowser {
		return browser.OpenURL(console.URL)
	}

	return output.CommandOutput(cmd, ConsoleSessionCollection([]ecloud.ConsoleSession{console}))
}
