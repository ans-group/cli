package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/registrar"
)

func registrarRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registrar",
		Short: "Commands relating to Registrar service",
	}

	// Child root commands
	cmd.AddCommand(registrarDomainRootCmd())

	return cmd
}

// OutputRegistrarDomains implements OutputDataProvider for outputting an array of Domains
type OutputRegistrarDomains struct {
	Domains []registrar.Domain
}

func outputRegistrarDomains(domains []registrar.Domain) {
	err := Output(&OutputRegistrarDomains{Domains: domains})
	if err != nil {
		output.Fatalf("Failed to output domains: %s", err)
	}
}

func (o *OutputRegistrarDomains) GetData() interface{} {
	return o.Domains
}

func (o *OutputRegistrarDomains) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, domain := range o.Domains {
		fields := o.getOrderedFields(domain)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputRegistrarDomains) getOrderedFields(domain registrar.Domain) *output.OrderedFields {
	fields := output.NewOrderedFields()

	fields.Set("name", output.NewFieldValue(domain.Name, true))
	fields.Set("status", output.NewFieldValue(domain.Status, true))
	fields.Set("registrar", output.NewFieldValue(domain.Registrar, true))
	fields.Set("registered_at", output.NewFieldValue(domain.RegisteredAt.String(), true))
	fields.Set("updated_at", output.NewFieldValue(domain.UpdatedAt.String(), true))
	fields.Set("renewal_at", output.NewFieldValue(domain.RenewalAt.String(), true))
	fields.Set("auto_renew", output.NewFieldValue(strconv.FormatBool(domain.AutoRenew), true))
	fields.Set("whois_privacy", output.NewFieldValue(strconv.FormatBool(domain.WHOISPrivacy), false))

	return fields
}

// OutputRegistrarNameservers implements OutputDataProvider for outputting an array of Nameservers
type OutputRegistrarNameservers struct {
	Nameservers []registrar.Nameserver
}

func outputRegistrarNameservers(domains []registrar.Nameserver) {
	err := Output(&OutputRegistrarNameservers{Nameservers: domains})
	if err != nil {
		output.Fatalf("Failed to output domains: %s", err)
	}
}

func (o *OutputRegistrarNameservers) GetData() interface{} {
	return o.Nameservers
}

func (o *OutputRegistrarNameservers) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, domain := range o.Nameservers {
		fields := o.getOrderedFields(domain)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputRegistrarNameservers) getOrderedFields(domain registrar.Nameserver) *output.OrderedFields {
	fields := output.NewOrderedFields()

	fields.Set("host", output.NewFieldValue(domain.Host, true))
	fields.Set("ip", output.NewFieldValue(domain.IP.String(), true))

	return fields
}
