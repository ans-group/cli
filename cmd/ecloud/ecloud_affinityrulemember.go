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

func ecloudAffinityRuleMemberRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "affinityrulemember",
		Short: "sub-commands relating to affinity rule members",
	}

	// Child commands
	cmd.AddCommand(ecloudAffinityRuleMemberShowCmd(f))
	cmd.AddCommand(ecloudAffinityRuleMemberCreateCmd(f))
	cmd.AddCommand(ecloudAffinityRuleMemberDeleteCmd(f))

	return cmd
}

func ecloudAffinityRuleMemberShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <rulemember: id>...",
		Short:   "Shows an affinity rule member",
		Long:    "This command shows one or more affinity rule members",
		Example: "ans ecloud affinityrulemember show arm-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing affinity rule member")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudAffinityRuleMemberShow),
	}
}

func ecloudAffinityRuleMemberShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var members []ecloud.AffinityRuleMember
	for _, arg := range args {
		rule, err := service.GetAffinityRuleMember(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving affinity rule member [%s]: %s", arg, err)
			continue
		}

		members = append(members, rule)
	}

	return output.CommandOutput(cmd, AffinityRuleMemberCollection(members))
}

func ecloudAffinityRuleMemberCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates an affinity rule member",
		Long:    "This command creates an affinity rule member",
		Example: "ans ecloud affinityrulemember create --affinity-rule ar-abcdef12 --instance i-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudAffinityRuleMemberCreate),
	}

	// Setup flags
	cmd.Flags().String("affinity-rule", "", "ID of affinity rule")
	_ = cmd.MarkFlagRequired("affinity-rule")
	cmd.Flags().String("instance", "", "ID of instance")
	_ = cmd.MarkFlagRequired("instance")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the affinity rule member has been completely created")

	return cmd
}

func ecloudAffinityRuleMemberCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateAffinityRuleMemberRequest{}
	createRequest.AffinityRuleID, _ = cmd.Flags().GetString("affinity-rule")
	createRequest.InstanceID, _ = cmd.Flags().GetString("instance")

	taskRef, err := service.CreateAffinityRuleMember(createRequest)
	if err != nil {
		return fmt.Errorf("error creating affinity rule member: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("error waiting for affinity rule member task to complete: %s", err)
		}
	}

	rule, err := service.GetAffinityRuleMember(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("error retrieving new affinity rule member: %s", err)
	}

	return output.CommandOutput(cmd, AffinityRuleMemberCollection([]ecloud.AffinityRuleMember{rule}))
}

func ecloudAffinityRuleMemberDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <rulemember: id>...",
		Short:   "Removes an affinity rule member",
		Long:    "This command removes one or more affinity rule members",
		Example: "ans ecloud affinityrulemember delete arm-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing affinity rule member")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudAffinityRuleMemberDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the affinity rule member has been completely removed")

	return cmd
}

func ecloudAffinityRuleMemberDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteAffinityRuleMember(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing affinity rule member [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for affinity rule member [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
