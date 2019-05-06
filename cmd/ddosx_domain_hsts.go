package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxDomainHSTSRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hsts",
		Short: "sub-commands relating to domain HSTS",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainHSTSEnableCmd())
	cmd.AddCommand(ddosxDomainHSTSDisableCmd())

	// Child root commands
	cmd.AddCommand(ddosxDomainHSTSRuleRootCmd())

	return cmd
}

func ddosxDomainHSTSEnableCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "enable <domain: name>...",
		Short:   "Enables HSTS for a domain",
		Long:    "This command enables HSTS for one or more domains",
		Example: "ukfast ddosx domain hsts enable example.com",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainHSTSEnable(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainHSTSEnable(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	var configurations []ddosx.HSTSConfiguration

	for _, arg := range args {
		err := service.AddDomainHSTSConfiguration(arg)
		if err != nil {
			OutputWithErrorLevelf("Error enabling HSTS for domain [%s]: %s", arg, err.Error())
			continue
		}

		configuration, err := service.GetDomainHSTSConfiguration(arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving updated HSTS configuration for domain [%s]: %s", arg, err)
			continue
		}

		configurations = append(configurations, configuration)
	}

	outputDDoSXHSTSConfiguration(configurations)
}

func ddosxDomainHSTSDisableCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "disable <domain: name>...",
		Short:   "Disables HSTS for a domain",
		Long:    "This command disables HSTS for one or more domains",
		Example: "ukfast ddosx domain hsts disable example.com",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainHSTSDisable(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainHSTSDisable(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	var domains []ddosx.Domain

	for _, arg := range args {
		err := service.DeleteDomainHSTSConfiguration(arg)
		if err != nil {
			OutputWithErrorLevelf("Error disabling HSTS for domain [%s]: %s", arg, err.Error())
			continue
		}

		domain, err := service.GetDomain(arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving updated HSTS configuration for domain [%s]: %s", arg, err)
			continue
		}

		domains = append(domains, domain)
	}

	outputDDoSXDomains(domains)
}
