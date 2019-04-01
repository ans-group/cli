package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ssl"
)

func sslCertificateRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certificate",
		Short: "sub-commands relating to certificates",
	}

	// Child commands
	cmd.AddCommand(sslCertificateListCmd())
	cmd.AddCommand(sslCertificateShowCmd())

	// Child root commands

	return cmd
}

func sslCertificateListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists certificates",
		Long:    "This command lists certificates",
		Example: "ukfast ssl certificate list",
		Run: func(cmd *cobra.Command, args []string) {
			sslCertificateList(getClient().SSLService(), cmd, args)
		},
	}
}

func sslCertificateList(service ssl.SSLService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	certificates, err := service.GetCertificates(params)
	if err != nil {
		output.Fatalf("Error retrieving certificates: %s", err)
		return
	}

	outputSSLCertificates(certificates)
}

func sslCertificateShowCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			sslCertificateShow(getClient().SSLService(), cmd, args)
		},
	}
}

func sslCertificateShow(service ssl.SSLService, cmd *cobra.Command, args []string) {
	var certificates []ssl.Certificate
	for _, arg := range args {
		certificateID, err := strconv.Atoi(arg)
		if err != nil {
			OutputWithErrorLevelf("Invalid certificate ID [%s]", arg)
			continue
		}

		certificate, err := service.GetCertificate(certificateID)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving certificate [%s]: %s", arg, err)
			continue
		}

		certificates = append(certificates, certificate)
	}

	outputSSLCertificates(certificates)
}
