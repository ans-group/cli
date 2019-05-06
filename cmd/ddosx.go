package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ddosx",
		Short: "Commands relating to DDoSX service",
	}

	// Child root commands
	cmd.AddCommand(ddosxDomainRootCmd())
	cmd.AddCommand(ddosxRecordRootCmd())
	cmd.AddCommand(ddosxSSLRootCmd())

	return cmd
}

// OutputDDoSXDomains implements OutputDataProvider for outputting an array of Domains
type OutputDDoSXDomains struct {
	Domains []ddosx.Domain
}

func outputDDoSXDomains(domains []ddosx.Domain) {
	err := Output(&OutputDDoSXDomains{Domains: domains})
	if err != nil {
		output.Fatalf("Failed to output domains: %s", err)
	}
}

func (o *OutputDDoSXDomains) GetData() interface{} {
	return o.Domains
}

func (o *OutputDDoSXDomains) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, domain := range o.Domains {
		fields := o.getOrderedFields(domain)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputDDoSXDomains) getOrderedFields(domain ddosx.Domain) *output.OrderedFields {
	fields := output.NewOrderedFields()

	var safednsZoneID string
	if domain.SafeDNSZoneID != nil {
		safednsZoneID = strconv.Itoa(*domain.SafeDNSZoneID)
	}

	fields.Set("name", output.NewFieldValue(domain.Name, true))
	fields.Set("status", output.NewFieldValue(domain.Status.String(), true))
	fields.Set("safedns_zone_id", output.NewFieldValue(safednsZoneID, true))
	fields.Set("dns_active", output.NewFieldValue(strconv.FormatBool(domain.DNSActive), true))
	fields.Set("cdn_active", output.NewFieldValue(strconv.FormatBool(domain.CDNActive), true))
	fields.Set("waf_active", output.NewFieldValue(strconv.FormatBool(domain.WAFActive), true))

	return fields
}

// OutputDDoSXRecords implements OutputDataProvider for outputting an array of Records
type OutputDDoSXRecords struct {
	Records []ddosx.Record
}

func outputDDoSXRecords(records []ddosx.Record) {
	err := Output(&OutputDDoSXRecords{Records: records})
	if err != nil {
		output.Fatalf("Failed to output records: %s", err)
	}
}

func (o *OutputDDoSXRecords) GetData() interface{} {
	return o.Records
}

func (o *OutputDDoSXRecords) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, record := range o.Records {
		fields := o.getOrderedFields(record)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputDDoSXRecords) getOrderedFields(record ddosx.Record) *output.OrderedFields {
	fields := output.NewOrderedFields()

	var safeDNSRecordID string
	if record.SafeDNSRecordID != nil {
		safeDNSRecordID = strconv.Itoa(*record.SafeDNSRecordID)
	}
	var sslID string
	if record.SSLID != nil {
		sslID = *record.SSLID
	}

	fields.Set("id", output.NewFieldValue(record.ID, true))
	fields.Set("safedns_record_id", output.NewFieldValue(safeDNSRecordID, true))
	fields.Set("ssl_id", output.NewFieldValue(sslID, true))
	fields.Set("domain_name", output.NewFieldValue(record.DomainName, true))
	fields.Set("name", output.NewFieldValue(record.Name, true))
	fields.Set("type", output.NewFieldValue(record.Type.String(), true))
	fields.Set("content", output.NewFieldValue(record.Content, true))

	return fields
}

// OutputDDoSXWAFs implements OutputDataProvider for outputting an array of WAFs
type OutputDDoSXWAFs struct {
	WAFs []ddosx.WAF
}

func outputDDoSXWAFs(wafs []ddosx.WAF) {
	err := Output(&OutputDDoSXWAFs{WAFs: wafs})
	if err != nil {
		output.Fatalf("Failed to output wafs: %s", err)
	}
}

func (o *OutputDDoSXWAFs) GetData() interface{} {
	return o.WAFs
}

func (o *OutputDDoSXWAFs) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, waf := range o.WAFs {
		fields := o.getOrderedFields(waf)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputDDoSXWAFs) getOrderedFields(waf ddosx.WAF) *output.OrderedFields {
	fields := output.NewOrderedFields()

	fields.Set("mode", output.NewFieldValue(waf.Mode.String(), true))
	fields.Set("paranoia_level", output.NewFieldValue(waf.ParanoiaLevel.String(), true))

	return fields
}

