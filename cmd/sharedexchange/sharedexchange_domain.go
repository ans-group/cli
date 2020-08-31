package sharedexchange

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/sharedexchange"
)

func sharedexchangeDomainRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain",
		Short: "sub-commands relating to domains",
	}

	// Child commands
	cmd.AddCommand(sharedexchangeDomainListCmd(f))
	cmd.AddCommand(sharedexchangeDomainShowCmd(f))

	return cmd
}

func sharedexchangeDomainListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists domains",
		Long:    "This command lists domains",
		Example: "ukfast sharedexchange domain list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return sharedexchangeDomainList(c.SharedExchangeService(), cmd, args)
		},
	}
}

func sharedexchangeDomainList(service sharedexchange.SharedExchangeService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	domains, err := service.GetDomains(params)
	if err != nil {
		return fmt.Errorf("Error retrieving domains: %s", err)
	}

	return output.CommandOutput(cmd, OutputSharedExchangeDomainsProvider(domains))
}

func sharedexchangeDomainShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: id>...",
		Short:   "Shows a domain",
		Long:    "This command shows one or more domains",
		Example: "ukfast sharedexchange domain show 123",
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

			return sharedexchangeDomainShow(c.SharedExchangeService(), cmd, args)
		},
	}
}

func sharedexchangeDomainShow(service sharedexchange.SharedExchangeService, cmd *cobra.Command, args []string) error {
	var domains []sharedexchange.Domain
	for _, arg := range args {
		domainID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid domain ID [%s]", arg)
			continue
		}

		domain, err := service.GetDomain(domainID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving domain [%s]: %s", arg, err)
			continue
		}

		domains = append(domains, domain)
	}

	return output.CommandOutput(cmd, OutputSharedExchangeDomainsProvider(domains))
}
