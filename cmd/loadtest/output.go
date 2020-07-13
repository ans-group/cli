package loadtest

import (
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func OutputLoadTestDomainsProvider(domains []ltaas.Domain) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(domains).WithDefaultFields([]string{"id", "name"})
}

func OutputLoadTestTestsProvider(tests []ltaas.Test) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(tests).WithDefaultFields([]string{"id", "name", "number_of_users", "duration", "protocol", "path"})
}

func OutputLoadTestJobsProvider(jobs []ltaas.Job) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(jobs).WithDefaultFields([]string{"id", "status", "job_start_timestamp", "job_end_timestamp"})
}

func OutputLoadTestJobResultsProvider(results []ltaas.JobResults) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(results).WithDefaultFields([]string{})
}

func OutputLoadTestJobSettingsProvider(settings []ltaas.JobSettings) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(settings).WithDefaultFields([]string{"date", "name", "duration", "max_users", "domain"})
}

func OutputLoadTestThresholdsProvider(thresholds []ltaas.Threshold) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(thresholds).WithDefaultFields([]string{"id", "name", "description", "query"})
}

func OutputLoadTestScenariosProvider(scenarios []ltaas.Scenario) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(scenarios).WithDefaultFields([]string{"id", "name", "description"})
}

func OutputLoadTestAgreementsProvider(agreements []ltaas.Agreement) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(agreements).WithDefaultFields([]string{"version"})
}

func OutputLoadTestAccountsProvider(accounts []ltaas.Account) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(accounts).WithDefaultFields([]string{"id"})
}
