package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
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
	cmd := &cobra.Command{
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
			raw, _ := cmd.Flags().GetBool("raw")
			if raw {
				registrarWhoisShowRaw(getClient().RegistrarService(), cmd, args)
			} else {
				registrarWhoisShow(getClient().RegistrarService(), cmd, args)
			}
		},
	}

	cmd.Flags().Bool("raw", false, "Specifies that whois content should be returned raw")

	return cmd
}

func registrarWhoisShow(service registrar.RegistrarService, cmd *cobra.Command, args []string) {
	var whoisArr []registrar.Whois
	for _, arg := range args {
		whois, err := service.GetWhois(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving whois for domain [%s]: %s", arg, err)
			continue
		}

		whoisArr = append(whoisArr, whois)
	}

	outputRegistrarWhois(whoisArr)
}

func registrarWhoisShowRaw(service registrar.RegistrarService, cmd *cobra.Command, args []string) {
	for _, arg := range args {
		whois, err := service.GetWhoisRaw(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving raw whois for domain [%s]: %s", arg, err)
			continue
		}

		fmt.Println(whois)
	}
}
