package cmd

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/cli/internal/pkg/resource"
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func loadtestRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "loadtest",
		Short: "Commands relating to load testing service",
	}

	// Child root commands
	cmd.AddCommand(loadtestDomainRootCmd())
	cmd.AddCommand(loadtestTestRootCmd())
	cmd.AddCommand(loadtestJobRootCmd())
	cmd.AddCommand(loadtestThresholdRootCmd())
	cmd.AddCommand(loadtestScenarioRootCmd())
	cmd.AddCommand(loadtestAgreementRootCmd())

	return cmd
}

// Currently non-functional, as domains aren't yet filterable server-side
type LoadTestDomainLocatorProvider struct {
	service ltaas.LTaaSService
}

func NewLoadTestDomainLocatorProvider(service ltaas.LTaaSService) *LoadTestDomainLocatorProvider {
	return &LoadTestDomainLocatorProvider{service: service}
}

func (p *LoadTestDomainLocatorProvider) SupportedProperties() []string {
	return []string{"name"}
}

func (p *LoadTestDomainLocatorProvider) Locate(property string, value string) (interface{}, error) {
	params := connection.APIRequestParameters{}
	params.WithFilter(connection.APIRequestFiltering{Property: property, Operator: connection.EQOperator, Value: []string{value}})

	return p.service.GetDomains(params)
}

func getLoadTestDomainByNameOrID(service ltaas.LTaaSService, nameOrID string) (ltaas.Domain, error) {
	_, err := uuid.Parse(nameOrID)
	if err != nil {
		locator := resource.NewResourceLocator(NewLoadTestDomainLocatorProvider(service))

		domain, err := locator.Invoke(nameOrID)
		if err != nil {
			return ltaas.Domain{}, fmt.Errorf("Error locating domain [%s]: %s", nameOrID, err)
		}

		return domain.(ltaas.Domain), nil
	}

	domain, err := service.GetDomain(nameOrID)
	if err != nil {
		return ltaas.Domain{}, fmt.Errorf("Error retrieving domain by ID [%s]: %s", nameOrID, err)
	}

	return domain, nil
}

func outputLoadTestDomains(domains []ltaas.Domain) error {
	err := Output(output.NewGenericOutputHandlerProvider(domains, []string{"id", "name"}, nil))
	if err != nil {
		return fmt.Errorf("Failed to output domains: %s", err)
	}

	return nil
}

func outputLoadTestTests(tests []ltaas.Test) error {
	err := Output(output.NewGenericOutputHandlerProvider(tests, []string{"id", "name", "number_of_users", "duration", "protocol", "path"}, nil))
	if err != nil {
		return fmt.Errorf("Failed to output tests: %s", err)
	}

	return nil
}

func outputLoadTestJobs(jobs []ltaas.Job) error {
	err := Output(output.NewGenericOutputHandlerProvider(jobs, []string{"id", "status", "job_start_timestamp", "job_end_timestamp"}, nil))
	if err != nil {
		return fmt.Errorf("Failed to output jobs: %s", err)
	}

	return nil
}

func outputLoadTestJobResults(results []ltaas.JobResults) error {
	err := Output(output.NewGenericOutputHandlerProvider(results, []string{"id", "status", "job_start_timestamp", "job_end_timestamp"}, nil))
	if err != nil {
		return fmt.Errorf("Failed to output job results: %s", err)
	}

	return nil
}

func outputLoadTestJobSettings(settings []ltaas.JobSettings) error {
	err := Output(output.NewGenericOutputHandlerProvider(settings, []string{"date", "name", "duration", "max_users", "domain"}, nil))
	if err != nil {
		return fmt.Errorf("Failed to output job settings: %s", err)
	}

	return nil
}

func outputLoadTestThresholds(thresholds []ltaas.Threshold) error {
	err := Output(output.NewGenericOutputHandlerProvider(thresholds, []string{"id", "name", "description", "query"}, nil))
	if err != nil {
		return fmt.Errorf("Failed to output thresholds: %s", err)
	}

	return nil
}

func outputLoadTestScenarios(scenarios []ltaas.Scenario) error {
	err := Output(output.NewGenericOutputHandlerProvider(scenarios, []string{"id", "name", "description"}, nil))
	if err != nil {
		return fmt.Errorf("Failed to output scenarios: %s", err)
	}

	return nil
}

func outputLoadTestAgreements(agreements []ltaas.Agreement) error {
	err := Output(output.NewGenericOutputHandlerProvider(agreements, []string{"version"}, nil))
	if err != nil {
		return fmt.Errorf("Failed to output agreements: %s", err)
	}

	return nil
}
