package safedns

import (
	"fmt"

	"github.com/ukfast/sdk-go/pkg/connection"
)

// GetZones retrieves a list of zones
func (s *Service) GetZones(parameters connection.APIRequestParameters) ([]Zone, error) {
	r := connection.RequestAll{}

	var zones []Zone
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getZonesPaginatedResponseBody(parameters)
		if err != nil {
			return nil, err
		}

		for _, zone := range response.Data {
			zones = append(zones, zone)
		}

		return response, nil
	}

	err := r.Invoke(parameters)

	return zones, err
}

// GetZonesPaginated retrieves a paginated list of zones
func (s *Service) GetZonesPaginated(parameters connection.APIRequestParameters) ([]Zone, error) {
	body, err := s.getZonesPaginatedResponseBody(parameters)

	return body.Data, err
}

func (s *Service) getZonesPaginatedResponseBody(parameters connection.APIRequestParameters) (*GetZonesResponseBody, error) {
	body := &GetZonesResponseBody{}

	response, err := s.connection.Get("/safedns/v1/zones", parameters)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse([]int{}, body)
}

// GetZone retrieves a single zone by name
func (s *Service) GetZone(zoneName string) (Zone, error) {
	body, err := s.getZoneResponseBody(zoneName)

	return body.Data, err
}

func (s *Service) getZoneResponseBody(zoneName string) (*GetZoneResponseBody, error) {
	body := &GetZoneResponseBody{}

	if zoneName == "" {
		return body, fmt.Errorf("invalid zone name")
	}

	response, err := s.connection.Get(fmt.Sprintf("/safedns/v1/zones/%s", zoneName), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &ZoneNotFoundError{ZoneName: zoneName}
	}

	return body, response.HandleResponse([]int{}, body)
}

// CreateZone creates a new SafeDNS zone
func (s *Service) CreateZone(req CreateZoneRequest) error {
	_, err := s.createZoneResponseBody(req)

	return err
}

func (s *Service) createZoneResponseBody(req CreateZoneRequest) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	response, err := s.connection.Post("/safedns/v1/zones", &req)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse([]int{}, body)
}

// PatchZone patches a SafeDNS zone
func (s *Service) PatchZone(zoneName string, req PatchZoneRequest) error {
	_, err := s.patchZoneResponseBody(zoneName, req)

	return err
}

func (s *Service) patchZoneResponseBody(zoneName string, req PatchZoneRequest) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if zoneName == "" {
		return body, fmt.Errorf("invalid zone name")
	}

	response, err := s.connection.Patch(fmt.Sprintf("/safedns/v1/zones/%s", zoneName), &req)
	if err != nil {
		return body, err
	}

	return body, response.HandleResponse([]int{}, body)
}

// DeleteZone removes a SafeDNS zone
func (s *Service) DeleteZone(zoneName string) error {
	_, err := s.deleteZoneResponseBody(zoneName)

	return err
}

func (s *Service) deleteZoneResponseBody(zoneName string) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if zoneName == "" {
		return body, fmt.Errorf("invalid zone name")
	}

	response, err := s.connection.Delete(fmt.Sprintf("/safedns/v1/zones/%s", zoneName), nil)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &ZoneNotFoundError{ZoneName: zoneName}
	}

	return body, response.HandleResponse([]int{}, body)
}

// GetZoneRecords retrieves a list of records
func (s *Service) GetZoneRecords(zoneName string, parameters connection.APIRequestParameters) ([]Record, error) {
	r := connection.RequestAll{}

	var records []Record
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getZoneRecordsPaginatedResponseBody(zoneName, parameters)
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

// GetZoneRecordsPaginated retrieves a paginated list of zones
func (s *Service) GetZoneRecordsPaginated(zoneName string, parameters connection.APIRequestParameters) ([]Record, error) {
	body, err := s.getZoneRecordsPaginatedResponseBody(zoneName, parameters)

	return body.Data, err
}

func (s *Service) getZoneRecordsPaginatedResponseBody(zoneName string, parameters connection.APIRequestParameters) (*GetRecordsResponseBody, error) {
	body := &GetRecordsResponseBody{}

	if zoneName == "" {
		return body, fmt.Errorf("invalid zone name")
	}

	response, err := s.connection.Get(fmt.Sprintf("/safedns/v1/zones/%s/records", zoneName), parameters)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &ZoneNotFoundError{ZoneName: zoneName}
	}

	return body, response.HandleResponse([]int{}, body)
}

