package billing

import (
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/billing"
	"github.com/spf13/cobra"
)

func billingDirectDebitRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "directdebit",
		Short: "sub-commands relating to direct debits",
	}

	// Child commands
	cmd.AddCommand(billingDirectDebitShowCmd(f))

	return cmd
}

func billingDirectDebitShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show",
		Short:   "Shows direct debit details",
		Long:    "This command shows direct debit details",
		Example: "ans billing directdebit show",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return billingDirectDebitShow(c.BillingService(), cmd, args)
		},
	}
}

func billingDirectDebitShow(service billing.BillingService, cmd *cobra.Command, args []string) error {
	details, err := service.GetDirectDebit()
	if err != nil {
		return fmt.Errorf("Error retrieving direct debit details: %s", err)
	}

	return output.CommandOutput(cmd, OutputBillingDirectDebitsProvider([]billing.DirectDebit{details}))
}
