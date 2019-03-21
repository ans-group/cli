package registrar

import (
	"fmt"

	"github.com/ukfast/sdk-go/pkg/connection"
)

// GetDomains retrieves a list of domains
func (s *Service) GetDomains(parameters connection.APIRequestParameters) ([]Domain, error) {
	r := connection.RequestAll{}

	var domains []Domain
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getDomainsPaginatedResponseBody(parameters)
		if err != nil {
			return nil, err
		}

		for _, domain := range response.Data {
			domains = append(domains, domain)
		}

		return response, nil
	}

	err := r.Invoke(parameters)

	return domains, err
}

// GetDomainsPaginated retrieves a paginated list of domains
func (s *Service) GetDomainsPaginated(parameters connection.APIRequestParameters) ([]Domain, error) {
	body, err := s.getDomainsPaginatedResponseBody(parameters)

	return body.Data, err
}

func (s *Service) getDomainsPaginatedResponseBody(parameters connection.APIRequestParameters) (*GetDomainsResponseBody, error) {
	body := &GetDomainsResponseBody{}

	response, err := s.connection.Get("/registrar/v1/domains", parameters)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse([]int{200}, body)
}

// GetDomain retrieves a single domain by name
func (s *Service) GetDomain(domainName string) (Domain, error) {
	body, err := s.getDomainResponseBody(domainName)

	return body.Data, err
}

func (s *Service) getDomainResponseBody(domainName string) (*GetDomainResponseBody, error) {
	body := &GetDomainResponseBody{}

	if domainName == "" {
		return body, fmt.Errorf("invalid domain name")
	}

	response, err := s.connection.Get(fmt.Sprintf("/registrar/v1/domains/%s", domainName), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &DomainNotFoundError{Name: domainName}
	}

	return body, response.HandleResponse([]int{200}, body)
}

// GetDomainNameservers retrieves the nameservers for a domain
func (s *Service) GetDomainNameservers(domainName string) ([]Nameserver, error) {
	body, err := s.getDomainNameserversResponseBody(domainName)

	return body.Data, err
}

func (s *Service) getDomainNameserversResponseBody(domainName string) (*GetNameserversResponseBody, error) {
	body := &GetNameserversResponseBody{}

	if domainName == "" {
		return body, fmt.Errorf("invalid domain name")
	}

	response, err := s.connection.Get(fmt.Sprintf("/registrar/v1/domains/%s/nameservers", domainName), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &DomainNotFoundError{Name: domainName}
	}

	return body, response.HandleResponse([]int{200}, body)
}
