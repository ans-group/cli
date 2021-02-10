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

func ecloudFirewallRuleFirewallRulePortRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "port",
		Short: "sub-commands relating to firewall rule ports",
	}

	// Child commands
	cmd.AddCommand(ecloudFirewallRuleFirewallRulePortListCmd(f))

	return cmd
}

func ecloudFirewallRuleFirewallRulePortListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists ports for firewall rule",
		Long:    "This command lists ports for firewall rule",
		Example: "ukfast ecloud firewallrule firewallport list fwp-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing firewall rule")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFirewallRuleFirewallRulePortList),
	}

	cmd.Flags().String("name", "", "Firewall rule port name for filtering")

	return cmd
}

func ecloudFirewallRuleFirewallRulePortList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	ports, err := service.GetFirewallRuleFirewallRulePorts(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving firewall rule ports: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudFirewallRulePortsProvider(ports))
}
