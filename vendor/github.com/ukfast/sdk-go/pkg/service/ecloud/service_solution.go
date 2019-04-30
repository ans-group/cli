package ecloud

import (
	"fmt"

	"github.com/ukfast/sdk-go/pkg/connection"
)

// GetSolutions retrieves a list of solutions
func (s *Service) GetSolutions(parameters connection.APIRequestParameters) ([]Solution, error) {
	r := connection.RequestAll{}

	var solutions []Solution
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getSolutionsPaginatedResponseBody(parameters)
		if err != nil {
			return nil, err
		}

		for _, solution := range response.Data {
			solutions = append(solutions, solution)
		}

		return response, nil
	}

	err := r.Invoke(parameters)

	return solutions, err
}

// GetSolutionsPaginated retrieves a paginated list of solutions
func (s *Service) GetSolutionsPaginated(parameters connection.APIRequestParameters) ([]Solution, error) {
	body, err := s.getSolutionsPaginatedResponseBody(parameters)

	return body.Data, err
}

func (s *Service) getSolutionsPaginatedResponseBody(parameters connection.APIRequestParameters) (*GetSolutionsResponseBody, error) {
	body := &GetSolutionsResponseBody{}

	response, err := s.connection.Get("/ecloud/v1/solutions", parameters)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse([]int{}, body)
}

// GetSolution retrieves a single Solution by ID
func (s *Service) GetSolution(solutionID int) (Solution, error) {
	body, err := s.getSolutionResponseBody(solutionID)

	return body.Data, err
}

func (s *Service) getSolutionResponseBody(solutionID int) (*GetSolutionResponseBody, error) {
	body := &GetSolutionResponseBody{}

	if solutionID < 1 {
		return body, fmt.Errorf("invalid solution id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/solutions/%d", solutionID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &SolutionNotFoundError{ID: solutionID}
	}

	return body, response.HandleResponse([]int{}, body)
}

// PatchSolution patches an eCloud solution
func (s *Service) PatchSolution(solutionID int, patch PatchSolutionRequest) (int, error) {
	body, err := s.patchSolutionResponseBody(solutionID, patch)

	return body.Data.ID, err
}

func (s *Service) patchSolutionResponseBody(solutionID int, patch PatchSolutionRequest) (*GetSolutionResponseBody, error) {
	body := &GetSolutionResponseBody{}

	if solutionID < 1 {
		return body, fmt.Errorf("invalid solution id")
	}

	response, err := s.connection.Patch(fmt.Sprintf("/ecloud/v1/solutions/%d", solutionID), &patch)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &SolutionNotFoundError{ID: solutionID}
	}

	return body, response.HandleResponse([]int{}, body)
}

// GetSolutionVirtualMachines retrieves a list of virtual machines within a solution
func (s *Service) GetSolutionVirtualMachines(solutionID int, parameters connection.APIRequestParameters) ([]VirtualMachine, error) {
	r := connection.RequestAll{}

	var vms []VirtualMachine
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getSolutionVirtualMachinesPaginatedResponseBody(solutionID, parameters)
		if err != nil {
			return nil, err
		}

		for _, vm := range response.Data {
			vms = append(vms, vm)
		}

		return response, nil
	}

	err := r.Invoke(parameters)

	return vms, err
}

// GetSolutionVirtualMachinesPaginated retrieves a paginated list of virtual machines within a solution
func (s *Service) GetSolutionVirtualMachinesPaginated(solutionID int, parameters connection.APIRequestParameters) ([]VirtualMachine, error) {
	body, err := s.getSolutionVirtualMachinesPaginatedResponseBody(solutionID, parameters)

	return body.Data, err
}

func (s *Service) getSolutionVirtualMachinesPaginatedResponseBody(solutionID int, parameters connection.APIRequestParameters) (*GetVirtualMachinesResponseBody, error) {
	body := &GetVirtualMachinesResponseBody{}

	if solutionID < 1 {
		return body, fmt.Errorf("invalid solution id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/solutions/%d/vms", solutionID), parameters)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &SolutionNotFoundError{ID: solutionID}
	}

	return body, response.HandleResponse([]int{}, body)
}

// GetSolutionSites retrieves a list of virtual machines within a solution
func (s *Service) GetSolutionSites(solutionID int, parameters connection.APIRequestParameters) ([]Site, error) {
	r := connection.RequestAll{}

	var sites []Site
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getSolutionSitesPaginatedResponseBody(solutionID, parameters)
		if err != nil {
			return nil, err
		}

		for _, site := range response.Data {
			sites = append(sites, site)
		}

		return response, nil
	}

	err := r.Invoke(parameters)

	return sites, err
}

