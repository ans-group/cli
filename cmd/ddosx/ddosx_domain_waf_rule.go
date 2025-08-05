package ddosx

import (
	"errors"
	"fmt"

	"github.com/ans-group/sdk-go/pkg/connection"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	"github.com/spf13/cobra"
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
		Example: "ans ddosx domain waf rule list",
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

			return ddosxDomainWAFRuleList(c.DDoSXService(), cmd, args)
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
		return fmt.Errorf("error retrieving domain WAF rules: %s", err)
	}

	return output.CommandOutput(cmd, WAFRuleCollection(domains))
}

func ddosxDomainWAFRuleShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: name> <rule: id>...",
		Short:   "Shows domain WAF rules",
		Long:    "This command shows a WAF rule",
		Example: "ans ddosx domain waf rule show example.com 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing domain")
			}
			if len(args) < 2 {
				return errors.New("missing rule")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxDomainWAFRuleShow(c.DDoSXService(), cmd, args)
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

	return output.CommandOutput(cmd, WAFRuleCollection(rules))
}

func ddosxDomainWAFRuleCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <domain: name>",
		Short:   "Creates domain WAF rules",
		Long:    "This command creates domain WAF rules",
		Example: "ans ddosx domain waf rule create example.com --uri example.html --ip 1.2.3.4",
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

			return ddosxDomainWAFRuleCreate(c.DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("uri", "", "URI for rule")
	_ = cmd.MarkFlagRequired("uri")
	cmd.Flags().String("ip", "", "IP for rule")
	_ = cmd.MarkFlagRequired("ip")

	return cmd
}

func ddosxDomainWAFRuleCreate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	createRequest := ddosx.CreateWAFRuleRequest{}
	createRequest.URI, _ = cmd.Flags().GetString("uri")

	ip, _ := cmd.Flags().GetString("ip")
	createRequest.IP = connection.IPAddress(ip)

	id, err := service.CreateDomainWAFRule(args[0], createRequest)
	if err != nil {
		return fmt.Errorf("error creating domain WAF rule: %s", err)
	}

	rule, err := service.GetDomainWAFRule(args[0], id)
	if err != nil {
		return fmt.Errorf("error retrieving new domain WAF rule [%s]: %s", id, err)
	}

	return output.CommandOutput(cmd, WAFRuleCollection([]ddosx.WAFRule{rule}))
}

func ddosxDomainWAFRuleUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <domain: name> <ruleset: id>...",
		Short:   "Updates WAF rules",
		Long:    "This command updates one or more domain WAF rules",
		Example: "ans ddosx domain waf ruleset update example.com 00000000-0000-0000-0000-000000000000 --active=true",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing domain")
			}
			if len(args) < 2 {
				return errors.New("missing rule")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxDomainWAFRuleUpdate(c.DDoSXService(), cmd, args)
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

	return output.CommandOutput(cmd, WAFRuleCollection(rules))
}

func ddosxDomainWAFRuleDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <domain: name> <rule: id>...",
		Short:   "Deletes WAF rules",
		Long:    "This command deletes one or more domain WAF rules",
		Example: "ans ddosx domain waf rule delete example.com 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing domain")
			}
			if len(args) < 2 {
				return errors.New("missing rule")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			ddosxDomainWAFRuleDelete(c.DDoSXService(), cmd, args)
			return nil
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
