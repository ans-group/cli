package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/ptr"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxDomainWAFRuleSetRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ruleset",
		Short: "sub-commands relating to domain web application firewall rule sets",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainWAFRuleSetListCmd())
	cmd.AddCommand(ddosxDomainWAFRuleSetShowCmd())
	cmd.AddCommand(ddosxDomainWAFRuleSetUpdateCmd())

	return cmd
}

func ddosxDomainWAFRuleSetListCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainWAFRuleSetList(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainWAFRuleSetList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	domains, err := service.GetDomainWAFRuleSets(args[0], params)
	if err != nil {
		output.Fatalf("Error retrieving domain waf rule sets: %s", err)
		return
	}

	outputDDoSXWAFRuleSets(domains)
}

func ddosxDomainWAFRuleSetShowCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainWAFRuleSetShow(getClient().DDoSXService(), cmd, args)
		},
	}
}
func ddosxDomainWAFRuleSetShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	var rulesets []ddosx.WAFRuleSet
	for _, arg := range args[1:] {
		ruleset, err := service.GetDomainWAFRuleSet(args[0], arg, connection.APIRequestParameters{})
		if err != nil {
			OutputWithErrorLevelf("Error retrieving domain WAF rule set [%s]: %s", arg, err)
			continue
		}

		rulesets = append(rulesets, ruleset)
	}

	outputDDoSXWAFRuleSets(rulesets)
}

func ddosxDomainWAFRuleSetUpdateCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainWAFRuleSetUpdate(getClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().Bool("active", false, "Specifies whether WAF rule set is active")

	return cmd
}

func ddosxDomainWAFRuleSetUpdate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	patchRequest := ddosx.PatchWAFRuleSetRequest{}

	if cmd.Flags().Changed("active") {
		activeBool, _ := cmd.Flags().GetBool("active")
		patchRequest.Active = ptr.Bool(activeBool)
	}

	var rulesets []ddosx.WAFRuleSet
	for _, arg := range args[1:] {
		err := service.PatchDomainWAFRuleSet(args[0], arg, patchRequest)
		if err != nil {
			OutputWithErrorLevelf("Error updating domain WAF rule set [%s]: %s", arg, err.Error())
			continue
		}

		ruleset, err := service.GetDomainWAFRuleSet(args[0], arg, connection.APIRequestParameters{})
		if err != nil {
			OutputWithErrorLevelf("Error retrieving updated domain WAF rule set [%s]: %s", arg, err)
			continue
		}

		rulesets = append(rulesets, ruleset)
	}

	outputDDoSXWAFRuleSets(rulesets)
}
