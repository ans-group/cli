package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestGetKeyValueFromStringFlag(t *testing.T) {
	t.Run("Valid_NoError", func(t *testing.T) {
		flag := "testkey=testvalue"

		key, value, err := GetKeyValueFromStringFlag(flag)

		assert.Nil(t, err)
		assert.Equal(t, "testkey", key)
		assert.Equal(t, "testvalue", value)
	})

	t.Run("Empty_Error", func(t *testing.T) {
		flag := ""

		_, _, err := GetKeyValueFromStringFlag(flag)

		assert.NotNil(t, err)
	})

	t.Run("OnlyKey_Error", func(t *testing.T) {
		flag := "testkey"

		_, _, err := GetKeyValueFromStringFlag(flag)

		assert.NotNil(t, err)
	})

	t.Run("MissingValue_Error", func(t *testing.T) {
		flag := "testkey="

		_, _, err := GetKeyValueFromStringFlag(flag)

		assert.NotNil(t, err)
	})

	t.Run("MissingKey_Error", func(t *testing.T) {
		flag := "=testvalue"

		_, _, err := GetKeyValueFromStringFlag(flag)

		assert.NotNil(t, err)
	})

	t.Run("MultiValue_Error", func(t *testing.T) {
		flag := "testkey=testvalue1=testvalue2"

		_, _, err := GetKeyValueFromStringFlag(flag)

		assert.NotNil(t, err)
	})
}

func TestGetCreateTagRequestFromStringArrayFlag(t *testing.T) {
	t.Run("None_NoError", func(t *testing.T) {
		var tagFlags []string

		r, err := GetCreateTagRequestFromStringArrayFlag(tagFlags)

		assert.Nil(t, err)
		assert.Len(t, r, 0)
	})

	t.Run("Single", func(t *testing.T) {
		var tagFlags []string
		tagFlags = append(tagFlags, "testkey1=testvalue1")

		r, err := GetCreateTagRequestFromStringArrayFlag(tagFlags)

		assert.Nil(t, err)
		assert.Len(t, r, 1)
		assert.Equal(t, "testkey1", r[0].Key)
		assert.Equal(t, "testvalue1", r[0].Value)
	})

	t.Run("Multiple", func(t *testing.T) {
		var tagFlags []string
		tagFlags = append(tagFlags, "testkey1=testvalue1")
		tagFlags = append(tagFlags, "testkey2=testvalue2")

		r, err := GetCreateTagRequestFromStringArrayFlag(tagFlags)

		assert.Nil(t, err)
		assert.Len(t, r, 2)
		assert.Equal(t, "testkey1", r[0].Key)
		assert.Equal(t, "testvalue1", r[0].Value)
		assert.Equal(t, "testkey2", r[1].Key)
		assert.Equal(t, "testvalue2", r[1].Value)
	})

	t.Run("Invalid_ReturnsError", func(t *testing.T) {
		tagFlags := []string{"invalid"}

		_, err := GetCreateTagRequestFromStringArrayFlag(tagFlags)

		assert.NotNil(t, err)
		assert.Equal(t, "Invalid format, expecting: key=value", err.Error())
	})
}

