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

func Test_sslCertificateContentShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := sslCertificateContentShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := sslCertificateContentShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing certificate", err.Error())
	})
}

func Test_sslCertificateContentShow(t *testing.T) {
	t.Run("SingleCertificateContent", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		service.EXPECT().GetCertificateContent(123).Return(ssl.CertificateContent{}, nil).Times(1)

		sslCertificateContentShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleCertificateContents", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetCertificateContent(123).Return(ssl.CertificateContent{}, nil),
			service.EXPECT().GetCertificateContent(456).Return(ssl.CertificateContent{}, nil),
		)

		sslCertificateContentShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetCertificateContentID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid certificate ID [abc]\n", func() {
			sslCertificateContentShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetCertificateContentError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		service.EXPECT().GetCertificateContent(123).Return(ssl.CertificateContent{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving certificate content [123]: test error\n", func() {
			sslCertificateContentShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}
