package ddosx

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	"github.com/spf13/cobra"
)

func ddosxDomainCDNRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cdn",
		Short: "sub-commands relating to domain CDN",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainCDNEnableCmd(f))
	cmd.AddCommand(ddosxDomainCDNDisableCmd(f))
	cmd.AddCommand(ddosxDomainCDNPurgeCmd(f))

	// Child root commands
	cmd.AddCommand(ddosxDomainCDNRuleRootCmd(f))

	return cmd
}

func ddosxDomainCDNEnableCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "enable <domain: name>...",
		Short:   "Enables CDN for a domain",
		Long:    "This command enables CDN for one or more domains",
		Example: "ukfast ddosx domain cdn enable example.com",
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

			return ddosxDomainCDNEnable(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainCDNEnable(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	var domains []ddosx.Domain

	for _, arg := range args {
		err := service.AddDomainCDNConfiguration(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error enabling CDN for domain [%s]: %s", arg, err.Error())
			continue
		}

		domain, err := service.GetDomain(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated domain [%s]: %s", arg, err)
			continue
		}

		domains = append(domains, domain)
	}

	return output.CommandOutput(cmd, OutputDDoSXDomainsProvider(domains))
}

func ddosxDomainCDNDisableCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "disable <domain: name>...",
		Short:   "Disables CDN for a domain",
		Long:    "This command disables CDN for one or more domains",
		Example: "ukfast ddosx domain cdn disable example.com",
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

			return ddosxDomainCDNDisable(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainCDNDisable(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	var domains []ddosx.Domain

	for _, arg := range args {
		err := service.DeleteDomainCDNConfiguration(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error disabling CDN for domain [%s]: %s", arg, err.Error())
			continue
		}

		domain, err := service.GetDomain(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated domain [%s]: %s", arg, err)
			continue
		}

		domains = append(domains, domain)
	}

	return output.CommandOutput(cmd, OutputDDoSXDomainsProvider(domains))
}

func ddosxDomainCDNPurgeCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "purge <domain: name>",
		Short:   "Purges CDN content for a domain",
		Long:    "This command purges CDN content for a domain",
		Example: "ukfast ddosx domain cdn purge example.com --record-name something.example.com --uri /test",
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

			return ddosxDomainCDNPurge(c.DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("record-name", "", "Record name for purging")
	cmd.MarkFlagRequired("record-name")
	cmd.Flags().String("uri", "", "URI for purging")
	cmd.MarkFlagRequired("uri")

	return cmd
}

func ddosxDomainCDNPurge(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	purgeRequest := ddosx.PurgeCDNRequest{}
	purgeRequest.RecordName, _ = cmd.Flags().GetString("record-name")
	purgeRequest.URI, _ = cmd.Flags().GetString("uri")

	err := service.PurgeDomainCDN(args[0], purgeRequest)
	if err != nil {
		return fmt.Errorf("Error purging CDN content for domain: %s", err.Error())
	}

	return nil
}
