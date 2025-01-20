package loadbalancer

import (
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/spf13/cobra"
)

func loadbalancerCertificateRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certificate",
		Short: "sub-commands relating to certificates",
	}

	// Child commands
	cmd.AddCommand(loadbalancerCertificateListCmd(f))

	return cmd
}

func loadbalancerCertificateListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists certificates",
		Long:    "This command lists certificates",
		Example: "ans loadbalancer certificate list",
		RunE:    loadbalancerCobraRunEFunc(f, loadbalancerCertificateList),
	}
}

func loadbalancerCertificateList(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	certificates, err := service.GetCertificates(params)
	if err != nil {
		return fmt.Errorf("Error retrieving certificates: %s", err)
	}

	return output.CommandOutput(cmd, CertificateCollection(certificates))
}
