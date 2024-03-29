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

func Test_ddosxDomainACLGeoIPRuleListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainACLGeoIPRuleListCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainACLGeoIPRuleListCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainACLGeoIPRuleList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainACLGeoIPRules("testdomain1.co.uk", gomock.Any()).Return([]ddosx.ACLGeoIPRule{}, nil).Times(1)

		ddosxDomainACLGeoIPRuleList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ddosxDomainACLGeoIPRuleList(service, cmd, []string{"testdomain1.co.uk"})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetDomainACLGeoIPRules_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainACLGeoIPRules("testdomain1.co.uk", gomock.Any()).Return([]ddosx.ACLGeoIPRule{}, errors.New("test error")).Times(1)

		err := ddosxDomainACLGeoIPRuleList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})

		assert.Equal(t, "Error retrieving domain ACL GeoIP rules: test error", err.Error())
	})
}

func Test_ddosxDomainACLGeoIPRuleShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainACLGeoIPRuleShowCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainACLGeoIPRuleShowCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingRule_Error", func(t *testing.T) {
		cmd := ddosxDomainACLGeoIPRuleShowCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing rule", err.Error())
	})
}

func Test_ddosxDomainACLGeoIPRuleShow(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainACLGeoIPRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.ACLGeoIPRule{}, nil)

		ddosxDomainACLGeoIPRuleShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("GetDomainACLGeoIPRule_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainACLGeoIPRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.ACLGeoIPRule{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving domain ACL GeoIP rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainACLGeoIPRuleShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_ddosxDomainACLGeoIPRuleCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainACLGeoIPRuleCreateCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainACLGeoIPRuleCreateCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainACLGeoIPRuleCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainACLGeoIPRuleCreateCmd(nil)
		cmd.Flags().Set("code", "GB")

		expectedRequest := ddosx.CreateACLGeoIPRuleRequest{
			Code: "GB",
		}

		gomock.InOrder(
			service.EXPECT().CreateDomainACLGeoIPRule("testdomain1.co.uk", gomock.Eq(expectedRequest)).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetDomainACLGeoIPRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.ACLGeoIPRule{}, nil),
		)

		ddosxDomainACLGeoIPRuleCreate(service, cmd, []string{"testdomain1.co.uk"})
	})

	t.Run("CreateDomainACLGeoIPRuleError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().CreateDomainACLGeoIPRule("testdomain1.co.uk", gomock.Any()).Return("00000000-0000-0000-0000-000000000000", errors.New("test error")).Times(1)

		err := ddosxDomainACLGeoIPRuleCreate(service, &cobra.Command{}, []string{"testdomain1.co.uk"})

		assert.Equal(t, "Error creating domain ACL GeoIP rule: test error", err.Error())
	})

	t.Run("GetDomainACLGeoIPRuleError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainACLGeoIPRuleCreateCmd(nil)
		cmd.Flags().Set("code", "GB")

		gomock.InOrder(
			service.EXPECT().CreateDomainACLGeoIPRule("testdomain1.co.uk", gomock.Any()).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetDomainACLGeoIPRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.ACLGeoIPRule{}, errors.New("test error")),
		)

		err := ddosxDomainACLGeoIPRuleCreate(service, &cobra.Command{}, []string{"testdomain1.co.uk"})

		assert.Equal(t, "Error retrieving new domain ACL GeoIP rule [00000000-0000-0000-0000-000000000000]: test error", err.Error())
	})
}

func Test_ddosxDomainACLGeoIPRuleUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainACLGeoIPRuleUpdateCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainACLGeoIPRuleUpdateCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingRule_Error", func(t *testing.T) {
		cmd := ddosxDomainACLGeoIPRuleUpdateCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing rule", err.Error())
	})
}

func Test_ddosxDomainACLGeoIPRuleUpdate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainACLGeoIPRuleUpdateCmd(nil)
		cmd.Flags().Set("code", "GB")

		expectedRequest := ddosx.PatchACLGeoIPRuleRequest{
			Code: "GB",
		}

		gomock.InOrder(
			service.EXPECT().PatchDomainACLGeoIPRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Eq(expectedRequest)).Return(nil),
			service.EXPECT().GetDomainACLGeoIPRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.ACLGeoIPRule{}, nil),
		)

		ddosxDomainACLGeoIPRuleUpdate(service, cmd, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("PatchDomainACLGeoIPRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().PatchDomainACLGeoIPRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating domain ACL GeoIP rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainACLGeoIPRuleUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})

	t.Run("GetDomainACLGeoIPRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchDomainACLGeoIPRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(nil),
			service.EXPECT().GetDomainACLGeoIPRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.ACLGeoIPRule{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated domain ACL GeoIP rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainACLGeoIPRuleUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_ddosxDomainACLGeoIPRuleDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainACLGeoIPRuleDeleteCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainACLGeoIPRuleDeleteCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingRule_Error", func(t *testing.T) {
		cmd := ddosxDomainACLGeoIPRuleDeleteCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing rule", err.Error())
	})
}

func Test_ddosxDomainACLGeoIPRuleDelete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomainACLGeoIPRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(nil).Times(1)

		ddosxDomainACLGeoIPRuleDelete(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("DeleteDomainACLGeoIPRule_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomainACLGeoIPRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error removing domain ACL GeoIP rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainACLGeoIPRuleDelete(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}
