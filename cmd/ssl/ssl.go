package ssl

import (
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ssl"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func SSLRootCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssl",
		Short: "Commands relating to SSL service",
	}

	// Child commands
	cmd.AddCommand(sslValidateCmd(f, fs))

	// Child root commands
	cmd.AddCommand(sslCertificateRootCmd(f))
	cmd.AddCommand(sslRecommendationsRootCmd(f))
	cmd.AddCommand(sslReportRootCmd(f))

	return cmd
}

func sslValidateCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "validate",
		Short:   "Validates a certificate",
		Long:    "This command validates an SSL certificate",
		Example: "ans ssl validate --certificate-file /tmp/cert.crt --key-file /tmp/cert.key --ca-bundle-file /tmp/ca.crt",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return sslValidate(c.SSLService(), fs, cmd, args)
		},
	}
	cmd.Flags().String("key", "", "Key contents for SSL to validate")
	cmd.Flags().String("key-file", "", "Path to file containing key contents for SSL to validate")
	cmd.Flags().String("certificate", "", "Certificate contents for SSL to validate")
	cmd.Flags().String("certificate-file", "", "Path to file containing certificate contents for SSL to validate")
	cmd.Flags().String("ca-bundle", "", "CA bundle contents for SSL to validate")
	cmd.Flags().String("ca-bundle-file", "", "Path to file containing CA bundle contents for SSL to validate")

	return cmd
}

func sslValidate(service ssl.SSLService, fs afero.Fs, cmd *cobra.Command, args []string) error {
	validateRequest := ssl.ValidateRequest{}

	var err error
	validateRequest.Key, err = helper.GetContentsFromLiteralOrFilePathFlag(cmd, fs, "key", "key-file")
	if err != nil {
		return err
	}

	validateRequest.Certificate, err = helper.GetContentsFromLiteralOrFilePathFlag(cmd, fs, "certificate", "certificate-file")
	if err != nil {
		return err
	}

	validateRequest.CABundle, err = helper.GetContentsFromLiteralOrFilePathFlag(cmd, fs, "ca-bundle", "ca-bundle-file")
	if err != nil {
		return err
	}

	validation, err := service.ValidateCertificate(validateRequest)
	if err != nil {
		return fmt.Errorf("Error validating certificate: %s", err)
	}

	return output.CommandOutput(cmd, CertificateValidationCollection([]ssl.CertificateValidation{validation}))
}
