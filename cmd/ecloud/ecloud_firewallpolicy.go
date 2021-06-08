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

func ecloudFirewallPolicyRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "firewallpolicy",
		Short: "sub-commands relating to policies",
	}

	// Child commands
	cmd.AddCommand(ecloudFirewallPolicyListCmd(f))
	cmd.AddCommand(ecloudFirewallPolicyShowCmd(f))
	cmd.AddCommand(ecloudFirewallPolicyCreateCmd(f))
	cmd.AddCommand(ecloudFirewallPolicyUpdateCmd(f))
	cmd.AddCommand(ecloudFirewallPolicyDeleteCmd(f))

	// Child root commands
	cmd.AddCommand(ecloudFirewallPolicyFirewallRuleRootCmd(f))
	cmd.AddCommand(ecloudFirewallPolicyTaskRootCmd(f))

	return cmd
}

func ecloudFirewallPolicyListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists firewall policies",
		Long:    "This command lists firewall policies",
		Example: "ukfast ecloud firewallpolicy list",
		RunE:    ecloudCobraRunEFunc(f, ecloudFirewallPolicyList),
	}

	cmd.Flags().String("name", "", "Firewall policy name for filtering")
	cmd.Flags().String("router", "", "Firewall policy router ID for filtering")

	return cmd
}

func ecloudFirewallPolicyList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
		helper.NewStringFilterFlagOption("router", "router_id"),
	)
	if err != nil {
		return err
	}

	policies, err := service.GetFirewallPolicies(params)
	if err != nil {
		return fmt.Errorf("Error retrieving firewall policies: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudFirewallPoliciesProvider(policies))
}

func ecloudFirewallPolicyShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <policy: id>...",
		Short:   "Shows a firewall policy",
		Long:    "This command shows one or more firewall policies",
		Example: "ukfast ecloud firewallpolicy show fwp-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing firewall policy")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFirewallPolicyShow),
	}
}

func ecloudFirewallPolicyShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var policies []ecloud.FirewallPolicy
	for _, arg := range args {
		policy, err := service.GetFirewallPolicy(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving firewall policy [%s]: %s", arg, err)
			continue
		}

		policies = append(policies, policy)
	}

	return output.CommandOutput(cmd, OutputECloudFirewallPoliciesProvider(policies))
}

func ecloudFirewallPolicyCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a firewall policy",
		Long:    "This command creates a firewall policy",
		Example: "ukfast ecloud firewallpolicy create --router rtr-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudFirewallPolicyCreate),
	}

	// Setup flags
	cmd.Flags().String("router", "", "ID of router")
	cmd.MarkFlagRequired("router")
	cmd.Flags().Int("sequence", 0, "Sequence for policy")
	cmd.MarkFlagRequired("sequence")
	cmd.Flags().String("name", "", "Name of policy")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the firewall policy has been completely created")

	return cmd
}

func ecloudFirewallPolicyCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateFirewallPolicyRequest{}
	createRequest.RouterID, _ = cmd.Flags().GetString("router")
	if cmd.Flags().Changed("name") {
		createRequest.Name, _ = cmd.Flags().GetString("name")
	}
	if cmd.Flags().Changed("sequence") {
		createRequest.Sequence, _ = cmd.Flags().GetInt("sequence")
	}

	taskRef, err := service.CreateFirewallPolicy(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating firewall policy: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for firewall policy task to complete: %s", err)
		}
	}

	policy, err := service.GetFirewallPolicy(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new firewall policy: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudFirewallPoliciesProvider([]ecloud.FirewallPolicy{policy}))
}

func ecloudFirewallPolicyUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <policy: id>...",
		Short:   "Updates a firewall policy",
		Long:    "This command updates one or more firewall policies",
		Example: "ukfast ecloud firewallpolicy update fwp-abcdef12 --name \"my policy\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing firewall policy")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFirewallPolicyUpdate),
	}

	cmd.Flags().String("name", "", "Name of policy")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the firewall policy has been completely updated")

	return cmd
}

func ecloudFirewallPolicyUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchFirewallPolicyRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var policies []ecloud.FirewallPolicy
	for _, arg := range args {
		task, err := service.PatchFirewallPolicy(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating firewall policy [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for firewall policy [%s]: %s", arg, err)
				continue
			}
		}

		policy, err := service.GetFirewallPolicy(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated firewall policy [%s]: %s", arg, err)
			continue
		}

		policies = append(policies, policy)
	}

	return output.CommandOutput(cmd, OutputECloudFirewallPoliciesProvider(policies))
}

func ecloudFirewallPolicyDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <policy: id>...",
		Short:   "Removes a firewall policy",
		Long:    "This command removes one or more firewall policies",
		Example: "ukfast ecloud firewallpolicy delete fwp-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing firewall policy")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFirewallPolicyDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the firewall policy has been completely removed")

	return cmd
}

func ecloudFirewallPolicyDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteFirewallPolicy(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing firewall policy [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for firewall policy [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
