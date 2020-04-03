package ddosx

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/ptr"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxDomainWAFRuleSetRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ruleset",
		Short: "sub-commands relating to domain web application firewall rule sets",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainWAFRuleSetListCmd(f))
	cmd.AddCommand(ddosxDomainWAFRuleSetShowCmd(f))
	cmd.AddCommand(ddosxDomainWAFRuleSetUpdateCmd(f))

	return cmd
}

func ddosxDomainWAFRuleSetListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <domain: name>",
		Short:   "Lists WAF rule sets",
		Long:    "This command lists WAF rule sets",
		Example: "ukfast ddosx domain waf ruleset list",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return ddosxDomainWAFRuleSetList(f.NewClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainWAFRuleSetList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	domains, err := service.GetDomainWAFRuleSets(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving domain waf rule sets: %s", err)
	}

	return output.CommandOutput(cmd, OutputDDoSXWAFRuleSetsProvider(domains))
}

func ddosxDomainWAFRuleSetShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: name> <ruleset: id>...",
		Short:   "Shows WAF rule sets",
		Long:    "This command shows one or more domain WAF rule sets",
		Example: "ukfast ddosx domain waf ruleset show example.com 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}
			if len(args) < 2 {
				return errors.New("Missing rule set")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return ddosxDomainWAFRuleSetShow(f.NewClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainWAFRuleSetShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	var rulesets []ddosx.WAFRuleSet
	for _, arg := range args[1:] {
		ruleset, err := service.GetDomainWAFRuleSet(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving domain WAF rule set [%s]: %s", arg, err)
			continue
		}

		rulesets = append(rulesets, ruleset)
	}

	return output.CommandOutput(cmd, OutputDDoSXWAFRuleSetsProvider(rulesets))
}

func ddosxDomainWAFRuleSetUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <domain: name> <ruleset: id>...",
		Short:   "Updates WAF rule sets",
		Long:    "This command updates one or more domain WAF rule sets",
		Example: "ukfast ddosx domain waf ruleset update example.com 00000000-0000-0000-0000-000000000000 --active=true",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}
			if len(args) < 2 {
				return errors.New("Missing rule set")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return ddosxDomainWAFRuleSetUpdate(f.NewClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().Bool("active", false, "Specifies whether WAF rule set is active")

	return cmd
}

func ddosxDomainWAFRuleSetUpdate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	patchRequest := ddosx.PatchWAFRuleSetRequest{}

	if cmd.Flags().Changed("active") {
		activeBool, _ := cmd.Flags().GetBool("active")
		patchRequest.Active = ptr.Bool(activeBool)
	}

	var rulesets []ddosx.WAFRuleSet
	for _, arg := range args[1:] {
		err := service.PatchDomainWAFRuleSet(args[0], arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating domain WAF rule set [%s]: %s", arg, err.Error())
			continue
		}

		ruleset, err := service.GetDomainWAFRuleSet(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated domain WAF rule set [%s]: %s", arg, err)
			continue
		}

		rulesets = append(rulesets, ruleset)
	}

	return output.CommandOutput(cmd, OutputDDoSXWAFRuleSetsProvider(rulesets))
}
