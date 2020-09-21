package ecloud_v2

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/ptr"
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
	cmd.AddCommand(ecloudNetworkCreateCmd(f))
	cmd.AddCommand(ecloudNetworkUpdateCmd(f))
	cmd.AddCommand(ecloudNetworkDeleteCmd(f))

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

	helper.HydrateAPIRequestParametersWithStringFilterFlag(&params, cmd, helper.NewStringFilterFlag("name", "name"))

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
		Example: "ukfast ecloud network show net-abcdef12",
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

func ecloudNetworkCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a network",
		Long:    "This command creates a network",
		Example: "ukfast ecloud network create",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudNetworkCreate(c.ECloudService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of network")
	cmd.Flags().String("router", "", "ID of router")
	cmd.MarkFlagRequired("router")

	return cmd
}

func ecloudNetworkCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {

	name, _ := cmd.Flags().GetString("name")
	routerID, _ := cmd.Flags().GetString("router")

	createRequest := ecloud.CreateNetworkRequest{
		Name:     ptr.String(name),
		RouterID: routerID,
	}

	networkID, err := service.CreateNetwork(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating network: %s", err)
	}

	network, err := service.GetNetwork(networkID)
	if err != nil {
		return fmt.Errorf("Error retrieving new network: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudNetworksProvider([]ecloud.Network{network}))
}

func ecloudNetworkUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <network: name>...",
		Short:   "Updates a network",
		Long:    "This command updates one or more networks",
		Example: "ukfast ecloud network update net-abcdef12 --name \"my network\"",
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

			return ecloudNetworkUpdate(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Name of network")

	return cmd
}

func ecloudNetworkUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchNetworkRequest{}

	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		patchRequest.Name = ptr.String(name)
	}

	var networks []ecloud.Network
	for _, arg := range args {
		err := service.PatchNetwork(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating network [%s]: %s", arg, err)
			continue
		}

		network, err := service.GetNetwork(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated network [%s]: %s", arg, err)
			continue
		}

		networks = append(networks, network)
	}

	return output.CommandOutput(cmd, OutputECloudNetworksProvider(networks))
}

func ecloudNetworkDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <network: name...>",
		Short:   "Removes a network",
		Long:    "This command removes one or more networks",
		Example: "ukfast ecloud network delete net-abcdef12",
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

			ecloudNetworkDelete(c.ECloudService(), cmd, args)
			return nil
		},
	}
}

func ecloudNetworkDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	for _, arg := range args {
		err := service.DeleteNetwork(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing network [%s]: %s", arg, err)
		}
	}
}
