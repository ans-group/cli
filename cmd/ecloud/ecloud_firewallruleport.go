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
		Example: "ans ecloud firewallruleport list",
		RunE:    ecloudCobraRunEFunc(f, ecloudFirewallRulePortList),
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

	return output.CommandOutput(cmd, FirewallRulePortCollection(rules))
}

func ecloudFirewallRulePortShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <port: id>...",
		Short:   "Shows a firewall rule port",
		Long:    "This command shows one or more firewall rule ports",
		Example: "ans ecloud firewallruleport show fwrp-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing firewall rule port")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFirewallRulePortShow),
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

	return output.CommandOutput(cmd, FirewallRulePortCollection(rules))
}

func ecloudFirewallRulePortCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a firewall rule port",
		Long:    "This command creates a firewall rule port",
		Example: "ans ecloud firewallruleport create --rule fwr-abcdef12",
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
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the firewall rule port has been completely created")

	return cmd
}

func ecloudFirewallRulePortCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateFirewallRulePortRequest{}
	createRequest.FirewallRuleID, _ = cmd.Flags().GetString("rule")
	createRequest.Source, _ = cmd.Flags().GetString("source")
	createRequest.Destination, _ = cmd.Flags().GetString("destination")

	protocol, _ := cmd.Flags().GetString("protocol")
	protocolParsed, err := ecloud.FirewallRulePortProtocolEnum.Parse(protocol)
	if err != nil {
		return err
	}
	createRequest.Protocol = protocolParsed

	if cmd.Flags().Changed("name") {
		createRequest.Name, _ = cmd.Flags().GetString("name")
	}

	taskRef, err := service.CreateFirewallRulePort(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating firewall rule port: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for firewall rule port task to complete: %s", err)
		}
	}

	rule, err := service.GetFirewallRulePort(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new firewall rule port: %s", err)
	}

	return output.CommandOutput(cmd, FirewallRulePortCollection([]ecloud.FirewallRulePort{rule}))
}

func ecloudFirewallRulePortUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <port: id>...",
		Short:   "Updates a firewall rule port",
		Long:    "This command updates one or more firewall rule ports",
		Example: "ans ecloud firewallruleport update fwrp-abcdef12 --name \"my port\"",
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
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the firewall rule port has been completely updated")

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
		protocolParsed, err := ecloud.FirewallRulePortProtocolEnum.Parse(protocol)
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
		task, err := service.PatchFirewallRulePort(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating firewall rule port [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for firewall rule port [%s]: %s", arg, err)
				continue
			}
		}

		rule, err := service.GetFirewallRulePort(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated firewall rule port [%s]: %s", arg, err)
			continue
		}

		rules = append(rules, rule)
	}

	return output.CommandOutput(cmd, FirewallRulePortCollection(rules))
}

func ecloudFirewallRulePortDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <port: id>...",
		Short:   "Removes a firewall rule port",
		Long:    "This command removes one or more firewall rule ports",
		Example: "ans ecloud firewallruleport delete fwrp-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing firewall rule port")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFirewallRulePortDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the firewall rule port has been completely removed")

	return cmd
}

func ecloudFirewallRulePortDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteFirewallRulePort(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing firewall rule port [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for firewall rule port [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
