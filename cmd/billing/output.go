package billing

import (
	"fmt"
	"strconv"

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
	return output.NewGenericOutputHandlerProvider(
		output.WithData(costs),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, cost := range costs {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(strconv.Itoa(cost.ID), true))
				fields.Set("server_id", output.NewFieldValue(strconv.Itoa(cost.ServerID), true))
				fields.Set("resource_type", output.NewFieldValue(cost.Resource.Type, true))
				fields.Set("resource_quantity", output.NewFieldValue(strconv.Itoa(cost.Resource.Quantity), true))
				fields.Set("resource_price", output.NewFieldValue(fmt.Sprintf("%.2f", cost.Resource.Price), true))
				fields.Set("resource_period", output.NewFieldValue(cost.Resource.Period, true))
				fields.Set("resource_usage_since_last_invoice", output.NewFieldValue(strconv.Itoa(cost.Resource.UsageSinceLastInvoice), false))
				fields.Set("resource_cost_since_last_invoice", output.NewFieldValue(fmt.Sprintf("%.2f", cost.Resource.CostSinceLastInvoice), false))
				fields.Set("resource_usage_for_period_estimate", output.NewFieldValue(strconv.Itoa(cost.Resource.UsageForPeriodEstimate), false))
				fields.Set("resource_billing_start", output.NewFieldValue(cost.Resource.BillingStart.String(), false))
				fields.Set("resource_billing_end", output.NewFieldValue(cost.Resource.BillingEnd.String(), false))
				fields.Set("resource_billing_due_date", output.NewFieldValue(cost.Resource.BillingDueDate.String(), false))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}
