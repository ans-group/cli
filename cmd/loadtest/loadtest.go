package loadtest

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/resource"
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func LoadTestRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "loadtest",
		Short: "Commands relating to load testing service",
	}

	// Child root commands
	cmd.AddCommand(loadtestDomainRootCmd(f))
	cmd.AddCommand(loadtestTestRootCmd(f))
	cmd.AddCommand(loadtestJobRootCmd(f))
	cmd.AddCommand(loadtestThresholdRootCmd(f))
	cmd.AddCommand(loadtestScenarioRootCmd(f))
	cmd.AddCommand(loadtestAgreementRootCmd(f))
	cmd.AddCommand(loadtestAccountRootCmd(f))

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
