package registrar

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/registrar"
	"github.com/spf13/cobra"
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
		Example: "ans registrar domain nameserver list example.com",
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

			return registrarDomainNameserverList(c.RegistrarService(), cmd, args)
		},
	}
}

func registrarDomainNameserverList(service registrar.RegistrarService, cmd *cobra.Command, args []string) error {
	nameservers, err := service.GetDomainNameservers(args[0])
	if err != nil {
		return fmt.Errorf("error retrieving domain nameservers: %s", err)
	}

	return output.CommandOutput(cmd, NameserverCollection(nameservers))
}
