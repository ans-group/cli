package ddosx

import (
	"github.com/ukfast/sdk-go/pkg/connection"
)

// GetRecords retrieves a list of records
func (s *Service) GetRecords(parameters connection.APIRequestParameters) ([]Record, error) {
	r := connection.RequestAll{}

	var records []Record
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getRecordsPaginatedResponseBody(parameters)
		if err != nil {
			return nil, err
		}

		for _, record := range response.Data {
			records = append(records, record)
		}

		return response, nil
	}

	err := r.Invoke(parameters)

	return records, err
}

// GetRecordsPaginated retrieves a paginated list of records
func (s *Service) GetRecordsPaginated(parameters connection.APIRequestParameters) ([]Record, error) {
	body, err := s.getRecordsPaginatedResponseBody(parameters)

	return body.Data, err
}

func (s *Service) getRecordsPaginatedResponseBody(parameters connection.APIRequestParameters) (*GetRecordsResponseBody, error) {
	body := &GetRecordsResponseBody{}

	response, err := s.connection.Get("/ddosx/v1/records", parameters)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse(body, nil)
}
