package ddosx

import (
	"errors"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	"github.com/spf13/cobra"
)

func ddosxSSLPrivateKeyRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "privatekey",
		Short: "sub-commands relating to private keys",
	}

	// Child commands
	cmd.AddCommand(ddosxSSLPrivateKeyShowCmd(f))

	return cmd
}

func ddosxSSLPrivateKeyShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <ssl: id>...",
		Short:   "Shows a ssl's private key",
		Long:    "This command shows one or more ssl's private key",
		Example: "ans ddosx ssl privatekey show 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ssl")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxSSLPrivateKeyShow(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxSSLPrivateKeyShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	var sslPrivateKeys []ddosx.SSLPrivateKey
	for _, arg := range args {
		sslPrivateKey, err := service.GetSSLPrivateKey(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving ssl [%s]: %s", arg, err)
			continue
		}

		sslPrivateKeys = append(sslPrivateKeys, sslPrivateKey)
	}

	return output.CommandOutput(cmd, OutputDDoSXSSLPrivateKeysProvider(sslPrivateKeys))
}
