package loadbalancer

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func loadbalancerListenerCertificateRootCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certificate",
		Short: "sub-commands relating to certificates",
	}

	// Child commands
	cmd.AddCommand(loadbalancerListenerCertificateListCmd(f))
	cmd.AddCommand(loadbalancerListenerCertificateShowCmd(f))
	cmd.AddCommand(loadbalancerListenerCertificateCreateCmd(f, fs))
	cmd.AddCommand(loadbalancerListenerCertificateUpdateCmd(f, fs))
	cmd.AddCommand(loadbalancerListenerCertificateDeleteCmd(f))

	return cmd
}

func loadbalancerListenerCertificateListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <listener: id>",
		Short:   "Lists certificates",
		Long:    "This command lists certificates",
		Example: "ukfast loadbalancer listener certificate list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing listener")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerListenerCertificateList),
	}
}

func loadbalancerListenerCertificateList(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	listenerID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid listener ID")
	}

	certificates, err := service.GetListenerCertificates(listenerID, params)
	if err != nil {
		return fmt.Errorf("Error retrieving certificates: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerCertificatesProvider(certificates))
}

func loadbalancerListenerCertificateShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <listener: id> <certificate: id>...",
		Short:   "Shows a certificate",
		Long:    "This command shows one or more certificates",
		Example: "ukfast loadbalancer listener certificate show 123 345",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing listener")
			}
			if len(args) < 2 {
				return errors.New("Missing certificate")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerListenerCertificateShow),
	}
}

func loadbalancerListenerCertificateShow(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	listenerID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid listener ID")
	}

	var certificates []loadbalancer.Certificate
	for _, arg := range args[1:] {

		certificateID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid certificate ID [%s]", arg)
			continue
		}

		certificate, err := service.GetListenerCertificate(listenerID, certificateID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving certificate [%d]: %s", certificateID, err)
			continue
		}

		certificates = append(certificates, certificate)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerCertificatesProvider(certificates))
}

func loadbalancerListenerCertificateCreateCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <listener: id>",
		Short:   "Creates a certificate",
		Long:    "This command creates a certificate",
		Example: "ukfast loadbalancer listener certificate create 123 --key-file /tmp/cert.key --certificate-file /tmp/cert.crt --ca-bundle-file /tmp/ca.crt",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing listener")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadbalancerListenerCertificateCreate(c.LoadBalancerService(), cmd, fs, args)
		},
	}

	cmd.Flags().String("name", "", "Name for certificate")
	cmd.Flags().String("key", "", "Key contents for certificate")
	cmd.Flags().String("key-file", "", "Path to file containing key contents for certificate")
	cmd.Flags().String("certificate", "", "Certificate contents for certificate")
	cmd.Flags().String("certificate-file", "", "Path to file containing certificate contents for certificate")
	cmd.Flags().String("ca-bundle", "", "CA bundle contents for certificate")
	cmd.Flags().String("ca-bundle-file", "", "Path to file containing CA bundle contents for certificate")

	return cmd
}

func loadbalancerListenerCertificateCreate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, fs afero.Fs, args []string) error {
	listenerID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid listener ID")
	}

	createRequest := loadbalancer.CreateCertificateRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.Key, err = helper.GetContentsFromLiteralOrFilePathFlag(cmd, fs, "key", "key-file")
	if err != nil {
		return err
	}
	createRequest.Certificate, err = helper.GetContentsFromLiteralOrFilePathFlag(cmd, fs, "certificate", "certificate-file")
	if err != nil {
		return err
	}
	createRequest.CABundle, err = helper.GetContentsFromLiteralOrFilePathFlag(cmd, fs, "ca-bundle", "ca-bundle-file")
	if err != nil {
		return err
	}

	certificateID, err := service.CreateListenerCertificate(listenerID, createRequest)
	if err != nil {
		return fmt.Errorf("Error creating certificate: %s", err)
	}

	certificate, err := service.GetListenerCertificate(listenerID, certificateID)
	if err != nil {
		return fmt.Errorf("Error retrieving new certificate: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerCertificatesProvider([]loadbalancer.Certificate{certificate}))
}

func loadbalancerListenerCertificateUpdateCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <listener: id> <certificate: id>...",
		Short:   "Updates a certificate",
		Long:    "This command updates one or more certificates",
		Example: "ukfast loadbalancer listener certificate update 123 456 --name mycertificate",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing listener")
			}
			if len(args) < 2 {
				return errors.New("Missing certificate")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadbalancerListenerCertificateUpdate(c.LoadBalancerService(), cmd, fs, args)
		},
	}

	cmd.Flags().String("name", "", "Name for certificate")
	cmd.Flags().String("key", "", "Key contents for certificate")
	cmd.Flags().String("key-file", "", "Path to file containing key contents for certificate")
	cmd.Flags().String("certificate", "", "Certificate contents for certificate")
	cmd.Flags().String("certificate-file", "", "Path to file containing certificate contents for certificate")
	cmd.Flags().String("ca-bundle", "", "CA bundle contents for certificate")
	cmd.Flags().String("ca-bundle-file", "", "Path to file containing CA bundle contents for certificate")

	return cmd
}

func loadbalancerListenerCertificateUpdate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, fs afero.Fs, args []string) error {
	listenerID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid listener ID")
	}

	patchRequest := loadbalancer.PatchCertificateRequest{}
	patchRequest.Name, _ = cmd.Flags().GetString("name")
	patchRequest.Key, err = helper.GetContentsFromLiteralOrFilePathFlag(cmd, fs, "key", "key-file")
	if err != nil {
		return err
	}
	patchRequest.Certificate, err = helper.GetContentsFromLiteralOrFilePathFlag(cmd, fs, "certificate", "certificate-file")
	if err != nil {
		return err
	}
	patchRequest.CABundle, err = helper.GetContentsFromLiteralOrFilePathFlag(cmd, fs, "ca-bundle", "ca-bundle-file")
	if err != nil {
		return err
	}

	var certificates []loadbalancer.Certificate
	for _, arg := range args[1:] {
		certificateID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid certificate ID [%s]", arg)
			continue
		}

		err = service.PatchListenerCertificate(listenerID, certificateID, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating certificate [%d]: %s", certificateID, err)
			continue
		}

		certificate, err := service.GetListenerCertificate(listenerID, certificateID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated certificate [%d]: %s", certificateID, err)
			continue
		}

		certificates = append(certificates, certificate)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerCertificatesProvider(certificates))
}

func loadbalancerListenerCertificateDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <listener: id> <certificate: id>...",
		Short:   "Removes a certificate",
		Long:    "This command removes one or more certificates",
		Example: "ukfast loadbalancer listener certificate delete 123 456",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing listener")
			}
			if len(args) < 2 {
				return errors.New("Missing certificate")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerListenerCertificateDelete),
	}
}

func loadbalancerListenerCertificateDelete(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	listenerID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid listener ID")
	}

	for _, arg := range args[1:] {
		certificateID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid certificate ID [%s]", arg)
			continue
		}

		err = service.DeleteListenerCertificate(listenerID, certificateID)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing certificate [%d]: %s", certificateID, err)
			continue
		}
	}

	return nil
}
