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

func ecloudNetworkNICRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nic",
		Short: "sub-commands relating to network NICs",
	}

	// Child commands
	cmd.AddCommand(ecloudNetworkNICListCmd(f))

	return cmd
}

func ecloudNetworkNICListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists network nics",
		Long:    "This command lists network nics",
		Example: "ans ecloud network nic list net-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing network")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNetworkNICList),
	}

	cmd.Flags().String("name", "", "NIC name for filtering")

	return cmd
}

func ecloudNetworkNICList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	nics, err := service.GetNetworkNICs(args[0], params)
	if err != nil {
		return fmt.Errorf("error retrieving network NICs: %s", err)
	}

	return output.CommandOutput(cmd, NICCollection(nics))
}
