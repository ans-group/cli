package ddosx

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxDomainWAFRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "waf",
		Short: "sub-commands relating to domain web application firewalls",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainWAFShowCmd(f))
	cmd.AddCommand(ddosxDomainWAFCreateCmd(f))
	cmd.AddCommand(ddosxDomainWAFUpdateCmd(f))
	cmd.AddCommand(ddosxDomainWAFDeleteCmd(f))

	// Child root commands
	cmd.AddCommand(ddosxDomainWAFRuleSetRootCmd(f))
	cmd.AddCommand(ddosxDomainWAFRuleRootCmd(f))
	cmd.AddCommand(ddosxDomainWAFAdvancedRuleRootCmd(f))

	return cmd
}

func ddosxDomainWAFShowCmd(f factory.ClientFactory) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return ddosxDomainWAFShow(f.NewClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainWAFShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	var wafs []ddosx.WAF
	for _, arg := range args {
		waf, err := service.GetDomainWAF(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving domain waf [%s]: %s", arg, err)
			continue
		}

		wafs = append(wafs, waf)
	}

	return output.CommandOutput(cmd, OutputDDoSXWAFsProvider(wafs))
}

func ddosxDomainWAFCreateCmd(f factory.ClientFactory) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return ddosxDomainWAFCreate(f.NewClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("mode", "", "Mode for WAF")
	cmd.MarkFlagRequired("mode")
	cmd.Flags().String("paranoia-level", "", "Paranoia level for WAF")
	cmd.MarkFlagRequired("paranoia-level")

	return cmd
}

func ddosxDomainWAFCreate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	mode, _ := cmd.Flags().GetString("mode")
	parsedMode, err := ddosx.ParseWAFMode(mode)
	if err != nil {
		return err
	}
	paranoiaLevel, _ := cmd.Flags().GetString("paranoia-level")
	parsedParanoiaLevel, err := ddosx.ParseWAFParanoiaLevel(paranoiaLevel)
	if err != nil {
		return err
	}

	createRequest := ddosx.CreateWAFRequest{
		Mode:          parsedMode,
		ParanoiaLevel: parsedParanoiaLevel,
	}

	err = service.CreateDomainWAF(args[0], createRequest)
	if err != nil {
		return fmt.Errorf("Error creating domain waf: %s", err)
	}

	waf, err := service.GetDomainWAF(args[0])
	if err != nil {
		return fmt.Errorf("Error retrieving domain waf: %s", err)
	}

	return output.CommandOutput(cmd, OutputDDoSXWAFsProvider([]ddosx.WAF{waf}))
}

func ddosxDomainWAFUpdateCmd(f factory.ClientFactory) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return ddosxDomainWAFUpdate(f.NewClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("mode", "", "Mode for WAF")
	cmd.Flags().String("paranoia-level", "", "Paranoia level for WAF")

	return cmd
}

func ddosxDomainWAFUpdate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	patchRequest := ddosx.PatchWAFRequest{}

	if cmd.Flags().Changed("mode") {
		mode, _ := cmd.Flags().GetString("mode")
		parsedMode, err := ddosx.ParseWAFMode(mode)
		if err != nil {
			return err
		}

		patchRequest.Mode = parsedMode
	}

	if cmd.Flags().Changed("paranoia-level") {
		paranoiaLevel, _ := cmd.Flags().GetString("paranoia-level")
		parsedParanoiaLevel, err := ddosx.ParseWAFParanoiaLevel(paranoiaLevel)
		if err != nil {
			return err
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

	return output.CommandOutput(cmd, OutputDDoSXWAFsProvider(wafs))
}

func ddosxDomainWAFDeleteCmd(f factory.ClientFactory) *cobra.Command {
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
			ddosxDomainWAFDelete(f.NewClient().DDoSXService(), cmd, args)
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
