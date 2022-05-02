package cloudflare

import (
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/cloudflare"
)

func OutputCloudflareAccountsProvider(accounts []cloudflare.Account) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(accounts).WithDefaultFields([]string{"id", "name", "status", "cloudflare_account_id"})
}

func OutputCloudflareSpendPlansProvider(plans []cloudflare.SpendPlan) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(plans).
		WithDefaultFields([]string{"id", "amount", "started_at", "ended_at"}).
		WithMonetaryFields([]string{"amount"})
}

func OutputCloudflareZonesProvider(zones []cloudflare.Zone) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(zones).WithDefaultFields([]string{"id", "name", "account_id"})
}

func OutputCloudflareSubscriptionsProvider(subscriptions []cloudflare.Subscription) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(subscriptions).WithDefaultFields([]string{"id", "name", "type", "price"})
}

func OutputCloudflareTotalSpendsProvider(spends []cloudflare.TotalSpend) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(spends).WithDefaultFields([]string{"id", "name", "type", "price"})
}
