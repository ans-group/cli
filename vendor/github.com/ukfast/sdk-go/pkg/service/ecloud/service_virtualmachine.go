package ecloud

import (
	"fmt"

	"github.com/ukfast/sdk-go/pkg/connection"
)

// GetVirtualMachines retrieves a list of virtual machines
func (s *Service) GetVirtualMachines(parameters connection.APIRequestParameters) ([]VirtualMachine, error) {
	r := connection.RequestAll{}

	var vms []VirtualMachine
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getVirtualMachinesPaginatedResponseBody(parameters)
		if err != nil {
			return nil, err
		}

		for _, virtualMachine := range response.Data {
			vms = append(vms, virtualMachine)
		}

		return response, nil
	}

	err := r.Invoke(parameters)

	return vms, err
}

// GetVirtualMachinesPaginated retrieves a paginated list of virtual machines
func (s *Service) GetVirtualMachinesPaginated(parameters connection.APIRequestParameters) ([]VirtualMachine, error) {
	body, err := s.getVirtualMachinesPaginatedResponseBody(parameters)

	return body.Data, err
}

func (s *Service) getVirtualMachinesPaginatedResponseBody(parameters connection.APIRequestParameters) (*GetVirtualMachinesResponseBody, error) {
	body := &GetVirtualMachinesResponseBody{}

	response, err := s.connection.Get("/ecloud/v1/vms", parameters)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse([]int{200}, body)
}

// GetVirtualMachine retrieves a single virtual machine by ID
func (s *Service) GetVirtualMachine(vmID int) (VirtualMachine, error) {
	body, err := s.getVirtualMachineResponseBody(vmID)

	return body.Data, err
}

func (s *Service) getVirtualMachineResponseBody(vmID int) (*GetVirtualMachineResponseBody, error) {
	body := &GetVirtualMachineResponseBody{}

	if vmID < 1 {
		return body, fmt.Errorf("invalid virtual machine id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/vms/%d", vmID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &VirtualMachineNotFoundError{ID: vmID}
	}

	return body, response.HandleResponse([]int{200}, body)
}

// DeleteVirtualMachine removes a virtual machine
func (s *Service) DeleteVirtualMachine(vmID int) error {
	_, err := s.deleteVirtualMachineResponseBody(vmID)

	return err
}

func (s *Service) deleteVirtualMachineResponseBody(vmID int) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if vmID < 1 {
		return body, fmt.Errorf("invalid virtual machine id")
	}

	response, err := s.connection.Delete(fmt.Sprintf("/ecloud/v1/vms/%d", vmID), nil)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &VirtualMachineNotFoundError{ID: vmID}
	}

	return body, response.HandleResponse([]int{202}, body)
}

// CreateVirtualMachine creates a new virtual machine
func (s *Service) CreateVirtualMachine(req CreateVirtualMachineRequest) (int, error) {
	body, err := s.createVirtualMachineResponseBody(req)

	return body.Data.ID, err
}

func (s *Service) createVirtualMachineResponseBody(req CreateVirtualMachineRequest) (*GetVirtualMachineResponseBody, error) {
	body := &GetVirtualMachineResponseBody{}

	response, err := s.connection.Post("/ecloud/v1/vms", &req)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse([]int{202}, body)
}

// PatchVirtualMachine patches an eCloud virtual machine
func (s *Service) PatchVirtualMachine(vmID int, patch PatchVirtualMachineRequest) error {
	_, err := s.patchVirtualMachineResponseBody(vmID, patch)

	return err
}

func (s *Service) patchVirtualMachineResponseBody(vmID int, patch PatchVirtualMachineRequest) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if vmID < 1 {
		return body, fmt.Errorf("invalid virtual machine id")
	}

	response, err := s.connection.Patch(fmt.Sprintf("/ecloud/v1/vms/%d", vmID), &patch)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &VirtualMachineNotFoundError{ID: vmID}
	}

	return body, response.HandleResponse([]int{200, 202}, body)
}

// CloneVirtualMachine clones a virtual machine
func (s *Service) CloneVirtualMachine(vmID int, req CloneVirtualMachineRequest) (int, error) {
	body, err := s.cloneVirtualMachineResponseBody(vmID, req)

	return body.Data.ID, err
}