// GetZoneRecord retrieves a single zone record by ID
func (s *Service) GetZoneRecord(zoneName string, recordID int) (Record, error) {
	body, err := s.getZoneRecordResponseBody(zoneName, recordID)

	return body.Data, err
}

func (s *Service) getZoneRecordResponseBody(zoneName string, recordID int) (*GetRecordResponseBody, error) {
	body := &GetRecordResponseBody{}

	if zoneName == "" {
		return body, fmt.Errorf("invalid zone name")
	}
	if recordID < 1 {
		return body, fmt.Errorf("invalid record id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/safedns/v1/zones/%s/records/%d", zoneName, recordID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &ZoneRecordNotFoundError{ZoneName: zoneName, RecordID: recordID}
	}

	return body, response.HandleResponse([]int{}, body)
}

// CreateZoneRecord creates a new SafeDNS zone record
func (s *Service) CreateZoneRecord(zoneName string, req CreateRecordRequest) (int, error) {
	body, err := s.createZoneRecordResponseBody(zoneName, req)

	return body.Data.ID, err
}

func (s *Service) createZoneRecordResponseBody(zoneName string, req CreateRecordRequest) (*GetRecordResponseBody, error) {
	body := &GetRecordResponseBody{}

	if zoneName == "" {
		return body, fmt.Errorf("invalid zone name")
	}

	response, err := s.connection.Post(fmt.Sprintf("/safedns/v1/zones/%s/records", zoneName), &req)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &ZoneNotFoundError{ZoneName: zoneName}
	}

	return body, response.HandleResponse([]int{}, body)
}

// UpdateZoneRecord updates a SafeDNS zone record
func (s *Service) UpdateZoneRecord(zoneName string, record Record) (int, error) {
	body, err := s.updateZoneRecordResponseBody(zoneName, record)

	return body.Data.ID, err
}

func (s *Service) updateZoneRecordResponseBody(zoneName string, record Record) (*GetRecordResponseBody, error) {
	body := &GetRecordResponseBody{}

	if zoneName == "" {
		return body, fmt.Errorf("invalid zone name")
	}
	if record.ID < 1 {
		return body, fmt.Errorf("invalid record id")
	}

	response, err := s.connection.Put(fmt.Sprintf("/safedns/v1/zones/%s/records/%d", zoneName, record.ID), &record)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &ZoneRecordNotFoundError{ZoneName: zoneName, RecordID: record.ID}
	}

	return body, response.HandleResponse([]int{}, body)
}

// PatchZoneRecord patches a SafeDNS zone record
func (s *Service) PatchZoneRecord(zoneName string, recordID int, patch PatchRecordRequest) (int, error) {
	body, err := s.patchZoneRecordResponseBody(zoneName, recordID, patch)

	return body.Data.ID, err
}

func (s *Service) patchZoneRecordResponseBody(zoneName string, recordID int, patch PatchRecordRequest) (*GetRecordResponseBody, error) {
	body := &GetRecordResponseBody{}

	if zoneName == "" {
		return body, fmt.Errorf("invalid zone name")
	}
	if recordID < 1 {
		return body, fmt.Errorf("invalid record id")
	}

	response, err := s.connection.Put(fmt.Sprintf("/safedns/v1/zones/%s/records/%d", zoneName, recordID), &patch)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &ZoneRecordNotFoundError{ZoneName: zoneName, RecordID: recordID}
	}

	return body, response.HandleResponse([]int{}, body)
}

// DeleteZoneRecord removes a SafeDNS zone record
func (s *Service) DeleteZoneRecord(zoneName string, recordID int) error {
	_, err := s.deleteZoneRecordResponseBody(zoneName, recordID)

	return err
}

