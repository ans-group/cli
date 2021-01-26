package account

import (
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/account"
)

func OutputAccountContactsProvider(contacts []account.Contact) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(contacts).
		WithDefaultFields([]string{"id", "type", "first_name", "last_name"})
}

func OutputAccountDetailsProvider(details []account.Details) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(details).
		WithDefaultFields([]string{"company_registration_number", "vat_identification_number", "primary_contact_id"})
}

func OutputAccountCreditsProvider(credits []account.Credit) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(credits).
		WithDefaultFields([]string{"type", "total", "remaining"})
}

func OutputAccountInvoicesProvider(invoices []account.Invoice) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(invoices).
		WithDefaultFields([]string{"id", "date", "paid", "net", "vat", "gross"}).
		WithMonetaryFields([]string{"net", "vat", "gross"})
}

func OutputAccountInvoiceQueriesProvider(queries []account.InvoiceQuery) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(queries).
		WithDefaultFields([]string{"id", "contact_id", "amount", "what_was_expected", "what_was_received", "proposed_solution", "invoice_ids"})
}
