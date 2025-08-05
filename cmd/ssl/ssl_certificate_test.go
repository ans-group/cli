package ssl

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/ssl"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_sslCertificateList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		service.EXPECT().GetCertificates(gomock.Any()).Return([]ssl.Certificate{}, nil).Times(1)

		sslCertificateList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := sslCertificateList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetCertificatesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		service.EXPECT().GetCertificates(gomock.Any()).Return([]ssl.Certificate{}, errors.New("test error")).Times(1)

		err := sslCertificateList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "error retrieving certificates: test error", err.Error())
	})
}

func Test_sslCertificateShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := sslCertificateShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := sslCertificateShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing certificate", err.Error())
	})
}

func Test_sslCertificateShow(t *testing.T) {
	t.Run("SingleCertificate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		service.EXPECT().GetCertificate(123).Return(ssl.Certificate{}, nil).Times(1)

		sslCertificateShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleCertificates", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetCertificate(123).Return(ssl.Certificate{}, nil),
			service.EXPECT().GetCertificate(456).Return(ssl.Certificate{}, nil),
		)

		sslCertificateShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetCertificateID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid certificate ID [abc]\n", func() {
			sslCertificateShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetCertificateError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		service.EXPECT().GetCertificate(123).Return(ssl.Certificate{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving certificate [123]: test error\n", func() {
			sslCertificateShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}
