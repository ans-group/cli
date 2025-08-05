package ddosx

import (
	"errors"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	"github.com/spf13/cobra"
)

func ddosxDomainDNSRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dns",
		Short: "sub-commands relating to domain DNS",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainDNSActivateCmd(f))
	cmd.AddCommand(ddosxDomainDNSDeactivateCmd(f))

	return cmd
}

func ddosxDomainDNSActivateCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "activate <domain: name>...",
		Short:   "Activates DNS routing for a domain",
		Long:    "This command activates DNS routing for one or more domains",
		Example: "ans ddosx domain dns activate example.com",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing domain")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxDomainDNSActivate(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainDNSActivate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	var domains []ddosx.Domain
	for _, arg := range args {
		err := service.ActivateDomainDNSRouting(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error activating DNS routing for domain [%s]: %s", arg, err)
			continue
		}

		domain, err := service.GetDomain(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving domain [%s]: %s", arg, err)
			continue
		}

		domains = append(domains, domain)
	}

	return output.CommandOutput(cmd, DomainCollection(domains))
}

func ddosxDomainDNSDeactivateCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "deactivate <domain: name>...",
		Short:   "Deactivates DNS routing for a domain",
		Long:    "This command deactivates DNS routing for one or more domains",
		Example: "ans ddosx domain dns deactivate example.com",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing domain")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxDomainDNSDeactivate(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainDNSDeactivate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	var domains []ddosx.Domain
	for _, arg := range args {
		err := service.DeactivateDomainDNSRouting(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error deactivating DNS routing for domain [%s]: %s", arg, err)
			continue
		}

		domain, err := service.GetDomain(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving domain [%s]: %s", arg, err)
			continue
		}

		domains = append(domains, domain)
	}

	return output.CommandOutput(cmd, DomainCollection(domains))
}
