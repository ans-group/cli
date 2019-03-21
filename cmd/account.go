package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/account"
)

func accountRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "Commands relating to Account service",
	}

	// Child root commands
	cmd.AddCommand(accountContactRootCmd())

	return cmd
}

// OutputAccountContacts implements OutputDataProvider for outputting an array of Contacts
type OutputAccountContacts struct {
	Contacts []account.Contact
}

func outputAccountContacts(contacts []account.Contact) {
	err := Output(&OutputAccountContacts{Contacts: contacts})
	if err != nil {
		output.Fatalf("Failed to output contacts: %s", err)
	}
}

func (o *OutputAccountContacts) GetData() interface{} {
	return o.Contacts
}

func (o *OutputAccountContacts) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, contact := range o.Contacts {
		fields := o.getOrderedFields(contact)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputAccountContacts) getOrderedFields(contact account.Contact) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("id", output.NewFieldValue(strconv.Itoa(contact.ID), true))
	fields.Set("type", output.NewFieldValue(contact.Type.String(), true))
	fields.Set("first_name", output.NewFieldValue(contact.FirstName, true))
	fields.Set("last_name", output.NewFieldValue(contact.LastName, true))

	return fields
}