// GetSolutionSitesPaginated retrieves a paginated list of sites within a solution
func (s *Service) GetSolutionSitesPaginated(solutionID int, parameters connection.APIRequestParameters) ([]Site, error) {
	body, err := s.getSolutionSitesPaginatedResponseBody(solutionID, parameters)

	return body.Data, err
}

func (s *Service) getSolutionSitesPaginatedResponseBody(solutionID int, parameters connection.APIRequestParameters) (*GetSitesResponseBody, error) {
	body := &GetSitesResponseBody{}

	if solutionID < 1 {
		return body, fmt.Errorf("invalid solution id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/solutions/%d/sites", solutionID), parameters)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &SolutionNotFoundError{ID: solutionID}
	}

	return body, response.HandleResponse([]int{}, body)
}

// GetSolutionDatastores retrieves a list of datastores within a solution
func (s *Service) GetSolutionDatastores(solutionID int, parameters connection.APIRequestParameters) ([]Datastore, error) {
	r := connection.RequestAll{}

	var datastores []Datastore
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getSolutionDatastoresPaginatedResponseBody(solutionID, parameters)
		if err != nil {
			return nil, err
		}

		for _, datastore := range response.Data {
			datastores = append(datastores, datastore)
		}

		return response, nil
	}

	err := r.Invoke(parameters)

	return datastores, err
}

// GetSolutionDatastoresPaginated retrieves a paginated list of datastores within a solution
func (s *Service) GetSolutionDatastoresPaginated(solutionID int, parameters connection.APIRequestParameters) ([]Datastore, error) {
	body, err := s.getSolutionDatastoresPaginatedResponseBody(solutionID, parameters)

	return body.Data, err
}

func (s *Service) getSolutionDatastoresPaginatedResponseBody(solutionID int, parameters connection.APIRequestParameters) (*GetDatastoresResponseBody, error) {
	body := &GetDatastoresResponseBody{}

	if solutionID < 1 {
		return body, fmt.Errorf("invalid solution id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/solutions/%d/datastores", solutionID), parameters)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &SolutionNotFoundError{ID: solutionID}
	}

	return body, response.HandleResponse([]int{}, body)
}

// GetSolutionHosts retrieves a list of hosts within a solution
func (s *Service) GetSolutionHosts(solutionID int, parameters connection.APIRequestParameters) ([]Host, error) {
	r := connection.RequestAll{}

	var hosts []Host
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getSolutionHostsPaginatedResponseBody(solutionID, parameters)
		if err != nil {
			return nil, err
		}

		for _, host := range response.Data {
			hosts = append(hosts, host)
		}

		return response, nil
	}

	err := r.Invoke(parameters)

	return hosts, err
}

// GetSolutionHostsPaginated retrieves a paginated list of hosts within a solution
func (s *Service) GetSolutionHostsPaginated(solutionID int, parameters connection.APIRequestParameters) ([]Host, error) {
	body, err := s.getSolutionHostsPaginatedResponseBody(solutionID, parameters)

	return body.Data, err
}

func (s *Service) getSolutionHostsPaginatedResponseBody(solutionID int, parameters connection.APIRequestParameters) (*GetHostsResponseBody, error) {
	body := &GetHostsResponseBody{}

	if solutionID < 1 {
		return body, fmt.Errorf("invalid solution id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/solutions/%d/hosts", solutionID), parameters)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &SolutionNotFoundError{ID: solutionID}
	}

	return body, response.HandleResponse([]int{}, body)
}

// GetSolutionNetworks retrieves a list of networks within a solution
func (s *Service) GetSolutionNetworks(solutionID int, parameters connection.APIRequestParameters) ([]Network, error) {
	r := connection.RequestAll{}

	var networks []Network
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getSolutionNetworksPaginatedResponseBody(solutionID, parameters)
		if err != nil {
			return nil, err
		}

		for _, network := range response.Data {
			networks = append(networks, network)
		}

		return response, nil
	}

	err := r.Invoke(parameters)

	return networks, err
}

// GetSolutionNetworksPaginated retrieves a paginated list of networks within a solution
func (s *Service) GetSolutionNetworksPaginated(solutionID int, parameters connection.APIRequestParameters) ([]Network, error) {
	body, err := s.getSolutionNetworksPaginatedResponseBody(solutionID, parameters)

	return body.Data, err
}

func (s *Service) getSolutionNetworksPaginatedResponseBody(solutionID int, parameters connection.APIRequestParameters) (*GetNetworksResponseBody, error) {
	body := &GetNetworksResponseBody{}

	if solutionID < 1 {
		return body, fmt.Errorf("invalid solution id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/solutions/%d/networks", solutionID), parameters)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &SolutionNotFoundError{ID: solutionID}
	}

	return body, response.HandleResponse([]int{}, body)
}