// OutputDDoSXWAFRuleSets implements OutputDataProvider for outputting an array of WAFRuleSets
type OutputDDoSXWAFRuleSets struct {
	WAFRuleSets []ddosx.WAFRuleSet
}

func outputDDoSXWAFRuleSets(wafRuleSets []ddosx.WAFRuleSet) {
	err := Output(&OutputDDoSXWAFRuleSets{WAFRuleSets: wafRuleSets})
	if err != nil {
		output.Fatalf("Failed to output waf rule sets: %s", err)
	}
}

func (o *OutputDDoSXWAFRuleSets) GetData() interface{} {
	return o.WAFRuleSets
}

func (o *OutputDDoSXWAFRuleSets) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, wafRuleSet := range o.WAFRuleSets {
		fields := o.getOrderedFields(wafRuleSet)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputDDoSXWAFRuleSets) getOrderedFields(wafRuleSet ddosx.WAFRuleSet) *output.OrderedFields {
	fields := output.NewOrderedFields()

	fields.Set("id", output.NewFieldValue(wafRuleSet.ID, true))
	fields.Set("name", output.NewFieldValue(wafRuleSet.Name.String(), true))
	fields.Set("active", output.NewFieldValue(strconv.FormatBool(wafRuleSet.Active), true))

	return fields
}

// OutputDDoSXWAFRules implements OutputDataProvider for outputting an array of WAFRules
type OutputDDoSXWAFRules struct {
	WAFRules []ddosx.WAFRule
}

func outputDDoSXWAFRules(wafRules []ddosx.WAFRule) {
	err := Output(&OutputDDoSXWAFRules{WAFRules: wafRules})
	if err != nil {
		output.Fatalf("Failed to output waf rules: %s", err)
	}
}

func (o *OutputDDoSXWAFRules) GetData() interface{} {
	return o.WAFRules
}

func (o *OutputDDoSXWAFRules) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, wafRule := range o.WAFRules {
		fields := o.getOrderedFields(wafRule)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputDDoSXWAFRules) getOrderedFields(wafRule ddosx.WAFRule) *output.OrderedFields {
	fields := output.NewOrderedFields()

	fields.Set("id", output.NewFieldValue(wafRule.ID, true))
	fields.Set("uri", output.NewFieldValue(wafRule.URI, true))
	fields.Set("ip", output.NewFieldValue(wafRule.IP.String(), true))

	return fields
}

// OutputDDoSXWAFAdvancedRules implements OutputDataProvider for outputting an array of WAFAdvancedRules
type OutputDDoSXWAFAdvancedRules struct {
	WAFAdvancedRules []ddosx.WAFAdvancedRule
}

func outputDDoSXWAFAdvancedRules(wafAdvancedRules []ddosx.WAFAdvancedRule) {
	err := Output(&OutputDDoSXWAFAdvancedRules{WAFAdvancedRules: wafAdvancedRules})
	if err != nil {
		output.Fatalf("Failed to output waf rules: %s", err)
	}
}

func (o *OutputDDoSXWAFAdvancedRules) GetData() interface{} {
	return o.WAFAdvancedRules
}

func (o *OutputDDoSXWAFAdvancedRules) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, wafAdvancedRule := range o.WAFAdvancedRules {
		fields := o.getOrderedFields(wafAdvancedRule)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputDDoSXWAFAdvancedRules) getOrderedFields(wafAdvancedRule ddosx.WAFAdvancedRule) *output.OrderedFields {
	fields := output.NewOrderedFields()

	fields.Set("id", output.NewFieldValue(wafAdvancedRule.ID, true))
	fields.Set("section", output.NewFieldValue(wafAdvancedRule.Section.String(), true))
	fields.Set("modifier", output.NewFieldValue(wafAdvancedRule.Modifier.String(), true))
	fields.Set("phrase", output.NewFieldValue(wafAdvancedRule.Phrase, true))
	fields.Set("ip", output.NewFieldValue(wafAdvancedRule.IP.String(), true))

	return fields
}

// OutputDDoSXSSLs implements OutputDataProvider for outputting an array of SSLs
type OutputDDoSXSSLs struct {
	SSLs []ddosx.SSL
}

func outputDDoSXSSLs(ssls []ddosx.SSL) {
	err := Output(&OutputDDoSXSSLs{SSLs: ssls})
	if err != nil {
		output.Fatalf("Failed to output ssls: %s", err)
	}
}

func (o *OutputDDoSXSSLs) GetData() interface{} {
	return o.SSLs
}

