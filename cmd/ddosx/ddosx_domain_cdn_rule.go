package ddosx

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	"github.com/spf13/cobra"
)

func ddosxDomainCDNRuleRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rule",
		Short: "sub-commands relating to CDN rules",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainCDNRuleListCmd(f))
	cmd.AddCommand(ddosxDomainCDNRuleShowCmd(f))
	cmd.AddCommand(ddosxDomainCDNRuleCreateCmd(f))
	cmd.AddCommand(ddosxDomainCDNRuleUpdateCmd(f))
	cmd.AddCommand(ddosxDomainCDNRuleDeleteCmd(f))

	return cmd
}

func ddosxDomainCDNRuleListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <domain: name>",
		Short:   "Lists domain CDN rules",
		Long:    "This command lists CDN rules",
		Example: "ans ddosx domain cdn rule list",
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

			return ddosxDomainCDNRuleList(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainCDNRuleList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	domains, err := service.GetDomainCDNRules(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving CDN rules: %s", err)
	}

	return output.CommandOutput(cmd, CDNRuleCollection(domains))
}

func ddosxDomainCDNRuleShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: name> <rule: id>...",
		Short:   "Shows CDN rules",
		Long:    "This command shows one or more CDN rules",
		Example: "ans ddosx domain cdn rule show example.com 00000000-0000-0000-0000-000000000000",
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

			return ddosxDomainCDNRuleShow(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainCDNRuleShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	var rules []ddosx.CDNRule
	for _, arg := range args[1:] {
		rule, err := service.GetDomainCDNRule(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving CDN rule [%s]: %s", arg, err)
			continue
		}

		rules = append(rules, rule)
	}

	return output.CommandOutput(cmd, CDNRuleCollection(rules))
}

func ddosxDomainCDNRuleCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <domain: name>",
		Short:   "Creates domain CDN rules",
		Long:    "This command creates domain CDN rules",
		Example: "ans ddosx domain cdn rule create example.com --uri example.html --cache-control custom --mime-type image/* --type global",
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

			return ddosxDomainCDNRuleCreate(c.DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("uri", "", "URI for rule")
	cmd.MarkFlagRequired("uri")
	cmd.Flags().String("cache-control", "", "Cache control configuration for rule")
	cmd.MarkFlagRequired("cache-control")
	cmd.Flags().String("cache-control-duration", "", "Cache control duration for rule (applicable with 'Custom' cache control), e.g. 1d4h")
	cmd.Flags().StringSlice("mime-type", []string{}, "Mime type for rule, can be repeated")
	cmd.MarkFlagRequired("mime-type")
	cmd.Flags().String("type", "", "Type of rule")
	cmd.MarkFlagRequired("type")

	return cmd
}

func ddosxDomainCDNRuleCreate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	cacheControl, _ := cmd.Flags().GetString("cache-control")
	parsedCacheControl, err := ddosx.CDNRuleCacheControlEnum.Parse(cacheControl)
	if err != nil {
		return clierrors.NewErrInvalidFlagValue("cache-control", cacheControl, err)
	}

	ruleType, _ := cmd.Flags().GetString("type")
	parsedRuleType, err := ddosx.CDNRuleTypeEnum.Parse(ruleType)
	if err != nil {
		return clierrors.NewErrInvalidFlagValue("type", ruleType, err)
	}

	createRequest := ddosx.CreateCDNRuleRequest{}
	createRequest.URI, _ = cmd.Flags().GetString("uri")
	createRequest.CacheControl = parsedCacheControl
	createRequest.MimeTypes, _ = cmd.Flags().GetStringSlice("mime-type")
	createRequest.Type = parsedRuleType

	if cmd.Flags().Changed("cache-control-duration") {
		cacheControlDuration, _ := cmd.Flags().GetString("cache-control-duration")
		parsedCacheControlDuration, err := ddosx.ParseCDNRuleCacheControlDuration(cacheControlDuration)
		if err != nil {
			return clierrors.NewErrInvalidFlagValue("cache-control-duration", cacheControlDuration, err)
		}

		createRequest.CacheControlDuration = parsedCacheControlDuration
	}

	id, err := service.CreateDomainCDNRule(args[0], createRequest)
	if err != nil {
		return fmt.Errorf("Error creating CDN rule: %s", err)
	}

	rule, err := service.GetDomainCDNRule(args[0], id)
	if err != nil {
		return fmt.Errorf("Error retrieving new CDN rule [%s]: %s", id, err.Error())
	}

	return output.CommandOutput(cmd, CDNRuleCollection([]ddosx.CDNRule{rule}))
}

func ddosxDomainCDNRuleUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <domain: name> <rule: id>...",
		Short:   "Updates CDN rules",
		Long:    "This command updates one or more domain CDN rules",
		Example: "ans ddosx domain cdn rule update example.com 00000000-0000-0000-0000-000000000000 --mime-type image/*",
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

			return ddosxDomainCDNRuleUpdate(c.DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("uri", "", "URI for rule")
	cmd.Flags().String("cache-control", "", "Cache control configuration for rule")
	cmd.Flags().String("cache-control-duration", "", "Cache control duration for rule")
	cmd.Flags().StringSlice("mime-type", []string{}, "Mime type for rule, can be repeated")
	cmd.Flags().String("type", "", "Type of rule")

	return cmd
}

func ddosxDomainCDNRuleUpdate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	patchRequest := ddosx.PatchCDNRuleRequest{}

	if cmd.Flags().Changed("uri") {
		patchRequest.URI, _ = cmd.Flags().GetString("uri")
	}

	if cmd.Flags().Changed("cache-control") {
		cacheControl, _ := cmd.Flags().GetString("cache-control")
		parsedCacheControl, err := ddosx.CDNRuleCacheControlEnum.Parse(cacheControl)
		if err != nil {
			return clierrors.NewErrInvalidFlagValue("cache-control", cacheControl, err)
		}

		patchRequest.CacheControl = parsedCacheControl
	}

	if cmd.Flags().Changed("cache-control-duration") {
		cacheControlDuration, _ := cmd.Flags().GetString("cache-control-duration")
		parsedCacheControlDuration, err := ddosx.ParseCDNRuleCacheControlDuration(cacheControlDuration)
		if err != nil {
			return clierrors.NewErrInvalidFlagValue("cache-control-duration", cacheControlDuration, err)
		}

		patchRequest.CacheControlDuration = parsedCacheControlDuration
	}

	if cmd.Flags().Changed("mime-type") {
		patchRequest.MimeTypes, _ = cmd.Flags().GetStringSlice("mime-type")
	}

	if cmd.Flags().Changed("type") {
		ruleType, _ := cmd.Flags().GetString("type")
		parsedRuleType, err := ddosx.CDNRuleTypeEnum.Parse(ruleType)
		if err != nil {
			return clierrors.NewErrInvalidFlagValue("type", ruleType, err)
		}

		patchRequest.Type = parsedRuleType
	}

	var rules []ddosx.CDNRule

	for _, arg := range args[1:] {
		err := service.PatchDomainCDNRule(args[0], arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating domain CDN rule [%s]: %s", arg, err.Error())
			continue
		}

		rule, err := service.GetDomainCDNRule(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated CDN rule [%s]: %s", arg, err.Error())
			continue
		}

		rules = append(rules, rule)
	}

	return output.CommandOutput(cmd, CDNRuleCollection(rules))
}

func ddosxDomainCDNRuleDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <domain: name> <rule: id>...",
		Short:   "Deletes CDN rules",
		Long:    "This command deletes one or more domain CDN rules",
		Example: "ans ddosx domain cdn rule delete example.com 00000000-0000-0000-0000-000000000000",
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

			ddosxDomainCDNRuleDelete(c.DDoSXService(), cmd, args)
			return nil
		},
	}
}

func ddosxDomainCDNRuleDelete(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	for _, arg := range args[1:] {
		err := service.DeleteDomainCDNRule(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing domain CDN rule [%s]: %s", arg, err.Error())
			continue
		}
	}
}
