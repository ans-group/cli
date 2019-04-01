package ecloud

import (
	"fmt"

	"github.com/ukfast/sdk-go/pkg/connection"
)

// GetPods retrieves a list of pods
func (s *Service) GetPods(parameters connection.APIRequestParameters) ([]Pod, error) {
	r := connection.RequestAll{}

	var pods []Pod
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getPodsPaginatedResponseBody(parameters)
		if err != nil {
			return nil, err
		}

		for _, pod := range response.Data {
			pods = append(pods, pod)
		}

		return response, nil
	}

	err := r.Invoke(parameters)

	return pods, err
}

// GetPodsPaginated retrieves a paginated list of pods
func (s *Service) GetPodsPaginated(parameters connection.APIRequestParameters) ([]Pod, error) {
	body, err := s.getPodsPaginatedResponseBody(parameters)

	return body.Data, err
}

func (s *Service) getPodsPaginatedResponseBody(parameters connection.APIRequestParameters) (*GetPodsResponseBody, error) {
	body := &GetPodsResponseBody{}

	response, err := s.connection.Get("/ecloud/v1/pods", parameters)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse([]int{200}, body)
}

// GetPod retrieves a single pod by ID
func (s *Service) GetPod(podID int) (Pod, error) {
	body, err := s.getPodResponseBody(podID)

	return body.Data, err
}

func (s *Service) getPodResponseBody(podID int) (*GetPodResponseBody, error) {
	body := &GetPodResponseBody{}

	if podID < 1 {
		return body, fmt.Errorf("invalid pod id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/pods/%d", podID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &PodNotFoundError{ID: podID}
	}

	return body, response.HandleResponse([]int{200}, body)
}

// GetPodTemplates retrieves a list of pod templates
func (s *Service) GetPodTemplates(podID int, parameters connection.APIRequestParameters) ([]Template, error) {
	r := connection.RequestAll{}

	var templates []Template
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getPodTemplatesPaginatedResponseBody(podID, parameters)
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

// GetPodTemplatesPaginated retrieves a paginated list of pod templates
func (s *Service) GetPodTemplatesPaginated(podID int, parameters connection.APIRequestParameters) ([]Template, error) {
	body, err := s.getPodTemplatesPaginatedResponseBody(podID, parameters)

	return body.Data, err
}

func (s *Service) getPodTemplatesPaginatedResponseBody(podID int, parameters connection.APIRequestParameters) (*GetTemplatesResponseBody, error) {
	body := &GetTemplatesResponseBody{}

	if podID < 1 {
		return body, fmt.Errorf("invalid pod id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/pods/%d/templates", podID), parameters)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse([]int{200}, body)
}

// GetPodTemplate retrieves a single pod template by name
func (s *Service) GetPodTemplate(podID int, templateName string) (Template, error) {
	body, err := s.getPodTemplateResponseBody(podID, templateName)

	return body.Data, err
}

func (s *Service) getPodTemplateResponseBody(podID int, templateName string) (*GetTemplateResponseBody, error) {
	body := &GetTemplateResponseBody{}

	if podID < 1 {
		return body, fmt.Errorf("invalid pod id")
	}
	if templateName == "" {
		return body, fmt.Errorf("invalid template name")
	}

	response, err := s.connection.Get(fmt.Sprintf("/ecloud/v1/pods/%d/templates/%s", podID, templateName), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TemplateNotFoundError{Name: templateName}
	}

	return body, response.HandleResponse([]int{200}, body)
}

// RenamePodTemplate renames a pod template
func (s *Service) RenamePodTemplate(podID int, templateName string, req RenameTemplateRequest) error {
	_, err := s.renamePodTemplateResponseBody(podID, templateName, req)

	return err
}

func (s *Service) renamePodTemplateResponseBody(podID int, templateName string, req RenameTemplateRequest) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if podID < 1 {
		return body, fmt.Errorf("invalid pod id")
	}
	if templateName == "" {
		return body, fmt.Errorf("invalid template name")
	}

	response, err := s.connection.Post(fmt.Sprintf("/ecloud/v1/pods/%d/templates/%s/move", podID, templateName), &req)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TemplateNotFoundError{Name: templateName}
	}

	return body, response.HandleResponse([]int{202}, body)
}

// DeletePodTemplate removes a pod template
func (s *Service) DeletePodTemplate(podID int, templateName string) error {
	_, err := s.deletePodTemplateResponseBody(podID, templateName)

	return err
}

func (s *Service) deletePodTemplateResponseBody(podID int, templateName string) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if podID < 1 {
		return body, fmt.Errorf("invalid pod id")
	}
	if templateName == "" {
		return body, fmt.Errorf("invalid template name")
	}

	response, err := s.connection.Delete(fmt.Sprintf("/ecloud/v1/pods/%d/templates/%s", podID, templateName), nil)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &TemplateNotFoundError{Name: templateName}
	}

	return body, response.HandleResponse([]int{202}, body)
}
