//go:generate go run ../../gen/model_paginated_gen.go -package ssl -typename Certificate -destination model_paginated.go

package ssl

import (
	"github.com/ukfast/sdk-go/pkg/connection"
)

type CertificateStatus string

func (s CertificateStatus) String() string {
	return string(s)
}

const (
	CertificateStatusCompleted      CertificateStatus = "Completed"
	CertificateStatusProcessing     CertificateStatus = "Processing"
	CertificateStatusExpired        CertificateStatus = "Expired"
	CertificateStatusExpiring       CertificateStatus = "Expiring"
	CertificateStatusPendingInstall CertificateStatus = "Pending Install"
)

// Certificate represents an SSL certificate
type Certificate struct {
	ID               int                 `json:"id"`
	Name             string              `json:"name"`
	Status           CertificateStatus   `json:"status"`
	CommonName       string              `json:"common_name"`
	AlternativeNames []string            `json:"alternative_names"`
	ValidDays        int                 `json:"valid_days"`
	OrderedDate      connection.DateTime `json:"ordered_date"`
	RenewalDate      connection.DateTime `json:"renewal_date"`
}

// CertificateContent represents the content of an SSL certificate
type CertificateContent struct {
	Server       string `json:"server"`
	Intermediate string `json:"intermediate"`
}

// CertificatePrivateKey represents an SSL certificate private key
type CertificatePrivateKey struct {
	Key string `json:"key"`
}
