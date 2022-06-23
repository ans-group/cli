package ddosx

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ddosxDomainHSTSEnableCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainHSTSEnableCmd(nil).Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		err := ddosxDomainHSTSEnableCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainHSTSEnable(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().AddDomainHSTSConfiguration("testdomain1.co.uk").Return(nil)
		service.EXPECT().GetDomainHSTSConfiguration("testdomain1.co.uk").Return(ddosx.HSTSConfiguration{}, nil)

		ddosxDomainHSTSEnable(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("AddDomainHSTSConfiguration_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().AddDomainHSTSConfiguration("testdomain1.co.uk").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error enabling HSTS for domain [testdomain1.co.uk]: test error\n", func() {
			ddosxDomainHSTSEnable(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})

	t.Run("GetDomainError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().AddDomainHSTSConfiguration("testdomain1.co.uk").Return(nil)
		service.EXPECT().GetDomainHSTSConfiguration("testdomain1.co.uk").Return(ddosx.HSTSConfiguration{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving updated HSTS configuration for domain [testdomain1.co.uk]: test error\n", func() {
			ddosxDomainHSTSEnable(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}

func Test_ddosxDomainHSTSDisableCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainHSTSDisableCmd(nil).Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		err := ddosxDomainHSTSDisableCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainHSTSDisable(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomainHSTSConfiguration("testdomain1.co.uk").Return(nil)
		service.EXPECT().GetDomainHSTSConfiguration("testdomain1.co.uk").Return(ddosx.HSTSConfiguration{}, nil)

		ddosxDomainHSTSDisable(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("DeleteDomainHSTSConfiguration_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomainHSTSConfiguration("testdomain1.co.uk").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error disabling HSTS for domain [testdomain1.co.uk]: test error\n", func() {
			ddosxDomainHSTSDisable(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})

	t.Run("GetDomainError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomainHSTSConfiguration("testdomain1.co.uk").Return(nil)
		service.EXPECT().GetDomainHSTSConfiguration("testdomain1.co.uk").Return(ddosx.HSTSConfiguration{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving updated HSTS configuration for domain [testdomain1.co.uk]: test error\n", func() {
			ddosxDomainHSTSDisable(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}
