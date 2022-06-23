package registrar

import (
	"strconv"
	"strings"

	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/registrar"
)

func OutputRegistrarDomainsProvider(domains []registrar.Domain) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(domains),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, domain := range domains {
				fields := output.NewOrderedFields()
				fields.Set("name", output.NewFieldValue(domain.Name, true))
				fields.Set("status", output.NewFieldValue(domain.Status, true))
				fields.Set("registrar", output.NewFieldValue(domain.Registrar, true))
				fields.Set("registered_at", output.NewFieldValue(domain.RegisteredAt.String(), true))
				fields.Set("updated_at", output.NewFieldValue(domain.UpdatedAt.String(), true))
				fields.Set("renewal_at", output.NewFieldValue(domain.RenewalAt.String(), true))
				fields.Set("auto_renew", output.NewFieldValue(strconv.FormatBool(domain.AutoRenew), true))
				fields.Set("whois_privacy", output.NewFieldValue(strconv.FormatBool(domain.WHOISPrivacy), false))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputRegistrarNameserversProvider(nameservers []registrar.Nameserver) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(nameservers),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, nameserver := range nameservers {
				fields := output.NewOrderedFields()
				fields.Set("host", output.NewFieldValue(nameserver.Host, true))
				fields.Set("ip", output.NewFieldValue(nameserver.IP.String(), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputRegistrarWhoisProvider(whoisArr []registrar.Whois) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(whoisArr),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, whois := range whoisArr {
				fields := output.NewOrderedFields()
				fields.Set("name", output.NewFieldValue(whois.Name, true))
				fields.Set("status", output.NewFieldValue(strings.Join(whois.Status, ", "), true))
				fields.Set("created_at", output.NewFieldValue(whois.CreatedAt.String(), true))
				fields.Set("updated_at", output.NewFieldValue(whois.UpdatedAt.String(), true))
				fields.Set("expires_at", output.NewFieldValue(whois.ExpiresAt.String(), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}
