//go:generate go run ../../gen/model_paginated_gen.go -package pss -typename Request,Reply -destination model_paginated.go

package pss

import "github.com/ukfast/sdk-go/pkg/connection"

type AuthorType string

func (s AuthorType) String() string {
	return string(s)
}

const (
	AuthorTypeClient  AuthorType = "Client"
	AuthorTypeAuto    AuthorType = "Auto"
	AuthorTypeSupport AuthorType = "Support"
)

type RequestPriority string

func (s RequestPriority) String() string {
	return string(s)
}

const (
	RequestPriorityNormal   RequestPriority = "Normal"
	RequestPriorityHigh     RequestPriority = "High"
	RequestPriorityCritical RequestPriority = "Critical"
)

var RequestPriorityEnum = []connection.Enum{RequestPriorityNormal, RequestPriorityHigh, RequestPriorityCritical}

// ParseRequestPriority attempts to parse a RequestPriority from string
func ParseRequestPriority(s string) (RequestPriority, error) {
	e, err := connection.ParseEnum(s, RequestPriorityEnum)
	if err != nil {
		return "", err
	}

	return e.(RequestPriority), err
}

type RequestStatus string

func (s RequestStatus) String() string {
	return string(s)
}

const (
	RequestStatusCompleted                RequestStatus = "Completed"
	RequestStatusAwaitingCustomerResponse RequestStatus = "Awaiting Customer Response"
	RequestStatusRepliedAndCompleted      RequestStatus = "Replied and Completed"
	RequestStatusSubmitted                RequestStatus = "Submitted"
)

// Request represents a PSS request
type Request struct {
	ID                int                 `json:"id"`
	Author            Author              `json:"author"`
	Type              string              `json:"type"`
	Secure            bool                `json:"secure"`
	Subject           string              `json:"subject"`
	CreatedAt         connection.DateTime `json:"created_at"`
	Priority          RequestPriority     `json:"priority"`
	Archived          bool                `json:"archived"`
	Status            RequestStatus       `json:"status"`
	RequestSMS        bool                `json:"request_sms"`
	Version           int                 `json:"version"`
	CustomerReference string              `json:"customer_reference"`
	Product           Product             `json:"product"`
	LastRepliedAt     connection.DateTime `json:"last_replied_at"`
}

// Author represents a PSS request author
type Author struct {
	ID   int        `json:"id"`
	Name string     `json:"name"`
	Type AuthorType `json:"type"`
}

// Reply represents a PSS reply
type Reply struct {
	ID          string              `json:"id"`
	Author      Author              `json:"author"`
	Description string              `json:"description"`
	CreatedAt   connection.DateTime `json:"created_at"`
	Attachments []Attachment        `json:"attachments"`
}

// Attachment represents a PSS attachment
type Attachment struct {
	Name string `json:"name"`
}

// Product represents a product to which the request applies to
type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}
