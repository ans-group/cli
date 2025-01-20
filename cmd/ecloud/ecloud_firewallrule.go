package ecloud

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/ptr"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/spf13/cobra"
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
		Example: "ans ecloud firewallrule list",
		RunE:    ecloudCobraRunEFunc(f, ecloudFirewallRuleList),
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

	return output.CommandOutput(cmd, FirewallRuleCollection(rules))
}

func ecloudFirewallRuleShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <rule: id>...",
		Short:   "Shows an firewall rule",
		Long:    "This command shows one or more firewall rules",
		Example: "ans ecloud firewallrule show fwr-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing firewall rule")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFirewallRuleShow),
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

	return output.CommandOutput(cmd, FirewallRuleCollection(rules))
}

func ecloudFirewallRuleCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a firewall rule",
		Long:    "This command creates a firewall rule",
		Example: "ans ecloud firewallrule create --policy fwp-abcdef12",
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
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the firewall rule has been completely created")

	return cmd
}

func ecloudFirewallRuleCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateFirewallRuleRequest{}
	createRequest.FirewallPolicyID, _ = cmd.Flags().GetString("policy")
	createRequest.Source, _ = cmd.Flags().GetString("source")
	createRequest.Destination, _ = cmd.Flags().GetString("destination")
	createRequest.Enabled, _ = cmd.Flags().GetBool("enabled")

	direction, _ := cmd.Flags().GetString("direction")
	directionParsed, err := ecloud.FirewallRuleDirectionEnum.Parse(direction)
	if err != nil {
		return err
	}
	createRequest.Direction = directionParsed

	action, _ := cmd.Flags().GetString("action")
	actionParsed, err := ecloud.FirewallRuleActionEnum.Parse(action)
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

	taskRef, err := service.CreateFirewallRule(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating firewall rule: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for firewall rule task to complete: %s", err)
		}
	}

	rule, err := service.GetFirewallRule(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new firewall rule: %s", err)
	}

	return output.CommandOutput(cmd, FirewallRuleCollection([]ecloud.FirewallRule{rule}))
}

func ecloudFirewallRuleUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <rule: id>...",
		Short:   "Updates a firewall rule",
		Long:    "This command updates one or more firewall rules",
		Example: "ans ecloud firewallrule update fwp-abcdef12 --name \"my rule\"",
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
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the firewall rule has been completely updated")

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
		directionParsed, err := ecloud.FirewallRuleDirectionEnum.Parse(direction)
		if err != nil {
			return err
		}
		patchRequest.Direction = directionParsed

	}

	if cmd.Flags().Changed("action") {
		action, _ := cmd.Flags().GetString("action")
		actionParsed, err := ecloud.FirewallRuleActionEnum.Parse(action)
		if err != nil {
			return err
		}
		patchRequest.Action = actionParsed
	}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	if cmd.Flags().Changed("sequence") {
		sequence, _ := cmd.Flags().GetInt("sequence")
		patchRequest.Sequence = ptr.Int(sequence)
	}

	if cmd.Flags().Changed("enabled") {
		enabled, _ := cmd.Flags().GetBool("enabled")
		patchRequest.Enabled = ptr.Bool(enabled)
	}

	var rules []ecloud.FirewallRule
	for _, arg := range args {
		task, err := service.PatchFirewallRule(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating firewall rule [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for firewall rule [%s]: %s", arg, err)
				continue
			}
		}

		rule, err := service.GetFirewallRule(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated firewall rule [%s]: %s", arg, err)
			continue
		}

		rules = append(rules, rule)
	}

	return output.CommandOutput(cmd, FirewallRuleCollection(rules))
}

func ecloudFirewallRuleDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <rule: id>...",
		Short:   "Removes a firewall rule",
		Long:    "This command removes one or more firewall rules",
		Example: "ans ecloud firewallrule delete fwr-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing firewall rule")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFirewallRuleDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the firewall rule has been completely removed")

	return cmd
}

func ecloudFirewallRuleDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteFirewallRule(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing firewall rule [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for firewall rule [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
