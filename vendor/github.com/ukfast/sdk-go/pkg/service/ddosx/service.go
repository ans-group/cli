package ddosx

import (
	"io"

	"github.com/ukfast/sdk-go/pkg/connection"
)

// DDoSXService is an interface for managing the DDoSX service
type DDoSXService interface {
	GetDomains(parameters connection.APIRequestParameters) ([]Domain, error)
	GetDomainsPaginated(parameters connection.APIRequestParameters) ([]Domain, error)
	GetDomain(domainName string) (Domain, error)
	CreateDomain(req CreateDomainRequest) error
	DeployDomain(domainName string) error
	GetDomainRecords(domainName string, parameters connection.APIRequestParameters) ([]Record, error)
	GetDomainRecordsPaginated(domainName string, parameters connection.APIRequestParameters) ([]Record, error)
	GetRecords(parameters connection.APIRequestParameters) ([]Record, error)
	GetRecordsPaginated(parameters connection.APIRequestParameters) ([]Record, error)
	CreateDomainRecord(domainName string, req CreateRecordRequest) (string, error)
	PatchDomainRecord(domainName string, recordID string, req PatchRecordRequest) error
	DeleteDomainRecord(domainName string, recordID string) error
	GetDomainProperties(domainName string, parameters connection.APIRequestParameters) ([]DomainProperty, error)
	GetDomainPropertiesPaginated(domainName string, parameters connection.APIRequestParameters) ([]DomainProperty, error)
	GetDomainProperty(domainName string, propertyID string) (DomainProperty, error)
	PatchDomainProperty(domainName string, propertyID string, req PatchDomainPropertyRequest) error
	GetDomainWAF(domainName string) (WAF, error)
	CreateDomainWAF(domainName string, req CreateWAFRequest) error
	PatchDomainWAF(domainName string, req PatchWAFRequest) error
	DeleteDomainWAF(domainName string) error
	GetDomainWAFRuleSets(domainName string, parameters connection.APIRequestParameters) ([]WAFRuleSet, error)
	GetDomainWAFRuleSetsPaginated(domainName string, parameters connection.APIRequestParameters) ([]WAFRuleSet, error)
	GetDomainWAFRuleSet(domainName string, ruleSetID string, parameters connection.APIRequestParameters) (WAFRuleSet, error)
	PatchDomainWAFRuleSet(domainName string, ruleSetID string, req PatchWAFRuleSetRequest) error
	GetDomainWAFRules(domainName string, parameters connection.APIRequestParameters) ([]WAFRule, error)
	GetDomainWAFRulesPaginated(domainName string, parameters connection.APIRequestParameters) ([]WAFRule, error)
	CreateDomainWAFRule(domainName string, req CreateWAFRuleRequest) (string, error)
	PatchDomainWAFRule(domainName string, ruleSetID string, req PatchWAFRuleRequest) error
	DeleteDomainWAFRule(domainName string, ruleID string) error
	GetDomainWAFAdvancedRules(domainName string, parameters connection.APIRequestParameters) ([]WAFAdvancedRule, error)
	GetDomainWAFAdvancedRulesPaginated(domainName string, parameters connection.APIRequestParameters) ([]WAFAdvancedRule, error)
	CreateDomainWAFAdvancedRule(domainName string, req CreateWAFAdvancedRuleRequest) (string, error)
	PatchDomainWAFAdvancedRule(domainName string, ruleID string, req PatchWAFAdvancedRuleRequest) error
	DeleteDomainWAFAdvancedRule(domainName string, ruleID string) error
	GetSSLs(parameters connection.APIRequestParameters) ([]SSL, error)
	GetSSLsPaginated(parameters connection.APIRequestParameters) ([]SSL, error)
	GetSSL(sslID string) (SSL, error)
	CreateSSL(req CreateSSLRequest) (string, error)
	PatchSSL(sslID string, req PatchSSLRequest) (string, error)
	DeleteSSL(sslID string) error
	GetSSLContent(sslID string) (SSLContent, error)
	GetSSLPrivateKey(sslID string) (SSLPrivateKey, error)
	GetDomainACLGeoIPRules(domainName string, parameters connection.APIRequestParameters) ([]ACLGeoIPRule, error)
	GetDomainACLGeoIPRulesPaginated(domainName string, parameters connection.APIRequestParameters) ([]ACLGeoIPRule, error)
	CreateDomainACLGeoIPRule(domainName string, req CreateACLGeoIPRuleRequest) (string, error)
	PatchDomainACLGeoIPRule(domainName string, ruleID string, req PatchACLGeoIPRuleRequest) error
	DeleteDomainACLGeoIPRule(domainName string, ruleID string) error
	GetDomainACLGeoIPRulesMode(domainName string) (ACLGeoIPRulesMode, error)
	PatchDomainACLGeoIPRulesMode(domainName string, req PatchACLGeoIPRulesModeRequest) error
	GetDomainACLIPRules(domainName string, parameters connection.APIRequestParameters) ([]ACLIPRule, error)
	GetDomainACLIPRulesPaginated(domainName string, parameters connection.APIRequestParameters) ([]ACLIPRule, error)
	CreateDomainACLIPRule(domainName string, req CreateACLIPRuleRequest) (string, error)
	PatchDomainACLIPRule(domainName string, ruleID string, req PatchACLIPRuleRequest) error
	DeleteDomainACLIPRule(domainName string, ruleID string) error
	DownloadDomainVerificationFile(domainName string) (string, string, error)
	DownloadDomainVerificationFileStream(domainName string) (io.ReadCloser, string, error)
}

// Service implements DDoSXService for managing
// DDoSX via the UKFast API
type Service struct {
	connection connection.Connection
}

// NewService returns a new instance of Service
func NewService(connection connection.Connection) *Service {
	return &Service{
		connection: connection,
	}
}
