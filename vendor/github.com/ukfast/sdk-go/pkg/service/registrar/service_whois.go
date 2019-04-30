package registrar

import (
	"fmt"

	"github.com/ukfast/sdk-go/pkg/connection"
)

// GetWhois retrieves WHOIS information for a single domain
func (s *Service) GetWhois(domainName string) (Whois, error) {
	body, err := s.getWhoisResponseBody(domainName)

	return body.Data, err
}

func (s *Service) getWhoisResponseBody(domainName string) (*GetWhoisResponseBody, error) {
	body := &GetWhoisResponseBody{}

	if domainName == "" {
		return body, fmt.Errorf("invalid domain name")
	}

	response, err := s.connection.Get(fmt.Sprintf("/registrar/v1/whois/%s", domainName), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &DomainNotFoundError{Name: domainName}
	}

	return body, response.HandleResponse([]int{}, body)
}

// GetWhoisRaw retrieves raw WHOIS information for a single domain
func (s *Service) GetWhoisRaw(domainName string) (string, error) {
	body, err := s.getWhoisRawResponseBody(domainName)

	return body.Data, err
}

func (s *Service) getWhoisRawResponseBody(domainName string) (*GetWhoisRawResponseBody, error) {
	body := &GetWhoisRawResponseBody{}

	if domainName == "" {
		return body, fmt.Errorf("invalid domain name")
	}

	response, err := s.connection.Get(fmt.Sprintf("/registrar/v1/whois/%s/raw", domainName), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &DomainNotFoundError{Name: domainName}
	}

	return body, response.HandleResponse([]int{}, body)
}
