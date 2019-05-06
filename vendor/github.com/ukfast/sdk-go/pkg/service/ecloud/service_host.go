package ecloud

import (
	"fmt"

	"github.com/ukfast/sdk-go/pkg/connection"
)

// GetHosts retrieves a list of hosts
func (s *Service) GetHosts(parameters connection.APIRequestParameters) ([]Host, error) {
	r := connection.RequestAll{}

	var hosts []Host
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getHostsPaginatedResponseBody(parameters)
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

// GetHostsPaginated retrieves a paginated list of hosts
func (s *Service) GetHostsPaginated(parameters connection.APIRequestParameters) ([]Host, error) {
	body, err := s.getHostsPaginatedResponseBody(parameters)

	return body.Data, err
}

func (s *Service) getHostsPaginatedResponseBody(parameters connection.APIRequestParameters) (*GetHostsResponseBody, error) {
	body := &GetHostsResponseBody{}

	response, err := s.connection.Get("/ecloud/v1/hosts", parameters)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, nil)
}

// GetHost retrieves a single host by ID
func (s *Service) GetHost(hostID int) (Host, error) {
	body, err := s.getHostResponseBody(hostID)

	return body.Data, err
}

func (s *Service) getHostResponseBody(hostID int) (*GetHostResponseBody, error) {
	body := &GetHostResponseBody{}

	if hostID < 1 {
		return body, fmt.Errorf("invalid host id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/hosts/%d", hostID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, func(resp *connection.APIResponse) error {
		if response.StatusCode == 404 {
			return &HostNotFoundError{ID: hostID}
		}

		return nil
	})
}
