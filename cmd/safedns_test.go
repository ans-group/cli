package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ukfast/sdk-go/pkg/service/safedns"
)

func TestOutputSafeDNSZones_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputSafeDNSZones{
			Zones: []safedns.Zone{
				safedns.Zone{
					Name: "testdomain1.com",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []safedns.Zone{}, data)
		assert.Equal(t, "testdomain1.com", data.([]safedns.Zone)[0].Name)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputSafeDNSZones{
			Zones: []safedns.Zone{
				safedns.Zone{
					Name: "testdomain1.com",
				},
				safedns.Zone{
					Name: "testdomain2.com",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []safedns.Zone{}, data)
		assert.Len(t, data, 2)
		assert.Equal(t, "testdomain1.com", data.([]safedns.Zone)[0].Name)
		assert.Equal(t, "testdomain2.com", data.([]safedns.Zone)[1].Name)
	})
}

func TestOutputSafeDNSZones_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputSafeDNSZones{
			Zones: []safedns.Zone{
				safedns.Zone{
					Name: "testdomain1.com",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testdomain1.com", data[0].Get("name").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputSafeDNSZones{
			Zones: []safedns.Zone{
				safedns.Zone{
					Name: "testdomain1.com",
				},
				safedns.Zone{
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

func TestOutputSafeDNSRecords_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputSafeDNSRecords{
			Records: []safedns.Record{
				safedns.Record{
					Name: "www.testdomain1.com",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []safedns.Record{}, data)
		assert.Equal(t, "www.testdomain1.com", data.([]safedns.Record)[0].Name)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputSafeDNSRecords{
			Records: []safedns.Record{
				safedns.Record{
					Name: "www.testdomain1.com",
				},
				safedns.Record{
					Name: "www.testdomain2.com",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []safedns.Record{}, data)
		assert.Len(t, data, 2)
		assert.Equal(t, "www.testdomain1.com", data.([]safedns.Record)[0].Name)
		assert.Equal(t, "www.testdomain2.com", data.([]safedns.Record)[1].Name)
	})
}

func TestOutputSafeDNSRecords_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputSafeDNSRecords{
			Records: []safedns.Record{
				safedns.Record{
					Name: "www.testdomain1.com",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "www.testdomain1.com", data[0].Get("name").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputSafeDNSRecords{
			Records: []safedns.Record{
				safedns.Record{
					Name: "www.testdomain1.com",
				},
				safedns.Record{
					Name: "www.testdomain2.com",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "www.testdomain1.com", data[0].Get("name").Value)
		assert.True(t, data[1].Exists("name"))
		assert.Equal(t, "www.testdomain2.com", data[1].Get("name").Value)
	})
}

func TestOutputSafeDNSNotes_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputSafeDNSNotes{
			Notes: []safedns.Note{
				safedns.Note{
					Notes: "testnote1",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []safedns.Note{}, data)
		assert.Equal(t, "testnote1", data.([]safedns.Note)[0].Notes)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputSafeDNSNotes{
			Notes: []safedns.Note{
				safedns.Note{
					Notes: "testnote1",
				},
				safedns.Note{
					Notes: "testnote2",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []safedns.Note{}, data)
		assert.Len(t, data, 2)
		assert.Equal(t, "testnote1", data.([]safedns.Note)[0].Notes)
		assert.Equal(t, "testnote2", data.([]safedns.Note)[1].Notes)
	})
}

func TestOutputSafeDNSNotes_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputSafeDNSNotes{
			Notes: []safedns.Note{
				safedns.Note{
					Notes: "testnote1",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("notes"))
		assert.Equal(t, "testnote1", data[0].Get("notes").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputSafeDNSNotes{
			Notes: []safedns.Note{
				safedns.Note{
					Notes: "testnote1",
				},
				safedns.Note{
					Notes: "testnote2",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("notes"))
		assert.Equal(t, "testnote1", data[0].Get("notes").Value)
		assert.True(t, data[1].Exists("notes"))
		assert.Equal(t, "testnote2", data[1].Get("notes").Value)
	})
}

func TestOutputSafeDNSTemplates_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputSafeDNSTemplates{
			Templates: []safedns.Template{
				safedns.Template{
					Name: "testtemplate1",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []safedns.Template{}, data)
		assert.Equal(t, "testtemplate1", data.([]safedns.Template)[0].Name)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputSafeDNSTemplates{
			Templates: []safedns.Template{
				safedns.Template{
					Name: "testtemplate1",
				},
				safedns.Template{
					Name: "testtemplate2",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []safedns.Template{}, data)
		assert.Len(t, data, 2)
		assert.Equal(t, "testtemplate1", data.([]safedns.Template)[0].Name)
		assert.Equal(t, "testtemplate2", data.([]safedns.Template)[1].Name)
	})
}

func TestOutputSafeDNSTemplates_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputSafeDNSTemplates{
			Templates: []safedns.Template{
				safedns.Template{
					Name: "testtemplate1",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testtemplate1", data[0].Get("name").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputSafeDNSTemplates{
			Templates: []safedns.Template{
				safedns.Template{
					Name: "testtemplate1",
				},
				safedns.Template{
					Name: "testtemplate2",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testtemplate1", data[0].Get("name").Value)
		assert.True(t, data[1].Exists("name"))
		assert.Equal(t, "testtemplate2", data[1].Get("name").Value)
	})
}
