package ssl

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	flaghelper "github.com/ukfast/cli/internal/pkg/helper/flag"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ssl"
)

func sslCertificateRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certificate",
		Short: "sub-commands relating to certificates",
	}

	// Child commands
	cmd.AddCommand(sslCertificateListCmd(f))
	cmd.AddCommand(sslCertificateShowCmd(f))

	// Child root commands
	cmd.AddCommand(sslCertificateContentRootCmd(f))
	cmd.AddCommand(sslCertificatePrivateKeyRootCmd(f))

	return cmd
}

func sslCertificateListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists certificates",
		Long:    "This command lists certificates",
		Example: "ukfast ssl certificate list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return sslCertificateList(c.SSLService(), cmd, args)
		},
	}
}

func sslCertificateList(service ssl.SSLService, cmd *cobra.Command, args []string) error {
	params, err := flaghelper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	certificates, err := service.GetCertificates(params)
	if err != nil {
		return fmt.Errorf("Error retrieving certificates: %s", err)
	}

	return output.CommandOutput(cmd, OutputSSLCertificatesProvider(certificates))
}

func sslCertificateShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <certificate: id>...",
		Short:   "Shows a certificate",
		Long:    "This command shows one or more certificates",
		Example: "ukfast ssl certificate show 123",
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

			return sslCertificateShow(c.SSLService(), cmd, args)
		},
	}
}

func sslCertificateShow(service ssl.SSLService, cmd *cobra.Command, args []string) error {
	var certificates []ssl.Certificate
	for _, arg := range args {
		certificateID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid certificate ID [%s]", arg)
			continue
		}

		certificate, err := service.GetCertificate(certificateID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving certificate [%s]: %s", arg, err)
			continue
		}

		certificates = append(certificates, certificate)
	}

	return output.CommandOutput(cmd, OutputSSLCertificatesProvider(certificates))
}