func TestOutputECloudVirtualMachines_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputECloudVirtualMachines{
			VirtualMachines: []ecloud.VirtualMachine{
				ecloud.VirtualMachine{
					Name: "testvm1",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ecloud.VirtualMachine{}, data)
		assert.Equal(t, "testvm1", data.([]ecloud.VirtualMachine)[0].Name)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputECloudVirtualMachines{
			VirtualMachines: []ecloud.VirtualMachine{
				ecloud.VirtualMachine{
					Name: "testvm1",
				},
				ecloud.VirtualMachine{
					Name: "testvm2",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ecloud.VirtualMachine{}, data)
		assert.Len(t, data, 2)
		assert.Equal(t, "testvm1", data.([]ecloud.VirtualMachine)[0].Name)
		assert.Equal(t, "testvm2", data.([]ecloud.VirtualMachine)[1].Name)
	})
}

func TestOutputECloudVirtualMachines_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputECloudVirtualMachines{
			VirtualMachines: []ecloud.VirtualMachine{
				ecloud.VirtualMachine{
					Name: "testvm1",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testvm1", data[0].Get("name").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputECloudVirtualMachines{
			VirtualMachines: []ecloud.VirtualMachine{
				ecloud.VirtualMachine{
					Name: "testvm1",
				},
				ecloud.VirtualMachine{
					Name: "testvm2",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testvm1", data[0].Get("name").Value)
		assert.True(t, data[1].Exists("name"))
		assert.Equal(t, "testvm2", data[1].Get("name").Value)
	})
}

func TestOutputECloudVirtualMachineDisks_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputECloudVirtualMachineDisks{
			VirtualMachineDisks: []ecloud.VirtualMachineDisk{
				ecloud.VirtualMachineDisk{
					Name: "testvmdisk1",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ecloud.VirtualMachineDisk{}, data)
		assert.Equal(t, "testvmdisk1", data.([]ecloud.VirtualMachineDisk)[0].Name)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputECloudVirtualMachineDisks{
			VirtualMachineDisks: []ecloud.VirtualMachineDisk{
				ecloud.VirtualMachineDisk{
					Name: "testvmdisk1",
				},
				ecloud.VirtualMachineDisk{
					Name: "testvmdisk2",
				},
			},
		}

		data := o.GetData()

		assert.Len(t, data, 2)
		assert.IsType(t, []ecloud.VirtualMachineDisk{}, data)
		assert.Equal(t, "testvmdisk1", data.([]ecloud.VirtualMachineDisk)[0].Name)
		assert.Equal(t, "testvmdisk2", data.([]ecloud.VirtualMachineDisk)[1].Name)
	})
}

func TestOutputECloudVirtualMachineDisks_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputECloudVirtualMachineDisks{
			VirtualMachineDisks: []ecloud.VirtualMachineDisk{
				ecloud.VirtualMachineDisk{
					Name: "testvmdisk1",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testvmdisk1", data[0].Get("name").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputECloudVirtualMachineDisks{
			VirtualMachineDisks: []ecloud.VirtualMachineDisk{
				ecloud.VirtualMachineDisk{
					Name: "testvmdisk1",
				},
				ecloud.VirtualMachineDisk{
					Name: "testvmdisk2",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testvmdisk1", data[0].Get("name").Value)
		assert.True(t, data[1].Exists("name"))
		assert.Equal(t, "testvmdisk2", data[1].Get("name").Value)
	})
}

func TestOutputECloudTags_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputECloudTags{
			Tags: []ecloud.Tag{
				ecloud.Tag{
					Key: "testkey1",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ecloud.Tag{}, data)
		assert.Equal(t, "testkey1", data.([]ecloud.Tag)[0].Key)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputECloudTags{
			Tags: []ecloud.Tag{
				ecloud.Tag{
					Key: "testkey1",
				},
				ecloud.Tag{
					Key: "testkey2",
				},
			},
		}

		data := o.GetData()

		assert.Len(t, data, 2)
		assert.IsType(t, []ecloud.Tag{}, data)
		assert.Equal(t, "testkey1", data.([]ecloud.Tag)[0].Key)
		assert.Equal(t, "testkey2", data.([]ecloud.Tag)[1].Key)
	})
}

func TestOutputECloudTags_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputECloudTags{
			Tags: []ecloud.Tag{
				ecloud.Tag{
					Key: "testkey1",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("key"))
		assert.Equal(t, "testkey1", data[0].Get("key").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputECloudTags{
			Tags: []ecloud.Tag{
				ecloud.Tag{
					Key: "testkey1",
				},
				ecloud.Tag{
					Key: "testkey2",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("key"))
		assert.Equal(t, "testkey1", data[0].Get("key").Value)
		assert.True(t, data[1].Exists("key"))
		assert.Equal(t, "testkey2", data[1].Get("key").Value)
	})
}

func TestOutputECloudSolutions_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputECloudSolutions{
			Solutions: []ecloud.Solution{
				ecloud.Solution{
					Name: "testsolution1",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ecloud.Solution{}, data)
		assert.Equal(t, "testsolution1", data.([]ecloud.Solution)[0].Name)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputECloudSolutions{
			Solutions: []ecloud.Solution{
				ecloud.Solution{
					Name: "testsolution1",
				},
				ecloud.Solution{
					Name: "testsolution2",
				},
			},
		}

		data := o.GetData()

		assert.Len(t, data, 2)
		assert.IsType(t, []ecloud.Solution{}, data)
		assert.Equal(t, "testsolution1", data.([]ecloud.Solution)[0].Name)
		assert.Equal(t, "testsolution2", data.([]ecloud.Solution)[1].Name)
	})
}

func TestOutputECloudSolutions_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputECloudSolutions{
			Solutions: []ecloud.Solution{
				ecloud.Solution{
					Name: "testsolution1",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testsolution1", data[0].Get("name").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputECloudSolutions{
			Solutions: []ecloud.Solution{
				ecloud.Solution{
					Name: "testsolution1",
				},
				ecloud.Solution{
					Name: "testsolution2",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testsolution1", data[0].Get("name").Value)
		assert.True(t, data[1].Exists("name"))
		assert.Equal(t, "testsolution2", data[1].Get("name").Value)
	})
}

func TestOutputECloudSites_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputECloudSites{
			Sites: []ecloud.Site{
				ecloud.Site{
					SolutionID: 123,
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ecloud.Site{}, data)
		assert.Equal(t, 123, data.([]ecloud.Site)[0].SolutionID)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputECloudSites{
			Sites: []ecloud.Site{
				ecloud.Site{
					SolutionID: 123,
				},
				ecloud.Site{
					SolutionID: 456,
				},
			},
		}

		data := o.GetData()

		assert.Len(t, data, 2)
		assert.IsType(t, []ecloud.Site{}, data)
		assert.Equal(t, 123, data.([]ecloud.Site)[0].SolutionID)
		assert.Equal(t, 456, data.([]ecloud.Site)[1].SolutionID)
	})
}

func TestOutputECloudSites_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputECloudSites{
			Sites: []ecloud.Site{
				ecloud.Site{
					SolutionID: 123,
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("solution_id"))
		assert.Equal(t, "123", data[0].Get("solution_id").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputECloudSites{
			Sites: []ecloud.Site{
				ecloud.Site{
					SolutionID: 123,
				},
				ecloud.Site{
					SolutionID: 456,
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("solution_id"))
		assert.Equal(t, "123", data[0].Get("solution_id").Value)
		assert.True(t, data[1].Exists("solution_id"))
		assert.Equal(t, "456", data[1].Get("solution_id").Value)
	})
}

func TestOutputECloudHosts_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputECloudHosts{
			Hosts: []ecloud.Host{
				ecloud.Host{
					Name: "testhost1",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ecloud.Host{}, data)
		assert.Equal(t, "testhost1", data.([]ecloud.Host)[0].Name)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputECloudHosts{
			Hosts: []ecloud.Host{
				ecloud.Host{
					Name: "testhost1",
				},
				ecloud.Host{
					Name: "testhost2",
				},
			},
		}

		data := o.GetData()

		assert.Len(t, data, 2)
		assert.IsType(t, []ecloud.Host{}, data)
		assert.Equal(t, "testhost1", data.([]ecloud.Host)[0].Name)
		assert.Equal(t, "testhost2", data.([]ecloud.Host)[1].Name)
	})
}

func TestOutputECloudHosts_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputECloudHosts{
			Hosts: []ecloud.Host{
				ecloud.Host{
					Name: "testhost1",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testhost1", data[0].Get("name").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputECloudHosts{
			Hosts: []ecloud.Host{
				ecloud.Host{
					Name: "testhost1",
				},
				ecloud.Host{
					Name: "testhost2",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testhost1", data[0].Get("name").Value)
		assert.True(t, data[1].Exists("name"))
		assert.Equal(t, "testhost2", data[1].Get("name").Value)
	})
}

func TestOutputECloudDatastores_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputECloudDatastores{
			Datastores: []ecloud.Datastore{
				ecloud.Datastore{
					Name: "testdatastore1",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ecloud.Datastore{}, data)
		assert.Equal(t, "testdatastore1", data.([]ecloud.Datastore)[0].Name)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputECloudDatastores{
			Datastores: []ecloud.Datastore{
				ecloud.Datastore{
					Name: "testdatastore1",
				},
				ecloud.Datastore{
					Name: "testdatastore2",
				},
			},
		}

		data := o.GetData()

		assert.Len(t, data, 2)
		assert.IsType(t, []ecloud.Datastore{}, data)
		assert.Equal(t, "testdatastore1", data.([]ecloud.Datastore)[0].Name)
		assert.Equal(t, "testdatastore2", data.([]ecloud.Datastore)[1].Name)
	})
}

func TestOutputECloudDatastores_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputECloudDatastores{
			Datastores: []ecloud.Datastore{
				ecloud.Datastore{
					Name: "testdatastore1",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testdatastore1", data[0].Get("name").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputECloudDatastores{
			Datastores: []ecloud.Datastore{
				ecloud.Datastore{
					Name: "testdatastore1",
				},
				ecloud.Datastore{
					Name: "testdatastore2",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testdatastore1", data[0].Get("name").Value)
		assert.True(t, data[1].Exists("name"))
		assert.Equal(t, "testdatastore2", data[1].Get("name").Value)
	})
}

func TestOutputECloudTemplates_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputECloudTemplates{
			Templates: []ecloud.Template{
				ecloud.Template{
					Name: "testtemplate1",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ecloud.Template{}, data)
		assert.Equal(t, "testtemplate1", data.([]ecloud.Template)[0].Name)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputECloudTemplates{
			Templates: []ecloud.Template{
				ecloud.Template{
					Name: "testtemplate1",
				},
				ecloud.Template{
					Name: "testtemplate2",
				},
			},
		}

		data := o.GetData()

		assert.Len(t, data, 2)
		assert.IsType(t, []ecloud.Template{}, data)
		assert.Equal(t, "testtemplate1", data.([]ecloud.Template)[0].Name)
		assert.Equal(t, "testtemplate2", data.([]ecloud.Template)[1].Name)
	})
}

func TestOutputECloudTemplates_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputECloudTemplates{
			Templates: []ecloud.Template{
				ecloud.Template{
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
		o := OutputECloudTemplates{
			Templates: []ecloud.Template{
				ecloud.Template{
					Name: "testtemplate1",
				},
				ecloud.Template{
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

func TestOutputECloudNetworks_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputECloudNetworks{
			Networks: []ecloud.Network{
				ecloud.Network{
					Name: "testnetwork1",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ecloud.Network{}, data)
		assert.Equal(t, "testnetwork1", data.([]ecloud.Network)[0].Name)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputECloudNetworks{
			Networks: []ecloud.Network{
				ecloud.Network{
					Name: "testnetwork1",
				},
				ecloud.Network{
					Name: "testnetwork2",
				},
			},
		}

		data := o.GetData()

		assert.Len(t, data, 2)
		assert.IsType(t, []ecloud.Network{}, data)
		assert.Equal(t, "testnetwork1", data.([]ecloud.Network)[0].Name)
		assert.Equal(t, "testnetwork2", data.([]ecloud.Network)[1].Name)
	})
}

func TestOutputECloudNetworks_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputECloudNetworks{
			Networks: []ecloud.Network{
				ecloud.Network{
					Name: "testnetwork1",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testnetwork1", data[0].Get("name").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputECloudNetworks{
			Networks: []ecloud.Network{
				ecloud.Network{
					Name: "testnetwork1",
				},
				ecloud.Network{
					Name: "testnetwork2",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testnetwork1", data[0].Get("name").Value)
		assert.True(t, data[1].Exists("name"))
		assert.Equal(t, "testnetwork2", data[1].Get("name").Value)
	})
}

func TestOutputECloudFirewalls_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputECloudFirewalls{
			Firewalls: []ecloud.Firewall{
				ecloud.Firewall{
					Name: "testfirewall1",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ecloud.Firewall{}, data)
		assert.Equal(t, "testfirewall1", data.([]ecloud.Firewall)[0].Name)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputECloudFirewalls{
			Firewalls: []ecloud.Firewall{
				ecloud.Firewall{
					Name: "testfirewall1",
				},
				ecloud.Firewall{
					Name: "testfirewall2",
				},
			},
		}

		data := o.GetData()

		assert.Len(t, data, 2)
		assert.IsType(t, []ecloud.Firewall{}, data)
		assert.Equal(t, "testfirewall1", data.([]ecloud.Firewall)[0].Name)
		assert.Equal(t, "testfirewall2", data.([]ecloud.Firewall)[1].Name)
	})
}

func TestOutputECloudFirewalls_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputECloudFirewalls{
			Firewalls: []ecloud.Firewall{
				ecloud.Firewall{
					Name: "testfirewall1",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testfirewall1", data[0].Get("name").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputECloudFirewalls{
			Firewalls: []ecloud.Firewall{
				ecloud.Firewall{
					Name: "testfirewall1",
				},
				ecloud.Firewall{
					Name: "testfirewall2",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testfirewall1", data[0].Get("name").Value)
		assert.True(t, data[1].Exists("name"))
		assert.Equal(t, "testfirewall2", data[1].Get("name").Value)
	})
}
