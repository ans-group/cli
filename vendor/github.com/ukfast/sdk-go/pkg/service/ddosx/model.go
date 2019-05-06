package ddosx

import (
	"errors"
	"strings"
	"time"

	"github.com/ukfast/go-durationstring"
	"github.com/ukfast/sdk-go/pkg/connection"
)

type DomainStatus string

func (s DomainStatus) String() string {
	return string(s)
}

const (
	DomainStatusConfigured    DomainStatus = "Configured"
	DomainStatusNotConfigured DomainStatus = "Not Configured"
	DomainStatusPending       DomainStatus = "Pending"
	DomainStatusFailed        DomainStatus = "Failed"
	DomainStatusCancelling    DomainStatus = "Cancelling"
	DomainStatusCancelled     DomainStatus = "Cancelled"
)

type DomainPropertyName string

func (e DomainPropertyName) String() string {
	return string(e)
}

const (
	DomainPropertyNameClientMaxBodySize DomainPropertyName = "client_max_body_size"
	DomainPropertyNameProxyTimeout      DomainPropertyName = "proxy_timeout"
	DomainPropertyNameIPv6Enabled       DomainPropertyName = "ipv6_enabled"
	DomainPropertyNameSecureOrigin      DomainPropertyName = "secure_origin"
)

var DomainPropertyNameEnum = []connection.Enum{
	DomainPropertyNameClientMaxBodySize,
	DomainPropertyNameProxyTimeout,
	DomainPropertyNameIPv6Enabled,
	DomainPropertyNameSecureOrigin,
}

// ParseDomainPropertyName attempts to parse a DomainPropertyName from string
func ParseDomainPropertyName(s string) (DomainPropertyName, error) {
	e, err := connection.ParseEnum(s, DomainPropertyNameEnum)
	if e != nil {
		return e.(DomainPropertyName), err
	}

	return "", err
}

type RecordType string

func (s RecordType) String() string {
	return string(s)
}

const (
	RecordTypeA    RecordType = "A"
	RecordTypeAAAA RecordType = "AAAA"
)

type WAFMode string

func (s WAFMode) String() string {
	return string(s)
}

const (
	WAFModeOn            WAFMode = "On"
	WAFModeOff           WAFMode = "Off"
	WAFModeDetectionOnly WAFMode = "DetectionOnly"
)

// ParseWAFMode attempts to parse a WAFMode from string
func ParseWAFMode(s string) (WAFMode, error) {
	switch strings.ToUpper(s) {
	case "ON":
		return WAFModeOn, nil
	case "OFF":
		return WAFModeOff, nil
	case "DETECTIONONLY":
		return WAFModeDetectionOnly, nil
	}

	return "", errors.New("Invalid WAF mode")
}

type WAFParanoiaLevel string

func (s WAFParanoiaLevel) String() string {
	return string(s)
}

const (
	WAFParanoiaLevelLow     WAFParanoiaLevel = "Low"
	WAFParanoiaLevelMedium  WAFParanoiaLevel = "Medium"
	WAFParanoiaLevelHigh    WAFParanoiaLevel = "High"
	WAFParanoiaLevelHighest WAFParanoiaLevel = "Highest"
)

// ParseWAFParanoiaLevel attempts to parse a WAFMode from string
func ParseWAFParanoiaLevel(s string) (WAFParanoiaLevel, error) {
	switch strings.ToUpper(s) {
	case "LOW":
		return WAFParanoiaLevelLow, nil
	case "MEDIUM":
		return WAFParanoiaLevelMedium, nil
	case "HIGH":
		return WAFParanoiaLevelHigh, nil
	case "HIGHEST":
		return WAFParanoiaLevelHighest, nil
	}

	return "", errors.New("Invalid WAF paranoia level")
}

type WAFRuleSetName string

func (s WAFRuleSetName) String() string {
	return string(s)
}

