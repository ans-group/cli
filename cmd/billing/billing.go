package billing

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
)

func BillingRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "billing",
		Short: "Commands relating to Billing service",
	}

	// Child root commands
	cmd.AddCommand(billingCardRootCmd(f))
	cmd.AddCommand(billingCloudCostRootCmd(f))
	cmd.AddCommand(billingInvoiceRootCmd(f))
	cmd.AddCommand(billingInvoiceQueryRootCmd(f))

	return cmd
}
