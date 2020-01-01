package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/account"
)

func accountInvoiceRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoice",
		Short: "sub-commands relating to invoices",
	}

	// Child commands
	cmd.AddCommand(accountInvoiceListCmd())
	cmd.AddCommand(accountInvoiceShowCmd())

	return cmd
}

func accountInvoiceListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists invoices",
		Long:    "This command lists invoices",
		Example: "ukfast account invoice list",
		Run: func(cmd *cobra.Command, args []string) {
			accountInvoiceList(getClient().AccountService(), cmd, args)
		},
	}
}

func accountInvoiceList(service account.AccountService, cmd *cobra.Command, args []string) {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	invoices, err := service.GetInvoices(params)
	if err != nil {
		output.Fatalf("Error retrieving invoices: %s", err)
		return
	}

	outputAccountInvoices(invoices)
}

func accountInvoiceShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <invoice: id>...",
		Short:   "Shows a invoice",
		Long:    "This command shows one or more invoices",
		Example: "ukfast account invoice show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing invoice")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			accountInvoiceShow(getClient().AccountService(), cmd, args)
		},
	}
}

func accountInvoiceShow(service account.AccountService, cmd *cobra.Command, args []string) {
	var invoices []account.Invoice
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

	outputAccountInvoices(invoices)
}
