package ecloud

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

func ecloudFirewallRuleRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "firewallrule",
		Short: "sub-commands relating to firewall rules",
	}

	// Child commands
	cmd.AddCommand(ecloudFirewallRuleListCmd(f))
	cmd.AddCommand(ecloudFirewallRuleShowCmd(f))
	cmd.AddCommand(ecloudFirewallRuleCreateCmd(f))
	cmd.AddCommand(ecloudFirewallRuleUpdateCmd(f))
	cmd.AddCommand(ecloudFirewallRuleDeleteCmd(f))

	// Child root commands
	cmd.AddCommand(ecloudFirewallRuleFirewallRulePortRootCmd(f))

	return cmd
}

func ecloudFirewallRuleListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists firewall rules",
		Long:    "This command lists firewall rules",
		Example: "ukfast ecloud firewallrule list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudFirewallRuleList(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("policy", "", "Firewall policy ID for filtering")

	return cmd
}

func ecloudFirewallRuleList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
		helper.NewStringFilterFlagOption("policy", "firewall_policy_id"),
	)
	if err != nil {
		return err
	}

	rules, err := service.GetFirewallRules(params)
	if err != nil {
		return fmt.Errorf("Error retrieving firewall rules: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudFirewallRulesProvider(rules))
}

func ecloudFirewallRuleShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <rule: id>...",
		Short:   "Shows an firewall rule",
		Long:    "This command shows one or more firewall rules",
		Example: "ukfast ecloud firewallrule show fwr-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing firewall rule")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudFirewallRuleShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudFirewallRuleShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var rules []ecloud.FirewallRule
	for _, arg := range args {
		rule, err := service.GetFirewallRule(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving firewall rule [%s]: %s", arg, err)
			continue
		}

		rules = append(rules, rule)
	}

	return output.CommandOutput(cmd, OutputECloudFirewallRulesProvider(rules))
}

func ecloudFirewallRuleCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a firewall rule",
		Long:    "This command creates a firewall rule",
		Example: "ukfast ecloud firewallrule create --policy fwp-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudFirewallRuleCreate),
	}

	// Setup flags
	cmd.Flags().String("policy", "", "ID of firewall policy")
	cmd.MarkFlagRequired("policy")
	cmd.Flags().String("source", "", "Source of rule. IP range/subnet or ANY")
	cmd.MarkFlagRequired("source")
	cmd.Flags().String("destination", "", "Destination of rule. IP range/subnet or ANY")
	cmd.MarkFlagRequired("destination")
	cmd.Flags().String("direction", "", "Direction of rule. One of: IN/OUT/IN_OUT")
	cmd.MarkFlagRequired("direction")
	cmd.Flags().String("action", "", "Action of rule. One of: ALLOW/DROP/REJECT")
	cmd.MarkFlagRequired("action")
	cmd.Flags().String("name", "", "Name of rule")
	cmd.Flags().Int("sequence", 0, "Sequence for rule")
	cmd.Flags().Bool("enabled", false, "Specifies whether rule is enabled")

	return cmd
}

func ecloudFirewallRuleCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateFirewallRuleRequest{}
	createRequest.FirewallPolicyID, _ = cmd.Flags().GetString("policy")
	createRequest.Source, _ = cmd.Flags().GetString("source")
	createRequest.Destination, _ = cmd.Flags().GetString("destination")
	createRequest.Enabled, _ = cmd.Flags().GetBool("enabled")

	direction, _ := cmd.Flags().GetString("direction")
	directionParsed, err := ecloud.ParseFirewallRuleDirection(direction)
	if err != nil {
		return err
	}
	createRequest.Direction = directionParsed

	action, _ := cmd.Flags().GetString("action")
	actionParsed, err := ecloud.ParseFirewallRuleAction(action)
	if err != nil {
		return err
	}
	createRequest.Action = actionParsed

	if cmd.Flags().Changed("name") {
		createRequest.Name, _ = cmd.Flags().GetString("name")
	}
	if cmd.Flags().Changed("sequence") {
		createRequest.Sequence, _ = cmd.Flags().GetInt("sequence")
	}

	ruleID, err := service.CreateFirewallRule(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating firewall rule: %s", err)
	}

	rule, err := service.GetFirewallRule(ruleID)
	if err != nil {
		return fmt.Errorf("Error retrieving new firewall rule: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudFirewallRulesProvider([]ecloud.FirewallRule{rule}))
}

func ecloudFirewallRuleUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <rule: id>...",
		Short:   "Updates a firewall rule",
		Long:    "This command updates one or more firewall rules",
		Example: "ukfast ecloud firewallrule update fwp-abcdef12 --name \"my rule\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing firewall rule")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFirewallRuleUpdate),
	}

	cmd.Flags().String("source", "", "Source of rule. IP range/subnet or ANY")
	cmd.Flags().String("destination", "", "Destination of rule. IP range/subnet or ANY")
	cmd.Flags().String("direction", "", "Direction of rule. One of: IN/OUT/IN_OUT")
	cmd.Flags().String("action", "", "Action of rule. One of: ALLOW/DROP/REJECT")
	cmd.Flags().String("name", "", "Name of rule")
	cmd.Flags().Int("sequence", 0, "Sequence for rule")
	cmd.Flags().Bool("enabled", false, "Specifies whether rule is enabled")

	return cmd
}

func ecloudFirewallRuleUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchFirewallRuleRequest{}

	if cmd.Flags().Changed("source") {
		patchRequest.Source, _ = cmd.Flags().GetString("source")
	}

	if cmd.Flags().Changed("destination") {
		patchRequest.Destination, _ = cmd.Flags().GetString("destination")
	}

	if cmd.Flags().Changed("direction") {
		direction, _ := cmd.Flags().GetString("direction")
		directionParsed, err := ecloud.ParseFirewallRuleDirection(direction)
		if err != nil {
			return err
		}
		patchRequest.Direction = directionParsed

	}

	if cmd.Flags().Changed("action") {
		action, _ := cmd.Flags().GetString("action")
		actionParsed, err := ecloud.ParseFirewallRuleAction(action)
		if err != nil {
			return err
		}
		patchRequest.Action = actionParsed
	}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	if cmd.Flags().Changed("sequence") {
		patchRequest.Name, _ = cmd.Flags().GetString("sequence")
	}

	if cmd.Flags().Changed("enabled") {
		enabled, _ := cmd.Flags().GetBool("enabled")
		patchRequest.Enabled = ptr.Bool(enabled)
	}

	var rules []ecloud.FirewallRule
	for _, arg := range args {
		err := service.PatchFirewallRule(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating firewall rule [%s]: %s", arg, err)
			continue
		}

		rule, err := service.GetFirewallRule(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated firewall rule [%s]: %s", arg, err)
			continue
		}

		rules = append(rules, rule)
	}

	return output.CommandOutput(cmd, OutputECloudFirewallRulesProvider(rules))
}

func ecloudFirewallRuleDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <rule: id>...",
		Short:   "Removes a firewall rule",
		Long:    "This command removes one or more firewall rules",
		Example: "ukfast ecloud firewallrule delete fwr-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing firewall rule")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFirewallRuleDelete),
	}
}

func ecloudFirewallRuleDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		err := service.DeleteFirewallRule(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing firewall rule [%s]: %s", arg, err)
		}
	}
	return nil
}
