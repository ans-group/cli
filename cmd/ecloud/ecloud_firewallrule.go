package ecloud

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	flaghelper "github.com/ukfast/cli/internal/pkg/helper/flag"
	"github.com/ukfast/cli/internal/pkg/output"
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

	return cmd
}

func ecloudFirewallRuleListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
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
}

func ecloudFirewallRuleList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := flaghelper.GetAPIRequestParametersFromFlags(cmd)
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
