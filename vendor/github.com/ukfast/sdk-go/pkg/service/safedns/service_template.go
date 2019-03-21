package safedns

import (
	"fmt"

	"github.com/ukfast/sdk-go/pkg/connection"
)

// GetTemplates retrieves a list of templates
func (s *Service) GetTemplates(parameters connection.APIRequestParameters) ([]Template, error) {
	r := connection.RequestAll{}

	var templates []Template
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getTemplatesPaginatedResponseBody(parameters)
		if err != nil {
			return nil, err
		}

		for _, template := range response.Data {
			templates = append(templates, template)
		}

		return response, nil
	}

	err := r.Invoke(parameters)

	return templates, err
}

// GetTemplatesPaginated retrieves a paginated list of templates
func (s *Service) GetTemplatesPaginated(parameters connection.APIRequestParameters) ([]Template, error) {
	body, err := s.getTemplatesPaginatedResponseBody(parameters)

	return body.Data, err
}

func (s *Service) getTemplatesPaginatedResponseBody(parameters connection.APIRequestParameters) (*GetTemplatesResponseBody, error) {
	body := &GetTemplatesResponseBody{}

	response, err := s.connection.Get("/safedns/v1/templates", parameters)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse([]int{200}, body)
}

// GetTemplate retrieves a single template by ID
func (s *Service) GetTemplate(templateID int) (Template, error) {
	body, err := s.getTemplateResponseBody(templateID)

	return body.Data, err
}

func (s *Service) getTemplateResponseBody(templateID int) (*GetTemplateResponseBody, error) {
	body := &GetTemplateResponseBody{}

	if templateID < 1 {
		return body, fmt.Errorf("invalid template id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/safedns/v1/templates/%d", templateID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TemplateNotFoundError{TemplateID: templateID}
	}

	return body, response.HandleResponse([]int{200}, body)
}

// CreateTemplate creates a new SafeDNS template
func (s *Service) CreateTemplate(req CreateTemplateRequest) (int, error) {
	body, err := s.createTemplateResponseBody(req)

	return body.Data.ID, err
}

func (s *Service) createTemplateResponseBody(req CreateTemplateRequest) (*GetTemplateResponseBody, error) {
	body := &GetTemplateResponseBody{}

	response, err := s.connection.Post("/safedns/v1/templates", &req)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse([]int{201}, body)
}

// UpdateTemplate updates a SafeDNS template
func (s *Service) UpdateTemplate(template Template) (int, error) {
	body, err := s.updateTemplateResponseBody(template)

	return body.Data.ID, err
}

func (s *Service) updateTemplateResponseBody(template Template) (*GetTemplateResponseBody, error) {
	body := &GetTemplateResponseBody{}

	if template.ID < 1 {
		return body, fmt.Errorf("invalid template id")
	}

	response, err := s.connection.Put(fmt.Sprintf("/safedns/v1/templates/%d", template.ID), &template)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TemplateNotFoundError{TemplateID: template.ID}
	}

	return body, response.HandleResponse([]int{200}, body)
}

// PatchTemplate patches a SafeDNS template
func (s *Service) PatchTemplate(templateID int, patch PatchTemplateRequest) (int, error) {
	body, err := s.patchTemplateResponseBody(templateID, patch)

	return body.Data.ID, err
}

func (s *Service) patchTemplateResponseBody(templateID int, patch PatchTemplateRequest) (*GetTemplateResponseBody, error) {
	body := &GetTemplateResponseBody{}

	if templateID < 1 {
		return body, fmt.Errorf("invalid template id")
	}

	// Currently uses PUT
	response, err := s.connection.Put(fmt.Sprintf("/safedns/v1/templates/%d", templateID), &patch)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TemplateNotFoundError{TemplateID: templateID}
	}

	return body, response.HandleResponse([]int{200}, body)
}

// DeleteTemplate removes a SafeDNS template
func (s *Service) DeleteTemplate(templateID int) error {
	_, err := s.deleteTemplateResponseBody(templateID)

	return err
}

func (s *Service) deleteTemplateResponseBody(templateID int) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if templateID < 1 {
		return body, fmt.Errorf("invalid template id")
	}

	response, err := s.connection.Delete(fmt.Sprintf("/safedns/v1/templates/%d", templateID), nil)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TemplateNotFoundError{TemplateID: templateID}
	}

	return body, response.HandleResponse([]int{204}, body)
}

