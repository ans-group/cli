package account

import (
	"strings"

	"github.com/ans-group/cli/internal/pkg/output"
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

type ApplicationCollection []account.Application

func (m ApplicationCollection) DefaultColumns() []string {
	return []string{"id", "name", "description", "created_at", "created_by"}
}

func (m ApplicationCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, app := range m {
		fields := output.NewOrderedFields()
		fields.Set("id", app.ID)
		fields.Set("name", app.Name)
		fields.Set("description", app.Description)
		fields.Set("created_at", app.CreatedAt.String())
		fields.Set("created_by", app.CreatedBy)
		data = append(data, fields)
	}
	return data
}

type ApplicationRestrictionCollection []ApplicationRestrictionWithID

func (m ApplicationRestrictionCollection) DefaultColumns() []string {
	return []string{"id", "ip_restriction_type", "ip_ranges"}
}

func (m ApplicationRestrictionCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, restriction := range m {
		fields := output.NewOrderedFields()
		fields.Set("id", restriction.ID)
		fields.Set("ip_restriction_type", restriction.IPRestrictionType)
		fields.Set("ip_ranges", strings.Join(restriction.IPRanges, ", "))
		data = append(data, fields)
	}
	return data
}
