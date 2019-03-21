package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/sdk-go/pkg/service/registrar"
)

func registrarWhoisRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whois",
		Short: "sub-commands relating to whois",
	}

	// Child commands
	cmd.AddCommand(registrarWhoisShowCmd())

	return cmd
}

func registrarWhoisShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: name>...",
		Short:   "Shows whois for a domain",
		Long:    "This command shows whois for one or more domains",
		Example: "ukfast registrar whois show example.com",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			registrarWhoisShow(getClient().RegistrarService(), cmd, args)
		},
	}
}

func registrarWhoisShow(service registrar.RegistrarService, cmd *cobra.Command, args []string) {
	var whoisArr []registrar.Whois
	for _, arg := range args {
		whois, err := service.GetWhois(arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving whois for domain [%s]: %s", arg, err)
			continue
		}

		whoisArr = append(whoisArr, whois)
	}

	outputRegistrarWhois(whoisArr)
}
