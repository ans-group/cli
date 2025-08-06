package ddosx

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/ptr"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	"github.com/spf13/cobra"
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
		Example: "ans ddosx domain waf ruleset list",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing domain")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxDomainWAFRuleSetList(c.DDoSXService(), cmd, args)
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
		return fmt.Errorf("error retrieving domain waf rule sets: %s", err)
	}

	return output.CommandOutput(cmd, WAFRuleSetCollection(domains))
}

func ddosxDomainWAFRuleSetShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: name> <ruleset: id>...",
		Short:   "Shows WAF rule sets",
		Long:    "This command shows one or more domain WAF rule sets",
		Example: "ans ddosx domain waf ruleset show example.com 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing domain")
			}
			if len(args) < 2 {
				return errors.New("missing rule set")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxDomainWAFRuleSetShow(c.DDoSXService(), cmd, args)
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

	return output.CommandOutput(cmd, WAFRuleSetCollection(rulesets))
}

func ddosxDomainWAFRuleSetUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <domain: name> <ruleset: id>...",
		Short:   "Updates WAF rule sets",
		Long:    "This command updates one or more domain WAF rule sets",
		Example: "ans ddosx domain waf ruleset update example.com 00000000-0000-0000-0000-000000000000 --active=true",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing domain")
			}
			if len(args) < 2 {
				return errors.New("missing rule set")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxDomainWAFRuleSetUpdate(c.DDoSXService(), cmd, args)
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

	return output.CommandOutput(cmd, WAFRuleSetCollection(rulesets))
}
