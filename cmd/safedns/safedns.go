package safedns

import (
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/resource"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/service/safedns"
	"github.com/spf13/cobra"
)

func SafeDNSRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "safedns",
		Short: "Commands relating to SafeDNS service",
	}

	// Child root commands
	cmd.AddCommand(safednsZoneRootCmd(f))
	cmd.AddCommand(safednsZoneRecordRootCmd(f))
	cmd.AddCommand(safednsZoneNoteRootCmd(f))
	cmd.AddCommand(safednsTemplateRootCmd(f))
	cmd.AddCommand(safednsSettingsRootCmd(f))

	return cmd
}

type SafeDNSTemplateLocatorProvider struct {
	service safedns.SafeDNSService
}

func NewSafeDNSTemplateLocatorProvider(service safedns.SafeDNSService) *SafeDNSTemplateLocatorProvider {
	return &SafeDNSTemplateLocatorProvider{service: service}
}

func (p *SafeDNSTemplateLocatorProvider) SupportedProperties() []string {
	return []string{"name"}
}

func (p *SafeDNSTemplateLocatorProvider) Locate(property string, value string) (any, error) {
	params := connection.APIRequestParameters{}
	params.WithFilter(connection.APIRequestFiltering{Property: property, Operator: connection.EQOperator, Value: []string{value}})

	return p.service.GetTemplates(params)
}

func getSafeDNSTemplateByNameOrID(service safedns.SafeDNSService, nameOrID string) (safedns.Template, error) {
	templateID, err := strconv.Atoi(nameOrID)
	if err != nil {
		locator := resource.NewResourceLocator(NewSafeDNSTemplateLocatorProvider(service))

		template, err := locator.Invoke(nameOrID)
		if err != nil {
			return safedns.Template{}, fmt.Errorf("error locating template [%s]: %s", nameOrID, err)
		}

		return template.(safedns.Template), nil
	}

	template, err := service.GetTemplate(templateID)
	if err != nil {
		return safedns.Template{}, fmt.Errorf("error retrieving template by ID [%d]: %s", templateID, err)
	}

	return template, nil
}

func getSafeDNSTemplateIDByNameOrID(service safedns.SafeDNSService, nameOrID string) (int, error) {
	templateID, err := strconv.Atoi(nameOrID)
	if err != nil {
		locator := resource.NewResourceLocator(NewSafeDNSTemplateLocatorProvider(service))

		template, err := locator.Invoke(nameOrID)
		if err != nil {
			return 0, fmt.Errorf("error locating template [%s]: %s", nameOrID, err)
		}

		return template.(safedns.Template).ID, nil
	}

	return templateID, nil
}
