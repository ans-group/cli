package cloudflare

import (
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/cloudflare"
)

type AccountCollection []cloudflare.Account

func (m AccountCollection) DefaultColumns() []string {
	return []string{"id", "name", "status", "cloudflare_account_id"}
}

type SpendPlanCollection []cloudflare.SpendPlan

func (m SpendPlanCollection) DefaultColumns() []string {
	return []string{"id", "amount", "started_at", "ended_at"}
}

func (m SpendPlanCollection) FieldValueHandlers() map[string]output.FieldValueHandlerFunc {
	return map[string]output.FieldValueHandlerFunc{
		"amount": output.MonetaryFieldValueHandler,
	}
}

type ZoneCollection []cloudflare.Zone

func (m ZoneCollection) DefaultColumns() []string {
	return []string{"id", "name", "account_id"}
}

type SubscriptionCollection []cloudflare.Subscription

func (m SubscriptionCollection) DefaultColumns() []string {
	return []string{"id", "name", "type", "price"}
}

type TotalSpendCollection []cloudflare.TotalSpend

func (m TotalSpendCollection) DefaultColumns() []string {
	return []string{"id", "name", "type", "price"}
}
