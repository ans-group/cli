package cloudflare

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/sdk-go/pkg/service/cloudflare"
	"github.com/spf13/cobra"
)

func CloudflareRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cloudflare",
		Short: "Commands relating to the Cloudflare service",
	}

	// Child root commands
	cmd.AddCommand(cloudflareAccountRootCmd(f))
	cmd.AddCommand(cloudflareSpendPlanRootCmd(f))
	cmd.AddCommand(cloudflareSubscriptionRootCmd(f))
	cmd.AddCommand(cloudflareTotalSpendRootCmd(f))
	cmd.AddCommand(cloudflareZoneRootCmd(f))

	return cmd
}

type cloudflareServiceCobraRunEFunc func(service cloudflare.CloudflareService, cmd *cobra.Command, args []string) error

func cloudflareCobraRunEFunc(f factory.ClientFactory, rf cloudflareServiceCobraRunEFunc) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		c, err := f.NewClient()
		if err != nil {
			return err
		}

		return rf(c.CloudflareService(), cmd, args)
	}
}
