package cmd

import (
	"errors"

	"github.com/spf13/cobra"
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
			OutputWithErrorLevelf("Error retrieving domain [%s]: %s", arg, err)
			continue
		}

		domains = append(domains, domain)
	}

	outputLoadTestDomains(domains)
}
