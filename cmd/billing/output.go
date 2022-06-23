package billing

import (
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/billing"
)

func OutputBillingCardsProvider(cards []billing.Card) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(cards).WithDefaultFields([]string{"id", "name", "card_number"})
}

func OutputBillingInvoicesProvider(invoices []billing.Invoice) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(invoices).WithDefaultFields([]string{"id", "date", "paid", "net", "vat", "gross"})
}

func OutputBillingInvoiceQueriesProvider(queries []billing.InvoiceQuery) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(queries).
		WithDefaultFields([]string{"id", "amount", "invoice_ids"}).
		WithMonetaryFields([]string{"amount"})
}

func OutputBillingCloudCostsProvider(costs []billing.CloudCost) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(costs).
		WithDefaultFields([]string{"id", "server_id", "resource_type", "resource_quantity", "resource_price", "resource_period"}).
		WithMonetaryFields([]string{"resource_price", "cost_since_last_invoice", "cost_for_period_estimate"})
}

func OutputBillingDirectDebitsProvider(details []billing.DirectDebit) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(details).WithDefaultFields([]string{"name", "status", "is_actived"})
}

func OutputBillingPaymentsProvider(payments []billing.Payment) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(payments).
		WithDefaultFields([]string{"id", "date", "cost", "vat", "gross", "discount"}).
		WithMonetaryFields([]string{"cost", "vat", "gross", "discount"})
}

func OutputBillingRecurringCostsProvider(costs []billing.RecurringCost) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(costs).
		WithDefaultFields([]string{"id", "status", "description", "cost", "interval", "created_at"}).
		WithMonetaryFields([]string{"cost"})
}
