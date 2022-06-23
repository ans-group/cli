package ssl

import (
	"errors"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ssl"
	"github.com/spf13/cobra"
)

func sslCertificatePrivateKeyRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "privatekey",
		Short: "sub-commands relating to certificate private keys",
	}

	// Child commands
	cmd.AddCommand(sslCertificatePrivateKeyShowCmd(f))

	return cmd
}

func sslCertificatePrivateKeyShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <certificate: id>...",
		Short:   "Shows a certificate private key",
		Long:    "This command shows one or more certificate private keys",
		Example: "ans ssl certificate privatekey show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing certificate")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return sslCertificatePrivateKeyShow(c.SSLService(), cmd, args)
		},
	}
}

func sslCertificatePrivateKeyShow(service ssl.SSLService, cmd *cobra.Command, args []string) error {
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

	return output.CommandOutput(cmd, OutputSSLCertificatesPrivateKeysProvider(certificatePrivateKeys))
}
