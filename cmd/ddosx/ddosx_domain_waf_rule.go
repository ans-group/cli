package ddosx

import (
	"errors"
	"fmt"

	"github.com/ukfast/sdk-go/pkg/connection"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxDomainWAFRuleRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rule",
		Short: "sub-commands relating to domain web application firewall rules",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainWAFRuleListCmd(f))
	cmd.AddCommand(ddosxDomainWAFRuleShowCmd(f))
	cmd.AddCommand(ddosxDomainWAFRuleCreateCmd(f))
	cmd.AddCommand(ddosxDomainWAFRuleUpdateCmd(f))
	cmd.AddCommand(ddosxDomainWAFRuleDeleteCmd(f))

	return cmd
}

func ddosxDomainWAFRuleListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <domain: name>",
		Short:   "Lists domain WAF rules",
		Long:    "This command lists WAF rules",
		Example: "ukfast ddosx domain waf rule list",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return ddosxDomainWAFRuleList(f.NewClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainWAFRuleList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	domains, err := service.GetDomainWAFRules(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving domain WAF rules: %s", err)
	}

	return output.CommandOutput(cmd, OutputDDoSXWAFRulesProvider(domains))
}

func ddosxDomainWAFRuleShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: name> <rule: id>...",
		Short:   "Shows domain WAF rules",
		Long:    "This command shows a WAF rule",
		Example: "ukfast ddosx domain waf rule show example.com 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}
			if len(args) < 2 {
				return errors.New("Missing rule")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return ddosxDomainWAFRuleShow(f.NewClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainWAFRuleShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {

	var rules []ddosx.WAFRule

	for _, arg := range args[1:] {
		rule, err := service.GetDomainWAFRule(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving domain WAF rule [%s]: %s", arg, err.Error())
			continue
		}

		rules = append(rules, rule)
	}

	return output.CommandOutput(cmd, OutputDDoSXWAFRulesProvider(rules))
}

func ddosxDomainWAFRuleCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <domain: name>",
		Short:   "Creates domain WAF rules",
		Long:    "This command creates domain WAF rules",
		Example: "ukfast ddosx domain waf rule create example.com --uri example.html --ip 1.2.3.4",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return ddosxDomainWAFRuleCreate(f.NewClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("uri", "", "URI for rule")
	cmd.MarkFlagRequired("uri")
	cmd.Flags().String("ip", "", "IP for rule")
	cmd.MarkFlagRequired("ip")

	return cmd
}

func ddosxDomainWAFRuleCreate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	createRequest := ddosx.CreateWAFRuleRequest{}
	createRequest.URI, _ = cmd.Flags().GetString("uri")

	ip, _ := cmd.Flags().GetString("ip")
	createRequest.IP = connection.IPAddress(ip)

	id, err := service.CreateDomainWAFRule(args[0], createRequest)
	if err != nil {
		return fmt.Errorf("Error creating domain WAF rule: %s", err)
	}

	rule, err := service.GetDomainWAFRule(args[0], id)
	if err != nil {
		return fmt.Errorf("Error retrieving new domain WAF rule [%s]: %s", id, err)
	}

	return output.CommandOutput(cmd, OutputDDoSXWAFRulesProvider([]ddosx.WAFRule{rule}))
}

func ddosxDomainWAFRuleUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <domain: name> <ruleset: id>...",
		Short:   "Updates WAF rules",
		Long:    "This command updates one or more domain WAF rules",
		Example: "ukfast ddosx domain waf ruleset update example.com 00000000-0000-0000-0000-000000000000 --active=true",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}
			if len(args) < 2 {
				return errors.New("Missing rule")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return ddosxDomainWAFRuleUpdate(f.NewClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("uri", "", "URI for rule")
	cmd.Flags().String("ip", "", "IP for rule")

	return cmd
}

func ddosxDomainWAFRuleUpdate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	patchRequest := ddosx.PatchWAFRuleRequest{}

	if cmd.Flags().Changed("uri") {
		patchRequest.URI, _ = cmd.Flags().GetString("uri")
	}

	if cmd.Flags().Changed("ip") {
		ip, _ := cmd.Flags().GetString("ip")
		patchRequest.IP = connection.IPAddress(ip)
	}

	var rules []ddosx.WAFRule
	for _, arg := range args[1:] {
		err := service.PatchDomainWAFRule(args[0], arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating domain WAF rule [%s]: %s", arg, err.Error())
			continue
		}

		rule, err := service.GetDomainWAFRule(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated domain WAF rule [%s]: %s", arg, err)
			continue
		}

		rules = append(rules, rule)
	}

	return output.CommandOutput(cmd, OutputDDoSXWAFRulesProvider(rules))
}

func ddosxDomainWAFRuleDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <domain: name> <rule: id>...",
		Short:   "Deletes WAF rules",
		Long:    "This command deletes one or more domain WAF rules",
		Example: "ukfast ddosx domain waf rule delete example.com 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}
			if len(args) < 2 {
				return errors.New("Missing rule")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainWAFRuleDelete(f.NewClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainWAFRuleDelete(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	for _, arg := range args[1:] {
		err := service.DeleteDomainWAFRule(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing domain WAF rule [%s]: %s", arg, err.Error())
			continue
		}
	}
}