// GetSolutionFirewalls retrieves a list of firewalls within a solution
func (s *Service) GetSolutionFirewalls(solutionID int, parameters connection.APIRequestParameters) ([]Firewall, error) {
	r := connection.RequestAll{}

	var firewalls []Firewall
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getSolutionFirewallsPaginatedResponseBody(solutionID, parameters)
		if err != nil {
			return nil, err
		}

		for _, firewall := range response.Data {
			firewalls = append(firewalls, firewall)
		}

		return response, nil
	}

	err := r.Invoke(parameters)

	return firewalls, err
}

// GetSolutionFirewallsPaginated retrieves a paginated list of firewalls within a solution
func (s *Service) GetSolutionFirewallsPaginated(solutionID int, parameters connection.APIRequestParameters) ([]Firewall, error) {
	body, err := s.getSolutionFirewallsPaginatedResponseBody(solutionID, parameters)

	return body.Data, err
}

func (s *Service) getSolutionFirewallsPaginatedResponseBody(solutionID int, parameters connection.APIRequestParameters) (*GetFirewallsResponseBody, error) {
	body := &GetFirewallsResponseBody{}

	if solutionID < 1 {
		return body, fmt.Errorf("invalid solution id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/solutions/%d/firewalls", solutionID), parameters)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &SolutionNotFoundError{ID: solutionID}
	}

	return body, response.HandleResponse([]int{}, body)
}

// GetSolutionTemplates retrieves a list of virtual machines within a solution
func (s *Service) GetSolutionTemplates(solutionID int, parameters connection.APIRequestParameters) ([]Template, error) {
	r := connection.RequestAll{}

	var templates []Template
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getSolutionTemplatesPaginatedResponseBody(solutionID, parameters)
		if err != nil {
			return nil, err
		}

		for _, template := range response.Data {
			templates = append(templates, template)
		}

		return response, nil
	}

	err := r.Invoke(parameters)

	return templates, err
}

// GetSolutionTemplatesPaginated retrieves a paginated list of templates within a solution
func (s *Service) GetSolutionTemplatesPaginated(solutionID int, parameters connection.APIRequestParameters) ([]Template, error) {
	body, err := s.getSolutionTemplatesPaginatedResponseBody(solutionID, parameters)

	return body.Data, err
}

func (s *Service) getSolutionTemplatesPaginatedResponseBody(solutionID int, parameters connection.APIRequestParameters) (*GetTemplatesResponseBody, error) {
	body := &GetTemplatesResponseBody{}

	if solutionID < 1 {
		return body, fmt.Errorf("invalid solution id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/solutions/%d/templates", solutionID), parameters)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &SolutionNotFoundError{ID: solutionID}
	}

	return body, response.HandleResponse([]int{}, body)
}

// GetSolutionTemplate retrieves a single solution template by name
func (s *Service) GetSolutionTemplate(solutionID int, templateName string) (Template, error) {
	body, err := s.getSolutionTemplateResponseBody(solutionID, templateName)

	return body.Data, err
}

func (s *Service) getSolutionTemplateResponseBody(solutionID int, templateName string) (*GetTemplateResponseBody, error) {
	body := &GetTemplateResponseBody{}

	if solutionID < 1 {
		return body, fmt.Errorf("invalid solution id")
	}
	if templateName == "" {
		return body, fmt.Errorf("invalid template name")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/solutions/%d/templates/%s", solutionID, templateName), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TemplateNotFoundError{Name: templateName}
	}

	return body, response.HandleResponse([]int{}, body)
}

// RenameSolutionTemplate renames a solution template
func (s *Service) RenameSolutionTemplate(solutionID int, templateName string, req RenameTemplateRequest) error {
	_, err := s.renameSolutionTemplateResponseBody(solutionID, templateName, req)

	return err
}

func (s *Service) renameSolutionTemplateResponseBody(solutionID int, templateName string, req RenameTemplateRequest) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if solutionID < 1 {
		return body, fmt.Errorf("invalid solution id")
	}
	if templateName == "" {
		return body, fmt.Errorf("invalid template name")
	}

	response, err := s.connection.Post(fmt.Sprintf("/ecloud/v1/solutions/%d/templates/%s/move", solutionID, templateName), &req)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TemplateNotFoundError{Name: templateName}
	}

	return body, response.HandleResponse([]int{}, body)
}

// DeleteSolutionTemplate removes a solution template
func (s *Service) DeleteSolutionTemplate(solutionID int, templateName string) error {
	_, err := s.deleteSolutionTemplateResponseBody(solutionID, templateName)

	return err
}

