package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxDomainRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain",
		Short: "sub-commands relating to domains",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainListCmd())
	cmd.AddCommand(ddosxDomainShowCmd())
	cmd.AddCommand(ddosxDomainCreateCmd())
	cmd.AddCommand(ddosxDomainDeployCmd())

	// Child root commands
	cmd.AddCommand(ddosxDomainRecordRootCmd())
	cmd.AddCommand(ddosxDomainWAFRootCmd())
	cmd.AddCommand(ddosxDomainACLRootCmd())
	cmd.AddCommand(ddosxDomainPropertyRootCmd())
	cmd.AddCommand(ddosxDomainVerificationRootCmd())
	cmd.AddCommand(ddosxDomainCDNRootCmd())

	return cmd
}

func ddosxDomainListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists domains",
		Long:    "This command lists domains",
		Example: "ukfast ddosx domain list",
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainList(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	domains, err := service.GetDomains(params)
	if err != nil {
		output.Fatalf("Error retrieving domains: %s", err)
		return
	}

	outputDDoSXDomains(domains)
}

func ddosxDomainShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: name>...",
		Short:   "Shows a domain",
		Long:    "This command shows one or more domains",
		Example: "ukfast ddosx domain show example.com",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainShow(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	var domains []ddosx.Domain
	for _, arg := range args {
		domain, err := service.GetDomain(arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving domain [%s]: %s", arg, err)
			continue
		}

		domains = append(domains, domain)
	}

	outputDDoSXDomains(domains)
}

func ddosxDomainCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a domain",
		Long:    "This command creates a new domain",
		Example: "ukfast ddosx domain create --name example.com",
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainCreate(getClient().DDoSXService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of domain")
	cmd.MarkFlagRequired("name")

	return cmd
}

func ddosxDomainCreate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	domainName, _ := cmd.Flags().GetString("name")

	createRequest := ddosx.CreateDomainRequest{
		Name: domainName,
	}

	err := service.CreateDomain(createRequest)
	if err != nil {
		output.Fatalf("Error creating domain: %s", err)
		return
	}

	domain, err := service.GetDomain(domainName)
	if err != nil {
		output.Fatalf("Error retrieving new domain: %s", err)
		return
	}

	outputDDoSXDomains([]ddosx.Domain{domain})
}

func ddosxDomainDeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deploy <domain: name>...",
		Short:   "Deploys a domain",
		Long:    "This command deploys one or more domains",
		Example: "ukfast ddosx domain deploy example.com",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainDeploy(getClient().DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the domain has been completely deployed before continuing on")

	return cmd
}

func ddosxDomainDeploy(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	var domains []ddosx.Domain
	for _, arg := range args {
		err := service.DeployDomain(arg)
		if err != nil {
			OutputWithErrorLevelf("Error deploying domain [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := WaitForCommand(DomainStatusWaitFunc(service, arg, ddosx.DomainStatusConfigured))
			if err != nil {
				OutputWithErrorLevelf("Error deploying domain [%s]: %s", arg, err)
				continue
			}
		}

		domain, err := service.GetDomain(arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving domain [%s]: %s", arg, err)
			continue
		}

		domains = append(domains, domain)
	}

	outputDDoSXDomains(domains)
}

func DomainStatusWaitFunc(service ddosx.DDoSXService, domainName string, status ddosx.DomainStatus) WaitFunc {
	return func() (finished bool, err error) {
		domain, err := service.GetDomain(domainName)
		if err != nil {
			return false, fmt.Errorf("Failed to retrieve domain [%s]: %s", domainName, err)
		}
		if domain.Status == ddosx.DomainStatusFailed {
			return false, fmt.Errorf("Domain [%s] in [%s] state", domainName, domain.Status.String())
		}
		if domain.Status == status {
			return true, nil
		}

		return false, nil
	}
}
