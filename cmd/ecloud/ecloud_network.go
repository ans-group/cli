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

	// Child root commands
	cmd.AddCommand(ecloudNetworkNICRootCmd(f))

	return cmd
}

func ecloudNetworkListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists networks",
		Long:    "This command lists networks",
		Example: "ans ecloud network list",
		RunE:    ecloudCobraRunEFunc(f, ecloudNetworkList),
	}

	cmd.Flags().String("name", "", "Network name for filtering")
	cmd.Flags().String("router", "", "Router ID for filtering")

	return cmd
}

func ecloudNetworkList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
		helper.NewStringFilterFlagOption("router", "router_id"),
	)
	if err != nil {
		return err
	}

	networks, err := service.GetNetworks(params)
	if err != nil {
		return fmt.Errorf("Error retrieving networks: %s", err)
	}

	return output.CommandOutput(cmd, NetworkCollection(networks))
}

func ecloudNetworkShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <network: id>...",
		Short:   "Shows a network",
		Long:    "This command shows one or more networks",
		Example: "ans ecloud network show net-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing network")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNetworkShow),
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

	return output.CommandOutput(cmd, NetworkCollection(networks))
}

func ecloudNetworkCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates an network",
		Long:    "This command creates an network",
		Example: "ans ecloud network create --router rtr-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudNetworkCreate),
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of network")
	cmd.Flags().String("router", "", "ID of router")
	cmd.MarkFlagRequired("router")
	cmd.Flags().String("subnet", "", "Subnet for network, e.g. 10.0.0.0/24")
	cmd.MarkFlagRequired("subnet")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the network has been completely created")

	return cmd
}

func ecloudNetworkCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateNetworkRequest{}
	if cmd.Flags().Changed("name") {
		createRequest.Name, _ = cmd.Flags().GetString("name")
	}
	createRequest.RouterID, _ = cmd.Flags().GetString("router")
	createRequest.Subnet, _ = cmd.Flags().GetString("subnet")

	networkID, err := service.CreateNetwork(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating network: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(NetworkResourceSyncStatusWaitFunc(service, networkID, ecloud.SyncStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for network sync: %s", err)
		}
	}

	network, err := service.GetNetwork(networkID)
	if err != nil {
		return fmt.Errorf("Error retrieving new network: %s", err)
	}

	return output.CommandOutput(cmd, NetworkCollection([]ecloud.Network{network}))
}

func ecloudNetworkUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <network: id>...",
		Short:   "Updates an network",
		Long:    "This command updates one or more networks",
		Example: "ans ecloud network update net-abcdef12 --name \"my network\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing network")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNetworkUpdate),
	}

	cmd.Flags().String("name", "", "Name of network")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the network has been completely updated")

	return cmd
}

func ecloudNetworkUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchNetworkRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var networks []ecloud.Network
	for _, arg := range args {
		err := service.PatchNetwork(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating network [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(NetworkResourceSyncStatusWaitFunc(service, arg, ecloud.SyncStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for network [%s] sync: %s", arg, err)
				continue
			}
		}

		network, err := service.GetNetwork(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated network [%s]: %s", arg, err)
			continue
		}

		networks = append(networks, network)
	}

	return output.CommandOutput(cmd, NetworkCollection(networks))
}

func ecloudNetworkDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <network: id>...",
		Short:   "Removes an network",
		Long:    "This command removes one or more networks",
		Example: "ans ecloud network delete net-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing network")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNetworkDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the network has been completely removed")

	return cmd
}

func ecloudNetworkDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		err := service.DeleteNetwork(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing network [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(NetworkNotFoundWaitFunc(service, arg))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for removal of network [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}

func NetworkResourceSyncStatusWaitFunc(service ecloud.ECloudService, networkID string, status ecloud.SyncStatus) helper.WaitFunc {
	return ResourceSyncStatusWaitFunc(func() (ecloud.SyncStatus, error) {
		network, err := service.GetNetwork(networkID)
		if err != nil {
			return "", err
		}
		return network.Sync.Status, nil
	}, status)
}

func NetworkNotFoundWaitFunc(service ecloud.ECloudService, networkID string) helper.WaitFunc {
	return func() (finished bool, err error) {
		_, err = service.GetNetwork(networkID)
		if err != nil {
			switch err.(type) {
			case *ecloud.NetworkNotFoundError:
				return true, nil
			default:
				return false, fmt.Errorf("Failed to retrieve network [%s]: %s", networkID, err)
			}
		}

		return false, nil
	}
}
