package ssl

import (
	"strconv"
	"strings"

	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ssl"
)

type CertificateCollection []ssl.Certificate

func (m CertificateCollection) DefaultColumns() []string {
	return []string{"id", "name", "status", "common_name", "valid_days", "ordered_date", "renewal_date"}
}

func (m CertificateCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, certificate := range m {
		fields := output.NewOrderedFields()
		fields.Set("id", strconv.Itoa(certificate.ID))
		fields.Set("name", certificate.Name)
		fields.Set("status", certificate.Status.String())
		fields.Set("common_name", certificate.CommonName)
		fields.Set("alternative_names", strings.Join(certificate.AlternativeNames, ", "))
		fields.Set("valid_days", strconv.Itoa(certificate.ValidDays))
		fields.Set("ordered_date", certificate.OrderedDate.String())
		fields.Set("renewal_date", certificate.RenewalDate.String())

		data = append(data, fields)
	}

	return data
}

type CertificateContentCollection []ssl.CertificateContent

func (m CertificateContentCollection) DefaultColumns() []string {
	return []string{"combined"}
}

func (m CertificateContentCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, certificateContent := range m {
		fields := output.NewOrderedFields()
		fields.Set("combined", certificateContent.Server+"\n"+certificateContent.Intermediate)
		fields.Set("server", certificateContent.Server)
		fields.Set("intermediate", certificateContent.Intermediate)

		data = append(data, fields)
	}

	return data
}

type CertificatePrivateKeyCollection []ssl.CertificatePrivateKey

func (m CertificatePrivateKeyCollection) DefaultColumns() []string {
	return []string{"key"}
}

func (m CertificatePrivateKeyCollection) Fields() []*output.OrderedFields {
	var data []*output.OrderedFields
	for _, certificatePrivateKey := range m {
		fields := output.NewOrderedFields()
		fields.Set("key", certificatePrivateKey.Key)

		data = append(data, fields)
	}

	return data
}

type CertificateValidationCollection []ssl.CertificateValidation

func (m CertificateValidationCollection) DefaultColumns() []string {
	return []string{"domains", "expires_at"}
}

type RecommendationsCollection []ssl.Recommendations

func (m RecommendationsCollection) DefaultColumns() []string {
	return []string{"level", "messages"}
}

type ReportCollection []ssl.Report

func (m ReportCollection) DefaultColumns() []string {
	return []string{"certificate_name", "certificate_expiring", "certificate_expired", "chain_intact"}
}
