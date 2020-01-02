package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func loadtestDomainVerificationDNSRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dns",
		Short: "sub-commands relating to DNS domain verification",
	}

	// Child commands
	cmd.AddCommand(loadtestDomainVerificationDNSVerifyCmd())

	return cmd
}

func loadtestDomainVerificationDNSVerifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "verify <domain: id>...",
		Short:   "Verifies a domain via DNS verification method",
		Long:    "This command verifies one or more domains via the DNS verification method",
		Example: "ukfast loadtest domain verification dns verify 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return loadtestDomainVerificationDNSVerify(getClient().LTaaSService(), cmd, args)
		},
	}
}

func loadtestDomainVerificationDNSVerify(service ltaas.LTaaSService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		err := service.VerifyDomainDNS(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error verifying domain [%s] via DNS verification method: %s", arg, err)
			continue
		}
	}

	return nil
}
