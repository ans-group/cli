package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ukfast/sdk-go/pkg/service/account"
)

func TestOutputAccountContacts_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputAccountContacts{
			Contacts: []account.Contact{
				account.Contact{
					FirstName: "testname",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []account.Contact{}, data)
		assert.Equal(t, "testname", data.([]account.Contact)[0].FirstName)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputAccountContacts{
			Contacts: []account.Contact{
				account.Contact{
					FirstName: "testname1",
				},
				account.Contact{
					FirstName: "testname2",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []account.Contact{}, data)
		assert.Len(t, data, 2)
		assert.Equal(t, "testname1", data.([]account.Contact)[0].FirstName)
		assert.Equal(t, "testname2", data.([]account.Contact)[1].FirstName)
	})
}

func TestOutputAccountContacts_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputAccountContacts{
			Contacts: []account.Contact{
				account.Contact{
					FirstName: "testname",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("first_name"))
		assert.Equal(t, "testname", data[0].Get("first_name").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputAccountContacts{
			Contacts: []account.Contact{
				account.Contact{
					FirstName: "testname1",
				},
				account.Contact{
					FirstName: "testname2",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("first_name"))
		assert.Equal(t, "testname1", data[0].Get("first_name").Value)
		assert.True(t, data[1].Exists("first_name"))
		assert.Equal(t, "testname2", data[1].Get("first_name").Value)
	})
}

func TestOutputAccountDetails_GetData_ExpectedData(t *testing.T) {
	o := OutputAccountDetails{
		Details: []account.Details{
			account.Details{
				VATIdentificationNumber: "GB123",
			},
		},
	}

	data := o.GetData()

	assert.IsType(t, []account.Details{}, data)
	assert.Equal(t, "GB123", data.([]account.Details)[0].VATIdentificationNumber)
}

func TestOutputAccountDetails_GetFieldData_ExpectedFieldData(t *testing.T) {
	o := OutputAccountDetails{
		Details: []account.Details{
			account.Details{
				VATIdentificationNumber: "GB123",
			},
			account.Details{
				VATIdentificationNumber: "GB456",
			},
		},
	}

	data, err := o.GetFieldData()

	assert.Nil(t, err)
	assert.Len(t, data, 2)
	assert.True(t, data[0].Exists("vat_identification_number"))
	assert.Equal(t, "GB123", data[0].Get("vat_identification_number").Value)
	assert.True(t, data[1].Exists("vat_identification_number"))
	assert.Equal(t, "GB456", data[1].Get("vat_identification_number").Value)
}
