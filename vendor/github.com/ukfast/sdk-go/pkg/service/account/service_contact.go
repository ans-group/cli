package account

import (
	"fmt"

	"github.com/ukfast/sdk-go/pkg/connection"
)

// GetContacts retrieves a list of contacts
func (s *Service) GetContacts(parameters connection.APIRequestParameters) ([]Contact, error) {
	r := connection.RequestAll{}

	var contacts []Contact
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getContactsPaginatedResponseBody(parameters)
		if err != nil {
			return nil, err
		}

		for _, contact := range response.Data {
			contacts = append(contacts, contact)
		}

		return response, nil
	}

	err := r.Invoke(parameters)

	return contacts, err
}

// GetContactsPaginated retrieves a paginated list of contacts
func (s *Service) GetContactsPaginated(parameters connection.APIRequestParameters) ([]Contact, error) {
	body, err := s.getContactsPaginatedResponseBody(parameters)

	return body.Data, err
}

func (s *Service) getContactsPaginatedResponseBody(parameters connection.APIRequestParameters) (*GetContactsResponseBody, error) {
	body := &GetContactsResponseBody{}

	response, err := s.connection.Get("/account/v1/contacts", parameters)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, nil)
}

// GetContact retrieves a single contact by id
func (s *Service) GetContact(contactID int) (Contact, error) {
	body, err := s.getContactResponseBody(contactID)

	return body.Data, err
}

func (s *Service) getContactResponseBody(contactID int) (*GetContactResponseBody, error) {
	body := &GetContactResponseBody{}

	if contactID < 1 {
		return body, fmt.Errorf("invalid contact id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/account/v1/contacts/%d", contactID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, func(resp *connection.APIResponse) error {
		if response.StatusCode == 404 {
			return &ContactNotFoundError{ID: contactID}
		}

		return nil
	})
}
