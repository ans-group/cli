package cmd

import (
	"errors"
	"time"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxDomainCDNRuleRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rule",
		Short: "sub-commands relating to CDN rules",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainCDNRuleListCmd())
	cmd.AddCommand(ddosxDomainCDNRuleShowCmd())
	cmd.AddCommand(ddosxDomainCDNRuleCreateCmd())
	cmd.AddCommand(ddosxDomainCDNRuleUpdateCmd())
	cmd.AddCommand(ddosxDomainCDNRuleDeleteCmd())

	return cmd
}

func ddosxDomainCDNRuleListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list <domain: name>",
		Short:   "Lists domain CDN rules",
		Long:    "This command lists CDN rules",
		Example: "ukfast ddosx domain cdn rule list",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainCDNRuleList(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainCDNRuleList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	domains, err := service.GetDomainCDNRules(args[0], params)
	if err != nil {
		output.Fatalf("Error retrieving CDN rules: %s", err)
		return
	}

	outputDDoSXCDNRules(domains)
}

func ddosxDomainCDNRuleShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: name> <rule: id>...",
		Short:   "Shows CDN rules",
		Long:    "This command shows one or more CDN rules",
		Example: "ukfast ddosx domain cdn rule show example.com 00000000-0000-0000-0000-000000000000",
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
			ddosxDomainCDNRuleShow(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainCDNRuleShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	var rules []ddosx.CDNRule
	for _, arg := range args[1:] {
		rule, err := service.GetDomainCDNRule(args[0], arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving CDN rule [%s]: %s", arg, err)
			continue
		}

		rules = append(rules, rule)
	}

	outputDDoSXCDNRules(rules)
}

func ddosxDomainCDNRuleCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <domain: name>",
		Short:   "Creates domain CDN rules",
		Long:    "This command creates domain CDN rules",
		Example: "ukfast ddosx domain cdn rule create --uri example.html --cache-control custom --mime-type image/* --type global",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainCDNRuleCreate(getClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("uri", "", "URI for rule")
	cmd.MarkFlagRequired("uri")
	cmd.Flags().String("cache-control", "", "Cache control configuration for rule")
	cmd.MarkFlagRequired("cache-control")
	cmd.Flags().String("cache-control-duration", "", "Cache control duration for rule (applicable with 'Custom' cache control), e.g. 4h")
	cmd.Flags().StringSlice("mime-type", []string{}, "Mime type for rule, can be repeated")
	cmd.MarkFlagRequired("mime-type")
	cmd.Flags().String("type", "", "Type of rule")
	cmd.MarkFlagRequired("type")

	return cmd
}

func ddosxDomainCDNRuleCreate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	cacheControl, _ := cmd.Flags().GetString("cache-control")
	parsedCacheControl, err := ddosx.ParseCDNRuleCacheControl(cacheControl)
	if err != nil {
		output.Fatal(clierrors.InvalidFlagValueString("cache-control", cacheControl, err))
		return
	}

	ruleType, _ := cmd.Flags().GetString("type")
	parsedRuleType, err := ddosx.ParseCDNRuleType(ruleType)
	if err != nil {
		output.Fatal(clierrors.InvalidFlagValueString("type", ruleType, err))
		return
	}

	createRequest := ddosx.CreateCDNRuleRequest{}
	createRequest.URI, _ = cmd.Flags().GetString("uri")
	createRequest.CacheControl = parsedCacheControl
	createRequest.MimeTypes, _ = cmd.Flags().GetStringSlice("mime-type")
	createRequest.Type = parsedRuleType

	if cmd.Flags().Changed("cache-control-duration") {
		cacheControlDuration, _ := cmd.Flags().GetString("cache-control-duration")
		parsedCacheControlDuration, err := time.ParseDuration(cacheControlDuration)
		if err != nil {
			output.Fatal(clierrors.InvalidFlagValueString("cache-control-duration", cacheControlDuration, err))
			return
		}

		createRequest.CacheControlDuration = ddosx.NewCDNRuleCacheControlDurationFromDuration(parsedCacheControlDuration)
	}

	id, err := service.CreateDomainCDNRule(args[0], createRequest)
	if err != nil {
		output.Fatalf("Error creating CDN rule: %s", err)
		return
	}

	rule, err := service.GetDomainCDNRule(args[0], id)
	if err != nil {
		output.Fatalf("Error retrieving new CDN rule [%s]: %s", id, err.Error())
		return
	}

	outputDDoSXCDNRules([]ddosx.CDNRule{rule})
}

func ddosxDomainCDNRuleUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <domain: name> <rule: id>...",
		Short:   "Updates CDN rules",
		Long:    "This command updates one or more domain CDN rules",
		Example: "ukfast ddosx domain cdn rule update example.com 00000000-0000-0000-0000-000000000000 --mime-type image/*",
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
			ddosxDomainCDNRuleUpdate(getClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("uri", "", "URI for rule")
	cmd.Flags().String("cache-control", "", "Cache control configuration for rule")
	cmd.Flags().String("cache-control-duration", "", "Cache control duration for rule")
	cmd.Flags().StringSlice("mime-type", []string{}, "Mime type for rule, can be repeated")
	cmd.Flags().String("type", "", "Type of rule")

	return cmd
}

func ddosxDomainCDNRuleUpdate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	patchRequest := ddosx.PatchCDNRuleRequest{}

	if cmd.Flags().Changed("uri") {
		patchRequest.URI, _ = cmd.Flags().GetString("uri")
	}

	if cmd.Flags().Changed("cache-control") {
		cacheControl, _ := cmd.Flags().GetString("cache-control")
		parsedCacheControl, err := ddosx.ParseCDNRuleCacheControl(cacheControl)
		if err != nil {
			output.Fatal(err.Error())
			return
		}

		patchRequest.CacheControl = parsedCacheControl
	}

	if cmd.Flags().Changed("cache-control-duration") {
		cacheControlDuration, _ := cmd.Flags().GetString("cache-control-duration")
		parsedCacheControlDuration, err := time.ParseDuration(cacheControlDuration)
		if err != nil {
			output.Fatal(err.Error())
			return
		}

		patchRequest.CacheControlDuration = ddosx.NewCDNRuleCacheControlDurationFromDuration(parsedCacheControlDuration)
	}

	if cmd.Flags().Changed("mime-type") {
		patchRequest.MimeTypes, _ = cmd.Flags().GetStringSlice("mime-type")
	}

	if cmd.Flags().Changed("type") {
		ruleType, _ := cmd.Flags().GetString("type")
		parsedRuleType, err := ddosx.ParseCDNRuleType(ruleType)
		if err != nil {
			output.Fatal(err.Error())
			return
		}

		patchRequest.Type = parsedRuleType
	}

	var rules []ddosx.CDNRule

	for _, arg := range args[1:] {
		err := service.PatchDomainCDNRule(args[0], arg, patchRequest)
		if err != nil {
			OutputWithErrorLevelf("Error updating domain CDN rule [%s]: %s", arg, err.Error())
			continue
		}

		rule, err := service.GetDomainCDNRule(args[0], arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving updated CDN rule [%s]: %s", arg, err.Error())
			return
		}

		rules = append(rules, rule)
	}

	outputDDoSXCDNRules(rules)
}

func ddosxDomainCDNRuleDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <domain: name> <rule: id>...",
		Short:   "Deletes CDN rules",
		Long:    "This command deletes one or more domain CDN rules",
		Example: "ukfast ddosx domain cdn rule delete example.com 00000000-0000-0000-0000-000000000000",
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
			ddosxDomainCDNRuleDelete(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainCDNRuleDelete(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	for _, arg := range args[1:] {
		err := service.DeleteDomainCDNRule(args[0], arg)
		if err != nil {
			OutputWithErrorLevelf("Error removing domain CDN rule [%s]: %s", arg, err.Error())
			continue
		}
	}
}
