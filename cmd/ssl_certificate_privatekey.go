package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ssl"
)

func sslCertificatePrivateKeyRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "privatekey",
		Short: "sub-commands relating to certificate private keys",
	}

	// Child commands
	cmd.AddCommand(sslCertificatePrivateKeyShowCmd())

	return cmd
}

func sslCertificatePrivateKeyShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <certificate: id>...",
		Short:   "Shows a certificate private key",
		Long:    "This command shows one or more certificate private keys",
		Example: "ukfast ssl certificate privatekey show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing certificate")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			sslCertificatePrivateKeyShow(getClient().SSLService(), cmd, args)
		},
	}
}

func sslCertificatePrivateKeyShow(service ssl.SSLService, cmd *cobra.Command, args []string) {
	var certificatePrivateKeys []ssl.CertificatePrivateKey
	for _, arg := range args {
		certificateID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid certificate ID [%s]", arg)
			continue
		}

		certificatePrivateKey, err := service.GetCertificatePrivateKey(certificateID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving certificate private key [%s]: %s", arg, err)
			continue
		}

		certificatePrivateKeys = append(certificatePrivateKeys, certificatePrivateKey)
	}

	outputSSLCertificatesPrivateKeys(certificatePrivateKeys)
}
