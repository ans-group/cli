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

func billingInvoiceQueryRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoicequery",
		Short: "sub-commands relating to invoice queries",
	}

	// Child commands
	cmd.AddCommand(billingInvoiceQueryListCmd(f))
	cmd.AddCommand(billingInvoiceQueryShowCmd(f))
	cmd.AddCommand(billingInvoiceQueryCreateCmd(f))

	return cmd
}

func billingInvoiceQueryListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists invoice queries",
		Long:    "This command lists invoice queries",
		Example: "ans billing invoicequery list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return billingInvoiceQueryList(c.BillingService(), cmd, args)
		},
	}
}

func billingInvoiceQueryList(service billing.BillingService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	queries, err := service.GetInvoiceQueries(params)
	if err != nil {
		return fmt.Errorf("Error retrieving invoice queries: %s", err)
	}

	return output.CommandOutput(cmd, OutputBillingInvoiceQueriesProvider(queries))
}

func billingInvoiceQueryShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <invoicequery: id>...",
		Short:   "Shows an invoice query",
		Long:    "This command shows one or more invoice queries",
		Example: "ans billing invoicequery show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing invoice query")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return billingInvoiceQueryShow(c.BillingService(), cmd, args)
		},
	}
}

func billingInvoiceQueryShow(service billing.BillingService, cmd *cobra.Command, args []string) error {
	var queries []billing.InvoiceQuery
	for _, arg := range args {
		queryID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid invoice query ID [%s]", arg)
			continue
		}

		query, err := service.GetInvoiceQuery(queryID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving invoice query [%s]: %s", arg, err)
			continue
		}

		queries = append(queries, query)
	}

	return output.CommandOutput(cmd, OutputBillingInvoiceQueriesProvider(queries))
}

func billingInvoiceQueryCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a invoice query",
		Long:    "This command creates a new invoice query",
		Example: "ans billing invoice query create",
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

			return billingInvoiceQueryCreate(c.BillingService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().Int("contact-id", 0, "The identifier of the contact that raised this query")
	cmd.MarkFlagRequired("contact-id")
	cmd.Flags().String("contact-method", "", "What is the preferred method to provide feedback on this query")
	cmd.MarkFlagRequired("contact-method")
	cmd.Flags().Float32("amount", 0, "The amount that is being queried")
	cmd.Flags().String("what-was-expected", "", "Text explaining what was expected")
	cmd.MarkFlagRequired("what-was-expected")
	cmd.Flags().String("what-was-received", "", "Text explaining what was actually received")
	cmd.MarkFlagRequired("what-was-received")
	cmd.Flags().String("proposed-solution", "", "What is the proposed solution")
	cmd.MarkFlagRequired("proposed-solution")
	cmd.Flags().IntSlice("invoice-id", []int{}, "Invoice ID for query")
	cmd.MarkFlagRequired("invoice-id")

	return cmd
}

func billingInvoiceQueryCreate(service billing.BillingService, cmd *cobra.Command, args []string) error {
	createRequest := billing.CreateInvoiceQueryRequest{}
	createRequest.ContactID, _ = cmd.Flags().GetInt("contact-id")
	createRequest.ContactMethod, _ = cmd.Flags().GetString("contact-method")
	createRequest.Amount, _ = cmd.Flags().GetFloat32("amount")
	createRequest.WhatWasExpected, _ = cmd.Flags().GetString("what-was-expected")
	createRequest.WhatWasReceived, _ = cmd.Flags().GetString("what-was-received")
	createRequest.ProposedSolution, _ = cmd.Flags().GetString("proposed-solution")
	createRequest.InvoiceIDs, _ = cmd.Flags().GetIntSlice("invoice-id")

	id, err := service.CreateInvoiceQuery(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating invoice query: %s", err)
	}

	query, err := service.GetInvoiceQuery(id)
	if err != nil {
		return fmt.Errorf("Error retrieving new invoice query [%d]: %s", id, err)
	}

	return output.CommandOutput(cmd, OutputBillingInvoiceQueriesProvider([]billing.InvoiceQuery{query}))
}
