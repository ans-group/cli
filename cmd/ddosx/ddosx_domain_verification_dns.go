package ddosx

import (
	"errors"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	"github.com/spf13/cobra"
)

func ddosxDomainVerificationDNSRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dns",
		Short: "sub-commands relating to DNS domain verification",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainVerificationDNSVerifyCmd(f))

	return cmd
}

func ddosxDomainVerificationDNSVerifyCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "verify <domain: name>...",
		Short:   "Verifies a domain via DNS verification method",
		Long:    "This command verifies one or more domains via the DNS verification method",
		Example: "ans ddosx domain verification dns verify example.com",
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

			return ddosxDomainVerificationDNSVerify(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainVerificationDNSVerify(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		err := service.VerifyDomainDNS(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error verifying domain [%s] via DNS verification method: %s", arg, err)
			continue
		}
	}

	return nil
}
