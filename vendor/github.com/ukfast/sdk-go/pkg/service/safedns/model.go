//go:generate go run ../../gen/model_paginated_gen.go -package safedns -typename Record,Zone,Note,Template -destination model_paginated.go

package safedns

import (
	"strconv"
	"time"

	"github.com/ukfast/sdk-go/pkg/connection"
)

// RecordTTL represents the record TTL time in seconds
type RecordTTL int

// Time returns the record TTL time
func (r RecordTTL) Time() time.Time {
	return time.Now().Add(r.Duration())
}

// Duration returns the record TTL duration (seconds)
func (r RecordTTL) Duration() time.Duration {
	return (time.Second * time.Duration(int(r)))
}

func (r RecordTTL) String() string {
	return strconv.Itoa(int(r))
}

type RecordType string

func (s RecordType) String() string {
	return string(s)
}

const (
	RecordTypeA     RecordType = "A"
	RecordTypeAAAA  RecordType = "AAAA"
	RecordTypeCAA   RecordType = "CAA"
	RecordTypeCNAME RecordType = "CNAME"
	RecordTypeMX    RecordType = "MX"
	RecordTypeSPF   RecordType = "SPF"
	RecordTypeSRV   RecordType = "SRV"
	RecordTypeTXT   RecordType = "TXT"
	RecordTypeNS    RecordType = "NS"
	RecordTypeSOA   RecordType = "SOA"
	RecordTypeAXFR  RecordType = "AXFR"
)

// Zone represents a SafeDNS zone
type Zone struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Record represents a SafeDNS record
type Record struct {
	connection.APIRequestBodyDefaultValidator

	ID         int                 `json:"id"`
	TemplateID int                 `json:"template_id"`
	Name       string              `json:"name" validate:"required"`
	Type       RecordType          `json:"type"`
	Content    string              `json:"content" validate:"required"`
	UpdatedAt  connection.DateTime `json:"updated_at"`
	TTL        RecordTTL           `json:"ttl"`
	Priority   int                 `json:"priority"`
}

// Validate returns an error if struct properties are missing/invalid
func (c *Record) Validate() *connection.ValidationError {
	return c.APIRequestBodyDefaultValidator.Validate(c)
}

// Note represents a SafeDNS note
type Note struct {
	ID        int                  `json:"id"`
	ContactID int                  `json:"contact_id"`
	Notes     string               `json:"notes"`
	CreatedAt connection.DateTime  `json:"created_at"`
	IP        connection.IPAddress `json:"ip"`
}

// Template represents a SafeDNS template
type Template struct {
	connection.APIRequestBodyDefaultValidator

	ID        int             `json:"id"`
	Name      string          `json:"name" validate:"required"`
	Default   bool            `json:"default"`
	CreatedAt connection.Date `json:"created_at"`
}

// Validate returns an error if struct properties are missing/invalid
func (c *Template) Validate() *connection.ValidationError {
	return c.APIRequestBodyDefaultValidator.Validate(c)
}
