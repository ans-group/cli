package ddosx

import (
	"errors"
	"fmt"

	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/ptr"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	"github.com/spf13/cobra"
)

func ddosxDomainACLIPRuleRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ip",
		Short: "sub-commands relating to domain ACL IP rules",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainACLIPRuleListCmd(f))
	cmd.AddCommand(ddosxDomainACLIPRuleShowCmd(f))
	cmd.AddCommand(ddosxDomainACLIPRuleCreateCmd(f))
	cmd.AddCommand(ddosxDomainACLIPRuleUpdateCmd(f))
	cmd.AddCommand(ddosxDomainACLIPRuleDeleteCmd(f))

	return cmd
}

func ddosxDomainACLIPRuleListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <domain: name>",
		Short:   "Lists ACL IP rules",
		Long:    "This command lists domain ACL IP rules",
		Example: "ans ddosx domain acl ip list",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxDomainACLIPRuleList(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainACLIPRuleList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	domains, err := service.GetDomainACLIPRules(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving domain ACL IP rules: %s", err)
	}

	return output.CommandOutput(cmd, OutputDDoSXACLIPRulesProvider(domains))
}

func ddosxDomainACLIPRuleShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: name> <rule: id>...",
		Short:   "Shows domain ACL IP rules",
		Long:    "This command shows an ACL IP rule",
		Example: "ans ddosx domain acl ip show example.com 00000000-0000-0000-0000-000000000000",
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
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxDomainACLIPRuleShow(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainACLIPRuleShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {

	var rules []ddosx.ACLIPRule

	for _, arg := range args[1:] {
		rule, err := service.GetDomainACLIPRule(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving domain ACL IP rule [%s]: %s", arg, err.Error())
			continue
		}

		rules = append(rules, rule)
	}

	return output.CommandOutput(cmd, OutputDDoSXACLIPRulesProvider(rules))
}

func ddosxDomainACLIPRuleCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <domain: name>",
		Short:   "Creates ACL IP rules",
		Long:    "This command creates domain ACL IP rules",
		Example: "ans ddosx domain acl ip create example.com --ip 1.2.3.4 --mode Deny --uri blog",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxDomainACLIPRuleCreate(c.DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("ip", "", "IP address for IP ACL rule")
	cmd.MarkFlagRequired("ip")
	cmd.Flags().String("uri", "", "Path for IP ACL rule, e.g. path/to/file.jpg")
	cmd.Flags().String("mode", "", "Mode for IP ACL rule. Valid values: "+ddosx.ACLIPModeEnum.String())
	cmd.MarkFlagRequired("mode")

	return cmd
}

func ddosxDomainACLIPRuleCreate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	ip, _ := cmd.Flags().GetString("ip")
	mode, _ := cmd.Flags().GetString("mode")
	parsedMode, err := ddosx.ParseACLIPMode(mode)
	if err != nil {
		return clierrors.NewErrInvalidFlagValue("mode", mode, err)
	}

	createRequest := ddosx.CreateACLIPRuleRequest{}
	createRequest.IP = connection.IPAddress(ip)
	createRequest.URI, _ = cmd.Flags().GetString("uri")
	createRequest.Mode = parsedMode

	id, err := service.CreateDomainACLIPRule(args[0], createRequest)
	if err != nil {
		return fmt.Errorf("Error creating domain ACL IP rule: %s", err)
	}

	rule, err := service.GetDomainACLIPRule(args[0], id)
	if err != nil {
		return fmt.Errorf("Error retrieving new domain ACL IP rule [%s]: %s", id, err)
	}

	return output.CommandOutput(cmd, OutputDDoSXACLIPRulesProvider([]ddosx.ACLIPRule{rule}))
}

func ddosxDomainACLIPRuleUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <domain: name> <rule: id>...",
		Short:   "Updates ACL IP rules",
		Long:    "This command updates one or more domain ACL IP rules",
		Example: "ans ddosx domain acl ip update example.com 00000000-0000-0000-0000-000000000000 --ip 1.2.3.4",
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
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxDomainACLIPRuleUpdate(c.DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("ip", "", "IP address for IP ACL rule")
	cmd.Flags().String("uri", "", "URI for IP ACL rule")
	cmd.Flags().String("mode", "", "Mode for IP ACL rule")

	return cmd
}

func ddosxDomainACLIPRuleUpdate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
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
			return err
		}
		patchRequest.Mode = parsedMode
	}

	var rules []ddosx.ACLIPRule
	for _, arg := range args[1:] {
		err := service.PatchDomainACLIPRule(args[0], arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating domain ACL IP rule [%s]: %s", arg, err.Error())
			continue
		}

		rule, err := service.GetDomainACLIPRule(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated domain ACL IP rule [%s]: %s", arg, err)
			continue
		}

		rules = append(rules, rule)
	}

	return output.CommandOutput(cmd, OutputDDoSXACLIPRulesProvider(rules))
}

func ddosxDomainACLIPRuleDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <domain: name> <rule: id>...",
		Short:   "Deletes ACL IP rules",
		Long:    "This command deletes one or more domain ACL IP rules",
		Example: "ans ddosx domain acl ip delete example.com 00000000-0000-0000-0000-000000000000",
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
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			ddosxDomainACLIPRuleDelete(c.DDoSXService(), cmd, args)
			return nil
		},
	}
}

func ddosxDomainACLIPRuleDelete(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	for _, arg := range args[1:] {
		err := service.DeleteDomainACLIPRule(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing domain ACL IP rule [%s]: %s", arg, err.Error())
			continue
		}
	}
}
