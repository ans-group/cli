package ddosx

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func ddosxDomainRootCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain",
		Short: "sub-commands relating to domains",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainListCmd(f))
	cmd.AddCommand(ddosxDomainShowCmd(f))
	cmd.AddCommand(ddosxDomainCreateCmd(f))
	cmd.AddCommand(ddosxDomainDeleteCmd(f))
	cmd.AddCommand(ddosxDomainDeployCmd(f))

	// Child root commands
	cmd.AddCommand(ddosxDomainRecordRootCmd(f))
	cmd.AddCommand(ddosxDomainWAFRootCmd(f))
	cmd.AddCommand(ddosxDomainACLRootCmd(f))
	cmd.AddCommand(ddosxDomainPropertyRootCmd(f, fs))
	cmd.AddCommand(ddosxDomainVerificationRootCmd(f, fs))
	cmd.AddCommand(ddosxDomainCDNRootCmd(f))
	cmd.AddCommand(ddosxDomainHSTSRootCmd(f))
	cmd.AddCommand(ddosxDomainDNSRootCmd(f))

	return cmd
}

func ddosxDomainListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists domains",
		Long:    "This command lists domains",
		Example: "ans ddosx domain list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxDomainList(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	domains, err := service.GetDomains(params)
	if err != nil {
		return fmt.Errorf("Error retrieving domains: %s", err)
	}

	return output.CommandOutput(cmd, DomainCollection(domains))
}

func ddosxDomainShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: name>...",
		Short:   "Shows a domain",
		Long:    "This command shows one or more domains",
		Example: "ans ddosx domain show example.com",
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

			return ddosxDomainShow(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	var domains []ddosx.Domain
	for _, arg := range args {
		domain, err := service.GetDomain(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving domain [%s]: %s", arg, err)
			continue
		}

		domains = append(domains, domain)
	}

	return output.CommandOutput(cmd, DomainCollection(domains))
}

func ddosxDomainCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a domain",
		Long:    "This command creates a new domain",
		Example: "ans ddosx domain create --name example.com",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxDomainCreate(c.DDoSXService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of domain")
	cmd.MarkFlagRequired("name")

	return cmd
}

func ddosxDomainCreate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	domainName, _ := cmd.Flags().GetString("name")

	createRequest := ddosx.CreateDomainRequest{
		Name: domainName,
	}

	err := service.CreateDomain(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating domain: %s", err)
	}

	domain, err := service.GetDomain(domainName)
	if err != nil {
		return fmt.Errorf("Error retrieving new domain: %s", err)
	}

	return output.CommandOutput(cmd, DomainCollection([]ddosx.Domain{domain}))
}

func ddosxDomainDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <domain: name>...",
		Short:   "Deletes a domain",
		Long:    "This command deletes one or more domains",
		Example: "ans ddosx domain delete example.com",
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

			ddosxDomainDelete(c.DDoSXService(), cmd, args)
			return nil
		},
	}

	cmd.Flags().String("summary", "", "Specifies summary for domain removal")
	cmd.MarkFlagRequired("summary")
	cmd.Flags().String("description", "", "Specifies description for domain removal")
	cmd.MarkFlagRequired("description")

	return cmd
}

func ddosxDomainDelete(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	req := ddosx.DeleteDomainRequest{}
	req.Summary, _ = cmd.Flags().GetString("summary")
	req.Description, _ = cmd.Flags().GetString("description")

	for _, arg := range args {
		err := service.DeleteDomain(arg, req)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing domain [%s]: %s", arg, err)
			continue
		}
	}
}

func ddosxDomainDeployCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deploy <domain: name>...",
		Short:   "Deploys a domain",
		Long:    "This command deploys one or more domains",
		Example: "ans ddosx domain deploy example.com",
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

			return ddosxDomainDeploy(c.DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the domain has been completely deployed before continuing on")

	return cmd
}

func ddosxDomainDeploy(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	var domains []ddosx.Domain
	for _, arg := range args {
		err := service.DeployDomain(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error deploying domain [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(DomainStatusWaitFunc(service, arg, ddosx.DomainStatusConfigured))
			if err != nil {
				output.OutputWithErrorLevelf("Error deploying domain [%s]: %s", arg, err)
				continue
			}
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

func DomainStatusWaitFunc(service ddosx.DDoSXService, domainName string, status ddosx.DomainStatus) helper.WaitFunc {
	return func() (finished bool, err error) {
		domain, err := service.GetDomain(domainName)
		if err != nil {
			return false, fmt.Errorf("Failed to retrieve domain [%s]: %s", domainName, err)
		}
		if domain.Status == ddosx.DomainStatusFailed {
			return false, fmt.Errorf("Domain [%s] in [%s] state", domainName, domain.Status)
		}
		if domain.Status == status {
			return true, nil
		}

		return false, nil
	}
}
