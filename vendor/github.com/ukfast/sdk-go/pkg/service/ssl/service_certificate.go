package ssl

import (
	"fmt"

	"github.com/ukfast/sdk-go/pkg/connection"
)

// GetCertificates retrieves a list of certificates
func (s *Service) GetCertificates(parameters connection.APIRequestParameters) ([]Certificate, error) {
	r := connection.RequestAll{}

	var certificates []Certificate
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getCertificatesPaginatedResponseBody(parameters)
		if err != nil {
			return nil, err
		}

		for _, certificate := range response.Data {
			certificates = append(certificates, certificate)
		}

		return response, nil
	}

	err := r.Invoke(parameters)

	return certificates, err
}

// GetCertificatesPaginated retrieves a paginated list of certificates
func (s *Service) GetCertificatesPaginated(parameters connection.APIRequestParameters) ([]Certificate, error) {
	body, err := s.getCertificatesPaginatedResponseBody(parameters)

	return body.Data, err
}

func (s *Service) getCertificatesPaginatedResponseBody(parameters connection.APIRequestParameters) (*GetCertificatesResponseBody, error) {
	body := &GetCertificatesResponseBody{}

	response, err := s.connection.Get("/ssl/v1/certificates", parameters)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse([]int{}, body)
}

// GetCertificate retrieves a single certificate by id
func (s *Service) GetCertificate(certificateID int) (Certificate, error) {
	body, err := s.getCertificateResponseBody(certificateID)

	return body.Data, err
}

func (s *Service) getCertificateResponseBody(certificateID int) (*GetCertificateResponseBody, error) {
	body := &GetCertificateResponseBody{}

	if certificateID < 1 {
		return body, fmt.Errorf("invalid certificate id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ssl/v1/certificates/%d", certificateID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &CertificateNotFoundError{ID: certificateID}
	}

	return body, response.HandleResponse([]int{}, body)
}

// GetCertificateContent retrieves the content of an SSL certificate
func (s *Service) GetCertificateContent(certificateID int) (CertificateContent, error) {
	body, err := s.getCertificateContentResponseBody(certificateID)

	return body.Data, err
}

func (s *Service) getCertificateContentResponseBody(certificateID int) (*GetCertificateContentResponseBody, error) {
	body := &GetCertificateContentResponseBody{}

	if certificateID < 1 {
		return body, fmt.Errorf("invalid certificate id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ssl/v1/certificates/%d/download", certificateID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &CertificateNotFoundError{ID: certificateID}
	}

	return body, response.HandleResponse([]int{}, body)
}

// GetCertificatePrivateKey retrieves an SSL certificate private key
func (s *Service) GetCertificatePrivateKey(certificateID int) (CertificatePrivateKey, error) {
	body, err := s.getCertificatePrivateKeyResponseBody(certificateID)

	return body.Data, err
}

func (s *Service) getCertificatePrivateKeyResponseBody(certificateID int) (*GetCertificateKeyResponseBody, error) {
	body := &GetCertificateKeyResponseBody{}

	if certificateID < 1 {
		return body, fmt.Errorf("invalid certificate id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ssl/v1/certificates/%d/private-key", certificateID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &CertificateNotFoundError{ID: certificateID}
	}

	return body, response.HandleResponse([]int{}, body)
}
