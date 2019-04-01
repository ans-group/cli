package pss

import "fmt"

// RequestNotFoundError indicates a request was not found
type RequestNotFoundError struct {
	ID int
}

func (e *RequestNotFoundError) Error() string {
	return fmt.Sprintf("Request not found with id [%d]", e.ID)
}
