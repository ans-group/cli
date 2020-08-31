package billing

import (
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/billing"
)

func OutputBillingCardsProvider(cards []billing.Card) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(cards).WithDefaultFields([]string{"id", "name", "card_number"})
}

func OutputBillingInvoicesProvider(invoices []billing.Invoice) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(invoices).WithDefaultFields([]string{"id", "date", "paid", "net", "vat", "gross"})
}

func OutputBillingInvoiceQueriesProvider(queries []billing.InvoiceQuery) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(queries).WithDefaultFields([]string{"id", "amount", "invoice_ids"})
}

func OutputBillingCloudCostsProvider(costs []billing.CloudCost) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(costs).WithDefaultFields([]string{"id", "server_id", "resource_type", "resource_quantity", "resource_price", "resource_period"})
}