func (o *OutputDDoSXSSLs) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, ssl := range o.SSLs {
		fields := o.getOrderedFields(ssl)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputDDoSXSSLs) getOrderedFields(ssl ddosx.SSL) *output.OrderedFields {
	fields := output.NewOrderedFields()

	var ukfastSSLID string
	if ssl.UKFastSSLID != nil {
		ukfastSSLID = strconv.Itoa(*ssl.UKFastSSLID)
	}

	fields.Set("id", output.NewFieldValue(ssl.ID, true))
	fields.Set("ukfast_ssl_id", output.NewFieldValue(ukfastSSLID, true))
	fields.Set("domains", output.NewFieldValue(strings.Join(ssl.Domains, ", "), true))
	fields.Set("friendly_name", output.NewFieldValue(ssl.FriendlyName, true))

	return fields
}

// OutputDDoSXSSLContents implements OutputDataProvider for outputting an array of SSLContentss
type OutputDDoSXSSLContents struct {
	SSLContents []ddosx.SSLContent
}

func outputDDoSXSSLContents(sslContents []ddosx.SSLContent) {
	err := Output(&OutputDDoSXSSLContents{SSLContents: sslContents})
	if err != nil {
		output.Fatalf("Failed to output ssl contents: %s", err)
	}
}

func (o *OutputDDoSXSSLContents) GetData() interface{} {
	return o.SSLContents
}

func (o *OutputDDoSXSSLContents) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, sslContent := range o.SSLContents {
		fields := o.getOrderedFields(sslContent)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputDDoSXSSLContents) getOrderedFields(sslContent ddosx.SSLContent) *output.OrderedFields {
	fields := output.NewOrderedFields()

	fields.Set("certificate", output.NewFieldValue(sslContent.Certificate, true))
	fields.Set("ca_bundle", output.NewFieldValue(sslContent.CABundle, true))
	return fields
}

// OutputDDoSXSSLPrivateKeys implements OutputDataProvider for outputting an array of SSLPrivateKeyss
type OutputDDoSXSSLPrivateKeys struct {
	SSLPrivateKeys []ddosx.SSLPrivateKey
}

func outputDDoSXSSLPrivateKeys(sslPrivateKeys []ddosx.SSLPrivateKey) {
	err := Output(&OutputDDoSXSSLPrivateKeys{SSLPrivateKeys: sslPrivateKeys})
	if err != nil {
		output.Fatalf("Failed to output ssl private key: %s", err)
	}
}

func (o *OutputDDoSXSSLPrivateKeys) GetData() interface{} {
	return o.SSLPrivateKeys
}

func (o *OutputDDoSXSSLPrivateKeys) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, sslPrivateKey := range o.SSLPrivateKeys {
		fields := o.getOrderedFields(sslPrivateKey)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputDDoSXSSLPrivateKeys) getOrderedFields(sslPrivateKey ddosx.SSLPrivateKey) *output.OrderedFields {
	fields := output.NewOrderedFields()

	fields.Set("key", output.NewFieldValue(sslPrivateKey.Key, true))
	return fields
}

// OutputDDoSXACLIPRules implements OutputDataProvider for outputting an array of ACLIPRules
type OutputDDoSXACLIPRules struct {
	ACLIPRules []ddosx.ACLIPRule
}

func outputDDoSXACLIPRules(rules []ddosx.ACLIPRule) {
	err := Output(&OutputDDoSXACLIPRules{ACLIPRules: rules})
	if err != nil {
		output.Fatalf("Failed to output domain ACL IP rules: %s", err)
	}
}

func (o *OutputDDoSXACLIPRules) GetData() interface{} {
	return o.ACLIPRules
}

func (o *OutputDDoSXACLIPRules) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, rule := range o.ACLIPRules {
		fields := o.getOrderedFields(rule)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputDDoSXACLIPRules) getOrderedFields(rule ddosx.ACLIPRule) *output.OrderedFields {
	fields := output.NewOrderedFields()

	fields.Set("id", output.NewFieldValue(rule.ID, true))
	fields.Set("mode", output.NewFieldValue(rule.Mode.String(), true))
	fields.Set("ip", output.NewFieldValue(rule.IP.String(), true))
	fields.Set("uri", output.NewFieldValue(rule.URI, true))

	return fields
}

// OutputDDoSXACLGeoIPRules implements OutputDataProvider for outputting an array of ACLGeoIPRules
type OutputDDoSXACLGeoIPRules struct {
	ACLGeoIPRules []ddosx.ACLGeoIPRule
}

