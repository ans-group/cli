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

func ecloudNetworkRuleNetworkRulePortRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "port",
		Short: "sub-commands relating to network rule ports",
	}

	// Child commands
	cmd.AddCommand(ecloudNetworkRuleNetworkRulePortListCmd(f))

	return cmd
}

func ecloudNetworkRuleNetworkRulePortListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists ports for network rule",
		Long:    "This command lists ports for network rule",
		Example: "ans ecloud networkrule networkport list np-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing network rule")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNetworkRuleNetworkRulePortList),
	}

	cmd.Flags().String("name", "", "Network rule port name for filtering")

	return cmd
}

func ecloudNetworkRuleNetworkRulePortList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	ports, err := service.GetNetworkRuleNetworkRulePorts(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving network rule ports: %s", err)
	}

	return output.CommandOutput(cmd, NetworkRulePortCollection(ports))
}
