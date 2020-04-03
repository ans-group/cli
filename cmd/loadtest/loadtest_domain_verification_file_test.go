package loadtest

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
)

func Test_loadtestDomainVerificationFileVerifyCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadtestDomainVerificationFileVerifyCmd(nil).Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadtestDomainVerificationFileVerifyCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_loadtestDomainVerificationFileVerify(t *testing.T) {
	t.Run("SingleDomain_NoError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().VerifyDomainFile("00000000-0000-0000-0000-000000000000").Return(nil)

		loadtestDomainVerificationFileVerify(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("DownloadDomainVerificationFileError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLTaaSService(mockCtrl)

		service.EXPECT().VerifyDomainFile("00000000-0000-0000-0000-000000000000").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error verifying domain [00000000-0000-0000-0000-000000000000] via File verification method: test error\n", func() {
			loadtestDomainVerificationFileVerify(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
		})
	})
}
