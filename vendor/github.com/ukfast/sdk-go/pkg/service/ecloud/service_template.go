package ecloud

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

	response, err := s.connection.Get("/ecloud/v1/templates", parameters)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse([]int{200}, body)
}

// GetTemplate retrieves a single template by name
func (s *Service) GetTemplate(templateName string) (Template, error) {
	body, err := s.getTemplateResponseBody(templateName)

	return body.Data, err
}

func (s *Service) getTemplateResponseBody(templateName string) (*GetTemplateResponseBody, error) {
	body := &GetTemplateResponseBody{}

	if templateName == "" {
		return body, fmt.Errorf("invalid template name")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/templates/%s", templateName), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TemplateNotFoundError{Name: templateName}
	}

	return body, response.HandleResponse([]int{200}, body)
}

// RenameTemplate renames an eCloud template
func (s *Service) RenameTemplate(templateName string, req RenameTemplateRequest) error {
	_, err := s.renameTemplateResponseBody(templateName, req)

	return err
}

func (s *Service) renameTemplateResponseBody(templateName string, req RenameTemplateRequest) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if templateName == "" {
		return body, fmt.Errorf("invalid template name")
	}

	response, err := s.connection.Post(fmt.Sprintf("/ecloud/v1/templates/%s/move", templateName), &req)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TemplateNotFoundError{Name: templateName}
	}

	return body, response.HandleResponse([]int{202}, body)
}

// DeleteTemplate removes a template
func (s *Service) DeleteTemplate(templateName string) error {
	_, err := s.deleteTemplateResponseBody(templateName)

	return err
}

func (s *Service) deleteTemplateResponseBody(templateName string) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if templateName == "" {
		return body, fmt.Errorf("invalid template name")
	}

	response, err := s.connection.Delete(fmt.Sprintf("/ecloud/v1/templates/%s", templateName), nil)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TemplateNotFoundError{Name: templateName}
	}

	return body, response.HandleResponse([]int{202}, body)
}