func (s *Service) deleteSolutionTemplateResponseBody(solutionID int, templateName string) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if solutionID < 1 {
		return body, fmt.Errorf("invalid solution id")
	}
	if templateName == "" {
		return body, fmt.Errorf("invalid template name")
	}

	response, err := s.connection.Delete(fmt.Sprintf("/ecloud/v1/solutions/%d/templates/%s", solutionID, templateName), nil)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TemplateNotFoundError{Name: templateName}
	}

	return body, response.HandleResponse([]int{}, body)
}

// GetSolutionTags retrieves a list of tags for a solution
func (s *Service) GetSolutionTags(solutionID int, parameters connection.APIRequestParameters) ([]Tag, error) {
	r := connection.RequestAll{}

	var tags []Tag
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getSolutionTagsPaginatedResponseBody(solutionID, parameters)
		if err != nil {
			return nil, err
		}

		for _, tag := range response.Data {
			tags = append(tags, tag)
		}

		return response, nil
	}

	err := r.Invoke(parameters)

	return tags, err
}

// GetSolutionTagsPaginated retrieves a paginated list of tags for a solution
func (s *Service) GetSolutionTagsPaginated(solutionID int, parameters connection.APIRequestParameters) ([]Tag, error) {
	body, err := s.getSolutionTagsPaginatedResponseBody(solutionID, parameters)

	return body.Data, err
}

func (s *Service) getSolutionTagsPaginatedResponseBody(solutionID int, parameters connection.APIRequestParameters) (*GetTagsResponseBody, error) {
	body := &GetTagsResponseBody{}

	if solutionID < 1 {
		return body, fmt.Errorf("invalid solution id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/solutions/%d/tags", solutionID), parameters)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &SolutionNotFoundError{ID: solutionID}
	}

	return body, response.HandleResponse([]int{}, body)
}

// GetSolutionTag retrieves a single solution tag by key
func (s *Service) GetSolutionTag(solutionID int, tagKey string) (Tag, error) {
	body, err := s.getSolutionTagResponseBody(solutionID, tagKey)

	return body.Data, err
}

func (s *Service) getSolutionTagResponseBody(solutionID int, tagKey string) (*GetTagResponseBody, error) {
	body := &GetTagResponseBody{}

	if solutionID < 1 {
		return body, fmt.Errorf("invalid solution id")
	}
	if tagKey == "" {
		return body, fmt.Errorf("invalid tag key")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/solutions/%d/tags/%s", solutionID, tagKey), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TagNotFoundError{Key: tagKey}
	}

	return body, response.HandleResponse([]int{}, body)
}

// CreateSolutionTag creates a new solution tag
func (s *Service) CreateSolutionTag(solutionID int, req CreateTagRequest) error {
	_, err := s.createSolutionTagResponseBody(solutionID, req)

	return err
}

func (s *Service) createSolutionTagResponseBody(solutionID int, req CreateTagRequest) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if solutionID < 1 {
		return body, fmt.Errorf("invalid solution id")
	}

	response, err := s.connection.Post(fmt.Sprintf("/ecloud/v1/solutions/%d/tags", solutionID), &req)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &SolutionNotFoundError{ID: solutionID}
	}

	return body, response.HandleResponse([]int{}, body)
}

// PatchSolutionTag patches an eCloud solution tag
func (s *Service) PatchSolutionTag(solutionID int, tagKey string, patch PatchTagRequest) error {
	_, err := s.patchSolutionTagResponseBody(solutionID, tagKey, patch)

	return err
}

func (s *Service) patchSolutionTagResponseBody(solutionID int, tagKey string, patch PatchTagRequest) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if solutionID < 1 {
		return body, fmt.Errorf("invalid solution id")
	}
	if tagKey == "" {
		return body, fmt.Errorf("invalid tag key")
	}

	response, err := s.connection.Patch(fmt.Sprintf("/ecloud/v1/solutions/%d/tags/%s", solutionID, tagKey), &patch)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TagNotFoundError{Key: tagKey}
	}

	return body, response.HandleResponse([]int{}, body)
}

// DeleteSolutionTag removes a solution tag
func (s *Service) DeleteSolutionTag(solutionID int, tagKey string) error {
	_, err := s.deleteSolutionTagResponseBody(solutionID, tagKey)

	return err
}

func (s *Service) deleteSolutionTagResponseBody(solutionID int, tagKey string) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if solutionID < 1 {
		return body, fmt.Errorf("invalid solution id")
	}
	if tagKey == "" {
		return body, fmt.Errorf("invalid tag key")
	}

	response, err := s.connection.Delete(fmt.Sprintf("/ecloud/v1/solutions/%d/tags/%s", solutionID, tagKey), nil)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TagNotFoundError{Key: tagKey}
	}

	return body, response.HandleResponse([]int{}, body)
}
