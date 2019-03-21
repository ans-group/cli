package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ukfast/sdk-go/pkg/service/ssl"
)

func TestOutputSSLCertificates_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputSSLCertificates{
			Certificates: []ssl.Certificate{
				ssl.Certificate{
					Name: "testcertificate1",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ssl.Certificate{}, data)
		assert.Equal(t, "testcertificate1", data.([]ssl.Certificate)[0].Name)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputSSLCertificates{
			Certificates: []ssl.Certificate{
				ssl.Certificate{
					Name: "testcertificate1",
				},
				ssl.Certificate{
					Name: "testcertificate2",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ssl.Certificate{}, data)
		assert.Len(t, data, 2)
		assert.Equal(t, "testcertificate1", data.([]ssl.Certificate)[0].Name)
		assert.Equal(t, "testcertificate2", data.([]ssl.Certificate)[1].Name)
	})
}

func TestOutputSSLCertificates_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputSSLCertificates{
			Certificates: []ssl.Certificate{
				ssl.Certificate{
					Name: "testcertificate1",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testcertificate1", data[0].Get("name").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputSSLCertificates{
			Certificates: []ssl.Certificate{
				ssl.Certificate{
					Name: "testcertificate1",
				},
				ssl.Certificate{
					Name: "testcertificate2",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("name"))
		assert.Equal(t, "testcertificate1", data[0].Get("name").Value)
		assert.True(t, data[1].Exists("name"))
		assert.Equal(t, "testcertificate2", data[1].Get("name").Value)
	})
}

func TestOutputSSLCertificateContents_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputSSLCertificateContents{
			CertificateContents: []ssl.CertificateContent{
				ssl.CertificateContent{
					Server: "testservercontent1",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ssl.CertificateContent{}, data)
		assert.Equal(t, "testservercontent1", data.([]ssl.CertificateContent)[0].Server)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputSSLCertificateContents{
			CertificateContents: []ssl.CertificateContent{
				ssl.CertificateContent{
					Server: "testservercontent1",
				},
				ssl.CertificateContent{
					Server: "testservercontent2",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ssl.CertificateContent{}, data)
		assert.Len(t, data, 2)
		assert.Equal(t, "testservercontent1", data.([]ssl.CertificateContent)[0].Server)
		assert.Equal(t, "testservercontent2", data.([]ssl.CertificateContent)[1].Server)
	})
}

func TestOutputSSLCertificateContents_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputSSLCertificateContents{
			CertificateContents: []ssl.CertificateContent{
				ssl.CertificateContent{
					Server: "testservercontent1",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.True(t, data[0].Exists("server"))
		assert.Equal(t, "testservercontent1", data[0].Get("server").Value)
	})

	t.Run("Multiple_ExpectedFieldData", func(t *testing.T) {
		o := OutputSSLCertificateContents{
			CertificateContents: []ssl.CertificateContent{
				ssl.CertificateContent{
					Server: "testservercontent1",
				},
				ssl.CertificateContent{
					Server: "testservercontent2",
				},
			},
		}

		data, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, data, 2)
		assert.True(t, data[0].Exists("server"))
		assert.Equal(t, "testservercontent1", data[0].Get("server").Value)
		assert.True(t, data[1].Exists("server"))
		assert.Equal(t, "testservercontent2", data[1].Get("server").Value)
	})
}

func TestOutputSSLCertificatePrivateKeys_GetData(t *testing.T) {
	t.Run("Single_ExpectedData", func(t *testing.T) {
		o := OutputSSLCertificatePrivateKeys{
			CertificatePrivateKeys: []ssl.CertificatePrivateKey{
				ssl.CertificatePrivateKey{
					Key: "testkey1",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ssl.CertificatePrivateKey{}, data)
		assert.Equal(t, "testkey1", data.([]ssl.CertificatePrivateKey)[0].Key)
	})

	t.Run("Multiple_ExpectedData", func(t *testing.T) {
		o := OutputSSLCertificatePrivateKeys{
			CertificatePrivateKeys: []ssl.CertificatePrivateKey{
				ssl.CertificatePrivateKey{
					Key: "testkey1",
				},
				ssl.CertificatePrivateKey{
					Key: "testkey2",
				},
			},
		}

		data := o.GetData()

		assert.IsType(t, []ssl.CertificatePrivateKey{}, data)
		assert.Len(t, data, 2)
		assert.Equal(t, "testkey1", data.([]ssl.CertificatePrivateKey)[0].Key)
		assert.Equal(t, "testkey2", data.([]ssl.CertificatePrivateKey)[1].Key)
	})
}

func TestOutputSSLCertificatePrivateKeys_GetFieldData(t *testing.T) {
	t.Run("Single_ExpectedFieldData", func(t *testing.T) {
		o := OutputSSLCertificatePrivateKeys{
			CertificatePrivateKeys: []ssl.CertificatePrivateKey{
				ssl.CertificatePrivateKey{
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
		o := OutputSSLCertificatePrivateKeys{
			CertificatePrivateKeys: []ssl.CertificatePrivateKey{
				ssl.CertificatePrivateKey{
					Key: "testkey1",
				},
				ssl.CertificatePrivateKey{
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
