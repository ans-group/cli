package cmd

import (
	"errors"

	"github.com/ukfast/sdk-go/pkg/connection"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxDomainWAFAdvancedRuleRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "advancedrule",
		Short: "sub-commands relating to domain web application firewall advanced rules",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainWAFAdvancedRuleListCmd())
	cmd.AddCommand(ddosxDomainWAFAdvancedRuleShowCmd())
	cmd.AddCommand(ddosxDomainWAFAdvancedRuleCreateCmd())
	cmd.AddCommand(ddosxDomainWAFAdvancedRuleUpdateCmd())
	cmd.AddCommand(ddosxDomainWAFAdvancedRuleDeleteCmd())

	return cmd
}

func ddosxDomainWAFAdvancedRuleListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list <domain: name>",
		Short:   "Lists domain WAF advanced rules",
		Long:    "This command lists domain WAF advanced rules",
		Example: "ukfast ddosx domain waf advancedrule list",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainWAFAdvancedRuleList(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainWAFAdvancedRuleList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	domains, err := service.GetDomainWAFAdvancedRules(args[0], params)
	if err != nil {
		output.Fatalf("Error retrieving domain WAF advanced rules: %s", err)
		return
	}

	outputDDoSXWAFAdvancedRules(domains)
}

func ddosxDomainWAFAdvancedRuleShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: name> <rule: id>...",
		Short:   "Shows domain WAF advanced rules",
		Long:    "This command shows a WAF advanced rule",
		Example: "ukfast ddosx domain waf advancedrule show example.com 00000000-0000-0000-0000-000000000000",
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
			ddosxDomainWAFAdvancedRuleShow(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainWAFAdvancedRuleShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {

	var rules []ddosx.WAFAdvancedRule

	for _, arg := range args[1:] {
		rule, err := service.GetDomainWAFAdvancedRule(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving domain WAF advanced rule [%s]: %s", arg, err.Error())
			continue
		}

		rules = append(rules, rule)
	}

	outputDDoSXWAFAdvancedRules(rules)
}

func ddosxDomainWAFAdvancedRuleCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <domain: name>",
		Short:   "Creates domain WAF advanced rules",
		Long:    "This command creates domain WAF advanced rules",
		Example: "ukfast ddosx domain waf advancedrule create --section REQUEST_URI --modifier beginswith --phrase test --ip 1.2.3.4",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainWAFAdvancedRuleCreate(getClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("section", "", "Section for advanced rule")
	cmd.MarkFlagRequired("section")
	cmd.Flags().String("modifier", "", "Modifier for advanced rule")
	cmd.MarkFlagRequired("modifier")
	cmd.Flags().String("phrase", "", "Phrase for advanced rule")
	cmd.MarkFlagRequired("phrase")
	cmd.Flags().String("ip", "", "IP for advanced rule")
	cmd.MarkFlagRequired("ip")

	return cmd
}

func ddosxDomainWAFAdvancedRuleCreate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	createRequest := ddosx.CreateWAFAdvancedRuleRequest{}

	modifier, _ := cmd.Flags().GetString("modifier")
	parsedModifier, err := ddosx.ParseWAFAdvancedRuleModifier(modifier)
	if err != nil {
		output.Fatal(err.Error())
		return
	}
	createRequest.Modifier = parsedModifier
	section, _ := cmd.Flags().GetString("section")
	createRequest.Section = ddosx.WAFAdvancedRuleSection(section)
	createRequest.Phrase, _ = cmd.Flags().GetString("phrase")
	ip, _ := cmd.Flags().GetString("ip")
	createRequest.IP = connection.IPAddress(ip)

	id, err := service.CreateDomainWAFAdvancedRule(args[0], createRequest)
	if err != nil {
		output.Fatalf("Error creating domain WAF advanced rule: %s", err)
		return
	}

	rule, err := service.GetDomainWAFAdvancedRule(args[0], id)
	if err != nil {
		output.Fatalf("Error retrieving new domain WAF advanced rule [%s]: %s", id, err)
		return
	}

	outputDDoSXWAFAdvancedRules([]ddosx.WAFAdvancedRule{rule})
}

func ddosxDomainWAFAdvancedRuleUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <domain: name> <advancedrule: id>...",
		Short:   "Updates domain WAF advanced rules",
		Long:    "This command updates one or more domain WAF advanced rules",
		Example: "ukfast ddosx domain waf advancedrule update example.com 00000000-0000-0000-0000-000000000000 --ip 1.2.3.4",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}
			if len(args) < 2 {
				return errors.New("Missing advanced rule")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainWAFAdvancedRuleUpdate(getClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("section", "", "Section for advanced rule")
	cmd.Flags().String("modifier", "", "Modifier for advanced rule")
	cmd.Flags().String("phrase", "", "Phrase for advanced rule")
	cmd.Flags().String("ip", "", "IP for advanced rule")

	return cmd
}

func ddosxDomainWAFAdvancedRuleUpdate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	patchRequest := ddosx.PatchWAFAdvancedRuleRequest{}

	if cmd.Flags().Changed("modifier") {
		modifier, _ := cmd.Flags().GetString("modifier")
		parsedModifier, err := ddosx.ParseWAFAdvancedRuleModifier(modifier)
		if err != nil {
			output.Fatal(err.Error())
			return
		}
		patchRequest.Modifier = parsedModifier
	}
	if cmd.Flags().Changed("section") {
		section, _ := cmd.Flags().GetString("section")
		patchRequest.Section = ddosx.WAFAdvancedRuleSection(section)
	}
	if cmd.Flags().Changed("phrase") {
		patchRequest.Phrase, _ = cmd.Flags().GetString("phrase")
	}
	if cmd.Flags().Changed("ip") {
		ip, _ := cmd.Flags().GetString("ip")
		patchRequest.IP = connection.IPAddress(ip)
	}

	var rules []ddosx.WAFAdvancedRule
	for _, arg := range args[1:] {
		err := service.PatchDomainWAFAdvancedRule(args[0], arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating domain WAF advanced rule [%s]: %s", arg, err.Error())
			continue
		}

		rule, err := service.GetDomainWAFAdvancedRule(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated domain WAF advanced rule [%s]: %s", arg, err)
			continue
		}

		rules = append(rules, rule)
	}

	outputDDoSXWAFAdvancedRules(rules)
}

func ddosxDomainWAFAdvancedRuleDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <domain: name> <advancedrule: id>...",
		Short:   "Deletes domain WAF advanced rules",
		Long:    "This command deletes one or more domain WAF advanced rules",
		Example: "ukfast ddosx domain waf advancedrule delete example.com 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}
			if len(args) < 2 {
				return errors.New("Missing advanced rule")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainWAFAdvancedRuleDelete(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainWAFAdvancedRuleDelete(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	for _, arg := range args[1:] {
		err := service.DeleteDomainWAFAdvancedRule(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing domain WAF advanced rule [%s]: %s", arg, err.Error())
			continue
		}
	}
}
