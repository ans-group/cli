package ecloud

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/connection"
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

	policyID, err := service.CreateFirewallPolicy(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating firewall policy: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(FirewallPolicyResourceSyncStatusWaitFunc(service, policyID, ecloud.SyncStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for firewall policy sync: %s", err)
		}
	}

	policy, err := service.GetFirewallPolicy(policyID)
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
		err := service.PatchFirewallPolicy(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating firewall policy [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(FirewallPolicyResourceSyncStatusWaitFunc(service, arg, ecloud.SyncStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for firewall policy [%s] sync: %s", arg, err)
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
		err := service.DeleteFirewallPolicy(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing firewall policy [%s]: %s", arg, err)
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(FirewallPolicyNotFoundWaitFunc(service, arg))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for removal of firewall policy [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}

func FirewallPolicyResourceSyncStatusWaitFunc(service ecloud.ECloudService, policyID string, status ecloud.SyncStatus) helper.WaitFunc {
	return ResourceSyncStatusWaitFunc(func() (ecloud.SyncStatus, error) {
		policy, err := service.GetFirewallPolicy(policyID)
		if err != nil {
			return "", err
		}
		return policy.Sync.Status, nil
	}, status)
}

func FirewallPolicyTaskStatusWaitFunc(service ecloud.ECloudService, policyID string, taskID string, status ecloud.TaskStatus) helper.WaitFunc {
	return TaskStatusFromResourceTaskListWaitFunc(service, taskID, FirewallPolicyTaskListFunc(service, policyID), status)
}

func FirewallPolicyTaskListFunc(service ecloud.ECloudService, policyID string) ResourceTaskListFunc {
	return func(params connection.APIRequestParameters) ([]ecloud.Task, error) {
		return service.GetFirewallPolicyTasks(policyID, params)
	}
}

func FirewallPolicyNotFoundWaitFunc(service ecloud.ECloudService, firewallPolicyID string) helper.WaitFunc {
	return func() (finished bool, err error) {
		_, err = service.GetFirewallPolicy(firewallPolicyID)
		if err != nil {
			switch err.(type) {
			case *ecloud.FirewallPolicyNotFoundError:
				return true, nil
			default:
				return false, fmt.Errorf("Failed to retrieve firewall policy [%s]: %s", firewallPolicyID, err)
			}
		}

		return false, nil
	}
}
