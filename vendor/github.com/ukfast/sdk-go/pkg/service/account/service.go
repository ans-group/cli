package account

import (
	"github.com/ukfast/sdk-go/pkg/connection"
)

// AccountService is an interface for managing account
type AccountService interface {
	GetContacts(parameters connection.APIRequestParameters) ([]Contact, error)
	GetContactsPaginated(parameters connection.APIRequestParameters) ([]Contact, error)
	GetContact(contactID int) (Contact, error)
}

// Service implements AccountService for managing
// Account certificates via the UKFast API
type Service struct {
	connection connection.Connection
}

// NewService returns a new instance of AccountService
func NewService(connection connection.Connection) *Service {
	return &Service{
		connection: connection,
	}
}
