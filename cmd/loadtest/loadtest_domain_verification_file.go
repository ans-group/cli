package loadtest

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func loadtestDomainVerificationFileRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "file",
		Short: "sub-commands relating to File domain verification",
	}

	// Child commands
	cmd.AddCommand(loadtestDomainVerificationFileVerifyCmd(f))

	return cmd
}

func loadtestDomainVerificationFileVerifyCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "verify <domain: id>...",
		Short:   "Verifies a domain via File verification method",
		Long:    "This command verifies one or more domains via the File verification method",
		Example: "ukfast loadtest domain verification file verify 00000000-0000-0000-0000-000000000000",
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

			return loadtestDomainVerificationFileVerify(c.LTaaSService(), cmd, args)
		},
	}
}

func loadtestDomainVerificationFileVerify(service ltaas.LTaaSService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		err := service.VerifyDomainFile(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error verifying domain [%s] via File verification method: %s", arg, err)
			continue
		}
	}

	return nil
}
