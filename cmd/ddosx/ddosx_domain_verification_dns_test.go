package ddosx

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ddosxDomainVerificationDNSVerifyCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainVerificationDNSVerifyCmd(nil).Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ddosxDomainVerificationDNSVerifyCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing domain", err.Error())
	})
}

func Test_ddosxDomainVerificationDNSVerify(t *testing.T) {
	t.Run("SingleDomain_NoError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().VerifyDomainDNS("testdomain1.co.uk").Return(nil)

		ddosxDomainVerificationDNSVerify(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("DownloadDomainVerificationFileError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().VerifyDomainDNS("testdomain1.co.uk").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error verifying domain [testdomain1.co.uk] via DNS verification method: test error\n", func() {
			ddosxDomainVerificationDNSVerify(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}
