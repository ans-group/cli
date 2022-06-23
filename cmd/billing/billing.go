package billing

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/spf13/cobra"
)

func BillingRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "billing",
		Short: "Commands relating to Billing service",
	}

	// Child root commands
	cmd.AddCommand(billingCardRootCmd(f))
	cmd.AddCommand(billingCloudCostRootCmd(f))
	cmd.AddCommand(billingDirectDebitRootCmd(f))
	cmd.AddCommand(billingInvoiceRootCmd(f))
	cmd.AddCommand(billingInvoiceQueryRootCmd(f))
	cmd.AddCommand(billingPaymentRootCmd(f))
	cmd.AddCommand(billingRecurringCostRootCmd(f))

	return cmd
}
