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

func Test_sslReportShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := sslReportShowCmd(nil).Args(nil, []string{"example.com"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := sslReportShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_sslReportShow(t *testing.T) {
	t.Run("SingleReport", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		service.EXPECT().GetReport("example.com").Return(ssl.Report{}, nil).Times(1)

		sslReportShow(service, &cobra.Command{}, []string{"example.com"})
	})

	t.Run("MultipleReport", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetReport("example.com").Return(ssl.Report{}, nil),
			service.EXPECT().GetReport("example2.com").Return(ssl.Report{}, nil),
		)

		sslReportShow(service, &cobra.Command{}, []string{"example.com", "example2.com"})
	})

	t.Run("GetReportError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockSSLService(mockCtrl)

		service.EXPECT().GetReport("example.com").Return(ssl.Report{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving SSL report for domain [example.com]: test error\n", func() {
			sslReportShow(service, &cobra.Command{}, []string{"example.com"})
		})
	})
}
