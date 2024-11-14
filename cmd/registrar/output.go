package registrar

import (
	"strconv"
	"strings"

	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/registrar"
)

type DomainCollection []registrar.Domain

func (d DomainCollection) DefaultColumns() []string {
	return []string{"name", "status", "registrar", "registered_at", "updated_at", "renewal_at", "auto_renew"}
}

func (d DomainCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, domain := range d {
		fields := output.NewOrderedFields()
		fields.Set("name", domain.Name)
		fields.Set("status", domain.Status)
		fields.Set("registrar", domain.Registrar)
		fields.Set("registered_at", domain.RegisteredAt.String())
		fields.Set("updated_at", domain.UpdatedAt.String())
		fields.Set("renewal_at", domain.RenewalAt.String())
		fields.Set("auto_renew", strconv.FormatBool(domain.AutoRenew))
		fields.Set("whois_privacy", strconv.FormatBool(domain.WHOISPrivacy))

		data = append(data, fields)
	}

	return data
}

type NameserverCollection []registrar.Nameserver

func (n NameserverCollection) DefaultColumns() []string {
	return []string{"host", "ip"}
}

func (n NameserverCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, nameserver := range n {
		fields := output.NewOrderedFields()
		fields.Set("host", nameserver.Host)
		fields.Set("ip", nameserver.IP.String())

		data = append(data, fields)
	}

	return data
}

type WhoisCollection []registrar.Whois

func (w WhoisCollection) DefaultColumns() []string {
	return []string{"name", "status", "created_at", "updated_at", "expires_at"}
}

func (w WhoisCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, whois := range w {
		fields := output.NewOrderedFields()
		fields.Set("name", whois.Name)
		fields.Set("status", strings.Join(whois.Status, ", "))
		fields.Set("created_at", whois.CreatedAt.String())
		fields.Set("updated_at", whois.UpdatedAt.String())
		fields.Set("expires_at", whois.ExpiresAt.String())

		data = append(data, fields)
	}

	return data
}
