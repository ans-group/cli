package safedns

import "github.com/ukfast/sdk-go/pkg/connection"

// GetZonesResponseBody represents the API response body from the GetZones resource
type GetZonesResponseBody struct {
	connection.APIResponseBody

	Data []Zone `json:"data"`
}

// GetZoneResponseBody represents the API response body from the GetZone resource
type GetZoneResponseBody struct {
	connection.APIResponseBody

	Data Zone `json:"data"`
}

// GetRecordsResponseBody represents the API response body from the GetRecords resource
type GetRecordsResponseBody struct {
	connection.APIResponseBody

	Data []Record `json:"data"`
}

// GetRecordResponseBody represents the API response body from the GetRecord resource
type GetRecordResponseBody struct {
	connection.APIResponseBody

	Data Record `json:"data"`
}

// GetZoneNotesResponseBody represents the API response body from the GetZoneNotes resource
type GetZoneNotesResponseBody struct {
	connection.APIResponseBody

	Data []Note `json:"data"`
}

// GetZoneNoteResponseBody represents the API response body from the GetZoneNote resource
type GetZoneNoteResponseBody struct {
	connection.APIResponseBody

	Data Note `json:"data"`
}

// GetTemplatesResponseBody represents the API response body from the GetTemplates resource
type GetTemplatesResponseBody struct {
	connection.APIResponseBody

	Data []Template `json:"data"`
}

// GetTemplateResponseBody represents the API response body from the GetTemplate resource
type GetTemplateResponseBody struct {
	connection.APIResponseBody

	Data Template `json:"data"`
}
