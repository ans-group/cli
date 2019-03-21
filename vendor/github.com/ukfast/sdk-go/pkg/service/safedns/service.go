package safedns

import (
	"github.com/ukfast/sdk-go/pkg/connection"
)

// SafeDNSService is an interface for managing SafeDNS
type SafeDNSService interface {
	GetZones(parameters connection.APIRequestParameters) ([]Zone, error)
	GetZonesPaginated(parameters connection.APIRequestParameters) ([]Zone, error)
	GetZone(zoneName string) (Zone, error)
	CreateZone(req CreateZoneRequest) error
	DeleteZone(zoneName string) error
	GetZoneRecords(zoneName string, parameters connection.APIRequestParameters) ([]Record, error)
	GetZoneRecordsPaginated(zoneName string, parameters connection.APIRequestParameters) ([]Record, error)
	GetZoneRecord(zoneName string, recordID int) (Record, error)
	CreateZoneRecord(zoneName string, req CreateRecordRequest) (int, error)
	UpdateZoneRecord(zoneName string, record Record) (int, error)
	PatchZoneRecord(zoneName string, recordID int, patch PatchRecordRequest) (int, error)
	DeleteZoneRecord(zoneName string, recordID int) error
	GetZoneNotes(zoneName string, parameters connection.APIRequestParameters) ([]Note, error)
	GetZoneNotesPaginated(zoneName string, parameters connection.APIRequestParameters) ([]Note, error)
	GetZoneNote(zoneName string, noteID int) (Note, error)
	CreateZoneNote(zoneName string, req CreateNoteRequest) (int, error)
	GetTemplates(parameters connection.APIRequestParameters) ([]Template, error)
	GetTemplatesPaginated(parameters connection.APIRequestParameters) ([]Template, error)
	GetTemplate(templateID int) (Template, error)
	CreateTemplate(req CreateTemplateRequest) (int, error)
	UpdateTemplate(template Template) (int, error)
	PatchTemplate(templateID int, patch PatchTemplateRequest) (int, error)
	DeleteTemplate(templateID int) error
	GetTemplateRecords(templateID int, parameters connection.APIRequestParameters) ([]Record, error)
	GetTemplateRecordsPaginated(templateID int, parameters connection.APIRequestParameters) ([]Record, error)
	GetTemplateRecord(templateID int, recordID int) (Record, error)
	CreateTemplateRecord(templateID int, req CreateRecordRequest) (int, error)
	UpdateTemplateRecord(templateID int, record Record) (int, error)
	PatchTemplateRecord(templateID int, recordID int, patch PatchRecordRequest) (int, error)
	DeleteTemplateRecord(templateID int, recordID int) error
}

// Service implements SafeDNSService for managing
// SafeDNS via the UKFast API
type Service struct {
	connection connection.Connection
}

// NewService returns a new instance of SafeDNSService
func NewService(connection connection.Connection) *Service {
	return &Service{
		connection: connection,
	}
}