const (
	WAFRuleSetNameIPRepution                             WAFRuleSetName = "IP Reputation"
	WAFRuleSetNameMethodEnforcement                      WAFRuleSetName = "Method Enforcement"
	WAFRuleSetNameScannerDetection                       WAFRuleSetName = "Scanner Detection"
	WAFRuleSetNameProtocolEnforcement                    WAFRuleSetName = "Protocol Enforcement"
	WAFRuleSetNameProtocolAttack                         WAFRuleSetName = "Protocol Attack"
	WAFRuleSetNameApplicationAttackLocalFileInclusion    WAFRuleSetName = "Application Attack (Local File Inclusion)"
	WAFRuleSetNameApplicationAttackRemoteFileInclusion   WAFRuleSetName = "Application Attack (Remote File Inclusion)"
	WAFRuleSetNameApplicationAttackRemoteCodeExecution   WAFRuleSetName = "Application Attack (Remote Code Execution)"
	WAFRuleSetNameApplicationAttackPHP                   WAFRuleSetName = "Application Attack PHP"
	WAFRuleSetNameApplicationAttackXSSCrossSiteScripting WAFRuleSetName = "Application Attack XSS (Cross Site Scripting)"
	WAFRuleSetNameApplicationAttackSQLISQLInjection      WAFRuleSetName = "Application Attack SQLI (SQL Injection)"
	WAFRuleSetNameApplicationAttackSessionFixation       WAFRuleSetName = "Application Attack Session Fixation"
	WAFRuleSetNameDataDeakages                           WAFRuleSetName = "Data Leakages"
	WAFRuleSetNameDataLeakageSQL                         WAFRuleSetName = "Data Leakage SQL"
	WAFRuleSetNameDataLeakageJava                        WAFRuleSetName = "Data Leakage Java"
	WAFRuleSetNameDataLeakagePHP                         WAFRuleSetName = "Data Leakage PHP"
	WAFRuleSetNameDataLeakageIIS                         WAFRuleSetName = "Data Leakage IIS"
)

type WAFAdvancedRuleSection string

func (s WAFAdvancedRuleSection) String() string {
	return string(s)
}

const (
	WAFAdvancedRuleSectionArgs           WAFAdvancedRuleSection = "ARGS"
	WAFAdvancedRuleSectionMatchedVars    WAFAdvancedRuleSection = "MATCHED_VARS"
	WAFAdvancedRuleSectionRemoteHost     WAFAdvancedRuleSection = "REMOTE_HOST"
	WAFAdvancedRuleSectionRequestBody    WAFAdvancedRuleSection = "REQUEST_BODY"
	WAFAdvancedRuleSectionRequestCookies WAFAdvancedRuleSection = "REQUEST_COOKIES"
	WAFAdvancedRuleSectionRequestHeaders WAFAdvancedRuleSection = "REQUEST_HEADERS"
	WAFAdvancedRuleSectionRequestURI     WAFAdvancedRuleSection = "REQUEST_URI"
)

// ParseWAFAdvancedRuleSection attempts to parse a WAFAdvancedRuleSection from string
func ParseWAFAdvancedRuleSection(s string) (WAFAdvancedRuleSection, error) {
	switch strings.ToUpper(s) {
	case "ARGS":
		return WAFAdvancedRuleSectionArgs, nil
	case "MATCHED_VARS":
		return WAFAdvancedRuleSectionMatchedVars, nil
	case "REMOTE_HOST":
		return WAFAdvancedRuleSectionRemoteHost, nil
	case "REQUEST_BODY":
		return WAFAdvancedRuleSectionRequestBody, nil
	case "REQUEST_COOKIES":
		return WAFAdvancedRuleSectionRequestCookies, nil
	case "REQUEST_HEADERS":
		return WAFAdvancedRuleSectionRequestHeaders, nil
	case "REQUEST_URI":
		return WAFAdvancedRuleSectionRequestURI, nil
	}

	return "", errors.New("Invalid advanced rule section")
}

type WAFAdvancedRuleModifier string

func (s WAFAdvancedRuleModifier) String() string {
	return string(s)
}

const (
	WAFAdvancedRuleModifierBeginsWith   WAFAdvancedRuleModifier = "beginsWith"
	WAFAdvancedRuleModifierEndsWith     WAFAdvancedRuleModifier = "endsWith"
	WAFAdvancedRuleModifierContains     WAFAdvancedRuleModifier = "contains"
	WAFAdvancedRuleModifierContainsWord WAFAdvancedRuleModifier = "containsWord"
)

// ParseWAFAdvancedRuleModifier attempts to parse a WAFAdvancedRuleModifier from string
func ParseWAFAdvancedRuleModifier(s string) (WAFAdvancedRuleModifier, error) {
	switch strings.ToUpper(s) {
	case "BEGINSWITH":
		return WAFAdvancedRuleModifierBeginsWith, nil
	case "ENDSWITH":
		return WAFAdvancedRuleModifierEndsWith, nil
	case "CONTAINS":
		return WAFAdvancedRuleModifierContains, nil
	case "CONTAINSWORD":
		return WAFAdvancedRuleModifierContainsWord, nil
	}

	return "", errors.New("Invalid advanced rule modifier")
}

type ACLIPMode string

func (s ACLIPMode) String() string {
	return string(s)
}

