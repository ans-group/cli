package ssl

import "github.com/ukfast/sdk-go/pkg/connection"

// GetCertificatesResponseBody represents the API response body from the GetCertificates resource
type GetCertificatesResponseBody struct {
	connection.APIResponseBody

	Data []Certificate `json:"data"`
}

// GetCertificateResponseBody represents the API response body from the GetCertificate resource
type GetCertificateResponseBody struct {
	connection.APIResponseBody

	Data Certificate `json:"data"`
}

// GetCertificateContentResponseBody represents the API response body from the GetCertificateContent resource
type GetCertificateContentResponseBody struct {
	connection.APIResponseBody

	Data CertificateContent `json:"data"`
}

// GetCertificateKeyResponseBody represents the API response body from the GetCertificateKey resource
type GetCertificateKeyResponseBody struct {
	connection.APIResponseBody

	Data CertificatePrivateKey `json:"data"`
}
