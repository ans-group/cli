package managedcloudflare

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/managedcloudflare"
)

func managedcloudflareTotalSpendRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "totalspend",
		Short: "sub-commands relating to total spend",
	}

	// Child commands
	cmd.AddCommand(managedcloudflareTotalSpendShowCmd(f))

	return cmd
}

func managedcloudflareTotalSpendShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show",
		Short:   "Shows total spend",
		Long:    "This command shows total spend",
		Example: "ukfast managedcloudflare totalspend show",
		RunE:    managedcloudflareCobraRunEFunc(f, managedcloudflareTotalSpendShow),
	}
}

func managedcloudflareTotalSpendShow(service managedcloudflare.ManagedCloudflareService, cmd *cobra.Command, args []string) error {
	spend, err := service.GetTotalSpendMonthToDate()
	if err != nil {
		return fmt.Errorf("Error retrieving total spend: %s", err)
	}

	return output.CommandOutput(cmd, OutputManagedCloudflareTotalSpendsProvider([]managedcloudflare.TotalSpend{spend}))
}
