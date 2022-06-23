package ddosx

import (
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
)

func OutputDDoSXDomainsProvider(domains []ddosx.Domain) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(domains).
		WithDefaultFields([]string{"name", "status", "safedns_zone_id", "dns_active", "cdn_active", "waf_active"})
}

func OutputDDoSXRecordsProvider(records []ddosx.Record) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(records).
		WithDefaultFields([]string{"id", "safedns_record_id", "ssl_id", "domain_name", "name", "type", "content"})
}

func OutputDDoSXWAFsProvider(wafs []ddosx.WAF) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(wafs).
		WithDefaultFields([]string{"mode", "paranoia_level"})
}

func OutputDDoSXWAFRuleSetsProvider(rulesets []ddosx.WAFRuleSet) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(rulesets).
		WithDefaultFields([]string{"id", "name", "active"})
}

func OutputDDoSXWAFRulesProvider(rules []ddosx.WAFRule) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(rules).
		WithDefaultFields([]string{"id", "uri", "ip"})
}

func OutputDDoSXWAFAdvancedRulesProvider(rules []ddosx.WAFAdvancedRule) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(rules).
		WithDefaultFields([]string{"id", "section", "modifier", "phrase", "ip"})
}

func OutputDDoSXSSLsProvider(ssls []ddosx.SSL) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(ssls).
		WithDefaultFields([]string{"id", "ukfast_ssl_id", "domains", "friendly_name"})
}

func OutputDDoSXSSLContentsProvider(sslContents []ddosx.SSLContent) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(sslContents).
		WithDefaultFields([]string{"certificate", "ca_bundle"})
}

func OutputDDoSXSSLPrivateKeysProvider(sslPrivateKeys []ddosx.SSLPrivateKey) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(sslPrivateKeys).
		WithDefaultFields([]string{"key"})
}

func OutputDDoSXACLIPRulesProvider(rules []ddosx.ACLIPRule) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(rules).
		WithDefaultFields([]string{"id", "mode", "ip", "uri"})
}

func OutputDDoSXACLGeoIPRulesProvider(rules []ddosx.ACLGeoIPRule) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(rules).
		WithDefaultFields([]string{"id", "name", "code"})
}

func OutputDDoSXACLGeoIPRulesModesProvider(modes []ddosx.ACLGeoIPRulesMode) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(modes).
		WithDefaultFields([]string{"mode"})
}

func OutputDDoSXDomainPropertiesProvider(properties []ddosx.DomainProperty) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(properties).
		WithDefaultFields([]string{"id", "name", "value"})
}

type OutputDDoSXDomainVerificationFilesFile struct {
	Name    string
	Content string
}

func OutputDDoSXDomainVerificationFilesProvider(files []OutputDDoSXDomainVerificationFilesFile) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(files).
		WithDefaultFields([]string{"name", "content"})
}

func OutputDDoSXCDNRulesProvider(rules []ddosx.CDNRule) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(rules).
		WithDefaultFields([]string{"id", "uri", "cache_control", "cache_control_duration", "mime_types", "type"})
}

func OutputDDoSXHSTSConfigurationsProvider(configurations []ddosx.HSTSConfiguration) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(configurations).
		WithDefaultFields([]string{"enabled"})
}

func OutputDDoSXHSTSRulesProvider(rules []ddosx.HSTSRule) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(rules).
		WithDefaultFields([]string{"id", "max_age", "preload", "include_subdomains", "type", "record_name"})
}

func OutputDDoSXWAFLogsProvider(logs []ddosx.WAFLog) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(logs).
		WithDefaultFields([]string{"id", "created_at", "client_ip", "request"})
}

func OutputDDoSXWAFLogMatchesProvider(matches []ddosx.WAFLogMatch) output.OutputHandlerDataProvider {
	return output.NewSerializedOutputHandlerDataProvider(matches).
		WithDefaultFields([]string{"id", "created_at", "country_code", "method", "message"})
}
