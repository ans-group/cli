package ddosx

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ddosxDomainWAFRuleListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleListCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleListCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainWAFRuleList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainWAFRules("testdomain1.co.uk", gomock.Any()).Return([]ddosx.WAFRule{}, nil).Times(1)

		ddosxDomainWAFRuleList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ddosxDomainWAFRuleList(service, cmd, []string{"testdomain1.co.uk"})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetDomainsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainWAFRules("testdomain1.co.uk", gomock.Any()).Return([]ddosx.WAFRule{}, errors.New("test error")).Times(1)

		err := ddosxDomainWAFRuleList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})

		assert.Equal(t, "Error retrieving domain WAF rules: test error", err.Error())
	})
}

func Test_ddosxDomainWAFRuleShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleShowCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleShowCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingRule_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleShowCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing rule", err.Error())
	})
}

func Test_ddosxDomainWAFRuleShow(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainWAFRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.WAFRule{}, nil)

		ddosxDomainWAFRuleShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("GetDomainWAFRule_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainWAFRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.WAFRule{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving domain WAF rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainWAFRuleShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_ddosxDomainWAFRuleCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleCreateCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleCreateCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainWAFRuleCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainWAFRuleCreateCmd(nil)
		cmd.Flags().Set("uri", "test.html")

		expectedRequest := ddosx.CreateWAFRuleRequest{
			URI: "test.html",
		}

		gomock.InOrder(
			service.EXPECT().CreateDomainWAFRule("testdomain1.co.uk", gomock.Eq(expectedRequest)).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetDomainWAFRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.WAFRule{}, nil),
		)

		ddosxDomainWAFRuleCreate(service, cmd, []string{"testdomain1.co.uk"})
	})

	t.Run("CreateDomainWAFRuleError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().CreateDomainWAFRule("testdomain1.co.uk", gomock.Any()).Return("00000000-0000-0000-0000-000000000000", errors.New("test error")).Times(1)

		err := ddosxDomainWAFRuleCreate(service, &cobra.Command{}, []string{"testdomain1.co.uk"})

		assert.Equal(t, "Error creating domain WAF rule: test error", err.Error())
	})

	t.Run("GetDomainWAFRuleError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().CreateDomainWAFRule("testdomain1.co.uk", gomock.Any()).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetDomainWAFRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.WAFRule{}, errors.New("test error")),
		)

		err := ddosxDomainWAFRuleCreate(service, &cobra.Command{}, []string{"testdomain1.co.uk"})

		assert.Equal(t, "Error retrieving new domain WAF rule [00000000-0000-0000-0000-000000000000]: test error", err.Error())
	})
}

func Test_ddosxDomainWAFRuleUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleUpdateCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleUpdateCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingRule_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleUpdateCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing rule", err.Error())
	})
}

func Test_ddosxDomainWAFRuleUpdate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainWAFRuleUpdateCmd(nil)
		cmd.Flags().Set("uri", "test.html")
		cmd.Flags().Set("ip", "1.2.3.4")

		expectedRequest := ddosx.PatchWAFRuleRequest{
			URI: "test.html",
			IP:  "1.2.3.4",
		}

		gomock.InOrder(
			service.EXPECT().PatchDomainWAFRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Eq(expectedRequest)).Return(nil),
			service.EXPECT().GetDomainWAFRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.WAFRule{}, nil),
		)

		ddosxDomainWAFRuleUpdate(service, cmd, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("PatchDomainWAFRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().PatchDomainWAFRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating domain WAF rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainWAFRuleUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})

	t.Run("GetDomainWAFRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchDomainWAFRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(nil),
			service.EXPECT().GetDomainWAFRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.WAFRule{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated domain WAF rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainWAFRuleUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_ddosxDomainWAFRuleDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleDeleteCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleDeleteCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingRule_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleDeleteCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing rule", err.Error())
	})
}

func Test_ddosxDomainWAFRuleDelete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomainWAFRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(nil).Times(1)

		ddosxDomainWAFRuleDelete(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("DeleteDomainWAFRule_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomainWAFRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error removing domain WAF rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainWAFRuleDelete(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}
