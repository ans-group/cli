package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ukfast/sdk-go/pkg/ptr"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func TestOutputDDoSXDomains_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputDDoSXDomains{
			Domains: []ddosx.Domain{
				ddosx.Domain{
					Name: "testdomain.com",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ddosx.Domain{}, data)
		assert.Equal(t, "testdomain.com", data.([]ddosx.Domain)[0].Name)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputDDoSXDomains{
			Domains: []ddosx.Domain{
				ddosx.Domain{
					Name: "testdomain1.com",
				},
				ddosx.Domain{
					Name: "testdomain2.com",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ddosx.Domain{}, data)
		assert.Len(t, data, 2)
		assert.Equal(t, "testdomain1.com", data.([]ddosx.Domain)[0].Name)
		assert.Equal(t, "testdomain2.com", data.([]ddosx.Domain)[1].Name)
	})
}

func TestOutputDDoSXDomains_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputDDoSXDomains{
			Domains: []ddosx.Domain{
				ddosx.Domain{
					Name: "testdomain.com",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testdomain.com", data[0].Get("name").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputDDoSXDomains{
			Domains: []ddosx.Domain{
				ddosx.Domain{
					Name: "testdomain1.com",
				},
				ddosx.Domain{
					Name: "testdomain2.com",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testdomain1.com", data[0].Get("name").Value)
		assert.True(t, data[1].Exists("name"))
		assert.Equal(t, "testdomain2.com", data[1].Get("name").Value)
	})
}

func TestOutputDDoSXRecords_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputDDoSXRecords{
			Records: []ddosx.Record{
				ddosx.Record{
					Name: "testdomain.com",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ddosx.Record{}, data)
		assert.Equal(t, "testdomain.com", data.([]ddosx.Record)[0].Name)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputDDoSXRecords{
			Records: []ddosx.Record{
				ddosx.Record{
					Name: "testdomain1.com",
				},
				ddosx.Record{
					Name: "testdomain2.com",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ddosx.Record{}, data)
		assert.Len(t, data, 2)
		assert.Equal(t, "testdomain1.com", data.([]ddosx.Record)[0].Name)
		assert.Equal(t, "testdomain2.com", data.([]ddosx.Record)[1].Name)
	})
}

func TestOutputDDoSXRecords_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputDDoSXRecords{
			Records: []ddosx.Record{
				ddosx.Record{
					Name: "testdomain.com",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testdomain.com", data[0].Get("name").Value)
	})

	t.Run("SingleWithSafeDNSRecordID_ExpectedFieldData", func(t *testing.T) {
		o := OutputDDoSXRecords{
			Records: []ddosx.Record{
				ddosx.Record{
					Name:            "testdomain.com",
					SafeDNSRecordID: ptr.Int(123),
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testdomain.com", data[0].Get("name").Value)
		assert.Equal(t, "123", data[0].Get("safedns_record_id").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputDDoSXRecords{
			Records: []ddosx.Record{
				ddosx.Record{
					Name: "testdomain1.com",
				},
				ddosx.Record{
					Name: "testdomain2.com",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testdomain1.com", data[0].Get("name").Value)
		assert.True(t, data[1].Exists("name"))
		assert.Equal(t, "testdomain2.com", data[1].Get("name").Value)
	})
}

func TestOutputDDoSXWAFs_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputDDoSXWAFs{
			WAFs: []ddosx.WAF{
				ddosx.WAF{
					Mode: ddosx.WAFModeOn,
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ddosx.WAF{}, data)
		assert.Equal(t, ddosx.WAFModeOn, data.([]ddosx.WAF)[0].Mode)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputDDoSXWAFs{
			WAFs: []ddosx.WAF{
				ddosx.WAF{
					Mode: ddosx.WAFModeOn,
				},
				ddosx.WAF{
					Mode: ddosx.WAFModeOff,
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ddosx.WAF{}, data)
		assert.Len(t, data, 2)
		assert.Equal(t, ddosx.WAFModeOn, data.([]ddosx.WAF)[0].Mode)
		assert.Equal(t, ddosx.WAFModeOff, data.([]ddosx.WAF)[1].Mode)
	})
}

func TestOutputDDoSXWAFs_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputDDoSXWAFs{
			WAFs: []ddosx.WAF{
				ddosx.WAF{
					Mode: ddosx.WAFModeOn,
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("mode"))
		assert.Equal(t, "On", data[0].Get("mode").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputDDoSXWAFs{
			WAFs: []ddosx.WAF{
				ddosx.WAF{
					Mode: ddosx.WAFModeOn,
				},
				ddosx.WAF{
					Mode: ddosx.WAFModeOff,
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("mode"))
		assert.Equal(t, "On", data[0].Get("mode").Value)
		assert.True(t, data[1].Exists("mode"))
		assert.Equal(t, "Off", data[1].Get("mode").Value)
	})
}

func TestOutputDDoSXWAFRuleSets_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputDDoSXWAFRuleSets{
			WAFRuleSets: []ddosx.WAFRuleSet{
				ddosx.WAFRuleSet{
					ID: "00000000-0000-0000-0000-000000000000",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ddosx.WAFRuleSet{}, data)
		assert.Equal(t, "00000000-0000-0000-0000-000000000000", data.([]ddosx.WAFRuleSet)[0].ID)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputDDoSXWAFRuleSets{
			WAFRuleSets: []ddosx.WAFRuleSet{
				ddosx.WAFRuleSet{
					ID: "00000000-0000-0000-0000-000000000000",
				},
				ddosx.WAFRuleSet{
					ID: "00000000-0000-0000-0000-000000000001",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ddosx.WAFRuleSet{}, data)
		assert.Len(t, data, 2)
		assert.Equal(t, "00000000-0000-0000-0000-000000000000", data.([]ddosx.WAFRuleSet)[0].ID)
		assert.Equal(t, "00000000-0000-0000-0000-000000000001", data.([]ddosx.WAFRuleSet)[1].ID)
	})
}

func TestOutputDDoSXWAFRuleSets_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputDDoSXWAFRuleSets{
			WAFRuleSets: []ddosx.WAFRuleSet{
				ddosx.WAFRuleSet{
					ID: "00000000-0000-0000-0000-000000000000",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("id"))
		assert.Equal(t, "00000000-0000-0000-0000-000000000000", data[0].Get("id").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputDDoSXWAFRuleSets{
			WAFRuleSets: []ddosx.WAFRuleSet{
				ddosx.WAFRuleSet{
					ID: "00000000-0000-0000-0000-000000000000",
				},
				ddosx.WAFRuleSet{
					ID: "00000000-0000-0000-0000-000000000001",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("id"))
		assert.Equal(t, "00000000-0000-0000-0000-000000000000", data[0].Get("id").Value)
		assert.True(t, data[1].Exists("id"))
		assert.Equal(t, "00000000-0000-0000-0000-000000000001", data[1].Get("id").Value)
	})
}

func TestOutputDDoSXSSLs_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputDDoSXSSLs{
			SSLs: []ddosx.SSL{
				ddosx.SSL{
					ID: "00000000-0000-0000-0000-000000000000",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ddosx.SSL{}, data)
		assert.Equal(t, "00000000-0000-0000-0000-000000000000", data.([]ddosx.SSL)[0].ID)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputDDoSXSSLs{
			SSLs: []ddosx.SSL{
				ddosx.SSL{
					ID: "00000000-0000-0000-0000-000000000000",
				},
				ddosx.SSL{
					ID: "00000000-0000-0000-0000-000000000001",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ddosx.SSL{}, data)
		assert.Len(t, data, 2)
		assert.Equal(t, "00000000-0000-0000-0000-000000000000", data.([]ddosx.SSL)[0].ID)
		assert.Equal(t, "00000000-0000-0000-0000-000000000001", data.([]ddosx.SSL)[1].ID)
	})
}

func TestOutputDDoSXSSLs_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputDDoSXSSLs{
			SSLs: []ddosx.SSL{
				ddosx.SSL{
					ID: "00000000-0000-0000-0000-000000000000",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("id"))
		assert.Equal(t, "00000000-0000-0000-0000-000000000000", data[0].Get("id").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputDDoSXSSLs{
			SSLs: []ddosx.SSL{
				ddosx.SSL{
					ID: "00000000-0000-0000-0000-000000000000",
				},
				ddosx.SSL{
					ID: "00000000-0000-0000-0000-000000000001",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("id"))
		assert.Equal(t, "00000000-0000-0000-0000-000000000000", data[0].Get("id").Value)
		assert.True(t, data[1].Exists("id"))
		assert.Equal(t, "00000000-0000-0000-0000-000000000001", data[1].Get("id").Value)
	})
}

func TestOutputDDoSXACLIPRules_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputDDoSXACLIPRules{
			ACLIPRules: []ddosx.ACLIPRule{
				ddosx.ACLIPRule{
					ID: "00000000-0000-0000-0000-000000000000",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ddosx.ACLIPRule{}, data)
		assert.Equal(t, "00000000-0000-0000-0000-000000000000", data.([]ddosx.ACLIPRule)[0].ID)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputDDoSXACLIPRules{
			ACLIPRules: []ddosx.ACLIPRule{
				ddosx.ACLIPRule{
					ID: "00000000-0000-0000-0000-000000000000",
				},
				ddosx.ACLIPRule{
					ID: "00000000-0000-0000-0000-000000000001",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ddosx.ACLIPRule{}, data)
		assert.Len(t, data, 2)
		assert.Equal(t, "00000000-0000-0000-0000-000000000000", data.([]ddosx.ACLIPRule)[0].ID)
		assert.Equal(t, "00000000-0000-0000-0000-000000000001", data.([]ddosx.ACLIPRule)[1].ID)
	})
}

func TestOutputDDoSXACLIPRules_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputDDoSXACLIPRules{
			ACLIPRules: []ddosx.ACLIPRule{
				ddosx.ACLIPRule{
					ID: "00000000-0000-0000-0000-000000000000",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("id"))
		assert.Equal(t, "00000000-0000-0000-0000-000000000000", data[0].Get("id").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputDDoSXACLIPRules{
			ACLIPRules: []ddosx.ACLIPRule{
				ddosx.ACLIPRule{
					ID: "00000000-0000-0000-0000-000000000000",
				},
				ddosx.ACLIPRule{
					ID: "00000000-0000-0000-0000-000000000001",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("id"))
		assert.Equal(t, "00000000-0000-0000-0000-000000000000", data[0].Get("id").Value)
		assert.True(t, data[1].Exists("id"))
		assert.Equal(t, "00000000-0000-0000-0000-000000000001", data[1].Get("id").Value)
	})
}

func TestOutputDDoSXACLGeoIPRules_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputDDoSXACLGeoIPRules{
			ACLGeoIPRules: []ddosx.ACLGeoIPRule{
				ddosx.ACLGeoIPRule{
					ID: "00000000-0000-0000-0000-000000000000",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ddosx.ACLGeoIPRule{}, data)
		assert.Equal(t, "00000000-0000-0000-0000-000000000000", data.([]ddosx.ACLGeoIPRule)[0].ID)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputDDoSXACLGeoIPRules{
			ACLGeoIPRules: []ddosx.ACLGeoIPRule{
				ddosx.ACLGeoIPRule{
					ID: "00000000-0000-0000-0000-000000000000",
				},
				ddosx.ACLGeoIPRule{
					ID: "00000000-0000-0000-0000-000000000001",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ddosx.ACLGeoIPRule{}, data)
		assert.Len(t, data, 2)
		assert.Equal(t, "00000000-0000-0000-0000-000000000000", data.([]ddosx.ACLGeoIPRule)[0].ID)
		assert.Equal(t, "00000000-0000-0000-0000-000000000001", data.([]ddosx.ACLGeoIPRule)[1].ID)
	})
}

func TestOutputDDoSXACLGeoIPRules_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputDDoSXACLGeoIPRules{
			ACLGeoIPRules: []ddosx.ACLGeoIPRule{
				ddosx.ACLGeoIPRule{
					ID: "00000000-0000-0000-0000-000000000000",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("id"))
		assert.Equal(t, "00000000-0000-0000-0000-000000000000", data[0].Get("id").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputDDoSXACLGeoIPRules{
			ACLGeoIPRules: []ddosx.ACLGeoIPRule{
				ddosx.ACLGeoIPRule{
					ID: "00000000-0000-0000-0000-000000000000",
				},
				ddosx.ACLGeoIPRule{
					ID: "00000000-0000-0000-0000-000000000001",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("id"))
		assert.Equal(t, "00000000-0000-0000-0000-000000000000", data[0].Get("id").Value)
		assert.True(t, data[1].Exists("id"))
		assert.Equal(t, "00000000-0000-0000-0000-000000000001", data[1].Get("id").Value)
	})
}

func TestOutputDDoSXACLGeoIPRulesModes_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputDDoSXACLGeoIPRulesModes{
			ACLGeoIPRulesModes: []ddosx.ACLGeoIPRulesMode{
				ddosx.ACLGeoIPRulesModeWhitelist,
			},
		}

		data := o.GetData()

		assert.IsType(t, []ddosx.ACLGeoIPRulesMode{}, data)
		assert.Equal(t, ddosx.ACLGeoIPRulesModeWhitelist, data.([]ddosx.ACLGeoIPRulesMode)[0])
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputDDoSXACLGeoIPRulesModes{
			ACLGeoIPRulesModes: []ddosx.ACLGeoIPRulesMode{
				ddosx.ACLGeoIPRulesModeWhitelist,
				ddosx.ACLGeoIPRulesModeBlacklist,
			},
		}

		data := o.GetData()

		assert.IsType(t, []ddosx.ACLGeoIPRulesMode{}, data)
		assert.Len(t, data, 2)
		assert.Equal(t, ddosx.ACLGeoIPRulesModeWhitelist, data.([]ddosx.ACLGeoIPRulesMode)[0])
		assert.Equal(t, ddosx.ACLGeoIPRulesModeBlacklist, data.([]ddosx.ACLGeoIPRulesMode)[1])
	})
}

func TestOutputDDoSXACLGeoIPRulesModes_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputDDoSXACLGeoIPRulesModes{
			ACLGeoIPRulesModes: []ddosx.ACLGeoIPRulesMode{
				ddosx.ACLGeoIPRulesModeWhitelist,
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("mode"))
		assert.Equal(t, "Whitelist", data[0].Get("mode").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputDDoSXACLGeoIPRulesModes{
			ACLGeoIPRulesModes: []ddosx.ACLGeoIPRulesMode{
				ddosx.ACLGeoIPRulesModeWhitelist,
				ddosx.ACLGeoIPRulesModeBlacklist,
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("mode"))
		assert.Equal(t, "Whitelist", data[0].Get("mode").Value)
		assert.True(t, data[1].Exists("mode"))
		assert.Equal(t, "Blacklist", data[1].Get("mode").Value)
	})
}

func TestOutputDDoSXCDNRules_GetData_ExpectedData(t *testing.T) {
	o := OutputDDoSXCDNRules{
		CDNRules: []ddosx.CDNRule{
			ddosx.CDNRule{
				ID: "00000000-0000-0000-0000-000000000000",
			},
		},
	}

	data := o.GetData()

	assert.IsType(t, []ddosx.CDNRule{}, data)
	assert.Equal(t, "00000000-0000-0000-0000-000000000000", data.([]ddosx.CDNRule)[0].ID)
}

func TestOutputDDoSXCDNRules_GetFieldData_ExpectedFieldData(t *testing.T) {
	o := OutputDDoSXCDNRules{
		CDNRules: []ddosx.CDNRule{
			ddosx.CDNRule{
				ID: "00000000-0000-0000-0000-000000000000",
			},
		},
	}

	data, err := o.GetFieldData()

	assert.Nil(t, err)
	assert.True(t, data[0].Exists("id"))
	assert.Equal(t, "00000000-0000-0000-0000-000000000000", data[0].Get("id").Value)
}

func TestOutputDDoSXHSTSConfiguration_GetData_ExpectedData(t *testing.T) {
	o := OutputDDoSXHSTSConfiguration{
		HSTSConfiguration: []ddosx.HSTSConfiguration{
			ddosx.HSTSConfiguration{
				Enabled: true,
			},
		},
	}

	data := o.GetData()

	assert.IsType(t, []ddosx.HSTSConfiguration{}, data)
	assert.Equal(t, "true", data.([]ddosx.HSTSConfiguration)[0].Enabled)
}

func TestOutputDDoSXHSTSConfiguration_GetFieldData_ExpectedFieldData(t *testing.T) {
	o := OutputDDoSXHSTSConfiguration{
		HSTSConfiguration: []ddosx.HSTSConfiguration{
			ddosx.HSTSConfiguration{
				Enabled: true,
			},
		},
	}

	data, err := o.GetFieldData()

	assert.Nil(t, err)
	assert.True(t, data[0].Exists("enabled"))
	assert.Equal(t, "true", data[0].Get("enabled").Value)
}

func TestOutputDDoSXHSTSRules_GetData_ExpectedData(t *testing.T) {
	o := OutputDDoSXHSTSRules{
		HSTSRules: []ddosx.HSTSRule{
			ddosx.HSTSRule{
				ID: "00000000-0000-0000-0000-000000000000",
			},
		},
	}

	data := o.GetData()

	assert.IsType(t, []ddosx.HSTSRule{}, data)
	assert.Equal(t, "00000000-0000-0000-0000-000000000000", data.([]ddosx.HSTSRule)[0].ID)
}

func TestOutputDDoSXHSTSRules_GetFieldData_ExpectedFieldData(t *testing.T) {
	o := OutputDDoSXHSTSRules{
		HSTSRules: []ddosx.HSTSRule{
			ddosx.HSTSRule{
				ID:         "00000000-0000-0000-0000-000000000000",
				RecordName: ptr.String("example.com"),
			},
		},
	}

	data, err := o.GetFieldData()

	assert.Nil(t, err)
	assert.True(t, data[0].Exists("id"))
	assert.Equal(t, "00000000-0000-0000-0000-000000000000", data[0].Get("id").Value)
	assert.True(t, data[0].Exists("record_name"))
	assert.Equal(t, "example.com", data[0].Get("record_name").Value)
}
