package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxDomainVerificationDNSRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dns",
		Short: "sub-commands relating to DNS domain verification",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainVerificationDNSVerifyCmd())

	return cmd
}

func ddosxDomainVerificationDNSVerifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "verify <domain: name>...",
		Short:   "Verifies a domain via DNS verification method",
		Long:    "This command verifies one or more domains via the DNS verification method",
		Example: "ukfast ddosx domain verification dns verify example.com",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return ddosxDomainVerificationDNSVerify(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainVerificationDNSVerify(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		err := service.VerifyDomainDNS(arg)
		if err != nil {
			OutputWithErrorLevelf("Error verifying domain [%s] via DNS verification method: %s", arg, err)
			continue
		}
	}

	return nil
}
