package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func loadtestDomainRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain",
		Short: "sub-commands relating to domains",
	}

	// Child commands
	cmd.AddCommand(loadtestDomainListCmd())
	cmd.AddCommand(loadtestDomainShowCmd())

	return cmd
}

func loadtestDomainListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists domains",
		Long:    "This command lists domains",
		Example: "ukfast loadtest domain list",
		Run: func(cmd *cobra.Command, args []string) {
			loadtestDomainList(getClient().LTaaSService(), cmd, args)
		},
	}
}

func loadtestDomainList(service ltaas.LTaaSService, cmd *cobra.Command, args []string) {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	domains, err := service.GetDomains(params)
	if err != nil {
		output.Fatalf("Error retrieving domains: %s", err)
		return
	}

	outputLoadTestDomains(domains)
}

func loadtestDomainShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: id>...",
		Short:   "Shows a domain",
		Long:    "This command shows one or more domains",
		Example: "ukfast loadtest domain show 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			loadtestDomainShow(getClient().LTaaSService(), cmd, args)
		},
	}
}

func loadtestDomainShow(service ltaas.LTaaSService, cmd *cobra.Command, args []string) {
	var domains []ltaas.Domain
	for _, arg := range args {
		domain, err := service.GetDomain(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving domain [%s]: %s", arg, err)
			continue
		}

		domains = append(domains, domain)
	}

	outputLoadTestDomains(domains)
}

func loadtestDomainCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a domain",
		Long:    "This command creates a domain ",
		Example: "ukfast loadtest domain create",
		RunE: func(cmd *cobra.Command, args []string) error {
			return loadtestDomainCreate(getClient().LTaaSService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Name of domain")
	cmd.MarkFlagRequired("domain")
	cmd.Flags().String("verification-method", "", "Verification method for domain")
	cmd.MarkFlagRequired("domain")

	return cmd
}

func loadtestDomainCreate(service ltaas.LTaaSService, cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	verificationMethod, _ := cmd.Flags().GetString("verification-method")
	parsedVerificationMethod, err := ltaas.ParseDomainVerificationMethod(verificationMethod)
	if err != nil {
		return clierrors.NewErrInvalidFlagValue("verification-method", verificationMethod, err)
	}

	createRequest := ltaas.CreateDomainRequest{
		Name:               name,
		VerificationMethod: parsedVerificationMethod,
	}

	domainID, err := service.CreateDomain(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating domain: %s", err)
	}

	domain, err := service.GetDomain(domainID)
	if err != nil {
		return fmt.Errorf("Error retrieving new domain: %s", err)
	}

	outputLoadTestDomains([]ltaas.Domain{domain})
	return nil
}

func loadtestDomainDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <domain: id>...",
		Short:   "Deletes a domain",
		Long:    "This command deletes one or more domains",
		Example: "ukfast loadtest domain delete 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			loadtestDomainDelete(getClient().LTaaSService(), cmd, args)
		},
	}
}

func loadtestDomainDelete(service ltaas.LTaaSService, cmd *cobra.Command, args []string) {
	var domains []ltaas.Domain
	for _, arg := range args {
		domain, err := service.GetDomain(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing domain [%s]: %s", arg, err)
			continue
		}

		domains = append(domains, domain)
	}

	outputLoadTestDomains(domains)
}
