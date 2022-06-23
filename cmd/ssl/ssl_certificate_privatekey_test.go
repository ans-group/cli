package ssl

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/ssl"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_sslCertificatePrivateKeyShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := sslCertificatePrivateKeyShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := sslCertificatePrivateKeyShowCmd(nil).Args(nil, []string{})

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

		test_output.AssertErrorOutput(t, "Invalid certificate ID [abc]\n", func() {
			sslCertificatePrivateKeyShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetCertificatePrivateKeyError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		service.EXPECT().GetCertificatePrivateKey(123).Return(ssl.CertificatePrivateKey{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving certificate private key [123]: test error\n", func() {
			sslCertificatePrivateKeyShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}
