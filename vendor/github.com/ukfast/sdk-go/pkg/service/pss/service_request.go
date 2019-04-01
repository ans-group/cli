package pss

import (
	"fmt"

	"github.com/ukfast/sdk-go/pkg/connection"
)

// GetRequests retrieves a list of requests
func (s *Service) GetRequests(parameters connection.APIRequestParameters) ([]Request, error) {
	r := connection.RequestAll{}

	var requests []Request
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getRequestsPaginatedResponseBody(parameters)
		if err != nil {
			return nil, err
		}

		for _, request := range response.Data {
			requests = append(requests, request)
		}

		return response, nil
	}

	err := r.Invoke(parameters)

	return requests, err
}

// GetRequestsPaginated retrieves a paginated list of requests
func (s *Service) GetRequestsPaginated(parameters connection.APIRequestParameters) ([]Request, error) {
	body, err := s.getRequestsPaginatedResponseBody(parameters)

	return body.Data, err
}

func (s *Service) getRequestsPaginatedResponseBody(parameters connection.APIRequestParameters) (*GetRequestsResponseBody, error) {
	body := &GetRequestsResponseBody{}

	response, err := s.connection.Get("/pss/v1/requests", parameters)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse([]int{200}, body)
}

// GetRequest retrieves a single request by id
func (s *Service) GetRequest(requestID int) (Request, error) {
	body, err := s.getRequestResponseBody(requestID)

	return body.Data, err
}

func (s *Service) getRequestResponseBody(requestID int) (*GetRequestResponseBody, error) {
	body := &GetRequestResponseBody{}

	if requestID < 1 {
		return body, fmt.Errorf("invalid request id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/pss/v1/requests/%d", requestID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &RequestNotFoundError{ID: requestID}
	}

	return body, response.HandleResponse([]int{200}, body)
}

// GetRequestConversation retrieves the conversation for a request
func (s *Service) GetRequestConversation(requestID int) ([]Reply, error) {
	body, err := s.getRequestConversationResponseBody(requestID)

	return body.Data, err
}

func (s *Service) getRequestConversationResponseBody(requestID int) (*GetRepliesResponseBody, error) {
	body := &GetRepliesResponseBody{}

	if requestID < 1 {
		return body, fmt.Errorf("invalid request id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/pss/v1/requests/%d/conversation", requestID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &RequestNotFoundError{ID: requestID}
	}

	return body, response.HandleResponse([]int{200}, body)
}
