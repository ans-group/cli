package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxDomainWAFRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "waf",
		Short: "sub-commands relating to domain web application firewalls",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainWAFShowCmd())
	cmd.AddCommand(ddosxDomainWAFCreateCmd())
	cmd.AddCommand(ddosxDomainWAFUpdateCmd())
	cmd.AddCommand(ddosxDomainWAFDeleteCmd())

	// Child root commands
	cmd.AddCommand(ddosxDomainWAFRuleSetRootCmd())
	cmd.AddCommand(ddosxDomainWAFRuleRootCmd())
	cmd.AddCommand(ddosxDomainWAFAdvancedRuleRootCmd())

	return cmd
}

func ddosxDomainWAFShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: name>...",
		Short:   "Shows a domain WAF",
		Long:    "This command shows one or more domain WAFs",
		Example: "ukfast ddosx domain waf show example.com",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainWAFShow(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainWAFShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	var wafs []ddosx.WAF
	for _, arg := range args {
		waf, err := service.GetDomainWAF(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving domain waf [%s]: %s", arg, err)
			continue
		}

		wafs = append(wafs, waf)
	}

	outputDDoSXWAFs(wafs)
}

func ddosxDomainWAFCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <domain: name>",
		Short:   "Creates a domain WAF",
		Long:    "This command creates a domain WAF",
		Example: "ukfast ddosx domain waf create example.com --mode on --paranoia-level high",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainWAFCreate(getClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("mode", "", "Mode for WAF")
	cmd.MarkFlagRequired("mode")
	cmd.Flags().String("paranoia-level", "", "Paranoia level for WAF")
	cmd.MarkFlagRequired("paranoia-level")

	return cmd
}

func ddosxDomainWAFCreate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	mode, _ := cmd.Flags().GetString("mode")
	parsedMode, err := ddosx.ParseWAFMode(mode)
	if err != nil {
		output.Fatal(err.Error())
		return
	}
	paranoiaLevel, _ := cmd.Flags().GetString("paranoia-level")
	parsedParanoiaLevel, err := ddosx.ParseWAFParanoiaLevel(paranoiaLevel)
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	createRequest := ddosx.CreateWAFRequest{
		Mode:          parsedMode,
		ParanoiaLevel: parsedParanoiaLevel,
	}

	err = service.CreateDomainWAF(args[0], createRequest)
	if err != nil {
		output.Fatalf("Error creating domain waf: %s", err)
		return
	}

	waf, err := service.GetDomainWAF(args[0])
	if err != nil {
		output.Fatalf("Error retrieving domain waf: %s", err)
		return
	}

	outputDDoSXWAFs([]ddosx.WAF{waf})
}

func ddosxDomainWAFUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <domain: name>..",
		Short:   "Updates a domain WAF",
		Long:    "This command updates one or more domain WAFs",
		Example: "ukfast ddosx domain waf update example.com --mode on",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainWAFUpdate(getClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("mode", "", "Mode for WAF")
	cmd.Flags().String("paranoia-level", "", "Paranoia level for WAF")

	return cmd
}

func ddosxDomainWAFUpdate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	patchRequest := ddosx.PatchWAFRequest{}

	if cmd.Flags().Changed("mode") {
		mode, _ := cmd.Flags().GetString("mode")
		parsedMode, err := ddosx.ParseWAFMode(mode)
		if err != nil {
			output.Fatal(err.Error())
			return
		}

		patchRequest.Mode = parsedMode
	}

	if cmd.Flags().Changed("paranoia-level") {
		paranoiaLevel, _ := cmd.Flags().GetString("paranoia-level")
		parsedParanoiaLevel, err := ddosx.ParseWAFParanoiaLevel(paranoiaLevel)
		if err != nil {
			output.Fatal(err.Error())
			return
		}

		patchRequest.ParanoiaLevel = parsedParanoiaLevel
	}

	var wafs []ddosx.WAF
	for _, arg := range args {
		err := service.PatchDomainWAF(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating domain waf [%s]: %s", arg, err)
			continue
		}

		waf, err := service.GetDomainWAF(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated domain waf [%s]: %s", arg, err)
			continue
		}

		wafs = append(wafs, waf)
	}

	outputDDoSXWAFs(wafs)
}

func ddosxDomainWAFDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <domain: name>...",
		Short:   "Deletes a domain WAF",
		Long:    "This command deletes one or more domain WAFs",
		Example: "ukfast ddosx domain waf delete example.com",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainWAFDelete(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainWAFDelete(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	for _, arg := range args {
		err := service.DeleteDomainWAF(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing domain waf [%s]: %s", arg, err)
			continue
		}
	}
}
