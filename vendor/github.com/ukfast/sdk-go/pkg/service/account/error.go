package account

import "fmt"

// ContactNotFoundError indicates a contact was not found
type ContactNotFoundError struct {
	ID int
}

func (e *ContactNotFoundError) Error() string {
	return fmt.Sprintf("Contact not found with ID [%d]", e.ID)
}
