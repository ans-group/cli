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

func ecloudNetworkPolicyRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "networkpolicy",
		Short: "sub-commands relating to policies",
	}

	// Child commands
	cmd.AddCommand(ecloudNetworkPolicyListCmd(f))
	cmd.AddCommand(ecloudNetworkPolicyShowCmd(f))
	cmd.AddCommand(ecloudNetworkPolicyCreateCmd(f))
	cmd.AddCommand(ecloudNetworkPolicyUpdateCmd(f))
	cmd.AddCommand(ecloudNetworkPolicyDeleteCmd(f))

	// Child root commands
	cmd.AddCommand(ecloudNetworkPolicyNetworkRuleRootCmd(f))
	cmd.AddCommand(ecloudNetworkPolicyTaskRootCmd(f))

	return cmd
}

func ecloudNetworkPolicyListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists network policies",
		Long:    "This command lists network policies",
		Example: "ans ecloud networkpolicy list",
		RunE:    ecloudCobraRunEFunc(f, ecloudNetworkPolicyList),
	}

	cmd.Flags().String("name", "", "Network policy name for filtering")
	cmd.Flags().String("network", "", "Network policy network ID for filtering")

	return cmd
}

func ecloudNetworkPolicyList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
		helper.NewStringFilterFlagOption("network", "network_id"),
	)
	if err != nil {
		return err
	}

	policies, err := service.GetNetworkPolicies(params)
	if err != nil {
		return fmt.Errorf("Error retrieving network policies: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudNetworkPoliciesProvider(policies))
}

func ecloudNetworkPolicyShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <policy: id>...",
		Short:   "Shows a network policy",
		Long:    "This command shows one or more network policies",
		Example: "ans ecloud networkpolicy show np-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing network policy")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNetworkPolicyShow),
	}
}

func ecloudNetworkPolicyShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var policies []ecloud.NetworkPolicy
	for _, arg := range args {
		policy, err := service.GetNetworkPolicy(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving network policy [%s]: %s", arg, err)
			continue
		}

		policies = append(policies, policy)
	}

	return output.CommandOutput(cmd, OutputECloudNetworkPoliciesProvider(policies))
}

func ecloudNetworkPolicyCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a network policy",
		Long:    "This command creates a network policy",
		Example: "ans ecloud networkpolicy create --network rtr-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudNetworkPolicyCreate),
	}

	// Setup flags
	cmd.Flags().String("network", "", "ID of network")
	cmd.MarkFlagRequired("network")
	cmd.Flags().String("name", "", "Name of policy")
	cmd.Flags().String("catchall-rule-action", "", "Action of catchall rule. One of: ALLOW/DROP/REJECT")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the network policy has been completely created")

	return cmd
}

func ecloudNetworkPolicyCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateNetworkPolicyRequest{}
	createRequest.NetworkID, _ = cmd.Flags().GetString("network")
	if cmd.Flags().Changed("name") {
		createRequest.Name, _ = cmd.Flags().GetString("name")
	}
	if cmd.Flags().Changed("catchall-rule-action") {
		catchallRuleAction, _ := cmd.Flags().GetString("catchall-rule-action")
		catchallRuleActionParsed, err := ecloud.NetworkPolicyCatchallRuleActionEnum.Parse(catchallRuleAction)
		if err != nil {
			return err
		}
		createRequest.CatchallRuleAction = catchallRuleActionParsed
	}

	taskRef, err := service.CreateNetworkPolicy(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating network policy: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for network policy task to complete: %s", err)
		}
	}

	policy, err := service.GetNetworkPolicy(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new network policy: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudNetworkPoliciesProvider([]ecloud.NetworkPolicy{policy}))
}

func ecloudNetworkPolicyUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <policy: id>...",
		Short:   "Updates a network policy",
		Long:    "This command updates one or more network policies",
		Example: "ans ecloud networkpolicy update np-abcdef12 --name \"my policy\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing network policy")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNetworkPolicyUpdate),
	}

	cmd.Flags().String("name", "", "Name of policy")
	cmd.Flags().String("catchall-rule-action", "", "Action of catchall rule. One of: ALLOW/DROP/REJECT")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the network policy has been completely updated")

	return cmd
}

func ecloudNetworkPolicyUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchNetworkPolicyRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}
	if cmd.Flags().Changed("catchall-rule-action") {
		catchallRuleAction, _ := cmd.Flags().GetString("catchall-rule-action")
		catchallRuleActionParsed, err := ecloud.NetworkPolicyCatchallRuleActionEnum.Parse(catchallRuleAction)
		if err != nil {
			return err
		}
		patchRequest.CatchallRuleAction = catchallRuleActionParsed
	}

	var policies []ecloud.NetworkPolicy
	for _, arg := range args {
		task, err := service.PatchNetworkPolicy(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating network policy [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for network policy [%s]: %s", arg, err)
				continue
			}
		}

		policy, err := service.GetNetworkPolicy(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated network policy [%s]: %s", arg, err)
			continue
		}

		policies = append(policies, policy)
	}

	return output.CommandOutput(cmd, OutputECloudNetworkPoliciesProvider(policies))
}

func ecloudNetworkPolicyDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <policy: id>...",
		Short:   "Removes a network policy",
		Long:    "This command removes one or more network policies",
		Example: "ans ecloud networkpolicy delete np-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing network policy")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNetworkPolicyDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the network policy has been completely removed")

	return cmd
}

func ecloudNetworkPolicyDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteNetworkPolicy(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing network policy [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for network policy [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
