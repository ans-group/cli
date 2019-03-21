package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/ssl"
)

func Test_sslCertificatePrivateKeyShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := sslCertificatePrivateKeyShowCmd().Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := sslCertificatePrivateKeyShowCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing certificate", err.Error())
	})
}

func Test_sslCertificatePrivateKeyShow(t *testing.T) {
	t.Run("SingleCertificatePrivateKey", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		service.EXPECT().GetCertificatePrivateKey(123).Return(ssl.CertificatePrivateKey{}, nil).Times(1)

		sslCertificatePrivateKeyShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleCertificatePrivateKeys", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetCertificatePrivateKey(123).Return(ssl.CertificatePrivateKey{}, nil),
			service.EXPECT().GetCertificatePrivateKey(456).Return(ssl.CertificatePrivateKey{}, nil),
		)

		sslCertificatePrivateKeyShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetCertificatePrivateKeyID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		output := test.CatchStdErr(t, func() {
			sslCertificatePrivateKeyShow(service, &cobra.Command{}, []string{"abc"})
		})

		assert.Equal(t, "Invalid certificate ID [abc]\n", output)
	})

	t.Run("GetCertificatePrivateKeyError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		service.EXPECT().GetCertificatePrivateKey(123).Return(ssl.CertificatePrivateKey{}, errors.New("test error"))

		output := test.CatchStdErr(t, func() {
			sslCertificatePrivateKeyShow(service, &cobra.Command{}, []string{"123"})
		})

		assert.Equal(t, "Error retrieving certificate private key [123]: test error\n", output)
	})
}