func (s *Service) deleteZoneRecordResponseBody(zoneName string, recordID int) (*connection.APIResponseBody, error) {
	body := &connection.APIResponseBody{}

	if zoneName == "" {
		return body, fmt.Errorf("invalid zone name")
	}
	if recordID < 1 {
		return body, fmt.Errorf("invalid record id")
	}

	response, err := s.connection.Delete(fmt.Sprintf("/safedns/v1/zones/%s/records/%d", zoneName, recordID), nil)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &ZoneRecordNotFoundError{ZoneName: zoneName, RecordID: recordID}
	}

	return body, response.HandleResponse([]int{}, body)
}

// GetZoneNotes retrieves a list of zone notes
func (s *Service) GetZoneNotes(zoneName string, parameters connection.APIRequestParameters) ([]Note, error) {
	r := connection.RequestAll{}

	var notes []Note
	r.GetNext = func(parameters connection.APIRequestParameters) (connection.ResponseBody, error) {
		response, err := s.getZoneNotesPaginatedResponseBody(zoneName, parameters)
		if err != nil {
			return nil, err
		}

		for _, note := range response.Data {
			notes = append(notes, note)
		}

		return response, nil
	}

	err := r.Invoke(parameters)

	return notes, err
}

// GetZoneNotesPaginated retrieves a paginated list of zone notes
func (s *Service) GetZoneNotesPaginated(zoneName string, parameters connection.APIRequestParameters) ([]Note, error) {
	body, err := s.getZoneNotesPaginatedResponseBody(zoneName, parameters)

	return body.Data, err
}

func (s *Service) getZoneNotesPaginatedResponseBody(zoneName string, parameters connection.APIRequestParameters) (*GetZoneNotesResponseBody, error) {
	body := &GetZoneNotesResponseBody{}

	if zoneName == "" {
		return body, fmt.Errorf("invalid zone name")
	}

	response, err := s.connection.Get(fmt.Sprintf("/safedns/v1/zones/%s/notes", zoneName), parameters)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &ZoneNotFoundError{ZoneName: zoneName}
	}

	return body, response.HandleResponse([]int{}, body)
}

// GetZoneNote retrieves a single zone note by ID
func (s *Service) GetZoneNote(zoneName string, noteID int) (Note, error) {
	body, err := s.getZoneNoteResponseBody(zoneName, noteID)

	return body.Data, err
}

func (s *Service) getZoneNoteResponseBody(zoneName string, noteID int) (*GetZoneNoteResponseBody, error) {
	body := &GetZoneNoteResponseBody{}

	if zoneName == "" {
		return body, fmt.Errorf("invalid zone name")
	}
	if noteID < 1 {
		return body, fmt.Errorf("invalid note id")
	}

	response, err := s.connection.Get(fmt.Sprintf("/safedns/v1/zones/%s/notes/%d", zoneName, noteID), connection.APIRequestParameters{})
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &ZoneNoteNotFoundError{ZoneName: zoneName, NoteID: noteID}
	}

	return body, response.HandleResponse([]int{}, body)
}

// CreateZoneNote creates a new SafeDNS zone note
func (s *Service) CreateZoneNote(zoneName string, req CreateNoteRequest) (int, error) {
	body, err := s.createZoneNote(zoneName, req)

	return body.Data.ID, err
}

func (s *Service) createZoneNote(zoneName string, req CreateNoteRequest) (*GetZoneNoteResponseBody, error) {
	body := &GetZoneNoteResponseBody{}

	if zoneName == "" {
		return body, fmt.Errorf("invalid zone name")
	}

	response, err := s.connection.Post(fmt.Sprintf("/safedns/v1/zones/%s/notes", zoneName), &req)
	if err != nil {
		return body, err
	}

	if response.StatusCode == 404 {
		return body, &ZoneNotFoundError{ZoneName: zoneName}
	}

	return body, response.HandleResponse([]int{}, body)
}
