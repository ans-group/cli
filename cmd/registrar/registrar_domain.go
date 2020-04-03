package registrar

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/registrar"
)

func registrarDomainRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain",
		Short: "sub-commands relating to domains",
	}

	// Child commands
	cmd.AddCommand(registrarDomainListCmd(f))
	cmd.AddCommand(registrarDomainShowCmd(f))

	// Child root commands
	cmd.AddCommand(registrarDomainNameserverRootCmd(f))

	return cmd
}

func registrarDomainListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists domains",
		Long:    "This command lists domains",
		Example: "ukfast registrar domain list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return registrarDomainList(f.NewClient().RegistrarService(), cmd, args)
		},
	}
}

func registrarDomainList(service registrar.RegistrarService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	domains, err := service.GetDomains(params)
	if err != nil {
		return fmt.Errorf("Error retrieving domains: %s", err)
	}

	return output.CommandOutput(cmd, OutputRegistrarDomainsProvider(domains))
}

func registrarDomainShowCmd(f factory.ClientFactory) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return registrarDomainShow(f.NewClient().RegistrarService(), cmd, args)
		},
	}
}

func registrarDomainShow(service registrar.RegistrarService, cmd *cobra.Command, args []string) error {
	var domains []registrar.Domain
	for _, arg := range args {
		domain, err := service.GetDomain(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving domain [%s]: %s", arg, err)
			continue
		}

		domains = append(domains, domain)
	}

	return output.CommandOutput(cmd, OutputRegistrarDomainsProvider(domains))
}
