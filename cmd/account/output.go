package account

import (
	"github.com/ans-group/sdk-go/pkg/service/account"
)

type ContactCollection []account.Contact

func (m ContactCollection) DefaultColumns() []string {
	return []string{"id", "type", "first_name", "last_name"}
}

type DetailsCollection []account.Details

func (m DetailsCollection) DefaultColumns() []string {
	return []string{"company_registration_number", "vat_identification_number", "primary_contact_id"}
}

type CreditCollection []account.Credit

func (m CreditCollection) DefaultColumns() []string {
	return []string{"type", "total", "remaining"}
}

type ClientCollection []account.Client

func (m ClientCollection) DefaultColumns() []string {
	return []string{"id", "company_name"}
}
