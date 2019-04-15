package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxDomainCDNRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cdn",
		Short: "sub-commands relating to domain CDN",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainCDNEnableCmd())
	cmd.AddCommand(ddosxDomainCDNDisableCmd())

	// Child root commands
	cmd.AddCommand(ddosxDomainCDNRuleRootCmd())

	return cmd
}

func ddosxDomainCDNEnableCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "enable <domain: name>...",
		Short:   "Enables CDN for a domain",
		Long:    "This command enables CDN for one or more domains",
		Example: "ukfast ddosx domain cdn enable example.com",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainCDNEnable(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainCDNEnable(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	var domains []ddosx.Domain

	for _, arg := range args {
		err := service.AddDomainCDNConfiguration(arg)
		if err != nil {
			OutputWithErrorLevelf("Error enabling CDN for domain [%s]: %s", arg, err.Error())
			continue
		}

		domain, err := service.GetDomain(arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving updated domain [%s]: %s", arg, err)
			continue
		}

		domains = append(domains, domain)
	}

	outputDDoSXDomains(domains)
}

func ddosxDomainCDNDisableCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "disable <domain: name>...",
		Short:   "Disables CDN for a domain",
		Long:    "This command disables CDN for one or more domains",
		Example: "ukfast ddosx domain cdn disable example.com",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainCDNDisable(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainCDNDisable(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	var domains []ddosx.Domain

	for _, arg := range args {
		err := service.DeleteDomainCDNConfiguration(arg)
		if err != nil {
			OutputWithErrorLevelf("Error disabling CDN for domain [%s]: %s", arg, err.Error())
			continue
		}

		domain, err := service.GetDomain(arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving updated domain [%s]: %s", arg, err)
			continue
		}

		domains = append(domains, domain)
	}

	outputDDoSXDomains(domains)
}
