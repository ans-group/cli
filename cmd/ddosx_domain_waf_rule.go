package cmd

import (
	"errors"

	"github.com/ukfast/sdk-go/pkg/connection"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxDomainWAFRuleRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rule",
		Short: "sub-commands relating to domain web application firewall rules",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainWAFRuleListCmd())
	cmd.AddCommand(ddosxDomainWAFRuleShowCmd())
	cmd.AddCommand(ddosxDomainWAFRuleCreateCmd())
	cmd.AddCommand(ddosxDomainWAFRuleUpdateCmd())
	cmd.AddCommand(ddosxDomainWAFRuleDeleteCmd())

	return cmd
}

func ddosxDomainWAFRuleListCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainWAFRuleList(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainWAFRuleList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	domains, err := service.GetDomainWAFRules(args[0], params)
	if err != nil {
		output.Fatalf("Error retrieving domain WAF rules: %s", err)
		return
	}

	outputDDoSXWAFRules(domains)
}

func ddosxDomainWAFRuleShowCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainWAFRuleShow(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainWAFRuleShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {

	var rules []ddosx.WAFRule

	for _, arg := range args[1:] {
		rule, err := service.GetDomainWAFRule(args[0], arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving domain WAF rule [%s]: %s", arg, err.Error())
			continue
		}

		rules = append(rules, rule)
	}

	outputDDoSXWAFRules(rules)
}

func ddosxDomainWAFRuleCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <domain: name>",
		Short:   "Creates domain WAF rules",
		Long:    "This command creates domain WAF rules",
		Example: "ukfast ddosx domain waf rule create --uri example.html --ip 1.2.3.4",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainWAFRuleCreate(getClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("uri", "", "URI for rule")
	cmd.MarkFlagRequired("uri")
	cmd.Flags().String("ip", "", "IP for rule")
	cmd.MarkFlagRequired("ip")

	return cmd
}

func ddosxDomainWAFRuleCreate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	createRequest := ddosx.CreateWAFRuleRequest{}
	createRequest.URI, _ = cmd.Flags().GetString("uri")

	ip, _ := cmd.Flags().GetString("ip")
	createRequest.IP = connection.IPAddress(ip)

	_, err := service.CreateDomainWAFRule(args[0], createRequest)
	if err != nil {
		output.Fatalf("Error creating domain WAF rule: %s", err)
		return
	}

	// TODO: add rule retrieval
}

func ddosxDomainWAFRuleUpdateCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainWAFRuleUpdate(getClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("uri", "", "URI for rule")
	cmd.Flags().String("ip", "", "IP for rule")

	return cmd
}

func ddosxDomainWAFRuleUpdate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	patchRequest := ddosx.PatchWAFRuleRequest{}

	if cmd.Flags().Changed("uri") {
		patchRequest.URI, _ = cmd.Flags().GetString("uri")
	}

	if cmd.Flags().Changed("ip") {
		ip, _ := cmd.Flags().GetString("ip")
		patchRequest.IP = connection.IPAddress(ip)
	}

	for _, arg := range args[1:] {
		err := service.PatchDomainWAFRule(args[0], arg, patchRequest)
		if err != nil {
			OutputWithErrorLevelf("Error updating domain WAF rule [%s]: %s", arg, err.Error())
			continue
		}

		// TODO: add rule retrieval
	}
}

func ddosxDomainWAFRuleDeleteCmd() *cobra.Command {
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
			ddosxDomainWAFRuleDelete(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainWAFRuleDelete(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	for _, arg := range args[1:] {
		err := service.DeleteDomainWAFRule(args[0], arg)
		if err != nil {
			OutputWithErrorLevelf("Error removing domain WAF rule [%s]: %s", arg, err.Error())
			continue
		}
	}
}