func outputDDoSXACLGeoIPRules(rules []ddosx.ACLGeoIPRule) {
	err := Output(&OutputDDoSXACLGeoIPRules{ACLGeoIPRules: rules})
	if err != nil {
		output.Fatalf("Failed to output domain ACL GeoIP rules: %s", err)
	}
}

func (o *OutputDDoSXACLGeoIPRules) GetData() interface{} {
	return o.ACLGeoIPRules
}

func (o *OutputDDoSXACLGeoIPRules) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, rule := range o.ACLGeoIPRules {
		fields := o.getOrderedFields(rule)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputDDoSXACLGeoIPRules) getOrderedFields(rule ddosx.ACLGeoIPRule) *output.OrderedFields {
	fields := output.NewOrderedFields()

	fields.Set("id", output.NewFieldValue(rule.ID, true))
	fields.Set("name", output.NewFieldValue(rule.Name, true))
	fields.Set("code", output.NewFieldValue(rule.Code, true))

	return fields
}

// OutputDDoSXACLGeoIPRulesModes implements OutputDataProvider for outputting an array of ACLGeoIPRulesModes
type OutputDDoSXACLGeoIPRulesModes struct {
	ACLGeoIPRulesModes []ddosx.ACLGeoIPRulesMode
}

func outputDDoSXACLGeoIPRulesModes(modes []ddosx.ACLGeoIPRulesMode) {
	err := Output(&OutputDDoSXACLGeoIPRulesModes{ACLGeoIPRulesModes: modes})
	if err != nil {
		output.Fatalf("Failed to output domain ACL GeoIP rules modes: %s", err)
	}
}

func (o *OutputDDoSXACLGeoIPRulesModes) GetData() interface{} {
	return o.ACLGeoIPRulesModes
}

func (o *OutputDDoSXACLGeoIPRulesModes) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, mode := range o.ACLGeoIPRulesModes {
		fields := o.getOrderedFields(mode)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputDDoSXACLGeoIPRulesModes) getOrderedFields(mode ddosx.ACLGeoIPRulesMode) *output.OrderedFields {
	fields := output.NewOrderedFields()

	fields.Set("mode", output.NewFieldValue(mode.String(), true))

	return fields
}

// OutputDDoSXDomainProperties implements OutputDataProvider for outputting an array of ACLGeoIPRulesModes
type OutputDDoSXDomainProperties struct {
	DomainProperties []ddosx.DomainProperty
}

func outputDDoSXDomainProperties(properties []ddosx.DomainProperty) {
	err := Output(&OutputDDoSXDomainProperties{DomainProperties: properties})
	if err != nil {
		output.Fatalf("Failed to output domain properties: %s", err)
	}
}

func (o *OutputDDoSXDomainProperties) GetData() interface{} {
	return o.DomainProperties
}

func (o *OutputDDoSXDomainProperties) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, property := range o.DomainProperties {
		fields := o.getOrderedFields(property)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputDDoSXDomainProperties) getOrderedFields(property ddosx.DomainProperty) *output.OrderedFields {
	fields := output.NewOrderedFields()

	fields.Set("id", output.NewFieldValue(property.ID, true))
	fields.Set("name", output.NewFieldValue(property.Name.String(), true))
	fields.Set("value", output.NewFieldValue(fmt.Sprintf("%v", property.Value), true))

	return fields
}

type OutputDDoSXDomainVerificationFilesFile struct {
	Name    string
	Content string
}

// OutputDDoSXDomainVerificationFiles implements OutputDataProvider for outputting an array of OutputDDoSXDomainVerificationFilesFile
type OutputDDoSXDomainVerificationFiles struct {
	DomainVerificationFiles []OutputDDoSXDomainVerificationFilesFile
}

func outputDDoSXDomainVerificationFiles(files []OutputDDoSXDomainVerificationFilesFile) {
	err := Output(&OutputDDoSXDomainVerificationFiles{DomainVerificationFiles: files})
	if err != nil {
		output.Fatalf("Failed to output domain verification files: %s", err)
	}
}

func (o *OutputDDoSXDomainVerificationFiles) GetData() interface{} {
	return o.DomainVerificationFiles
}