// GetTemplateRecords retrieves a list of records
func (s *Service) GetTemplateRecords(templateID int, parameters connection.APIRequestParameters) ([]Record, error) {
	r := connection.RequestAll{}

	var records []Record
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getTemplateRecordsPaginatedResponseBody(templateID, parameters)
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

// GetTemplateRecordsPaginated retrieves a paginated list of template records
func (s *Service) GetTemplateRecordsPaginated(templateID int, parameters connection.APIRequestParameters) ([]Record, error) {
	body, err := s.getTemplateRecordsPaginatedResponseBody(templateID, parameters)

	return body.Data, err
}

func (s *Service) getTemplateRecordsPaginatedResponseBody(templateID int, parameters connection.APIRequestParameters) (*GetRecordsResponseBody, error) {
	body := &GetRecordsResponseBody{}

	if templateID < 1 {
		return body, fmt.Errorf("invalid template id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/safedns/v1/templates/%d/records", templateID), parameters)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TemplateNotFoundError{TemplateID: templateID}
	}

	return body, response.HandleResponse([]int{200}, body)
}

// GetTemplateRecord retrieves a single zone record by ID
func (s *Service) GetTemplateRecord(templateID int, recordID int) (Record, error) {
	body, err := s.getTemplateRecordResponseBody(templateID, recordID)

	return body.Data, err
}

func (s *Service) getTemplateRecordResponseBody(templateID int, recordID int) (*GetRecordResponseBody, error) {
	body := &GetRecordResponseBody{}

	if templateID < 1 {
		return body, fmt.Errorf("invalid template id")
	}
	if recordID < 1 {
		return body, fmt.Errorf("invalid record id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/safedns/v1/templates/%d/records/%d", templateID, recordID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TemplateRecordNotFoundError{TemplateID: templateID, RecordID: recordID}
	}

	return body, response.HandleResponse([]int{200}, body)
}

// CreateTemplateRecord creates a new SafeDNS zone record
func (s *Service) CreateTemplateRecord(templateID int, req CreateRecordRequest) (int, error) {
	body, err := s.createTemplateRecordResponseBody(templateID, req)

	return body.Data.ID, err
}

func (s *Service) createTemplateRecordResponseBody(templateID int, req CreateRecordRequest) (*GetTemplateResponseBody, error) {
	body := &GetTemplateResponseBody{}

	if templateID < 1 {
		return body, fmt.Errorf("invalid template id")
	}

	response, err := s.connection.Post(fmt.Sprintf("/safedns/v1/templates/%d/records", templateID), &req)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TemplateNotFoundError{TemplateID: templateID}
	}

	return body, response.HandleResponse([]int{201}, body)
}

// UpdateTemplateRecord updates a SafeDNS template record
func (s *Service) UpdateTemplateRecord(templateID int, record Record) (int, error) {
	body, err := s.updateTemplateRecordResponseBody(templateID, record)

	return body.Data.ID, err
}

func (s *Service) updateTemplateRecordResponseBody(templateID int, record Record) (*GetTemplateResponseBody, error) {
	body := &GetTemplateResponseBody{}

	if templateID < 1 {
		return body, fmt.Errorf("invalid template id")
	}
	if record.ID < 1 {
		return body, fmt.Errorf("invalid record id")
	}

	response, err := s.connection.Put(fmt.Sprintf("/safedns/v1/templates/%d/records/%d", templateID, record.ID), &record)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TemplateRecordNotFoundError{TemplateID: templateID, RecordID: record.ID}
	}

	return body, response.HandleResponse([]int{200}, body)
}

// PatchTemplateRecord patches a SafeDNS template record
func (s *Service) PatchTemplateRecord(templateID int, recordID int, patch PatchRecordRequest) (int, error) {
	body, err := s.patchTemplateRecordResponseBody(templateID, recordID, patch)

	return body.Data.ID, err
}

func (s *Service) patchTemplateRecordResponseBody(templateID int, recordID int, patch PatchRecordRequest) (*GetTemplateResponseBody, error) {
	body := &GetTemplateResponseBody{}

	if templateID < 1 {
		return body, fmt.Errorf("invalid template id")
	}
	if recordID < 1 {
		return body, fmt.Errorf("invalid record id")
	}

	// Currently uses PUT
	response, err := s.connection.Put(fmt.Sprintf("/safedns/v1/templates/%d/records/%d", templateID, recordID), &patch)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TemplateRecordNotFoundError{TemplateID: templateID, RecordID: recordID}
	}

	return body, response.HandleResponse([]int{200}, body)
}

// DeleteTemplateRecord removes a SafeDNS template record
func (s *Service) DeleteTemplateRecord(templateID int, recordID int) error {
	_, err := s.deleteTemplateRecordResponseBody(templateID, recordID)

	return err
}

func (s *Service) deleteTemplateRecordResponseBody(templateID int, recordID int) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if templateID < 1 {
		return body, fmt.Errorf("invalid template id")
	}
	if recordID < 1 {
		return body, fmt.Errorf("invalid record id")
	}

	response, err := s.connection.Delete(fmt.Sprintf("/safedns/v1/templates/%d/records/%d", templateID, recordID), nil)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TemplateRecordNotFoundError{TemplateID: templateID, RecordID: recordID}
	}

	return body, response.HandleResponse([]int{204}, body)
}
