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

func TestOutputPSSRequests_GetFieldData(t *testing.T) {
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
}

func TestOutputPSSRequestReplies_GetData_ExpectedData(t *testing.T) {
	o := OutputPSSReplies{
		Replies: []pss.Reply{
			pss.Reply{
				Description: "test reply",
			},
		},
	}

	data := o.GetData()

	assert.IsType(t, []pss.Reply{}, data)
	assert.Equal(t, "test reply", data.([]pss.Reply)[0].Description)
}

func TestOutputPSSReplies_GetFieldData_ExpectedFieldData(t *testing.T) {
	o := OutputPSSReplies{
		Replies: []pss.Reply{
			pss.Reply{
				Description: "test reply",
			},
		},
	}

	data, err := o.GetFieldData()

	assert.Nil(t, err)
	assert.True(t, data[0].Exists("description"))
	assert.Equal(t, "test reply", data[0].Get("description").Value)
}
