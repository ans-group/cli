package cmd

import (
	"errors"

	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/ptr"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxDomainACLIPRuleRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ip",
		Short: "sub-commands relating to domain ACL IP rules",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainACLIPRuleListCmd())
	cmd.AddCommand(ddosxDomainACLIPRuleShowCmd())
	cmd.AddCommand(ddosxDomainACLIPRuleCreateCmd())
	cmd.AddCommand(ddosxDomainACLIPRuleUpdateCmd())
	cmd.AddCommand(ddosxDomainACLIPRuleDeleteCmd())

	return cmd
}

func ddosxDomainACLIPRuleListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list <domain: name>",
		Short:   "Lists ACL IP rules",
		Long:    "This command lists domain ACL IP rules",
		Example: "ukfast ddosx domain acl ip list",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainACLIPRuleList(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainACLIPRuleList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	domains, err := service.GetDomainACLIPRules(args[0], params)
	if err != nil {
		output.Fatalf("Error retrieving domain ACL IP rules: %s", err)
		return
	}

	outputDDoSXACLIPRules(domains)
}

func ddosxDomainACLIPRuleShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: name> <rule: id>...",
		Short:   "Shows domain ACL IP rules",
		Long:    "This command shows an ACL IP rule",
		Example: "ukfast ddosx domain acl ip show example.com 00000000-0000-0000-0000-000000000000",
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
			ddosxDomainACLIPRuleShow(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainACLIPRuleShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {

	var rules []ddosx.ACLIPRule

	for _, arg := range args[1:] {
		rule, err := service.GetDomainACLIPRule(args[0], arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving domain ACL IP rule [%s]: %s", arg, err.Error())
			continue
		}

		rules = append(rules, rule)
	}

	outputDDoSXACLIPRules(rules)
}

func ddosxDomainACLIPRuleCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <domain: name>",
		Short:   "Creates ACL IP rules",
		Long:    "This command creates domain ACL IP rules",
		Example: "ukfast ddosx domain acl ip create",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainACLIPRuleCreate(getClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("ip", "", "IP address for IP ACL rule")
	cmd.MarkFlagRequired("ip")
	cmd.Flags().String("uri", "", "URI for IP ACL rule")
	cmd.MarkFlagRequired("uri")
	cmd.Flags().String("mode", "", "Mode for IP ACL rule")
	cmd.MarkFlagRequired("mode")

	return cmd
}

func ddosxDomainACLIPRuleCreate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	ip, _ := cmd.Flags().GetString("ip")
	mode, _ := cmd.Flags().GetString("mode")
	parsedMode, err := ddosx.ParseACLIPMode(mode)
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	createRequest := ddosx.CreateACLIPRuleRequest{}
	createRequest.IP = connection.IPAddress(ip)
	createRequest.URI, _ = cmd.Flags().GetString("uri")
	createRequest.Mode = parsedMode

	id, err := service.CreateDomainACLIPRule(args[0], createRequest)
	if err != nil {
		output.Fatalf("Error creating domain ACL IP rule: %s", err)
		return
	}

	rule, err := service.GetDomainACLIPRule(args[0], id)
	if err != nil {
		output.Fatalf("Error retrieving new domain ACL IP rule [%s]: %s", id, err)
		return
	}

	outputDDoSXACLIPRules([]ddosx.ACLIPRule{rule})
}

func ddosxDomainACLIPRuleUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <domain: name> <rule: id>...",
		Short:   "Updates ACL IP rules",
		Long:    "This command updates one or more domain ACL IP rules",
		Example: "ukfast ddosx domain acl ip update example.com 00000000-0000-0000-0000-000000000000 --ip 1.2.3.4",
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
			ddosxDomainACLIPRuleUpdate(getClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("ip", "", "IP address for IP ACL rule")
	cmd.Flags().String("uri", "", "URI for IP ACL rule")
	cmd.Flags().String("mode", "", "Mode for IP ACL rule")

	return cmd
}

func ddosxDomainACLIPRuleUpdate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	patchRequest := ddosx.PatchACLIPRuleRequest{}

	if cmd.Flags().Changed("ip") {
		ipAddress, _ := cmd.Flags().GetString("ip")
		patchRequest.IP = connection.IPAddress(ipAddress)
	}

	if cmd.Flags().Changed("uri") {
		uri, _ := cmd.Flags().GetString("uri")
		patchRequest.URI = ptr.String(uri)
	}

	if cmd.Flags().Changed("mode") {
		mode, _ := cmd.Flags().GetString("mode")
		parsedMode, err := ddosx.ParseACLIPMode(mode)
		if err != nil {
			output.Fatal(err.Error())
			return
		}
		patchRequest.Mode = parsedMode
	}

	var rules []ddosx.ACLIPRule
	for _, arg := range args[1:] {
		err := service.PatchDomainACLIPRule(args[0], arg, patchRequest)
		if err != nil {
			OutputWithErrorLevelf("Error updating domain ACL IP rule [%s]: %s", arg, err.Error())
			continue
		}

		rule, err := service.GetDomainACLIPRule(args[0], arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving updated domain ACL IP rule [%s]: %s", arg, err)
			continue
		}

		rules = append(rules, rule)
	}

	outputDDoSXACLIPRules(rules)
}

func ddosxDomainACLIPRuleDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <domain: name> <rule: id>...",
		Short:   "Deletes ACL IP rules",
		Long:    "This command deletes one or more domain ACL IP rules",
		Example: "ukfast ddosx domain acl ip delete example.com 00000000-0000-0000-0000-000000000000",
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
			ddosxDomainACLIPRuleDelete(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainACLIPRuleDelete(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	for _, arg := range args[1:] {
		err := service.DeleteDomainACLIPRule(args[0], arg)
		if err != nil {
			OutputWithErrorLevelf("Error removing domain ACL IP rule [%s]: %s", arg, err.Error())
			continue
		}
	}
}
