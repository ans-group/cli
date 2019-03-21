package cmd

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ssl"
)

func sslRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssl",
		Short: "Commands relating to SSL service",
	}

	// Child root commands
	cmd.AddCommand(sslCertificateRootCmd())

	return cmd
}

// OutputSSLCertificates implements OutputDataProvider for outputting an array of Certificates
type OutputSSLCertificates struct {
	Certificates []ssl.Certificate
}

func outputSSLCertificates(certificates []ssl.Certificate) {
	err := Output(&OutputSSLCertificates{Certificates: certificates})
	if err != nil {
		output.Fatalf("Failed to output certificates: %s", err)
	}
}

func (o *OutputSSLCertificates) GetData() interface{} {
	return o.Certificates
}

func (o *OutputSSLCertificates) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, certificate := range o.Certificates {
		fields := o.getOrderedFields(certificate)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputSSLCertificates) getOrderedFields(certificate ssl.Certificate) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("id", output.NewFieldValue(strconv.Itoa(certificate.ID), true))
	fields.Set("name", output.NewFieldValue(certificate.Name, true))
	fields.Set("status", output.NewFieldValue(certificate.Status.String(), true))
	fields.Set("common_name", output.NewFieldValue(certificate.CommonName, true))
	fields.Set("alternative_names", output.NewFieldValue(strings.Join(certificate.AlternativeNames, ", "), false))
	fields.Set("valid_days", output.NewFieldValue(strconv.Itoa(certificate.ValidDays), true))
	fields.Set("ordered_date", output.NewFieldValue(certificate.OrderedDate.String(), true))
	fields.Set("renewal_date", output.NewFieldValue(certificate.RenewalDate.String(), true))

	return fields
}

// OutputSSLCertificateContents implements OutputDataProvider for outputting an array of certificates contents
type OutputSSLCertificateContents struct {
	CertificateContents []ssl.CertificateContent
}

func outputSSLCertificatesContents(certificatesContent []ssl.CertificateContent) {
	err := Output(&OutputSSLCertificateContents{CertificateContents: certificatesContent})
	if err != nil {
		output.Fatalf("Failed to output certificate contents: %s", err)
	}
}

func (o *OutputSSLCertificateContents) GetData() interface{} {
	return o.CertificateContents
}

func (o *OutputSSLCertificateContents) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, certificateContent := range o.CertificateContents {
		fields := o.getOrderedFields(certificateContent)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputSSLCertificateContents) getOrderedFields(certificateContent ssl.CertificateContent) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("combined", output.NewFieldValue(certificateContent.Server+"\n"+certificateContent.Intermediate, true))
	fields.Set("server", output.NewFieldValue(certificateContent.Server, false))
	fields.Set("intermediate", output.NewFieldValue(certificateContent.Intermediate, false))

	return fields
}

// OutputSSLCertificatePrivateKeys implements OutputDataProvider for outputting an array of Certificates
type OutputSSLCertificatePrivateKeys struct {
	CertificatePrivateKeys []ssl.CertificatePrivateKey
}

func outputSSLCertificatesPrivateKeys(certificatesPrivateKey []ssl.CertificatePrivateKey) {
	err := Output(&OutputSSLCertificatePrivateKeys{
		CertificatePrivateKeys: certificatesPrivateKey,
	})
	if err != nil {
		output.Fatalf("Failed to output certificate private keys: %s", err)
	}
}

func (o *OutputSSLCertificatePrivateKeys) GetData() interface{} {
	return o.CertificatePrivateKeys
}

func (o *OutputSSLCertificatePrivateKeys) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, certificatePrivateKey := range o.CertificatePrivateKeys {
		fields := o.getOrderedFields(certificatePrivateKey)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputSSLCertificatePrivateKeys) getOrderedFields(certificatePrivateKey ssl.CertificatePrivateKey) *output.OrderedFields {
	fields := output.NewOrderedFields()
	fields.Set("key", output.NewFieldValue(certificatePrivateKey.Key, true))

	return fields
}
