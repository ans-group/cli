package ddosx

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func Test_ddosxDomainWAFShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainWAFShowCmd(nil).Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ddosxDomainWAFShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainWAFShow(t *testing.T) {
	t.Run("SingleDomainWAF", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainWAF("testdomain1.co.uk").Return(ddosx.WAF{}, nil).Times(1)

		ddosxDomainWAFShow(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("MultipleDomainWAFs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetDomainWAF("testdomain1.co.uk").Return(ddosx.WAF{}, nil),
			service.EXPECT().GetDomainWAF("testdomain2.co.uk").Return(ddosx.WAF{}, nil),
		)

		ddosxDomainWAFShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "testdomain2.co.uk"})
	})

	t.Run("GetDomainWAFError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainWAF("testdomain1.co.uk").Return(ddosx.WAF{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving domain waf [testdomain1.co.uk]: test error\n", func() {
			ddosxDomainWAFShow(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}

func Test_ddosxDomainWAFCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainWAFCreateCmd(nil).Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ddosxDomainWAFCreateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainWAFCreate(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainWAFCreateCmd(nil)
		cmd.Flags().Set("mode", "on")
		cmd.Flags().Set("paranoia-level", "high")

		expectedRequest := ddosx.CreateWAFRequest{
			Mode:          ddosx.WAFModeOn,
			ParanoiaLevel: ddosx.WAFParanoiaLevelHigh,
		}

		gomock.InOrder(
			service.EXPECT().CreateDomainWAF("testdomain1.co.uk", gomock.Eq(expectedRequest)).Return(nil),
			service.EXPECT().GetDomainWAF("testdomain1.co.uk").Return(ddosx.WAF{}, nil),
		)

		ddosxDomainWAFCreate(service, cmd, []string{"testdomain1.co.uk"})
	})

	t.Run("InvalidMode_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainWAFCreateCmd(nil)
		cmd.Flags().Set("mode", "invalidmode")
		cmd.Flags().Set("paranoia-level", "high")

		err := ddosxDomainWAFCreate(service, cmd, []string{"testdomain1.co.uk"})

		assert.IsType(t, &connection.ErrInvalidEnumValue{}, err)
	})

	t.Run("InvalidParanoiaLevel_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainWAFCreateCmd(nil)
		cmd.Flags().Set("mode", "on")
		cmd.Flags().Set("paranoia-level", "invalidparanoialevel")

		err := ddosxDomainWAFCreate(service, cmd, []string{"testdomain1.co.uk"})

		assert.IsType(t, &connection.ErrInvalidEnumValue{}, err)
	})

	t.Run("CreateDomainWAFError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainWAFCreateCmd(nil)
		cmd.Flags().Set("mode", "on")
		cmd.Flags().Set("paranoia-level", "high")

		service.EXPECT().CreateDomainWAF("testdomain1.co.uk", gomock.Any()).Return(errors.New("test error"))

		err := ddosxDomainWAFCreate(service, cmd, []string{"testdomain1.co.uk"})

		assert.Equal(t, "Error creating domain waf: test error", err.Error())
	})

	t.Run("GetDomainWAFError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainWAFCreateCmd(nil)
		cmd.Flags().Set("mode", "on")
		cmd.Flags().Set("paranoia-level", "high")

		gomock.InOrder(
			service.EXPECT().CreateDomainWAF("testdomain1.co.uk", gomock.Any()).Return(nil),
			service.EXPECT().GetDomainWAF("testdomain1.co.uk").Return(ddosx.WAF{}, errors.New("test error")),
		)

		err := ddosxDomainWAFCreate(service, cmd, []string{"testdomain1.co.uk"})

		assert.Equal(t, "Error retrieving domain waf: test error", err.Error())
	})
}

func Test_ddosxDomainWAFUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainWAFUpdateCmd(nil).Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ddosxDomainWAFUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainWAFUpdate(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainWAFUpdateCmd(nil)
		cmd.Flags().Set("mode", "on")
		cmd.Flags().Set("paranoia-level", "high")

		expectedRequest := ddosx.PatchWAFRequest{
			Mode:          ddosx.WAFModeOn,
			ParanoiaLevel: ddosx.WAFParanoiaLevelHigh,
		}

		gomock.InOrder(
			service.EXPECT().PatchDomainWAF("testdomain1.co.uk", gomock.Eq(expectedRequest)).Return(nil),
			service.EXPECT().GetDomainWAF("testdomain1.co.uk").Return(ddosx.WAF{}, nil),
		)

		ddosxDomainWAFUpdate(service, cmd, []string{"testdomain1.co.uk"})
	})

	t.Run("InvalidMode_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainWAFUpdateCmd(nil)
		cmd.Flags().Set("mode", "invalidmode")
		cmd.Flags().Set("paranoia-level", "high")

		err := ddosxDomainWAFUpdate(service, cmd, []string{"testdomain1.co.uk"})

		assert.IsType(t, &connection.ErrInvalidEnumValue{}, err)
	})

	t.Run("InvalidParanoiaLevel_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainWAFUpdateCmd(nil)
		cmd.Flags().Set("mode", "on")
		cmd.Flags().Set("paranoia-level", "invalidparanoialevel")

		err := ddosxDomainWAFUpdate(service, cmd, []string{"testdomain1.co.uk"})

		assert.IsType(t, &connection.ErrInvalidEnumValue{}, err)
	})

	t.Run("UpdateDomainWAFError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainWAFUpdateCmd(nil)
		cmd.Flags().Set("mode", "on")
		cmd.Flags().Set("paranoia-level", "high")

		service.EXPECT().PatchDomainWAF("testdomain1.co.uk", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating domain waf [testdomain1.co.uk]: test error\n", func() {
			ddosxDomainWAFUpdate(service, cmd, []string{"testdomain1.co.uk"})
		})
	})

	t.Run("GetDomainWAFError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainWAFUpdateCmd(nil)
		cmd.Flags().Set("mode", "on")
		cmd.Flags().Set("paranoia-level", "high")

		gomock.InOrder(
			service.EXPECT().PatchDomainWAF("testdomain1.co.uk", gomock.Any()).Return(nil),
			service.EXPECT().GetDomainWAF("testdomain1.co.uk").Return(ddosx.WAF{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated domain waf [testdomain1.co.uk]: test error\n", func() {
			ddosxDomainWAFUpdate(service, cmd, []string{"testdomain1.co.uk"})
		})
	})
}

func Test_ddosxDomainWAFDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainWAFDeleteCmd(nil).Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ddosxDomainWAFDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainWAFDelete(t *testing.T) {
	t.Run("SingleDomainWAF", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomainWAF("testdomain1.co.uk").Return(nil).Times(1)

		ddosxDomainWAFDelete(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("MultipleDomainWAFs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteDomainWAF("testdomain1.co.uk").Return(nil),
			service.EXPECT().DeleteDomainWAF("testdomain2.co.uk").Return(nil),
		)

		ddosxDomainWAFDelete(service, &cobra.Command{}, []string{"testdomain1.co.uk", "testdomain2.co.uk"})
	})

	t.Run("DeleteDomainWAFError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomainWAF("testdomain1.co.uk").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error removing domain waf [testdomain1.co.uk]: test error\n", func() {
			ddosxDomainWAFDelete(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}
