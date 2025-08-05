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

func ecloudNICRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nic",
		Short: "sub-commands relating to NICs",
	}

	// Child commands
	cmd.AddCommand(ecloudNICListCmd(f))
	cmd.AddCommand(ecloudNICShowCmd(f))

	// Child root commands
	cmd.AddCommand(ecloudNICIPAddressRootCmd(f))

	return cmd
}

func ecloudNICListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists NICs",
		Long:    "This command lists NICs",
		Example: "ans ecloud nic list",
		RunE:    ecloudCobraRunEFunc(f, ecloudNICList),
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
		return fmt.Errorf("error retrieving NICs: %s", err)
	}

	return output.CommandOutput(cmd, NICCollection(nics))
}

func ecloudNICShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <nic: id>...",
		Short:   "Shows a NIC",
		Long:    "This command shows one or more NICs",
		Example: "ans ecloud nic show nic-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing nic")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNICShow),
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

	return output.CommandOutput(cmd, NICCollection(nics))
}