func (s *Service) cloneVirtualMachineResponseBody(vmID int, req CloneVirtualMachineRequest) (*GetVirtualMachineResponseBody, error) {
	body := &GetVirtualMachineResponseBody{}

	if vmID < 1 {
		return body, fmt.Errorf("invalid virtual machine id")
	}

	response, err := s.connection.Post(fmt.Sprintf("/ecloud/v1/vms/%d/clone", vmID), &req)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &VirtualMachineNotFoundError{ID: vmID}
	}

	return body, response.HandleResponse([]int{202}, body)
}

// PowerOnVirtualMachine powers on a virtual machine
func (s *Service) PowerOnVirtualMachine(vmID int) error {
	_, err := s.powerOnVirtualMachineResponseBody(vmID)

	return err
}

func (s *Service) powerOnVirtualMachineResponseBody(vmID int) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if vmID < 1 {
		return body, fmt.Errorf("invalid virtual machine id")
	}

	response, err := s.connection.Put(fmt.Sprintf("/ecloud/v1/vms/%d/power-on", vmID), nil)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &VirtualMachineNotFoundError{ID: vmID}
	}

	return body, response.HandleResponse([]int{204}, body)
}

// PowerOffVirtualMachine powers off a virtual machine
func (s *Service) PowerOffVirtualMachine(vmID int) error {
	_, err := s.powerOffVirtualMachineResponseBody(vmID)

	return err
}

func (s *Service) powerOffVirtualMachineResponseBody(vmID int) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if vmID < 1 {
		return body, fmt.Errorf("invalid virtual machine id")
	}

	response, err := s.connection.Put(fmt.Sprintf("/ecloud/v1/vms/%d/power-off", vmID), nil)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &VirtualMachineNotFoundError{ID: vmID}
	}

	return body, response.HandleResponse([]int{204}, body)
}

// PowerResetVirtualMachine resets a virtual machine (hard power off)
func (s *Service) PowerResetVirtualMachine(vmID int) error {
	_, err := s.powerResetVirtualMachineResponseBody(vmID)

	return err
}

func (s *Service) powerResetVirtualMachineResponseBody(vmID int) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if vmID < 1 {
		return body, fmt.Errorf("invalid virtual machine id")
	}

	response, err := s.connection.Put(fmt.Sprintf("/ecloud/v1/vms/%d/power-reset", vmID), nil)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &VirtualMachineNotFoundError{ID: vmID}
	}

	return body, response.HandleResponse([]int{200}, body)
}

// PowerShutdownVirtualMachine shuts down a virtual machine
func (s *Service) PowerShutdownVirtualMachine(vmID int) error {
	_, err := s.powerShutdownVirtualMachineResponseBody(vmID)

	return err
}

func (s *Service) powerShutdownVirtualMachineResponseBody(vmID int) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if vmID < 1 {
		return body, fmt.Errorf("invalid virtual machine id")
	}

	response, err := s.connection.Put(fmt.Sprintf("/ecloud/v1/vms/%d/power-shutdown", vmID), nil)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &VirtualMachineNotFoundError{ID: vmID}
	}

	return body, response.HandleResponse([]int{200}, body)
}

// PowerRestartVirtualMachine resets a virtual machine (graceful power off)
func (s *Service) PowerRestartVirtualMachine(vmID int) error {
	_, err := s.powerRestartVirtualMachineResponseBody(vmID)

	return err
}

func (s *Service) powerRestartVirtualMachineResponseBody(vmID int) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if vmID < 1 {
		return body, fmt.Errorf("invalid virtual machine id")
	}

	response, err := s.connection.Put(fmt.Sprintf("/ecloud/v1/vms/%d/power-restart", vmID), nil)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &VirtualMachineNotFoundError{ID: vmID}
	}

	return body, response.HandleResponse([]int{200}, body)
}

// CreateVirtualMachineTemplate creates a virtual machine template
func (s *Service) CreateVirtualMachineTemplate(vmID int, req CreateVirtualMachineTemplateRequest) error {
	_, err := s.createVirtualMachineTemplateResponseBody(vmID, req)

	return err
}

func (s *Service) createVirtualMachineTemplateResponseBody(vmID int, req CreateVirtualMachineTemplateRequest) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if vmID < 1 {
		return body, fmt.Errorf("invalid virtual machine id")
	}

	response, err := s.connection.Post(fmt.Sprintf("/ecloud/v1/vms/%d/clone-to-template", vmID), &req)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &VirtualMachineNotFoundError{ID: vmID}
	}

	return body, response.HandleResponse([]int{202}, body)
}

