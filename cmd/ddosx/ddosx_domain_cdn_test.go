package ddosx

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func Test_ddosxDomainCDNEnableCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainCDNEnableCmd(nil).Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		err := ddosxDomainCDNEnableCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainCDNEnable(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().AddDomainCDNConfiguration("testdomain1.co.uk").Return(nil)
		service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{}, nil)

		ddosxDomainCDNEnable(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("AddDomainCDNConfiguration_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().AddDomainCDNConfiguration("testdomain1.co.uk").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error enabling CDN for domain [testdomain1.co.uk]: test error\n", func() {
			ddosxDomainCDNEnable(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})

	t.Run("GetDomainError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().AddDomainCDNConfiguration("testdomain1.co.uk").Return(nil)
		service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving updated domain [testdomain1.co.uk]: test error\n", func() {
			ddosxDomainCDNEnable(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}

func Test_ddosxDomainCDNDisableCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainCDNDisableCmd(nil).Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		err := ddosxDomainCDNDisableCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainCDNDisable(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomainCDNConfiguration("testdomain1.co.uk").Return(nil)
		service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{}, nil)

		ddosxDomainCDNDisable(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("DeleteDomainCDNConfiguration_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomainCDNConfiguration("testdomain1.co.uk").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error disabling CDN for domain [testdomain1.co.uk]: test error\n", func() {
			ddosxDomainCDNDisable(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})

	t.Run("GetDomainError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomainCDNConfiguration("testdomain1.co.uk").Return(nil)
		service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving updated domain [testdomain1.co.uk]: test error\n", func() {
			ddosxDomainCDNDisable(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}

func Test_ddosxDomainCDNPurgeCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainCDNPurgeCmd(nil).Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		err := ddosxDomainCDNPurgeCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainCDNPurge(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ddosxDomainCDNPurgeCmd(nil)
		cmd.Flags().Set("record-name", "example.com")
		cmd.Flags().Set("uri", "test.html")

		expectedRequest := ddosx.PurgeCDNRequest{
			RecordName: "example.com",
			URI:        "test.html",
		}

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().PurgeDomainCDN("testdomain1.co.uk", gomock.Eq(expectedRequest)).Return(nil)

		ddosxDomainCDNPurge(service, cmd, []string{"testdomain1.co.uk"})
	})

	t.Run("DeleteDomainCDNConfiguration_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ddosxDomainCDNPurgeCmd(nil)
		cmd.Flags().Set("record-name", "example.com")
		cmd.Flags().Set("uri", "test.html")

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().PurgeDomainCDN("testdomain1.co.uk", gomock.Any()).Return(errors.New("test error"))

		err := ddosxDomainCDNPurge(service, cmd, []string{"testdomain1.co.uk"})

		assert.Equal(t, "Error purging CDN content for domain: test error", err.Error())
	})
}
