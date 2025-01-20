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

func ecloudAffinityRuleRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "affinityrule",
		Short: "sub-commands relating to affinity rules",
	}

	// Child commands
	cmd.AddCommand(ecloudAffinityRuleListCmd(f))
	cmd.AddCommand(ecloudAffinityRuleShowCmd(f))
	cmd.AddCommand(ecloudAffinityRuleCreateCmd(f))
	cmd.AddCommand(ecloudAffinityRuleUpdateCmd(f))
	cmd.AddCommand(ecloudAffinityRuleDeleteCmd(f))

	// Child root commands
	cmd.AddCommand(ecloudAffinityRuleAffinityMemberRootCmd(f))

	return cmd
}

func ecloudAffinityRuleListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists affinity rules",
		Long:    "This command lists affinity rules",
		Example: "ans ecloud affinityrule list",
		RunE:    ecloudCobraRunEFunc(f, ecloudAffinityRuleList),
	}

	cmd.Flags().String("name", "", "Affinity rule name for filtering")

	return cmd
}

func ecloudAffinityRuleList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	rules, err := service.GetAffinityRules(params)
	if err != nil {
		return fmt.Errorf("Error retrieving affinity rules: %s", err)
	}

	return output.CommandOutput(cmd, AffinityRuleCollection(rules))
}

func ecloudAffinityRuleShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <rule: id>...",
		Short:   "Shows an affinity rule",
		Long:    "This command shows one or more affinity rules",
		Example: "ans ecloud affinityrule show ar-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing affinity rule")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudAffinityRuleShow),
	}
}

func ecloudAffinityRuleShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var rules []ecloud.AffinityRule
	for _, arg := range args {
		rule, err := service.GetAffinityRule(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving affinity rule [%s]: %s", arg, err)
			continue
		}

		rules = append(rules, rule)
	}

	return output.CommandOutput(cmd, AffinityRuleCollection(rules))
}

func ecloudAffinityRuleCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates an affinity rule",
		Long:    "This command creates an affinity rule",
		Example: "ans ecloud affinityrule create --vpc vpc-abcdef12 --availability-zone az-abcdef12 --type anti-affinity",
		RunE:    ecloudCobraRunEFunc(f, ecloudAffinityRuleCreate),
	}

	// Setup flags
	cmd.Flags().String("vpc", "", "ID of VPC")
	cmd.MarkFlagRequired("vpc")
	cmd.Flags().String("availability-zone", "", "ID of AZ")
	cmd.MarkFlagRequired("availability-zone")
	cmd.Flags().String("type", "", "Type of rule. One of: affinity/anti-affinity")
	cmd.MarkFlagRequired("type")
	cmd.Flags().String("name", "", "Name of affinity rule")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the affinity rule has been completely created")

	return cmd
}

func ecloudAffinityRuleCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateAffinityRuleRequest{}
	createRequest.VPCID, _ = cmd.Flags().GetString("vpc")
	createRequest.AvailabilityZoneID, _ = cmd.Flags().GetString("availability-zone")

	ruleType, _ := cmd.Flags().GetString("type")
	parsedType, err := ecloud.AffinityRuleTypeEnum.Parse(ruleType)
	if err != nil {
		return err
	}
	createRequest.Type = parsedType

	if cmd.Flags().Changed("name") {
		createRequest.Name, _ = cmd.Flags().GetString("name")
	}

	taskRef, err := service.CreateAffinityRule(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating affinity rule: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for affinity rule task to complete: %s", err)
		}
	}

	rule, err := service.GetAffinityRule(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new affinity rule: %s", err)
	}

	return output.CommandOutput(cmd, AffinityRuleCollection([]ecloud.AffinityRule{rule}))
}

func ecloudAffinityRuleUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <rule: id>...",
		Short:   "Updates an affinity rule",
		Long:    "This command updates one or more affinity rules",
		Example: "ans ecloud affinityrule update ar-abcdef12 --name \"my rule\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing affinity rule")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudAffinityRuleUpdate),
	}

	cmd.Flags().String("name", "", "Name of affinity rule")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the affinity rule has been completely updated")

	return cmd
}

func ecloudAffinityRuleUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchAffinityRuleRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var rules []ecloud.AffinityRule
	for _, arg := range args {
		task, err := service.PatchAffinityRule(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating affinity rule [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for affinity rule [%s]: %s", arg, err)
				continue
			}
		}

		rule, err := service.GetAffinityRule(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated affinity rule [%s]: %s", arg, err)
			continue
		}

		rules = append(rules, rule)
	}

	return output.CommandOutput(cmd, AffinityRuleCollection(rules))
}

func ecloudAffinityRuleDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <rule: id>...",
		Short:   "Removes an affinity rule",
		Long:    "This command removes one or more affinity rules",
		Example: "ans ecloud affinityrule delete ar-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing affinity rule")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudAffinityRuleDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the affinity rule has been completely removed")

	return cmd
}

func ecloudAffinityRuleDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteAffinityRule(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing affinity rule [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for affinity rule [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
