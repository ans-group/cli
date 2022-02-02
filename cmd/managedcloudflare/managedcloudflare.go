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

	return cmd
}

func OutputManagedCloudflareAccountsProvider(projects []managedcloudflare.Account) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(projects).WithDefaultFields([]string{"id", "name", "status", "cloudflare_account_id"})
}
