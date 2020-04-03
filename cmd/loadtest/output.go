package loadtest

import (
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func OutputLoadTestDomainsProvider(domains []ltaas.Domain) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(domains, []string{"id", "name"}, nil)
}

func OutputLoadTestTestsProvider(tests []ltaas.Test) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(tests, []string{"id", "name", "number_of_users", "duration", "protocol", "path"}, nil)
}

func OutputLoadTestJobsProvider(jobs []ltaas.Job) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(jobs, []string{"id", "status", "job_start_timestamp", "job_end_timestamp"}, nil)
}

func OutputLoadTestJobResultsProvider(results []ltaas.JobResults) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(results, []string{}, nil)
}

func OutputLoadTestJobSettingsProvider(settings []ltaas.JobSettings) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(settings, []string{"date", "name", "duration", "max_users", "domain"}, nil)
}

func OutputLoadTestThresholdsProvider(thresholds []ltaas.Threshold) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(thresholds, []string{"id", "name", "description", "query"}, nil)
}

func OutputLoadTestScenariosProvider(scenarios []ltaas.Scenario) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(scenarios, []string{"id", "name", "description"}, nil)
}

func OutputLoadTestAgreementsProvider(agreements []ltaas.Agreement) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(agreements, []string{"version"}, nil)
}

func OutputLoadTestAccountsProvider(accounts []ltaas.Account) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(accounts, []string{"id"}, nil)
}
