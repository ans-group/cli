package sharedexchange

import (
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/sharedexchange"
)

func OutputSharedExchangeDomainsProvider(domains []sharedexchange.Domain) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(domains).WithDefaultFields([]string{"id", "name", "version", "created_at"})
}
