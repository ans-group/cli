package cmd

import (
	"fmt"
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
	cmd.AddCommand(accountDetailsRootCmd())
	cmd.AddCommand(accountCreditRootCmd())
	cmd.AddCommand(accountInvoiceRootCmd())

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

// OutputAccountDetails implements OutputDataProvider for outputting an array of Details
type OutputAccountDetails struct {
	Details []account.Details
}

func outputAccountDetails(details []account.Details) {
	err := Output(&OutputAccountDetails{Details: details})
	if err != nil {
		output.Fatalf("Failed to output details: %s", err)
	}
}

func (o *OutputAccountDetails) GetData() interface{} {
	return o.Details
}

func (o *OutputAccountDetails) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, detail := range o.Details {
		fields := o.getOrderedFields(detail)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputAccountDetails) getOrderedFields(details account.Details) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("company_registration_number", output.NewFieldValue(details.CompanyRegistrationNumber, true))
	fields.Set("vat_identification_number", output.NewFieldValue(details.VATIdentificationNumber, true))
	fields.Set("primary_contact_id", output.NewFieldValue(strconv.Itoa(details.PrimaryContactID), true))

	return fields
}

// OutputAccountCredits implements OutputDataProvider for outputting an array of Credits
type OutputAccountCredits struct {
	Credits []account.Credit
}

func outputAccountCredits(credits []account.Credit) {
	err := Output(&OutputAccountCredits{Credits: credits})
	if err != nil {
		output.Fatalf("Failed to output credits: %s", err)
	}
}

func (o *OutputAccountCredits) GetData() interface{} {
	return o.Credits
}

func (o *OutputAccountCredits) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, credit := range o.Credits {
		fields := o.getOrderedFields(credit)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputAccountCredits) getOrderedFields(credit account.Credit) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("type", output.NewFieldValue(credit.Type, true))
	fields.Set("total", output.NewFieldValue(strconv.Itoa(credit.Total), true))
	fields.Set("remaining", output.NewFieldValue(strconv.Itoa(credit.Remaining), true))

	return fields
}

// OutputAccountInvoices implements OutputDataProvider for outputting an array of Invoices
type OutputAccountInvoices struct {
	Invoices []account.Invoice
}

func outputAccountInvoices(invoices []account.Invoice) {
	err := Output(&OutputAccountInvoices{Invoices: invoices})
	if err != nil {
		output.Fatalf("Failed to output invoices: %s", err)
	}
}

func (o *OutputAccountInvoices) GetData() interface{} {
	return o.Invoices
}

func (o *OutputAccountInvoices) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, invoice := range o.Invoices {
		fields := o.getOrderedFields(invoice)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputAccountInvoices) getOrderedFields(invoice account.Invoice) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("id", output.NewFieldValue(strconv.Itoa(invoice.ID), true))
	fields.Set("date", output.NewFieldValue(invoice.Date.String(), true))
	fields.Set("paid", output.NewFieldValue(strconv.FormatBool(invoice.Paid), true))
	fields.Set("net", output.NewFieldValue(fmt.Sprintf("%f", invoice.Net), true))
	fields.Set("vat", output.NewFieldValue(fmt.Sprintf("%f", invoice.VAT), true))
	fields.Set("gross", output.NewFieldValue(fmt.Sprintf("%f", invoice.Gross), true))

	return fields
}
