package ecloud

import (
	"fmt"

	"github.com/ukfast/sdk-go/pkg/connection"
)

// GetFirewalls retrieves a list of firewalls
func (s *Service) GetFirewalls(parameters connection.APIRequestParameters) ([]Firewall, error) {
	r := connection.RequestAll{}

	var firewalls []Firewall
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getFirewallsPaginatedResponseBody(parameters)
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

// GetFirewallsPaginated retrieves a paginated list of firewalls
func (s *Service) GetFirewallsPaginated(parameters connection.APIRequestParameters) ([]Firewall, error) {
	body, err := s.getFirewallsPaginatedResponseBody(parameters)

	return body.Data, err
}

func (s *Service) getFirewallsPaginatedResponseBody(parameters connection.APIRequestParameters) (*GetFirewallsResponseBody, error) {
	body := &GetFirewallsResponseBody{}

	response, err := s.connection.Get("/ecloud/v1/firewalls", parameters)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse([]int{200}, body)
}

// GetFirewall retrieves a single firewall by ID
func (s *Service) GetFirewall(firewallID int) (Firewall, error) {
	body, err := s.getFirewallResponseBody(firewallID)

	return body.Data, err
}

func (s *Service) getFirewallResponseBody(firewallID int) (*GetFirewallResponseBody, error) {
	body := &GetFirewallResponseBody{}

	if firewallID < 1 {
		return body, fmt.Errorf("invalid firewall id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/firewalls/%d", firewallID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &FirewallNotFoundError{ID: firewallID}
	}

	return body, response.HandleResponse([]int{200}, body)
}

// GetFirewallConfig retrieves a single firewall config by ID
func (s *Service) GetFirewallConfig(firewallID int) (FirewallConfig, error) {
	body, err := s.getFirewallConfigResponseBody(firewallID)

	return body.Data, err
}

func (s *Service) getFirewallConfigResponseBody(firewallID int) (*GetFirewallConfigResponseBody, error) {
	body := &GetFirewallConfigResponseBody{}

	if firewallID < 1 {
		return body, fmt.Errorf("invalid firewall id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/firewalls/%d/config", firewallID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &FirewallNotFoundError{ID: firewallID}
	}

	return body, response.HandleResponse([]int{200}, body)
}
