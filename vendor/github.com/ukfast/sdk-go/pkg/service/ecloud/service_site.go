package ecloud

import (
	"fmt"

	"github.com/ukfast/sdk-go/pkg/connection"
)

// GetSites retrieves a list of sites
func (s *Service) GetSites(parameters connection.APIRequestParameters) ([]Site, error) {
	r := connection.RequestAll{}

	var sites []Site
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getSitesPaginatedResponseBody(parameters)
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

// GetSitesPaginated retrieves a paginated list of sites
func (s *Service) GetSitesPaginated(parameters connection.APIRequestParameters) ([]Site, error) {
	body, err := s.getSitesPaginatedResponseBody(parameters)

	return body.Data, err
}

func (s *Service) getSitesPaginatedResponseBody(parameters connection.APIRequestParameters) (*GetSitesResponseBody, error) {
	body := &GetSitesResponseBody{}

	response, err := s.connection.Get("/ecloud/v1/sites", parameters)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, nil)
}

// GetSite retrieves a single site by ID
func (s *Service) GetSite(siteID int) (Site, error) {
	body, err := s.getSiteResponseBody(siteID)

	return body.Data, err
}

func (s *Service) getSiteResponseBody(siteID int) (*GetSiteResponseBody, error) {
	body := &GetSiteResponseBody{}

	if siteID < 1 {
		return body, fmt.Errorf("invalid site id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/sites/%d", siteID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, func(resp *connection.APIResponse) error {
		if response.StatusCode == 404 {
			return &SiteNotFoundError{ID: siteID}
		}

		return nil
	})
}
