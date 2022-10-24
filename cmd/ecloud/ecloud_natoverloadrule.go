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

func ecloudNATOverloadRuleRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "natoverloadrule",
		Short: "sub-commands relating to NAT overload rules",
	}

	// Child commands
	cmd.AddCommand(ecloudNATOverloadRuleListCmd(f))
	cmd.AddCommand(ecloudNATOverloadRuleShowCmd(f))
	cmd.AddCommand(ecloudNATOverloadRuleCreateCmd(f))
	cmd.AddCommand(ecloudNATOverloadRuleUpdateCmd(f))
	cmd.AddCommand(ecloudNATOverloadRuleDeleteCmd(f))

	return cmd
}

func ecloudNATOverloadRuleListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists NAT overload rules",
		Long:    "This command lists NAT overload rules",
		Example: "ans ecloud natoverloadrule list",
		RunE:    ecloudCobraRunEFunc(f, ecloudNATOverloadRuleList),
	}

	cmd.Flags().String("name", "", "NAT overload rule name for filtering")

	return cmd
}

func ecloudNATOverloadRuleList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	rules, err := service.GetNATOverloadRules(params)
	if err != nil {
		return fmt.Errorf("Error retrieving NAT overload rules: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudNATOverloadRulesProvider(rules))
}

func ecloudNATOverloadRuleShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <rule: id>...",
		Short:   "Shows a NAT overload rule",
		Long:    "This command shows one or more NAT overload rules",
		Example: "ans ecloud natoverloadrule show nor-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing NAT overload rule")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNATOverloadRuleShow),
	}
}

func ecloudNATOverloadRuleShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var rules []ecloud.NATOverloadRule
	for _, arg := range args {
		rule, err := service.GetNATOverloadRule(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving NAT overload rule [%s]: %s", arg, err)
			continue
		}

		rules = append(rules, rule)
	}

	return output.CommandOutput(cmd, OutputECloudNATOverloadRulesProvider(rules))
}

func ecloudNATOverloadRuleCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a NAT overload rule",
		Long:    "This command creates a NAT overload rule",
		Example: "ans ecloud natoverloadrule create --network net-abcdef12 --subnet 10.0.0.0/24 --action ALLOW",
		RunE:    ecloudCobraRunEFunc(f, ecloudNATOverloadRuleCreate),
	}

	// Setup flags
	cmd.Flags().String("network", "", "ID of network")
	cmd.MarkFlagRequired("network")
	cmd.Flags().String("subnet", "", "Subnet for rule")
	cmd.MarkFlagRequired("subnet")
	cmd.Flags().String("floating-ip", "", "ID of floating IP for rule")
	cmd.MarkFlagRequired("floating-ip")
	cmd.Flags().String("action", "", "Action for rule - allow/deny")
	cmd.MarkFlagRequired("action")
	cmd.Flags().String("name", "", "Name of rule")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the NAT overload rule has been completely created")

	return cmd
}

func ecloudNATOverloadRuleCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateNATOverloadRuleRequest{}
	createRequest.NetworkID, _ = cmd.Flags().GetString("network")
	createRequest.Subnet, _ = cmd.Flags().GetString("subnet")
	createRequest.FloatingIPID, _ = cmd.Flags().GetString("floating-ip")
	createRequest.Name, _ = cmd.Flags().GetString("name")

	action, _ := cmd.Flags().GetString("action")
	actionParsed, err := ecloud.ParseNATOverloadRuleAction(action)
	if err != nil {
		return err
	}
	createRequest.Action = actionParsed

	taskRef, err := service.CreateNATOverloadRule(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating NAT overload rule: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for NAT overload rule task to complete: %s", err)
		}
	}

	rule, err := service.GetNATOverloadRule(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new NAT overload rule: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudNATOverloadRulesProvider([]ecloud.NATOverloadRule{rule}))
}

func ecloudNATOverloadRuleUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <rule: id>...",
		Short:   "Updates a NAT overload rule",
		Long:    "This command updates one or more NAT overload rules",
		Example: "ans ecloud natoverloadrule update nor-abcdef12 --name \"my\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing NAT overload rule")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNATOverloadRuleUpdate),
	}

	cmd.Flags().String("name", "", "Name of rule")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the NAT overload rule has been completely updated")

	return cmd
}

func ecloudNATOverloadRuleUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchNATOverloadRuleRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var rules []ecloud.NATOverloadRule
	for _, arg := range args {
		task, err := service.PatchNATOverloadRule(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating NAT overload rule [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for NAT overload rule [%s]: %s", arg, err)
				continue
			}
		}

		rule, err := service.GetNATOverloadRule(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated NAT overload rule [%s]: %s", arg, err)
			continue
		}

		rules = append(rules, rule)
	}

	return output.CommandOutput(cmd, OutputECloudNATOverloadRulesProvider(rules))
}

func ecloudNATOverloadRuleDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <rule: id>...",
		Short:   "Removes a NAT overload rule",
		Long:    "This command removes one or more NAT overload rules",
		Example: "ans ecloud natoverloadrule delete nor-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing NAT overload rule")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNATOverloadRuleDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the NAT overload rule has been completely removed")

	return cmd
}

func ecloudNATOverloadRuleDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteNATOverloadRule(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing NAT overload rule [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for NAT overload rule [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
