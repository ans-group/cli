package ecloud_v2

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudNetworkRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "sub-commands relating to networks",
	}

	// Child commands
	cmd.AddCommand(ecloudNetworkListCmd(f))
	cmd.AddCommand(ecloudNetworkShowCmd(f))

	return cmd
}

func ecloudNetworkListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists networks",
		Long:    "This command lists networks",
		Example: "ukfast ecloud network list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudNetworkList(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Network name for filtering")

	return cmd
}

func ecloudNetworkList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	helper.HydrateAPIRequestParametersWithStringFilterFlag(&params, cmd, "name", "name")

	networks, err := service.GetNetworks(params)
	if err != nil {
		return fmt.Errorf("Error retrieving networks: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudNetworksProvider(networks))
}

func ecloudNetworkShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <network: id>...",
		Short:   "Shows a network",
		Long:    "This command shows one or more networks",
		Example: "ukfast ecloud network show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing network")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudNetworkShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudNetworkShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var networks []ecloud.Network
	for _, arg := range args {
		network, err := service.GetNetwork(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving network [%s]: %s", arg, err)
			continue
		}

		networks = append(networks, network)
	}

	return output.CommandOutput(cmd, OutputECloudNetworksProvider(networks))
}
