package managedcloudflare

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/output"
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

	return cmd
}

func OutputManagedCloudflareAccountsProvider(accounts []managedcloudflare.Account) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(accounts).WithDefaultFields([]string{"id", "name", "status", "cloudflare_account_id"})
}

func OutputManagedCloudflareSpendPlansProvider(plans []managedcloudflare.SpendPlan) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(plans).WithDefaultFields([]string{"id", "amount", "started_at", "ended_at"})
}
