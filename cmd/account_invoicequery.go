package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/account"
)

func accountInvoiceQueryRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoicequery",
		Short: "sub-commands relating to invoice queries",
	}

	// Child commands
	cmd.AddCommand(accountInvoiceQueryListCmd())
	cmd.AddCommand(accountInvoiceQueryShowCmd())

	return cmd
}

func accountInvoiceQueryListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists invoice queries",
		Long:    "This command lists invoice queries",
		Example: "ukfast account invoicequery list",
		Run: func(cmd *cobra.Command, args []string) {
			accountInvoiceQueryList(getClient().AccountService(), cmd, args)
		},
	}
}

func accountInvoiceQueryList(service account.AccountService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	queries, err := service.GetInvoiceQueries(params)
	if err != nil {
		output.Fatalf("Error retrieving invoice queries: %s", err)
		return
	}

	outputAccountInvoiceQueries(queries)
}

func accountInvoiceQueryShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <invoicequery: id>...",
		Short:   "Shows an invoice query",
		Long:    "This command shows one or more invoice queries",
		Example: "ukfast account invoicequery show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing invoice query")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			accountInvoiceQueryShow(getClient().AccountService(), cmd, args)
		},
	}
}

func accountInvoiceQueryShow(service account.AccountService, cmd *cobra.Command, args []string) {
	var queries []account.InvoiceQuery
	for _, arg := range args {
		queryID, err := strconv.Atoi(arg)
		if err != nil {
			OutputWithErrorLevelf("Invalid invoice query ID [%s]", arg)
			continue
		}

		query, err := service.GetInvoiceQuery(queryID)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving invoice query [%s]: %s", arg, err)
			continue
		}

		queries = append(queries, query)
	}

	outputAccountInvoiceQueries(queries)
}
