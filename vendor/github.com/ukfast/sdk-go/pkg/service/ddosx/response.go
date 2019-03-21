package ddosx

import "github.com/ukfast/sdk-go/pkg/connection"

// GetDomainsResponseBody represents an API response body containing []Domain data
type GetDomainsResponseBody struct {
	connection.APIResponseBody

	Data []Domain `json:"data"`
}

// GetDomainResponseBody represents an API response body containing Domain data
type GetDomainResponseBody struct {
	connection.APIResponseBody

	Data Domain `json:"data"`
}

// GetRecordsResponseBody represents an API response body containing []Record data
type GetRecordsResponseBody struct {
	connection.APIResponseBody

	Data []Record `json:"data"`
}

// GetRecordResponseBody represents an API response body containing Record data
type GetRecordResponseBody struct {
	connection.APIResponseBody

	Data Record `json:"data"`
}

// GetDomainPropertiesResponseBody represents an API response body containing []DomainProperty data
type GetDomainPropertiesResponseBody struct {
	connection.APIResponseBody

	Data []DomainProperty `json:"data"`
}

// GetDomainPropertyResponseBody represents an API response body containing DomainProperty data
type GetDomainPropertyResponseBody struct {
	connection.APIResponseBody

	Data DomainProperty `json:"data"`
}

// GetWAFResponseBody represents an API response body containing WAF data
type GetWAFResponseBody struct {
	connection.APIResponseBody

	Data WAF `json:"data"`
}

// GetWAFRuleSetsResponseBody represents an API response body containing []WAFRuleSet data
type GetWAFRuleSetsResponseBody struct {
	connection.APIResponseBody

	Data []WAFRuleSet `json:"data"`
}

// GetWAFRuleSetResponseBody represents an API response body containing WAFRuleSet data
type GetWAFRuleSetResponseBody struct {
	connection.APIResponseBody

	Data WAFRuleSet `json:"data"`
}

// GetWAFRulesResponseBody represents an API response body containing []WAFRule
type GetWAFRulesResponseBody struct {
	connection.APIResponseBody

	Data []WAFRule `json:"data"`
}

// GetWAFRuleResponseBody represents an API response body containing WAFRule
type GetWAFRuleResponseBody struct {
	connection.APIResponseBody

	Data WAFRule `json:"data"`
}

// GetWAFAdvancedRulesResponseBody represents an API response body containing []WAFAdvancedRule
type GetWAFAdvancedRulesResponseBody struct {
	connection.APIResponseBody

	Data []WAFAdvancedRule `json:"data"`
}

// GetWAFAdvancedRuleResponseBody represents an API response body containing WAFAdvancedRule
type GetWAFAdvancedRuleResponseBody struct {
	connection.APIResponseBody

	Data WAFAdvancedRule `json:"data"`
}

// GetSSLsResponseBody represents an API response body containing []SSL data
type GetSSLsResponseBody struct {
	connection.APIResponseBody

	Data []SSL `json:"data"`
}

// GetSSLResponseBody represents an API response body containing SSL data
type GetSSLResponseBody struct {
	connection.APIResponseBody

	Data SSL `json:"data"`
}

// GetSSLContentResponseBody represents an API response body containing SSLContent data
type GetSSLContentResponseBody struct {
	connection.APIResponseBody

	Data SSLContent `json:"data"`
}

// GetSSLPrivateKeyResponseBody represents an API response body containing SSLPrivateKey data
type GetSSLPrivateKeyResponseBody struct {
	connection.APIResponseBody

	Data SSLPrivateKey `json:"data"`
}

// GetACLGeoIPRulesResponseBody represents an API response body containing []ACLGeoIPRule data
type GetACLGeoIPRulesResponseBody struct {
	connection.APIResponseBody

	Data []ACLGeoIPRule `json:"data"`
}

// GetACLGeoIPRuleResponseBody represents an API response body containing ACLGeoIPRule data
type GetACLGeoIPRuleResponseBody struct {
	connection.APIResponseBody

	Data ACLGeoIPRule `json:"data"`
}

// GetACLIPRulesResponseBody represents an API response body containing []ACLIPRule data
type GetACLIPRulesResponseBody struct {
	connection.APIResponseBody

	Data []ACLIPRule `json:"data"`
}

// GetACLIPRuleResponseBody represents an API response body containing ACLIPRule data
type GetACLIPRuleResponseBody struct {
	connection.APIResponseBody

	Data ACLIPRule `json:"data"`
}

// GetACLGeoIPRulesModeResponseBody represents an API response body containing ACLGeoIPRulesMode data
type GetACLGeoIPRulesModeResponseBody struct {
	connection.APIResponseBody

	Data struct {
		Mode ACLGeoIPRulesMode `json:"mode"`
	} `json:"data"`
}
