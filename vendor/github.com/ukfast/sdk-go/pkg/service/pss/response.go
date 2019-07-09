package pss

import "github.com/ukfast/sdk-go/pkg/connection"

// GetRequestsResponseBody represents the API response body from the GetRequests resource
type GetRequestsResponseBody struct {
	connection.APIResponseBody

	Data []Request `json:"data"`
}

// GetRequestResponseBody represents the API response body from the GetRequest resource
type GetRequestResponseBody struct {
	connection.APIResponseBody

	Data Request `json:"data"`
}

// GetReplyResponseBody represents the API response body from the GetReply resource
type GetReplyResponseBody struct {
	connection.APIResponseBody

	Data Reply `json:"data"`
}

// GetRepliesResponseBody represents the API response body from the GetReplies resource
type GetRepliesResponseBody struct {
	connection.APIResponseBody

	Data []Reply `json:"data"`
}
