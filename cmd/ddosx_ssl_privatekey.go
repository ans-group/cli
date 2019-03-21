package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxSSLPrivateKeyRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "privatekey",
		Short: "sub-commands relating to private keys",
	}

	// Child commands
	cmd.AddCommand(ddosxSSLPrivateKeyShowCmd())

	return cmd
}

func ddosxSSLPrivateKeyShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <ssl: id>...",
		Short:   "Shows a ssl's private key",
		Long:    "This command shows one or more ssl's private key",
		Example: "ukfast ddosx ssl privatekey show 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing ssl")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxSSLPrivateKeyShow(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxSSLPrivateKeyShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	var sslPrivateKeys []ddosx.SSLPrivateKey
	for _, arg := range args {
		sslPrivateKey, err := service.GetSSLPrivateKey(arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving ssl [%s]: %s", arg, err)
			continue
		}

		sslPrivateKeys = append(sslPrivateKeys, sslPrivateKey)
	}

	outputDDoSXSSLPrivateKeys(sslPrivateKeys)
}
