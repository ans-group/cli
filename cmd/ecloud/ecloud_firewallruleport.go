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

func ecloudFirewallRulePortRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "firewallruleport",
		Short: "sub-commands relating to firewall rule ports",
	}

	// Child commands
	cmd.AddCommand(ecloudFirewallRulePortListCmd(f))
	cmd.AddCommand(ecloudFirewallRulePortShowCmd(f))
	cmd.AddCommand(ecloudFirewallRulePortCreateCmd(f))
	cmd.AddCommand(ecloudFirewallRulePortUpdateCmd(f))
	cmd.AddCommand(ecloudFirewallRulePortDeleteCmd(f))

	return cmd
}

func ecloudFirewallRulePortListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists firewall rule ports",
		Long:    "This command lists firewall rule ports",
		Example: "ukfast ecloud firewallruleport list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudFirewallRulePortList(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("rule", "", "Firewall rule ID for filtering")

	return cmd
}

func ecloudFirewallRulePortList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
		helper.NewStringFilterFlagOption("rule", "firewall_rule_id"),
	)
	if err != nil {
		return err
	}

	rules, err := service.GetFirewallRulePorts(params)
	if err != nil {
		return fmt.Errorf("Error retrieving firewall rule ports: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudFirewallRulePortsProvider(rules))
}

func ecloudFirewallRulePortShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <port: id>...",
		Short:   "Shows a firewall rule port",
		Long:    "This command shows one or more firewall rule ports",
		Example: "ukfast ecloud firewallruleport show fwrp-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing firewall rule port")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudFirewallRulePortShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudFirewallRulePortShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var rules []ecloud.FirewallRulePort
	for _, arg := range args {
		rule, err := service.GetFirewallRulePort(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving firewall rule port [%s]: %s", arg, err)
			continue
		}

		rules = append(rules, rule)
	}

	return output.CommandOutput(cmd, OutputECloudFirewallRulePortsProvider(rules))
}

func ecloudFirewallRulePortCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a firewall rule port",
		Long:    "This command creates a firewall rule port",
		Example: "ukfast ecloud firewallruleport create --rule fwr-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudFirewallRulePortCreate),
	}

	// Setup flags
	cmd.Flags().String("rule", "", "ID of firewall rule")
	cmd.MarkFlagRequired("rule")
	cmd.Flags().String("source", "", "Source port. Single port, port range, or ANY")
	cmd.Flags().String("destination", "", "Destination port. Single port, port range, or ANY")
	cmd.Flags().String("protocol", "", "Protocol of port. One of: TCP/UDP/ICMPv4")
	cmd.MarkFlagRequired("protocol")
	cmd.Flags().String("name", "", "Name of port")

	return cmd
}

func ecloudFirewallRulePortCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateFirewallRulePortRequest{}
	createRequest.FirewallRuleID, _ = cmd.Flags().GetString("rule")
	createRequest.Source, _ = cmd.Flags().GetString("source")
	createRequest.Destination, _ = cmd.Flags().GetString("destination")

	protocol, _ := cmd.Flags().GetString("protocol")
	protocolParsed, err := ecloud.ParseFirewallRulePortProtocol(protocol)
	if err != nil {
		return err
	}
	createRequest.Protocol = protocolParsed

	if cmd.Flags().Changed("name") {
		createRequest.Name, _ = cmd.Flags().GetString("name")
	}

	ruleID, err := service.CreateFirewallRulePort(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating firewall rule port: %s", err)
	}

	rule, err := service.GetFirewallRulePort(ruleID)
	if err != nil {
		return fmt.Errorf("Error retrieving new firewall rule port: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudFirewallRulePortsProvider([]ecloud.FirewallRulePort{rule}))
}

func ecloudFirewallRulePortUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <port: id>...",
		Short:   "Updates a firewall rule port",
		Long:    "This command updates one or more firewall rule ports",
		Example: "ukfast ecloud firewallruleport update fwrp-abcdef12 --name \"my port\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing firewall rule port")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFirewallRulePortUpdate),
	}

	cmd.Flags().String("source", "", "Source port. Single port, port range, or ANY")
	cmd.Flags().String("destination", "", "Destination port. Single port, port range, or ANY")
	cmd.Flags().String("protocol", "", "Protocol of port. One of: TCP/UDP/ICMPv4")
	cmd.Flags().String("name", "", "Name of port")

	return cmd
}

func ecloudFirewallRulePortUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchFirewallRulePortRequest{}

	if cmd.Flags().Changed("source") {
		patchRequest.Source, _ = cmd.Flags().GetString("source")
	}

	if cmd.Flags().Changed("destination") {
		patchRequest.Destination, _ = cmd.Flags().GetString("destination")
	}

	if cmd.Flags().Changed("protocol") {

		protocol, _ := cmd.Flags().GetString("protocol")
		protocolParsed, err := ecloud.ParseFirewallRulePortProtocol(protocol)
		if err != nil {
			return err
		}
		patchRequest.Protocol = protocolParsed
	}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var rules []ecloud.FirewallRulePort
	for _, arg := range args {
		err := service.PatchFirewallRulePort(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating firewall rule port [%s]: %s", arg, err)
			continue
		}

		rule, err := service.GetFirewallRulePort(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated firewall rule port [%s]: %s", arg, err)
			continue
		}

		rules = append(rules, rule)
	}

	return output.CommandOutput(cmd, OutputECloudFirewallRulePortsProvider(rules))
}

func ecloudFirewallRulePortDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <port: id>...",
		Short:   "Removes a firewall rule port",
		Long:    "This command removes one or more firewall rule ports",
		Example: "ukfast ecloud firewallruleport delete fwrp-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing firewall rule port")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFirewallRulePortDelete),
	}
}

func ecloudFirewallRulePortDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		err := service.DeleteFirewallRulePort(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing firewall rule port [%s]: %s", arg, err)
		}
	}
	return nil
}
