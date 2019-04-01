package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ukfast/sdk-go/pkg/service/pss"
)

func TestOutputPSSRequests_GetData_ExpectedData(t *testing.T) {
	o := OutputPSSRequests{
		Requests: []pss.Request{
			pss.Request{
				ID: 123,
			},
		},
	}

	data := o.GetData()

	assert.IsType(t, []pss.Request{}, data)
	assert.Equal(t, 123, data.([]pss.Request)[0].ID)
}

func TestOutputPSSRequests_GetFieldData_ExpectedFieldData(t *testing.T) {
	t.Run("ExpectedFieldData", func(t *testing.T) {
		o := OutputPSSRequests{
			Requests: []pss.Request{
				pss.Request{
					ID: 123,
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("id"))
		assert.Equal(t, "123", data[0].Get("id").Value)
	})

	t.Run("NonNilAssignee_ExpectedAssignee", func(t *testing.T) {
		o := OutputPSSRequests{
			Requests: []pss.Request{
				pss.Request{
					Assignee: &pss.SupportUser{
						Name: "test a",
					},
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("assignee_name"))
		assert.Equal(t, "test a", data[0].Get("assignee_name").Value)
	})
}
