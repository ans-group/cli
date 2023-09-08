package account

import (
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/account"
)

func OutputAccountContactsProvider(contacts []account.Contact) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(contacts).
		WithDefaultFields([]string{"id", "type", "first_name", "last_name"})
}

func OutputAccountDetailsProvider(details []account.Details) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(details).
		WithDefaultFields([]string{"company_registration_number", "vat_identification_number", "primary_contact_id"})
}

func OutputAccountCreditsProvider(credits []account.Credit) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(credits).
		WithDefaultFields([]string{"type", "total", "remaining"})
}

func OutputAccountClientsProvider(clients []account.Client) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(clients).
		WithDefaultFields([]string{"id", "company_name"})
}
