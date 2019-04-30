package ecloud

import (
	"fmt"

	"github.com/ukfast/sdk-go/pkg/connection"
)

// GetAppliances retrieves a list of appliances
func (s *Service) GetAppliances(parameters connection.APIRequestParameters) ([]Appliance, error) {
	r := connection.RequestAll{}

	var appliances []Appliance
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getAppliancesPaginatedResponseBody(parameters)
		if err != nil {
			return nil, err
		}

		for _, appliance := range response.Data {
			appliances = append(appliances, appliance)
		}

		return response, nil
	}

	err := r.Invoke(parameters)

	return appliances, err
}

// GetAppliancesPaginated retrieves a paginated list of appliances
func (s *Service) GetAppliancesPaginated(parameters connection.APIRequestParameters) ([]Appliance, error) {
	body, err := s.getAppliancesPaginatedResponseBody(parameters)

	return body.Data, err
}

func (s *Service) getAppliancesPaginatedResponseBody(parameters connection.APIRequestParameters) (*GetAppliancesResponseBody, error) {
	body := &GetAppliancesResponseBody{}

	response, err := s.connection.Get("/ecloud/v1/appliances", parameters)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse([]int{}, body)
}

// GetAppliance retrieves a single Appliance by ID
func (s *Service) GetAppliance(applianceID string) (Appliance, error) {
	body, err := s.getApplianceResponseBody(applianceID)

	return body.Data, err
}

func (s *Service) getApplianceResponseBody(applianceID string) (*GetApplianceResponseBody, error) {
	body := &GetApplianceResponseBody{}

	if applianceID == "" {
		return body, fmt.Errorf("invalid appliance id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/appliances/%s", applianceID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &ApplianceNotFoundError{ID: applianceID}
	}

	return body, response.HandleResponse([]int{}, body)
}

// GetApplianceParameters retrieves a list of appliance parameters
func (s *Service) GetApplianceParameters(applianceID string, reqParameters connection.APIRequestParameters) ([]ApplianceParameter, error) {
	r := connection.RequestAll{}

	var parameters []ApplianceParameter
	r.GetNext = func(reqParameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getApplianceParametersPaginatedResponseBody(applianceID, reqParameters)
		if err != nil {
			return nil, err
		}

		for _, parameter := range response.Data {
			parameters = append(parameters, parameter)
		}

		return response, nil
	}

	err := r.Invoke(reqParameters)

	return parameters, err
}

// GetApplianceParametersPaginated retrieves a paginated list of appliance parameters
func (s *Service) GetApplianceParametersPaginated(applianceID string, parameters connection.APIRequestParameters) ([]ApplianceParameter, error) {
	body, err := s.getApplianceParametersPaginatedResponseBody(applianceID, parameters)

	return body.Data, err
}

func (s *Service) getApplianceParametersPaginatedResponseBody(applianceID string, parameters connection.APIRequestParameters) (*GetApplianceParametersResponseBody, error) {
	body := &GetApplianceParametersResponseBody{}

	if applianceID == "" {
		return body, fmt.Errorf("invalid appliance id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/appliances/%s/parameters", applianceID), parameters)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &ApplianceNotFoundError{ID: applianceID}
	}

	return body, response.HandleResponse([]int{}, body)
}
