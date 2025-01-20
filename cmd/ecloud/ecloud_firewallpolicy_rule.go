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

func ecloudFirewallPolicyFirewallRuleRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rule",
		Short: "sub-commands relating to firewall policy firewall rules",
	}

	// Child commands
	cmd.AddCommand(ecloudFirewallPolicyFirewallRuleListCmd(f))

	return cmd
}

func ecloudFirewallPolicyFirewallRuleListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists firewall rules for firewall policy",
		Long:    "This command lists firewall rules for firewall policy",
		Example: "ans ecloud firewallpolicy firewallrule list fwp-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing firewall policy")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudFirewallPolicyFirewallRuleList),
	}

	cmd.Flags().String("name", "", "Firewall rule name for filtering")

	return cmd
}

func ecloudFirewallPolicyFirewallRuleList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	rules, err := service.GetFirewallPolicyFirewallRules(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving firewall policy firewall rules: %s", err)
	}

	return output.CommandOutput(cmd, FirewallRuleCollection(rules))
}
