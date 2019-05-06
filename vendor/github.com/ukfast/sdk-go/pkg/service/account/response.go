package account

import "github.com/ukfast/sdk-go/pkg/connection"

// GetContactsResponseBody represents the API response body from the GetContacts resource
type GetContactsResponseBody struct {
	connection.APIResponseBody

	Data []Contact `json:"data"`
}

// GetContactResponseBody represents the API response body from the GetContact resource
type GetContactResponseBody struct {
	connection.APIResponseBody

	Data Contact `json:"data"`
}

// GetDetailsResponseBody represents the API response body from the GetDetails resource
type GetDetailsResponseBody struct {
	connection.APIResponseBody

	Data Details `json:"data"`
}

// GetCreditsResponseBody represents the API response body from the GetCredits resource
type GetCreditsResponseBody struct {
	connection.APIResponseBody

	Data []Credit `json:"data"`
}