func (o *OutputDDoSXDomainVerificationFiles) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, file := range o.DomainVerificationFiles {
		fields := o.getOrderedFields(file)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputDDoSXDomainVerificationFiles) getOrderedFields(file OutputDDoSXDomainVerificationFilesFile) *output.OrderedFields {
	fields := output.NewOrderedFields()

	fields.Set("name", output.NewFieldValue(file.Name, true))
	fields.Set("content", output.NewFieldValue(file.Content, true))

	return fields
}

// OutputDDoSXCDNRules implements OutputDataProvider for outputting an array of CDNRules
type OutputDDoSXCDNRules struct {
	CDNRules []ddosx.CDNRule
}

func outputDDoSXCDNRules(rules []ddosx.CDNRule) {
	err := Output(&OutputDDoSXCDNRules{CDNRules: rules})
	if err != nil {
		output.Fatalf("Failed to output domain ACL GeoIP rules: %s", err)
	}
}

func (o *OutputDDoSXCDNRules) GetData() interface{} {
	return o.CDNRules
}

func (o *OutputDDoSXCDNRules) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, rule := range o.CDNRules {
		fields := o.getOrderedFields(rule)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputDDoSXCDNRules) getOrderedFields(rule ddosx.CDNRule) *output.OrderedFields {
	fields := output.NewOrderedFields()

	fields.Set("id", output.NewFieldValue(rule.ID, true))
	fields.Set("uri", output.NewFieldValue(rule.URI, true))
	fields.Set("cache_control", output.NewFieldValue(rule.CacheControl.String(), true))
	fields.Set("cache_control_duration", output.NewFieldValue(rule.CacheControlDuration.String(), true))
	fields.Set("mime_types", output.NewFieldValue(strings.Join(rule.MimeTypes, ", "), true))
	fields.Set("type", output.NewFieldValue(rule.Type.String(), true))

	return fields
}

// OutputDDoSXHSTSConfiguration implements OutputDataProvider for outputting an array of HSTSConfiguration
type OutputDDoSXHSTSConfiguration struct {
	HSTSConfiguration []ddosx.HSTSConfiguration
}

func outputDDoSXHSTSConfiguration(configurations []ddosx.HSTSConfiguration) {
	err := Output(&OutputDDoSXHSTSConfiguration{HSTSConfiguration: configurations})
	if err != nil {
		output.Fatalf("Failed to output domain HSTS configurations: %s", err)
	}
}

func (o *OutputDDoSXHSTSConfiguration) GetData() interface{} {
	return o.HSTSConfiguration
}

func (o *OutputDDoSXHSTSConfiguration) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, configuration := range o.HSTSConfiguration {
		fields := o.getOrderedFields(configuration)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputDDoSXHSTSConfiguration) getOrderedFields(configuration ddosx.HSTSConfiguration) *output.OrderedFields {
	fields := output.NewOrderedFields()

	fields.Set("enabled", output.NewFieldValue(strconv.FormatBool(configuration.Enabled), true))

	return fields
}

// OutputDDoSXHSTSRules implements OutputDataProvider for outputting an array of HSTSRules
type OutputDDoSXHSTSRules struct {
	HSTSRules []ddosx.HSTSRule
}

func outputDDoSXHSTSRules(rules []ddosx.HSTSRule) {
	err := Output(&OutputDDoSXHSTSRules{HSTSRules: rules})
	if err != nil {
		output.Fatalf("Failed to output domain HSTS rules: %s", err)
	}
}

func (o *OutputDDoSXHSTSRules) GetData() interface{} {
	return o.HSTSRules
}

func (o *OutputDDoSXHSTSRules) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	for _, rule := range o.HSTSRules {
		fields := o.getOrderedFields(rule)
		data = append(data, fields)
	}

	return data, nil
}

func (o *OutputDDoSXHSTSRules) getOrderedFields(rule ddosx.HSTSRule) *output.OrderedFields {
	fields := output.NewOrderedFields()

	recordName := ""
	if rule.RecordName != nil {
		recordName = *rule.RecordName
	}

	fields.Set("id", output.NewFieldValue(rule.ID, true))
	fields.Set("max_age", output.NewFieldValue(strconv.Itoa(rule.MaxAge), true))
	fields.Set("preload", output.NewFieldValue(strconv.FormatBool(rule.Preload), true))
	fields.Set("include_subdomains", output.NewFieldValue(strconv.FormatBool(rule.IncludeSubdomains), true))
	fields.Set("rule_type", output.NewFieldValue(rule.RuleType.String(), true))
	fields.Set("record_name", output.NewFieldValue(recordName, true))

	return fields
}
