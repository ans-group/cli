package managedcloudflare

import (
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/managedcloudflare"
)

func OutputManagedCloudflareAccountsProvider(accounts []managedcloudflare.Account) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(accounts).WithDefaultFields([]string{"id", "name", "status", "cloudflare_account_id"})
}

func OutputManagedCloudflareSpendPlansProvider(plans []managedcloudflare.SpendPlan) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(plans).
		WithDefaultFields([]string{"id", "amount", "started_at", "ended_at"}).
		WithMonetaryFields([]string{"amount"})
}

func OutputManagedCloudflareZonesProvider(zones []managedcloudflare.Zone) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(zones).WithDefaultFields([]string{"id", "name", "account_id"})
}
