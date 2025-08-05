package registrar

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/registrar"
	"github.com/spf13/cobra"
)

func registrarWhoisRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whois",
		Short: "sub-commands relating to whois",
	}

	// Child commands
	cmd.AddCommand(registrarWhoisShowCmd(f))

	return cmd
}

func registrarWhoisShowCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show <domain: name>...",
		Short:   "Shows whois for a domain",
		Long:    "This command shows whois for one or more domains",
		Example: "ans registrar whois show example.com",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing domain")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Flags().GetBool("raw")
			if raw {
				c, err := f.NewClient()
				if err != nil {
					return err
				}

				registrarWhoisShowRaw(c.RegistrarService(), cmd, args)
				return nil
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return registrarWhoisShow(c.RegistrarService(), cmd, args)
		},
	}

	cmd.Flags().Bool("raw", false, "Specifies that whois content should be returned raw")

	return cmd
}

func registrarWhoisShow(service registrar.RegistrarService, cmd *cobra.Command, args []string) error {
	var whoisArr []registrar.Whois
	for _, arg := range args {
		whois, err := service.GetWhois(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving whois for domain [%s]: %s", arg, err)
			continue
		}

		whoisArr = append(whoisArr, whois)
	}

	return output.CommandOutput(cmd, WhoisCollection(whoisArr))
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