// GetVirtualMachineTags retrieves a list of tags for a virtual machine
func (s *Service) GetVirtualMachineTags(vmID int, parameters connection.APIRequestParameters) ([]Tag, error) {
	r := connection.RequestAll{}

	var vms []Tag
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getVirtualMachineTagsPaginatedResponseBody(vmID, parameters)
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

// GetVirtualMachineTagsPaginated retrieves a paginated list of tags for a virtual machine
func (s *Service) GetVirtualMachineTagsPaginated(vmID int, parameters connection.APIRequestParameters) ([]Tag, error) {
	body, err := s.getVirtualMachineTagsPaginatedResponseBody(vmID, parameters)

	return body.Data, err
}

func (s *Service) getVirtualMachineTagsPaginatedResponseBody(vmID int, parameters connection.APIRequestParameters) (*GetTagsResponseBody, error) {
	body := &GetTagsResponseBody{}

	if vmID < 1 {
		return body, fmt.Errorf("invalid virtual machine id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/vms/%d/tags", vmID), parameters)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &VirtualMachineNotFoundError{ID: vmID}
	}

	return body, response.HandleResponse([]int{200}, body)
}

// GetVirtualMachineTag retrieves a single virtual machine tag by key
func (s *Service) GetVirtualMachineTag(vmID int, tagKey string) (Tag, error) {
	body, err := s.getVirtualMachineTagResponseBody(vmID, tagKey)

	return body.Data, err
}

func (s *Service) getVirtualMachineTagResponseBody(vmID int, tagKey string) (*GetTagResponseBody, error) {
	body := &GetTagResponseBody{}

	if vmID < 1 {
		return body, fmt.Errorf("invalid virtual machine id")
	}
	if tagKey == "" {
		return body, fmt.Errorf("invalid tag key")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/vms/%d/tags/%s", vmID, tagKey), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TagNotFoundError{Key: tagKey}
	}

	return body, response.HandleResponse([]int{200}, body)
}

// CreateVirtualMachineTag creates a new virtual machine tag
func (s *Service) CreateVirtualMachineTag(vmID int, req CreateTagRequest) error {
	_, err := s.createVirtualMachineTagResponseBody(vmID, req)

	return err
}

func (s *Service) createVirtualMachineTagResponseBody(vmID int, req CreateTagRequest) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if vmID < 1 {
		return body, fmt.Errorf("invalid virtual machine id")
	}

	response, err := s.connection.Post(fmt.Sprintf("/ecloud/v1/vms/%d/tags", vmID), &req)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &VirtualMachineNotFoundError{ID: vmID}
	}

	return body, response.HandleResponse([]int{201}, body)
}

// PatchVirtualMachineTag patches an eCloud virtual machine tag
func (s *Service) PatchVirtualMachineTag(vmID int, tagKey string, patch PatchTagRequest) error {
	_, err := s.patchVirtualMachineTagResponseBody(vmID, tagKey, patch)

	return err
}

func (s *Service) patchVirtualMachineTagResponseBody(vmID int, tagKey string, patch PatchTagRequest) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if vmID < 1 {
		return body, fmt.Errorf("invalid virtual machine id")
	}
	if tagKey == "" {
		return body, fmt.Errorf("invalid tag key")
	}

	response, err := s.connection.Patch(fmt.Sprintf("/ecloud/v1/vms/%d/tags/%s", vmID, tagKey), &patch)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TagNotFoundError{Key: tagKey}
	}
	return body, response.HandleResponse([]int{200}, body)
}

// DeleteVirtualMachineTag removes a virtual machine tag
func (s *Service) DeleteVirtualMachineTag(vmID int, tagKey string) error {
	_, err := s.deleteVirtualMachineTagResponseBody(vmID, tagKey)

	return err
}

func (s *Service) deleteVirtualMachineTagResponseBody(vmID int, tagKey string) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if vmID < 1 {
		return body, fmt.Errorf("invalid virtual machine id")
	}
	if tagKey == "" {
		return body, fmt.Errorf("invalid tag key")
	}

	response, err := s.connection.Delete(fmt.Sprintf("/ecloud/v1/vms/%d/tags/%s", vmID, tagKey), nil)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TagNotFoundError{Key: tagKey}
	}

	return body, response.HandleResponse([]int{204}, body)
}
