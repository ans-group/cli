package registrar

import "github.com/ukfast/sdk-go/pkg/connection"

// GetDomainsResponseBody represents the API response body from the GetDomains resource
type GetDomainsResponseBody struct {
	connection.APIResponseBody

	Data []Domain `json:"data"`
}

// GetDomainResponseBody represents the API response body from the GetDomain resource
type GetDomainResponseBody struct {
	connection.APIResponseBody

	Data Domain `json:"data"`
}

// GetNameserversResponseBody represents the API response body from the GetNameservers resource
type GetNameserversResponseBody struct {
	connection.APIResponseBody

	Data []Nameserver `json:"data"`
}

// GetWhoisResponseBody represents the API response body from the GetWhois resource
type GetWhoisResponseBody struct {
	connection.APIResponseBody

	Data Whois `json:"data"`
}

// GetWhoisRawResponseBody represents the API response body from the GetWhoisRaw resource
type GetWhoisRawResponseBody struct {
	connection.APIResponseBody

	Data string `json:"data"`
}
