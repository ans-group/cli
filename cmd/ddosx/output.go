package ddosx

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func OutputDDoSXDomainsProvider(domains []ddosx.Domain) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(domains),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, domain := range domains {
				var safednsZoneID string
				if domain.SafeDNSZoneID != nil {
					safednsZoneID = strconv.Itoa(*domain.SafeDNSZoneID)
				}

				fields := output.NewOrderedFields()
				fields.Set("name", output.NewFieldValue(domain.Name, true))
				fields.Set("status", output.NewFieldValue(domain.Status.String(), true))
				fields.Set("safedns_zone_id", output.NewFieldValue(safednsZoneID, true))
				fields.Set("dns_active", output.NewFieldValue(strconv.FormatBool(domain.DNSActive), true))
				fields.Set("cdn_active", output.NewFieldValue(strconv.FormatBool(domain.CDNActive), true))
				fields.Set("waf_active", output.NewFieldValue(strconv.FormatBool(domain.WAFActive), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDDoSXRecordsProvider(records []ddosx.Record) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(records),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, record := range records {
				var safeDNSRecordID string
				if record.SafeDNSRecordID != nil {
					safeDNSRecordID = strconv.Itoa(*record.SafeDNSRecordID)
				}
				var sslID string
				if record.SSLID != nil {
					sslID = *record.SSLID
				}

				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(record.ID, true))
				fields.Set("safedns_record_id", output.NewFieldValue(safeDNSRecordID, true))
				fields.Set("ssl_id", output.NewFieldValue(sslID, true))
				fields.Set("domain_name", output.NewFieldValue(record.DomainName, true))
				fields.Set("name", output.NewFieldValue(record.Name, true))
				fields.Set("type", output.NewFieldValue(record.Type.String(), true))
				fields.Set("content", output.NewFieldValue(record.Content, true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDDoSXWAFsProvider(wafs []ddosx.WAF) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(wafs),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, waf := range wafs {
				fields := output.NewOrderedFields()
				fields.Set("mode", output.NewFieldValue(waf.Mode.String(), true))
				fields.Set("paranoia_level", output.NewFieldValue(waf.ParanoiaLevel.String(), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDDoSXWAFRuleSetsProvider(rulesets []ddosx.WAFRuleSet) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(rulesets),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, ruleset := range rulesets {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(ruleset.ID, true))
				fields.Set("name", output.NewFieldValue(ruleset.Name.String(), true))
				fields.Set("active", output.NewFieldValue(strconv.FormatBool(ruleset.Active), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDDoSXWAFRulesProvider(rules []ddosx.WAFRule) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(rules),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, rule := range rules {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(rule.ID, true))
				fields.Set("uri", output.NewFieldValue(rule.URI, true))
				fields.Set("ip", output.NewFieldValue(rule.IP.String(), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDDoSXWAFAdvancedRulesProvider(rules []ddosx.WAFAdvancedRule) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(rules),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, rule := range rules {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(rule.ID, true))
				fields.Set("section", output.NewFieldValue(rule.Section.String(), true))
				fields.Set("modifier", output.NewFieldValue(rule.Modifier.String(), true))
				fields.Set("phrase", output.NewFieldValue(rule.Phrase, true))
				fields.Set("ip", output.NewFieldValue(rule.IP.String(), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDDoSXSSLsProvider(ssls []ddosx.SSL) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(ssls),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, ssl := range ssls {
				var ukfastSSLID string
				if ssl.UKFastSSLID != nil {
					ukfastSSLID = strconv.Itoa(*ssl.UKFastSSLID)
				}

				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(ssl.ID, true))
				fields.Set("ukfast_ssl_id", output.NewFieldValue(ukfastSSLID, true))
				fields.Set("domains", output.NewFieldValue(strings.Join(ssl.Domains, ", "), true))
				fields.Set("friendly_name", output.NewFieldValue(ssl.FriendlyName, true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDDoSXSSLContentsProvider(sslContents []ddosx.SSLContent) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(sslContents),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, sslContent := range sslContents {
				fields := output.NewOrderedFields()
				fields.Set("certificate", output.NewFieldValue(sslContent.Certificate, true))
				fields.Set("ca_bundle", output.NewFieldValue(sslContent.CABundle, true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDDoSXSSLPrivateKeysProvider(sslPrivateKeys []ddosx.SSLPrivateKey) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(sslPrivateKeys),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, sslPrivateKey := range sslPrivateKeys {
				fields := output.NewOrderedFields()
				fields.Set("key", output.NewFieldValue(sslPrivateKey.Key, true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDDoSXACLIPRulesProvider(rules []ddosx.ACLIPRule) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(rules),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, rule := range rules {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(rule.ID, true))
				fields.Set("mode", output.NewFieldValue(rule.Mode.String(), true))
				fields.Set("ip", output.NewFieldValue(rule.IP.String(), true))
				fields.Set("uri", output.NewFieldValue(rule.URI, true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDDoSXACLGeoIPRulesProvider(rules []ddosx.ACLGeoIPRule) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(rules),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, rule := range rules {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(rule.ID, true))
				fields.Set("name", output.NewFieldValue(rule.Name, true))
				fields.Set("code", output.NewFieldValue(rule.Code, true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDDoSXACLGeoIPRulesModesProvider(modes []ddosx.ACLGeoIPRulesMode) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(modes),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, mode := range modes {
				fields := output.NewOrderedFields()
				fields.Set("mode", output.NewFieldValue(mode.String(), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDDoSXDomainPropertiesProvider(properties []ddosx.DomainProperty) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(properties),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, property := range properties {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(property.ID, true))
				fields.Set("name", output.NewFieldValue(property.Name.String(), true))
				fields.Set("value", output.NewFieldValue(fmt.Sprintf("%v", property.Value), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

type OutputDDoSXDomainVerificationFilesFile struct {
	Name    string
	Content string
}

func OutputDDoSXDomainVerificationFilesProvider(files []OutputDDoSXDomainVerificationFilesFile) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(files),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, file := range files {
				fields := output.NewOrderedFields()
				fields.Set("name", output.NewFieldValue(file.Name, true))
				fields.Set("content", output.NewFieldValue(file.Content, true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDDoSXCDNRulesProvider(rules []ddosx.CDNRule) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(rules),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, rule := range rules {
				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(rule.ID, true))
				fields.Set("uri", output.NewFieldValue(rule.URI, true))
				fields.Set("cache_control", output.NewFieldValue(rule.CacheControl.String(), true))
				fields.Set("cache_control_duration", output.NewFieldValue(rule.CacheControlDuration.String(), true))
				fields.Set("mime_types", output.NewFieldValue(strings.Join(rule.MimeTypes, ", "), true))
				fields.Set("type", output.NewFieldValue(rule.Type.String(), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDDoSXHSTSConfigurationsProvider(configurations []ddosx.HSTSConfiguration) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(configurations),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, configuration := range configurations {
				fields := output.NewOrderedFields()
				fields.Set("enabled", output.NewFieldValue(strconv.FormatBool(configuration.Enabled), true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}

func OutputDDoSXHSTSRulesProvider(rules []ddosx.HSTSRule) output.OutputHandlerProvider {
	return output.NewGenericOutputHandlerProvider(
		output.WithData(rules),
		output.WithFieldDataFunc(func() ([]*output.OrderedFields, error) {
			var data []*output.OrderedFields
			for _, rule := range rules {
				recordName := ""
				if rule.RecordName != nil {
					recordName = *rule.RecordName
				}

				fields := output.NewOrderedFields()
				fields.Set("id", output.NewFieldValue(rule.ID, true))
				fields.Set("max_age", output.NewFieldValue(strconv.Itoa(rule.MaxAge), true))
				fields.Set("preload", output.NewFieldValue(strconv.FormatBool(rule.Preload), true))
				fields.Set("include_subdomains", output.NewFieldValue(strconv.FormatBool(rule.IncludeSubdomains), true))
				fields.Set("type", output.NewFieldValue(rule.Type.String(), true))
				fields.Set("record_name", output.NewFieldValue(recordName, true))

				data = append(data, fields)
			}

			return data, nil
		}),
	)
}
