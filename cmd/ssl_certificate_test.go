package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/cli/test"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/ssl"
)

func Test_sslCertificateList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		service.EXPECT().GetCertificates(gomock.Any()).Return([]ssl.Certificate{}, nil).Times(1)

		sslCertificateList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()
		defer func() { flagFilter = nil }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)
		flagFilter = []string{"invalidfilter"}

		output := test.CatchStdErr(t, func() {
			sslCertificateList(service, &cobra.Command{}, []string{})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Missing value for filtering\n", output)
	})

	t.Run("GetCertificatesError_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		service.EXPECT().GetCertificates(gomock.Any()).Return([]ssl.Certificate{}, errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			sslCertificateList(service, &cobra.Command{}, []string{})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error retrieving certificates: test error\n", output)
	})
}

func Test_sslCertificateShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := sslCertificateShowCmd().Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := sslCertificateShowCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing certificate", err.Error())
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

		output := test.CatchStdErr(t, func() {
			sslCertificateShow(service, &cobra.Command{}, []string{"abc"})
		})

		assert.Equal(t, "Invalid certificate ID [abc]\n", output)
	})

	t.Run("GetCertificateError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		service.EXPECT().GetCertificate(123).Return(ssl.Certificate{}, errors.New("test error"))

		output := test.CatchStdErr(t, func() {
			sslCertificateShow(service, &cobra.Command{}, []string{"123"})
		})

		assert.Equal(t, "Error retrieving certificate [123]: test error\n", output)
	})
}
