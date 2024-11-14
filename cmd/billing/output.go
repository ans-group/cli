package billing

import (
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/billing"
)

type CardCollection []billing.Card

func (m CardCollection) DefaultColumns() []string {
	return []string{"id", "name", "card_number"}
}

type InvoiceCollection []billing.Invoice

func (m InvoiceCollection) DefaultColumns() []string {
	return []string{"id", "date", "paid", "net", "vat", "gross"}
}

type InvoiceQueryCollection []billing.InvoiceQuery

func (m InvoiceQueryCollection) DefaultColumns() []string {
	return []string{"id", "amount", "invoice_ids"}
}

func (m InvoiceQueryCollection) FieldValueHandlers() map[string]output.FieldValueHandlerFunc {
	return map[string]output.FieldValueHandlerFunc{
		"amount": output.MonetaryFieldValueHandler,
	}
}

type CloudCostCollection []billing.CloudCost

func (m CloudCostCollection) DefaultColumns() []string {
	return []string{"id", "server_id", "resource_type", "resource_quantity", "resource_price", "resource_period"}
}

func (m CloudCostCollection) FieldValueHandlers() map[string]output.FieldValueHandlerFunc {
	return map[string]output.FieldValueHandlerFunc{
		"resource_price":           output.MonetaryFieldValueHandler,
		"cost_since_last_invoice":  output.MonetaryFieldValueHandler,
		"cost_for_period_estimate": output.MonetaryFieldValueHandler,
	}
}

type DirectDebitCollection []billing.DirectDebit

func (m DirectDebitCollection) DefaultColumns() []string {
	return []string{"name", "status", "is_actived"}
}

type PaymentCollection []billing.Payment

func (m PaymentCollection) DefaultColumns() []string {
	return []string{"id", "date", "cost", "vat", "gross", "discount"}
}

func (m PaymentCollection) FieldValueHandlers() map[string]output.FieldValueHandlerFunc {
	return map[string]output.FieldValueHandlerFunc{
		"cost":     output.MonetaryFieldValueHandler,
		"vat":      output.MonetaryFieldValueHandler,
		"gross":    output.MonetaryFieldValueHandler,
		"discount": output.MonetaryFieldValueHandler,
	}
}

type RecurringCostCollection []billing.RecurringCost

func (m RecurringCostCollection) DefaultColumns() []string {
	return []string{"id", "status", "description", "cost", "interval", "created_at"}
}

func (m RecurringCostCollection) FieldValueHandlers() map[string]output.FieldValueHandlerFunc {
	return map[string]output.FieldValueHandlerFunc{
		"cost": output.MonetaryFieldValueHandler,
	}
}
