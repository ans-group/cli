package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/registrar"
)

func registrarDomainNameserverRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nameserver",
		Short: "sub-commands relating to domain nameservers",
	}

	// Child commands
	cmd.AddCommand(registrarDomainNameserverListCmd())

	return cmd
}

func registrarDomainNameserverListCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			registrarDomainNameserverList(getClient().RegistrarService(), cmd, args)
		},
	}
}

func registrarDomainNameserverList(service registrar.RegistrarService, cmd *cobra.Command, args []string) {
	nameservers, err := service.GetDomainNameservers(args[0])
	if err != nil {
		output.Fatalf("Error retrieving domain nameservers: %s", err)
		return
	}

	outputRegistrarNameservers(nameservers)
}
