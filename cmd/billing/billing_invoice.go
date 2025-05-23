package billing

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/billing"
	"github.com/spf13/cobra"
)

func billingInvoiceRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoice",
		Short: "sub-commands relating to invoices",
	}

	// Child commands
	cmd.AddCommand(billingInvoiceListCmd(f))
	cmd.AddCommand(billingInvoiceShowCmd(f))

	return cmd
}

func billingInvoiceListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists invoices",
		Long:    "This command lists invoices",
		Example: "ans billing invoice list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return billingInvoiceList(c.BillingService(), cmd, args)
		},
	}
}

func billingInvoiceList(service billing.BillingService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	invoices, err := service.GetInvoices(params)
	if err != nil {
		return fmt.Errorf("Error retrieving invoices: %s", err)
	}

	return output.CommandOutput(cmd, InvoiceCollection(invoices))
}

func billingInvoiceShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <invoice: id>...",
		Short:   "Shows a invoice",
		Long:    "This command shows one or more invoices",
		Example: "ans billing invoice show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing invoice")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return billingInvoiceShow(c.BillingService(), cmd, args)
		},
	}
}

func billingInvoiceShow(service billing.BillingService, cmd *cobra.Command, args []string) error {
	var invoices []billing.Invoice
	for _, arg := range args {
		invoiceID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid invoice ID [%s]", arg)
			continue
		}

		invoice, err := service.GetInvoice(invoiceID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving invoice [%s]: %s", arg, err)
			continue
		}

		invoices = append(invoices, invoice)
	}

	return output.CommandOutput(cmd, InvoiceCollection(invoices))
}
