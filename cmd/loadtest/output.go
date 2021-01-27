package loadtest

import (
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func OutputLoadTestDomainsProvider(domains []ltaas.Domain) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(domains).WithDefaultFields([]string{"id", "name"})
}

func OutputLoadTestTestsProvider(tests []ltaas.Test) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(tests).WithDefaultFields([]string{"id", "name", "number_of_users", "duration", "protocol", "path"})
}

func OutputLoadTestJobsProvider(jobs []ltaas.Job) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(jobs).WithDefaultFields([]string{"id", "status", "job_start_timestamp", "job_end_timestamp"})
}

func OutputLoadTestJobResultsProvider(results []ltaas.JobResults) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(results).WithDefaultFields([]string{})
}

func OutputLoadTestJobSettingsProvider(settings []ltaas.JobSettings) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(settings).WithDefaultFields([]string{"date", "name", "duration", "max_users", "domain"})
}

func OutputLoadTestThresholdsProvider(thresholds []ltaas.Threshold) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(thresholds).WithDefaultFields([]string{"id", "name", "description", "query"})
}

func OutputLoadTestScenariosProvider(scenarios []ltaas.Scenario) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(scenarios).WithDefaultFields([]string{"id", "name", "description"})
}

func OutputLoadTestAgreementsProvider(agreements []ltaas.Agreement) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(agreements).WithDefaultFields([]string{"version"})
}

func OutputLoadTestAccountsProvider(accounts []ltaas.Account) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(accounts).WithDefaultFields([]string{"id"})
}
