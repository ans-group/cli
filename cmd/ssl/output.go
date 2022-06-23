package ssl

import (
	"strconv"
	"strings"

	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ssl"
)

func OutputSSLCertificatesProvider(certificates []ssl.Certificate) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(certificates),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, certificate := range certificates {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(strconv.Itoa(certificate.ID), true))
				fields.Set("name", output.NewFieldValue(certificate.Name, true))
				fields.Set("status", output.NewFieldValue(certificate.Status.String(), true))
				fields.Set("common_name", output.NewFieldValue(certificate.CommonName, true))
				fields.Set("alternative_names", output.NewFieldValue(strings.Join(certificate.AlternativeNames, ", "), false))
				fields.Set("valid_days", output.NewFieldValue(strconv.Itoa(certificate.ValidDays), true))
				fields.Set("ordered_date", output.NewFieldValue(certificate.OrderedDate.String(), true))
				fields.Set("renewal_date", output.NewFieldValue(certificate.RenewalDate.String(), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputSSLCertificatesContentsProvider(certificatesContent []ssl.CertificateContent) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(certificatesContent),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, certificateContent := range certificatesContent {
				fields := output.NewOrderedFields()
				fields.Set("combined", output.NewFieldValue(certificateContent.Server+"\n"+certificateContent.Intermediate, true))
				fields.Set("server", output.NewFieldValue(certificateContent.Server, false))
				fields.Set("intermediate", output.NewFieldValue(certificateContent.Intermediate, false))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputSSLCertificatesPrivateKeysProvider(certificatesPrivateKey []ssl.CertificatePrivateKey) output.OutputHandlerDataProvider {
	return output.NewGenericOutputHandlerDataProvider(
		output.WithData(certificatesPrivateKey),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, certificatePrivateKey := range certificatesPrivateKey {
				fields := output.NewOrderedFields()
				fields.Set("key", output.NewFieldValue(certificatePrivateKey.Key, true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputSSLCertificateValidationsProvider(validations []ssl.CertificateValidation) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(validations).WithDefaultFields([]string{"domains", "expires_at"})
}

func OutputSSLRecommendationsProvider(recommendationsSlice []ssl.Recommendations) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(recommendationsSlice).WithDefaultFields([]string{"level", "messages"})
}

func OutputSSLReportsProvider(reports []ssl.Report) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(reports).WithDefaultFields([]string{"certificate_name", "certificate_expiring", "certificate_expired", "chain_intact"})
}
