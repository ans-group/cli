package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxDomainHSTSRuleRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rule",
		Short: "sub-commands relating to HSTS rules",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainHSTSRuleListCmd())
	cmd.AddCommand(ddosxDomainHSTSRuleShowCmd())
	cmd.AddCommand(ddosxDomainHSTSRuleCreateCmd())
	cmd.AddCommand(ddosxDomainHSTSRuleUpdateCmd())
	cmd.AddCommand(ddosxDomainHSTSRuleDeleteCmd())

	return cmd
}

func ddosxDomainHSTSRuleListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list <domain: name>",
		Short:   "Lists domain HSTS rules",
		Long:    "This command lists HSTS rules",
		Example: "ukfast ddosx domain hsts rule list",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainHSTSRuleList(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainHSTSRuleList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	domains, err := service.GetDomainHSTSRules(args[0], params)
	if err != nil {
		output.Fatalf("Error retrieving HSTS rules: %s", err)
		return
	}

	outputDDoSXHSTSRules(domains)
}

func ddosxDomainHSTSRuleShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: name> <rule: id>...",
		Short:   "Shows HSTS rules",
		Long:    "This command shows one or more HSTS rules",
		Example: "ukfast ddosx domain hsts rule show example.com 00000000-0000-0000-0000-000000000000",
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
			ddosxDomainHSTSRuleShow(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainHSTSRuleShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	var rules []ddosx.HSTSRule
	for _, arg := range args[1:] {
		rule, err := service.GetDomainHSTSRule(args[0], arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving HSTS rule [%s]: %s", arg, err)
			continue
		}

		rules = append(rules, rule)
	}

	outputDDoSXHSTSRules(rules)
}

func ddosxDomainHSTSRuleCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <domain: name>",
		Short:   "Creates domain HSTS rules",
		Long:    "This command creates domain HSTS rules",
		Example: "ukfast ddosx domain hsts rule create --uri example.html --cache-control custom --mime-type image/* --type global",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainHSTSRuleCreate(getClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().Int("max-age", 0, "Max age for rule")
	cmd.Flags().Bool("preload", false, "Specifies preload should be enabled")
	cmd.Flags().Bool("include-subdomains", false, "Specifies subdomains should be included")
	cmd.Flags().String("type", "", "Type of rule")
	cmd.MarkFlagRequired("type")
	cmd.Flags().String("record-name", "", "Specifies name of record")

	return cmd
}

func ddosxDomainHSTSRuleCreate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	ruleType, _ := cmd.Flags().GetString("type")
	parsedRuleType, err := ddosx.ParseHSTSRuleType(ruleType)
	if err != nil {
		output.Fatal(clierrors.InvalidFlagValueString("type", ruleType, err))
		return
	}

	createRequest := ddosx.CreateHSTSRuleRequest{}
	createRequest.MaxAge, _ = cmd.Flags().GetInt("max-age")
	createRequest.Preload, _ = cmd.Flags().GetBool("preload")
	createRequest.IncludeSubdomains, _ = cmd.Flags().GetBool("include-subdomains")
	createRequest.Type = parsedRuleType

	if cmd.Flags().Changed("record-name") {
		recordName, _ := cmd.Flags().GetString("record-name")
		createRequest.RecordName = &recordName
	}

	id, err := service.CreateDomainHSTSRule(args[0], createRequest)
	if err != nil {
		output.Fatalf("Error creating HSTS rule: %s", err)
		return
	}

	rule, err := service.GetDomainHSTSRule(args[0], id)
	if err != nil {
		output.Fatalf("Error retrieving new HSTS rule [%s]: %s", id, err.Error())
		return
	}

	outputDDoSXHSTSRules([]ddosx.HSTSRule{rule})
}

func ddosxDomainHSTSRuleUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <domain: name> <rule: id>...",
		Short:   "Updates HSTS rules",
		Long:    "This command updates one or more domain HSTS rules",
		Example: "ukfast ddosx domain hsts rule update example.com 00000000-0000-0000-0000-000000000000 --mime-type image/*",
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
			ddosxDomainHSTSRuleUpdate(getClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().Int("max-age", 0, "Max age for rule")
	cmd.Flags().Bool("preload", false, "Specifies preload should be enabled")
	cmd.Flags().Bool("include-subdomains", false, "Specifies subdomains should be included")

	return cmd
}

func ddosxDomainHSTSRuleUpdate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	patchRequest := ddosx.PatchHSTSRuleRequest{}

	if cmd.Flags().Changed("max-age") {
		maxAge, _ := cmd.Flags().GetInt("max-age")
		patchRequest.MaxAge = &maxAge
	}

	if cmd.Flags().Changed("preload") {
		preload, _ := cmd.Flags().GetBool("preload")
		patchRequest.Preload = &preload
	}
	if cmd.Flags().Changed("include-subdomains") {
		includeSubdomains, _ := cmd.Flags().GetBool("include-subdomains")
		patchRequest.IncludeSubdomains = &includeSubdomains
	}

	var rules []ddosx.HSTSRule

	for _, arg := range args[1:] {
		err := service.PatchDomainHSTSRule(args[0], arg, patchRequest)
		if err != nil {
			OutputWithErrorLevelf("Error updating domain HSTS rule [%s]: %s", arg, err.Error())
			continue
		}

		rule, err := service.GetDomainHSTSRule(args[0], arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving updated HSTS rule [%s]: %s", arg, err.Error())
			return
		}

		rules = append(rules, rule)
	}

	outputDDoSXHSTSRules(rules)
}

func ddosxDomainHSTSRuleDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <domain: name> <rule: id>...",
		Short:   "Deletes HSTS rules",
		Long:    "This command deletes one or more domain HSTS rules",
		Example: "ukfast ddosx domain hsts rule delete example.com 00000000-0000-0000-0000-000000000000",
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
			ddosxDomainHSTSRuleDelete(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainHSTSRuleDelete(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	for _, arg := range args[1:] {
		err := service.DeleteDomainHSTSRule(args[0], arg)
		if err != nil {
			OutputWithErrorLevelf("Error removing domain HSTS rule [%s]: %s", arg, err.Error())
			continue
		}
	}
}
