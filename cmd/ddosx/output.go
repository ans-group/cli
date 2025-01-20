package ddosx

import (
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
)

type DomainCollection []ddosx.Domain

func (m DomainCollection) DefaultColumns() []string {
	return []string{"name", "status", "safedns_zone_id", "dns_active", "cdn_active", "waf_active"}
}

type RecordCollection []ddosx.Record

func (m RecordCollection) DefaultColumns() []string {
	return []string{"id", "safedns_record_id", "ssl_id", "domain_name", "name", "type", "content"}
}

type WAFCollection []ddosx.WAF

func (m WAFCollection) DefaultColumns() []string {
	return []string{"mode", "paranoia_level"}
}

type WAFRuleSetCollection []ddosx.WAFRuleSet

func (m WAFRuleSetCollection) DefaultColumns() []string {
	return []string{"id", "name", "active"}
}

type WAFRuleCollection []ddosx.WAFRule

func (m WAFRuleCollection) DefaultColumns() []string {
	return []string{"id", "uri", "ip"}
}

type WAFAdvancedRuleCollection []ddosx.WAFAdvancedRule

func (m WAFAdvancedRuleCollection) DefaultColumns() []string {
	return []string{"id", "section", "modifier", "phrase", "ip"}
}

type SSLCollection []ddosx.SSL

func (m SSLCollection) DefaultColumns() []string {
	return []string{"id", "ans_ssl_id", "domains", "friendly_name"}
}

type SSLContentCollection []ddosx.SSLContent

func (m SSLContentCollection) DefaultColumns() []string {
	return []string{"certificate", "ca_bundle"}
}

type SSLPrivateKeyCollection []ddosx.SSLPrivateKey

func (m SSLPrivateKeyCollection) DefaultColumns() []string {
	return []string{"key"}
}

type ACLIPRuleCollection []ddosx.ACLIPRule

func (m ACLIPRuleCollection) DefaultColumns() []string {
	return []string{"id", "mode", "ip", "uri"}
}

type ACLGeoIPRuleCollection []ddosx.ACLGeoIPRule

func (m ACLGeoIPRuleCollection) DefaultColumns() []string {
	return []string{"id", "name", "code"}
}

type ACLGeoIPRulesModeCollection []ddosx.ACLGeoIPRulesMode

func (m ACLGeoIPRulesModeCollection) DefaultColumns() []string {
	return []string{"mode"}
}

type DomainPropertyCollection []ddosx.DomainProperty

func (m DomainPropertyCollection) DefaultColumns() []string {
	return []string{"id", "name", "value"}
}

type OutputDDoSXDomainVerificationFilesFile struct {
	Name    string
	Content string
}

type OutputDDoSXDomainVerificationFilesFileCollection []OutputDDoSXDomainVerificationFilesFile

func (m OutputDDoSXDomainVerificationFilesFile) DefaultColumns() []string {
	return []string{"name", "content"}
}

type CDNRuleCollection []ddosx.CDNRule

func (m CDNRuleCollection) DefaultColumns() []string {
	return []string{"id", "uri", "cache_control", "cache_control_duration", "mime_types", "type"}
}

type HSTSConfigurationCollection []ddosx.HSTSConfiguration

func (m HSTSConfigurationCollection) DefaultColumns() []string {
	return []string{"enabled"}
}

type HSTSRuleCollection []ddosx.HSTSRule

func (m HSTSRuleCollection) DefaultColumns() []string {
	return []string{"id", "max_age", "preload", "include_subdomains", "type", "record_name"}
}

type WAFLogCollection []ddosx.WAFLog

func (m WAFLogCollection) DefaultColumns() []string {
	return []string{"id", "created_at", "client_ip", "request"}
}

type WAFLogMatchCollection []ddosx.WAFLogMatch

func (m WAFLogMatchCollection) DefaultColumns() []string {
	return []string{"id", "created_at", "country_code", "method", "message"}
}
