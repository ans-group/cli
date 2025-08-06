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

func ecloudNetworkRuleRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "networkrule",
		Short: "sub-commands relating to network rules",
	}

	// Child commands
	cmd.AddCommand(ecloudNetworkRuleListCmd(f))
	cmd.AddCommand(ecloudNetworkRuleShowCmd(f))
	cmd.AddCommand(ecloudNetworkRuleCreateCmd(f))
	cmd.AddCommand(ecloudNetworkRuleUpdateCmd(f))
	cmd.AddCommand(ecloudNetworkRuleDeleteCmd(f))

	// Child root commands
	cmd.AddCommand(ecloudNetworkRuleNetworkRulePortRootCmd(f))

	return cmd
}

func ecloudNetworkRuleListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists network rules",
		Long:    "This command lists network rules",
		Example: "ans ecloud networkrule list",
		RunE:    ecloudCobraRunEFunc(f, ecloudNetworkRuleList),
	}

	cmd.Flags().String("policy", "", "Network policy ID for filtering")

	return cmd
}

func ecloudNetworkRuleList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
		helper.NewStringFilterFlagOption("policy", "network_policy_id"),
	)
	if err != nil {
		return err
	}

	rules, err := service.GetNetworkRules(params)
	if err != nil {
		return fmt.Errorf("error retrieving network rules: %s", err)
	}

	return output.CommandOutput(cmd, NetworkRuleCollection(rules))
}

func ecloudNetworkRuleShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <rule: id>...",
		Short:   "Shows an network rule",
		Long:    "This command shows one or more network rules",
		Example: "ans ecloud networkrule show nr-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing network rule")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNetworkRuleShow),
	}
}

func ecloudNetworkRuleShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var rules []ecloud.NetworkRule
	for _, arg := range args {
		rule, err := service.GetNetworkRule(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving network rule [%s]: %s", arg, err)
			continue
		}

		rules = append(rules, rule)
	}

	return output.CommandOutput(cmd, NetworkRuleCollection(rules))
}

func ecloudNetworkRuleCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a network rule",
		Long:    "This command creates a network rule",
		Example: "ans ecloud networkrule create --policy np-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudNetworkRuleCreate),
	}

	// Setup flags
	cmd.Flags().String("policy", "", "ID of network policy")
	_ = cmd.MarkFlagRequired("policy")
	cmd.Flags().String("source", "", "Source of rule. IP range/subnet or ANY")
	_ = cmd.MarkFlagRequired("source")
	cmd.Flags().String("destination", "", "Destination of rule. IP range/subnet or ANY")
	_ = cmd.MarkFlagRequired("destination")
	cmd.Flags().String("direction", "", "Direction of rule. One of: IN/OUT/IN_OUT")
	_ = cmd.MarkFlagRequired("direction")
	cmd.Flags().String("action", "", "Action of rule. One of: ALLOW/DROP/REJECT")
	_ = cmd.MarkFlagRequired("action")
	cmd.Flags().String("name", "", "Name of rule")
	cmd.Flags().Int("sequence", 0, "Sequence for rule")
	cmd.Flags().Bool("enabled", false, "Specifies whether rule is enabled")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the network rule has been completely created")

	return cmd
}

func ecloudNetworkRuleCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateNetworkRuleRequest{}
	createRequest.NetworkPolicyID, _ = cmd.Flags().GetString("policy")
	createRequest.Source, _ = cmd.Flags().GetString("source")
	createRequest.Destination, _ = cmd.Flags().GetString("destination")
	createRequest.Enabled, _ = cmd.Flags().GetBool("enabled")

	direction, _ := cmd.Flags().GetString("direction")
	directionParsed, err := ecloud.NetworkRuleDirectionEnum.Parse(direction)
	if err != nil {
		return err
	}
	createRequest.Direction = directionParsed

	action, _ := cmd.Flags().GetString("action")
	actionParsed, err := ecloud.NetworkRuleActionEnum.Parse(action)
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

	taskRef, err := service.CreateNetworkRule(createRequest)
	if err != nil {
		return fmt.Errorf("error creating network rule: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("error waiting for network rule task to complete: %s", err)
		}
	}

	rule, err := service.GetNetworkRule(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("error retrieving new network rule: %s", err)
	}

	return output.CommandOutput(cmd, NetworkRuleCollection([]ecloud.NetworkRule{rule}))
}

func ecloudNetworkRuleUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <rule: id>...",
		Short:   "Updates a network rule",
		Long:    "This command updates one or more network rules",
		Example: "ans ecloud networkrule update np-abcdef12 --name \"my rule\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing network rule")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNetworkRuleUpdate),
	}

	cmd.Flags().String("source", "", "Source of rule. IP range/subnet or ANY")
	cmd.Flags().String("destination", "", "Destination of rule. IP range/subnet or ANY")
	cmd.Flags().String("direction", "", "Direction of rule. One of: IN/OUT/IN_OUT")
	cmd.Flags().String("action", "", "Action of rule. One of: ALLOW/DROP/REJECT")
	cmd.Flags().String("name", "", "Name of rule")
	cmd.Flags().Int("sequence", 0, "Sequence for rule")
	cmd.Flags().Bool("enabled", false, "Specifies whether rule is enabled")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the network rule has been completely updated")

	return cmd
}

func ecloudNetworkRuleUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchNetworkRuleRequest{}

	if cmd.Flags().Changed("source") {
		patchRequest.Source, _ = cmd.Flags().GetString("source")
	}

	if cmd.Flags().Changed("destination") {
		patchRequest.Destination, _ = cmd.Flags().GetString("destination")
	}

	if cmd.Flags().Changed("direction") {
		direction, _ := cmd.Flags().GetString("direction")
		directionParsed, err := ecloud.NetworkRuleDirectionEnum.Parse(direction)
		if err != nil {
			return err
		}
		patchRequest.Direction = directionParsed

	}

	if cmd.Flags().Changed("action") {
		action, _ := cmd.Flags().GetString("action")
		actionParsed, err := ecloud.NetworkRuleActionEnum.Parse(action)
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

	var rules []ecloud.NetworkRule
	for _, arg := range args {
		task, err := service.PatchNetworkRule(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating network rule [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for network rule [%s]: %s", arg, err)
				continue
			}
		}

		rule, err := service.GetNetworkRule(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated network rule [%s]: %s", arg, err)
			continue
		}

		rules = append(rules, rule)
	}

	return output.CommandOutput(cmd, NetworkRuleCollection(rules))
}

func ecloudNetworkRuleDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <rule: id>...",
		Short:   "Removes a network rule",
		Long:    "This command removes one or more network rules",
		Example: "ans ecloud networkrule delete nr-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing network rule")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNetworkRuleDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the network rule has been completely removed")

	return cmd
}

func ecloudNetworkRuleDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteNetworkRule(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing network rule [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for network rule [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