const (
	ACLIPModeAllow ACLIPMode = "Allow"
	ACLIPModeDeny  ACLIPMode = "Deny"
)

// ParseACLIPMode attempts to parse a ACLIPMode from string
func ParseACLIPMode(s string) (ACLIPMode, error) {
	switch strings.ToUpper(s) {
	case "ALLOW":
		return ACLIPModeAllow, nil
	case "DENY":
		return ACLIPModeDeny, nil
	}

	return "", errors.New("Invalid ACL IP mode")

}

type ACLGeoIPRulesMode string

func (s ACLGeoIPRulesMode) String() string {
	return string(s)
}

const (
	ACLGeoIPRulesModeWhitelist ACLGeoIPRulesMode = "Whitelist"
	ACLGeoIPRulesModeBlacklist ACLGeoIPRulesMode = "Blacklist"
)

// ParseACLGeoIPRulesMode attempts to parse a ACLGeoIPRulesMode from string
func ParseACLGeoIPRulesMode(s string) (ACLGeoIPRulesMode, error) {
	switch strings.ToUpper(s) {
	case "WHITELIST":
		return ACLGeoIPRulesModeWhitelist, nil
	case "BLACKLIST":
		return ACLGeoIPRulesModeBlacklist, nil
	}

	return "", errors.New("Invalid ACL GeoIP rules filtering mode")
}

type CDNRuleCacheControl string

func (e CDNRuleCacheControl) String() string {
	return string(e)
}

const (
	CDNRuleCacheControlCustom CDNRuleCacheControl = "Custom"
	CDNRuleCacheControlOrigin CDNRuleCacheControl = "Origin"
)

var CDNRuleCacheControlEnum = []connection.Enum{CDNRuleCacheControlCustom, CDNRuleCacheControlOrigin}

// ParseCDNRuleCacheControl attempts to parse a CDNRuleCacheControl from string
func ParseCDNRuleCacheControl(s string) (CDNRuleCacheControl, error) {
	e, err := connection.ParseEnum(s, CDNRuleCacheControlEnum)
	if e != nil {
		return e.(CDNRuleCacheControl), err
	}

	return "", err
}

type CDNRuleType string

func (e CDNRuleType) String() string {
	return string(e)
}

const (
	CDNRuleTypeGlobal CDNRuleType = "global"
	CDNRuleTypePerURI CDNRuleType = "per-uri"
)

var CDNRuleTypeEnum = []connection.Enum{CDNRuleTypeGlobal, CDNRuleTypePerURI}

// ParseCDNRuleType attempts to parse a CDNRuleType from string
func ParseCDNRuleType(s string) (CDNRuleType, error) {
	e, err := connection.ParseEnum(s, CDNRuleTypeEnum)
	if e != nil {
		return e.(CDNRuleType), err
	}

	return "", err
}

type HSTSRuleType string

func (e HSTSRuleType) String() string {
	return string(e)
}

const (
	HSTSRuleTypeDomain HSTSRuleType = "domain"
	HSTSRuleTypeRecord HSTSRuleType = "record"
)

var HSTSRuleTypeEnum = []connection.Enum{HSTSRuleTypeDomain, HSTSRuleTypeRecord}

// ParseHSTSRuleType attempts to parse a HSTSRuleType from string
func ParseHSTSRuleType(s string) (HSTSRuleType, error) {
	e, err := connection.ParseEnum(s, HSTSRuleTypeEnum)
	if e != nil {
		return e.(HSTSRuleType), err
	}

	return "", err
}

// Domain represents a DDoSX domain
type Domain struct {
	SafeDNSZoneID *int               `json:"safedns_zone_id"`
	Name          string             `json:"name"`
	Status        DomainStatus       `json:"status"`
	DNSActive     bool               `json:"dns_active"`
	CDNActive     bool               `json:"cdn_active"`
	WAFActive     bool               `json:"waf_active"`
	ExternalDNS   *DomainExternalDNS `json:"external_dns"`
}

// DomainExternalDNS represents a DDoSX domain external DNS configuration
type DomainExternalDNS struct {
	Verified           bool   `json:"verified"`
	VerificationString string `json:"verification_string"`
	Target             string `json:"target"`
}

// DomainProperty represents a DDoSX domain property
type DomainProperty struct {
	ID    string             `json:"id"`
	Name  DomainPropertyName `json:"name"`
	Value interface{}        `json:"value"`
}

