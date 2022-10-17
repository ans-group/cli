package loadbalancer

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func LoadBalancerRootCmd(f factory.ClientFactory, fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "loadbalancer",
		Short: "Commands relating to load balancer service",
	}

	// Child root commands
	cmd.AddCommand(loadbalancerAccessIPRootCmd(f))
	cmd.AddCommand(loadbalancerACLRootCmd(f))
	cmd.AddCommand(loadbalancerBindRootCmd(f))
	cmd.AddCommand(loadbalancerCertificateRootCmd(f))
	cmd.AddCommand(loadbalancerClusterRootCmd(f))
	cmd.AddCommand(loadbalancerDeploymentRootCmd(f))
	cmd.AddCommand(loadbalancerListenerRootCmd(f, fs))
	cmd.AddCommand(loadbalancerTargetGroupRootCmd(f))

	return cmd
}

type loadbalancerServiceCobraRunEFunc func(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error

func loadbalancerCobraRunEFunc(f factory.ClientFactory, rf loadbalancerServiceCobraRunEFunc) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		c, err := f.NewClient()
		if err != nil {
			return err
		}

		return rf(c.LoadBalancerService(), cmd, args)
	}
}
