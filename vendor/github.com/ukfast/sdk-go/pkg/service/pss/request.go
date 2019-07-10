package pss

import "github.com/ukfast/sdk-go/pkg/connection"

// CreateRequestRequest represents a request to create a PSS request
type CreateRequestRequest struct {
	connection.APIRequestBodyDefaultValidator

	Author            Author          `json:"author" validate:"required"`
	Secure            bool            `json:"secure"`
	Priority          RequestPriority `json:"priority" validate:"required"`
	Subject           string          `json:"subject" validate:"required"`
	Details           string          `json:"details" validate:"required"`
	CC                []string        `json:"cc,omitempty"`
	RequestSMS        bool            `json:"request_sms"`
	CustomerReference string          `json:"customer_reference,omitempty"`
	Product           *Product        `json:"product,omitempty"`
}

// Validate returns an error if struct properties are missing/invalid
func (c *CreateRequestRequest) Validate() *connection.ValidationError {
	return c.APIRequestBodyDefaultValidator.Validate(c)
}

// PatchRequestRequest represents a PSS Request patch request
type PatchRequestRequest struct {
	Secure     *bool           `json:"secure,omitempty"`
	Read       *bool           `json:"read,omitempty"`
	Priority   RequestPriority `json:"priority,omitempty"`
	RequestSMS *bool           `json:"request_sms,omitempty"`
	Archived   *bool           `json:"archived,omitempty"`
}

// CreateReplyRequest represents a request to create a PSS request reply
type CreateReplyRequest struct {
	connection.APIRequestBodyDefaultValidator

	Author      Author `json:"author" validate:"required"`
	Description string `json:"description" validate:"required"`
}

// Validate returns an error if struct properties are missing/invalid
func (c *CreateReplyRequest) Validate() *connection.ValidationError {
	return c.APIRequestBodyDefaultValidator.Validate(c)
}
