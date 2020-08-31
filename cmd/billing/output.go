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
	return output.NewSerializedOutputHandlerProvider(queries).
		WithDefaultFields([]string{"id", "amount", "invoice_ids"}).
		WithMonetaryFields([]string{"amount"})
}

func OutputBillingCloudCostsProvider(costs []billing.CloudCost) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(costs).
		WithDefaultFields([]string{"id", "server_id", "resource_type", "resource_quantity", "resource_price", "resource_period"}).
		WithMonetaryFields([]string{"resource_price", "cost_since_last_invoice", "cost_for_period_estimate"})
}

func OutputBillingDirectDebitsProvider(details []billing.DirectDebit) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(details).WithDefaultFields([]string{"name", "status", "is_actived"})
}

func OutputBillingPaymentsProvider(payments []billing.Payment) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(payments).
		WithDefaultFields([]string{"id", "date", "cost", "vat", "gross", "discount"}).
		WithMonetaryFields([]string{"cost", "vat", "gross", "discount"})
}

func OutputBillingRecurringCostsProvider(costs []billing.RecurringCost) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(costs).
		WithDefaultFields([]string{"id", "status", "description", "cost", "interval", "created_at"}).
		WithMonetaryFields([]string{"cost"})
}
