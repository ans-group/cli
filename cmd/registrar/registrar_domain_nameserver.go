package registrar

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/registrar"
)

func registrarDomainNameserverRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nameserver",
		Short: "sub-commands relating to domain nameservers",
	}

	// Child commands
	cmd.AddCommand(registrarDomainNameserverListCmd(f))

	return cmd
}

func registrarDomainNameserverListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists domain nameservers",
		Long:    "This command lists a domain's nameservers",
		Example: "ukfast registrar domain nameserver list example.com",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return registrarDomainNameserverList(f.NewClient().RegistrarService(), cmd, args)
		},
	}
}

func registrarDomainNameserverList(service registrar.RegistrarService, cmd *cobra.Command, args []string) error {
	nameservers, err := service.GetDomainNameservers(args[0])
	if err != nil {
		return fmt.Errorf("Error retrieving domain nameservers: %s", err)
	}

	return output.CommandOutput(cmd, OutputRegistrarNameserversProvider(nameservers))
}
