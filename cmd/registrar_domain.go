package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/registrar"
)

func registrarDomainRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain",
		Short: "sub-commands relating to domains",
	}

	// Child commands
	cmd.AddCommand(registrarDomainListCmd())
	cmd.AddCommand(registrarDomainShowCmd())

	// Child root commands
	cmd.AddCommand(registrarDomainNameserverRootCmd())

	return cmd
}

func registrarDomainListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists domains",
		Long:    "This command lists domains",
		Example: "ukfast registrar domain list",
		Run: func(cmd *cobra.Command, args []string) {
			registrarDomainList(getClient().RegistrarService(), cmd, args)
		},
	}
}

func registrarDomainList(service registrar.RegistrarService, cmd *cobra.Command, args []string) {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	domains, err := service.GetDomains(params)
	if err != nil {
		output.Fatalf("Error retrieving domains: %s", err)
		return
	}

	outputRegistrarDomains(domains)
}

func registrarDomainShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: name>...",
		Short:   "Shows a domain",
		Long:    "This command shows one or more domains",
		Example: "ukfast registrar domain show example.com",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			registrarDomainShow(getClient().RegistrarService(), cmd, args)
		},
	}
}

func registrarDomainShow(service registrar.RegistrarService, cmd *cobra.Command, args []string) {
	var domains []registrar.Domain
	for _, arg := range args {
		domain, err := service.GetDomain(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving domain [%s]: %s", arg, err)
			continue
		}

		domains = append(domains, domain)
	}

	outputRegistrarDomains(domains)
}
