package account

import (
	"fmt"
	"strconv"

	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/account"
)

func OutputAccountContactsProvider(contacts []account.Contact) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(contacts),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, contact := range contacts {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(strconv.Itoa(contact.ID), true))
				fields.Set("type", output.NewFieldValue(contact.Type.String(), true))
				fields.Set("first_name", output.NewFieldValue(contact.FirstName, true))
				fields.Set("last_name", output.NewFieldValue(contact.LastName, true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputAccountDetailsProvider(details []account.Details) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(details),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, detail := range details {
				fields := output.NewOrderedFields()
				fields.Set("company_registration_number", output.NewFieldValue(detail.CompanyRegistrationNumber, true))
				fields.Set("vat_identification_number", output.NewFieldValue(detail.VATIdentificationNumber, true))
				fields.Set("primary_contact_id", output.NewFieldValue(strconv.Itoa(detail.PrimaryContactID), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputAccountCreditsProvider(credits []account.Credit) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(credits),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, credit := range credits {
				fields := output.NewOrderedFields()
				fields.Set("type", output.NewFieldValue(credit.Type, true))
				fields.Set("total", output.NewFieldValue(strconv.Itoa(credit.Total), true))
				fields.Set("remaining", output.NewFieldValue(strconv.Itoa(credit.Remaining), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputAccountInvoicesProvider(invoices []account.Invoice) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(invoices),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, invoice := range invoices {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(strconv.Itoa(invoice.ID), true))
				fields.Set("date", output.NewFieldValue(invoice.Date.String(), true))
				fields.Set("paid", output.NewFieldValue(strconv.FormatBool(invoice.Paid), true))
				fields.Set("net", output.NewFieldValue(fmt.Sprintf("%.2f", invoice.Net), true))
				fields.Set("vat", output.NewFieldValue(fmt.Sprintf("%.2f", invoice.VAT), true))
				fields.Set("gross", output.NewFieldValue(fmt.Sprintf("%.2f", invoice.Gross), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputAccountInvoiceQueriesProvider(queries []account.InvoiceQuery) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(queries),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, query := range queries {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(strconv.Itoa(query.ID), true))
				fields.Set("contact_id", output.NewFieldValue(strconv.Itoa(query.ContactID), false))
				fields.Set("amount", output.NewFieldValue(fmt.Sprintf("%f", query.Amount), true))
				fields.Set("what_was_expected", output.NewFieldValue(query.WhatWasExpected, false))
				fields.Set("what_was_received", output.NewFieldValue(query.WhatWasReceived, false))
				fields.Set("proposed_solution", output.NewFieldValue(query.ProposedSolution, false))
				fields.Set("invoice_ids", output.NewFieldValue(helper.JoinInt(query.InvoiceIDs, ", "), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}
