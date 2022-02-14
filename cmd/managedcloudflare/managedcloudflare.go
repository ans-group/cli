package managedcloudflare

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/sdk-go/pkg/service/managedcloudflare"
)

func ManagedCloudflareRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "managedcloudflare",
		Short: "Commands relating to the Managed Cloudflare service",
	}

	// Child root commands
	cmd.AddCommand(managedcloudflareAccountRootCmd(f))
	cmd.AddCommand(managedcloudflareSpendPlanRootCmd(f))
	cmd.AddCommand(managedcloudflareSubscriptionRootCmd(f))
	cmd.AddCommand(managedcloudflareTotalSpendRootCmd(f))
	cmd.AddCommand(managedcloudflareZoneRootCmd(f))

	return cmd
}

type managedcloudflareServiceCobraRunEFunc func(service managedcloudflare.ManagedCloudflareService, cmd *cobra.Command, args []string) error

func managedcloudflareCobraRunEFunc(f factory.ClientFactory, rf managedcloudflareServiceCobraRunEFunc) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		c, err := f.NewClient()
		if err != nil {
			return err
		}

		return rf(c.ManagedCloudflareService(), cmd, args)
	}
}
