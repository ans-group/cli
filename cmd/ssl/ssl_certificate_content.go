package ssl

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ssl"
)

func sslCertificateContentRootCmd(f factory.ClientFactory) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "content",
		Short: "sub-commands relating to certificate contents",
	}

	// Child commands
	cmd.AddCommand(sslCertificateContentShowCmd(f))

	return cmd
}

func sslCertificateContentShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <certificate: id>...",
		Short:   "Shows a certificate content",
		Long:    "This command shows one or more certificate contents",
		Example: "ukfast ssl certificate content show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing certificate")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return sslCertificateContentShow(f.NewClient().SSLService(), cmd, args)
		},
	}
}

func sslCertificateContentShow(service ssl.SSLService, cmd *cobra.Command, args []string) error {
	var certificateContents []ssl.CertificateContent
	for _, arg := range args {
		certificateID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid certificate ID [%s]", arg)
			continue
		}

		certificateContent, err := service.GetCertificateContent(certificateID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving certificate content [%s]: %s", arg, err)
			continue
		}

		certificateContents = append(certificateContents, certificateContent)
	}

	return output.CommandOutput(cmd, OutputSSLCertificatesContentsProvider(certificateContents))
}