// Record represents a DDoSX record
type Record struct {
	ID              string     `json:"id"`
	DomainName      string     `json:"domain_name"`
	SafeDNSRecordID *int       `json:"safedns_record_id"`
	SSLID           *string    `json:"ssl_id"`
	Name            string     `json:"name"`
	Type            RecordType `json:"type"`
	Content         string     `json:"content"`
}

// WAF represents a DDoSX WAF configuration
type WAF struct {
	Mode          WAFMode          `json:"mode"`
	ParanoiaLevel WAFParanoiaLevel `json:"paranoia_level"`
}

// WAFRuleSet represents a DDoSX WAF rule set
type WAFRuleSet struct {
	ID     string         `json:"id"`
	Name   WAFRuleSetName `json:"name"`
	Active bool           `json:"active"`
}

// WAFRule represents a DDoSX WAF rule
type WAFRule struct {
	ID  string               `json:"id"`
	URI string               `json:"uri"`
	IP  connection.IPAddress `json:"ip"`
}

// WAFAdvancedRule represents a DDoSX WAF advanced rule
type WAFAdvancedRule struct {
	ID       string                  `json:"id"`
	Section  WAFAdvancedRuleSection  `json:"section"`
	Modifier WAFAdvancedRuleModifier `json:"modifier"`
	Phrase   string                  `json:"phrase"`
	IP       connection.IPAddress    `json:"ip"`
}

// SSL represents a DDoSX SSL
type SSL struct {
	ID           string   `json:"id"`
	UKFastSSLID  *int     `json:"ukfast_ssl_id"`
	Domains      []string `json:"domains"`
	FriendlyName string   `json:"friendly_name"`
}

// SSLContent represents a DDoSX SSL content
type SSLContent struct {
	Certificate string `json:"certificate"`
	CABundle    string `json:"ca_bundle"`
}

// SSLPrivateKey represents a DDoSX SSL private key
type SSLPrivateKey struct {
	Key string `json:"key"`
}

// ACLGeoIPRule represents a DDoSX ACL GeoIP rule
type ACLGeoIPRule struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// ACLIPRule represents a DDoSX ACL IP rule
type ACLIPRule struct {
	ID   string               `json:"id"`
	IP   connection.IPAddress `json:"ip"`
	URI  string               `json:"uri"`
	Mode ACLIPMode            `json:"mode"`
}

// CDNRule represents a DDoSX CDN rule
type CDNRule struct {
	ID           string              `json:"id"`
	URI          string              `json:"uri"`
	CacheControl CDNRuleCacheControl `json:"cache_control"`
	// CacheControlDuration specifies the cache control duration. May be nil if duration not applicable
	CacheControlDuration CDNRuleCacheControlDuration `json:"cache_control_duration"`
	MimeTypes            []string                    `json:"mime_types"`
	Type                 CDNRuleType                 `json:"type"`
}

// CDNRuleCacheControlDuration represents a DDoSX CDN rule duration
type CDNRuleCacheControlDuration struct {
	Years   int `json:"years"`
	Months  int `json:"months"`
	Days    int `json:"days"`
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
}

// Duration returns the cache control duration as time.Duration
func (d *CDNRuleCacheControlDuration) Duration() time.Duration {
	day := time.Hour * 24
	td := time.Duration(d.Years) * day * 365
	td = td + time.Duration(d.Months)*day*30
	td = td + time.Duration(d.Days)*day
	td = td + time.Duration(d.Hours)*time.Hour
	return td + time.Duration(d.Minutes)*time.Minute
}

func (d *CDNRuleCacheControlDuration) String() string {
	return durationstring.String(d.Years, d.Months, d.Days, d.Hours, d.Minutes, 0, 0, 0, 0)
}

// ParseCDNRuleCacheControlDuration parses string s and returns a pointer to an
// initialised CDNRuleCacheControlDuration
func ParseCDNRuleCacheControlDuration(s string) (*CDNRuleCacheControlDuration, error) {
	years, months, days, hours, minutes, _, _, _, _, err := durationstring.Parse(s)
	if err != nil {
		return nil, err
	}

	return &CDNRuleCacheControlDuration{
		Years:   years,
		Months:  months,
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
	}, nil
}

// HSTSConfiguration represents HSTS configuration for a DDoSX domain
type HSTSConfiguration struct {
	Enabled bool `json:"enabled"`
}

// HSTSRule represents HSTS rule for a DDoSX domain
type HSTSRule struct {
	ID                string       `json:"id"`
	MaxAge            int          `json:"max_age"`
	Preload           bool         `json:"preload"`
	IncludeSubdomains bool         `json:"include_subdomains"`
	Type              HSTSRuleType `json:"type"`
	RecordName        *string      `json:"record_name"`
}
