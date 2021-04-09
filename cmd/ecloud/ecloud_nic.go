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

func ecloudNICRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nic",
		Short: "sub-commands relating to NICs",
	}

	// Child commands
	cmd.AddCommand(ecloudNICListCmd(f))
	cmd.AddCommand(ecloudNICShowCmd(f))

	return cmd
}

func ecloudNICListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists NICs",
		Long:    "This command lists NICs",
		Example: "ukfast ecloud nic list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudNICList(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "NIC name for filtering")

	return cmd
}

func ecloudNICList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	nics, err := service.GetNICs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving NICs: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudNICsProvider(nics))
}

func ecloudNICShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <nic: id>...",
		Short:   "Shows a NIC",
		Long:    "This command shows one or more NICs",
		Example: "ukfast ecloud nic show nic-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing nic")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudNICShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudNICShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var nics []ecloud.NIC
	for _, arg := range args {
		nic, err := service.GetNIC(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving NIC [%s]: %s", arg, err)
			continue
		}

		nics = append(nics, nic)
	}

	return output.CommandOutput(cmd, OutputECloudNICsProvider(nics))
}