package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/ptr"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func Test_ddosxDomainACLIPRuleListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainACLIPRuleListCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainACLIPRuleListCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainACLIPRuleList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainACLIPRules("testdomain1.co.uk", gomock.Any()).Return([]ddosx.ACLIPRule{}, nil).Times(1)

		ddosxDomainACLIPRuleList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("MalformedFlag_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		test_output.AssertFatalOutput(t, "Missing value for filtering\n", func() {
			ddosxDomainACLIPRuleList(service, cmd, []string{"testdomain1.co.uk"})
		})
	})

	t.Run("GetDomainsError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainACLIPRules("testdomain1.co.uk", gomock.Any()).Return([]ddosx.ACLIPRule{}, errors.New("test error")).Times(1)

		test_output.AssertFatalOutput(t, "Error retrieving domain ACL IP rules: test error\n", func() {
			ddosxDomainACLIPRuleList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}

func Test_ddosxDomainACLIPRuleShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainACLIPRuleShowCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainACLIPRuleShowCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingRule_Error", func(t *testing.T) {
		cmd := ddosxDomainACLIPRuleShowCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing rule", err.Error())
	})
}

func Test_ddosxDomainACLIPRuleShow(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainACLIPRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.ACLIPRule{}, nil)

		ddosxDomainACLIPRuleShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("GetDomainACLIPRule_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainACLIPRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.ACLIPRule{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving domain ACL IP rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainACLIPRuleShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_ddosxDomainACLIPRuleCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainACLIPRuleCreateCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainACLIPRuleCreateCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainACLIPRuleCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainACLIPRuleCreateCmd()
		cmd.Flags().Set("ip", "1.2.3.4")
		cmd.Flags().Set("uri", "testuri")
		cmd.Flags().Set("mode", "allow")

		expectedRequest := ddosx.CreateACLIPRuleRequest{
			IP:   "1.2.3.4",
			Mode: ddosx.ACLIPModeAllow,
			URI:  "testuri",
		}

		gomock.InOrder(
			service.EXPECT().CreateDomainACLIPRule("testdomain1.co.uk", gomock.Eq(expectedRequest)).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetDomainACLIPRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.ACLIPRule{}, nil),
		)

		ddosxDomainACLIPRuleCreate(service, cmd, []string{"testdomain1.co.uk"})
	})

	t.Run("InvalidMode_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainACLIPRuleCreateCmd()
		cmd.Flags().Set("ip", "1.2.3.4")
		cmd.Flags().Set("uri", "testuri")
		cmd.Flags().Set("mode", "invalidmode")

		err := ddosxDomainACLIPRuleCreate(service, cmd, []string{"testdomain1.co.uk"})
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("CreateDomainACLIPRuleError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainACLIPRuleCreateCmd()
		cmd.Flags().Set("ip", "1.2.3.4")
		cmd.Flags().Set("uri", "testuri")
		cmd.Flags().Set("mode", "allow")

		service.EXPECT().CreateDomainACLIPRule("testdomain1.co.uk", gomock.Any()).Return("00000000-0000-0000-0000-000000000000", errors.New("test error")).Times(1)

		err := ddosxDomainACLIPRuleCreate(service, cmd, []string{"testdomain1.co.uk"})
		assert.Equal(t, "Error creating domain ACL IP rule: test error", err.Error())
	})

	t.Run("CreateDomainACLIPRuleError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainACLIPRuleCreateCmd()
		cmd.Flags().Set("ip", "1.2.3.4")
		cmd.Flags().Set("uri", "testuri")
		cmd.Flags().Set("mode", "allow")

		gomock.InOrder(
			service.EXPECT().CreateDomainACLIPRule("testdomain1.co.uk", gomock.Any()).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetDomainACLIPRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.ACLIPRule{}, errors.New("test error")),
		)

		err := ddosxDomainACLIPRuleCreate(service, cmd, []string{"testdomain1.co.uk"})
		assert.Equal(t, "Error retrieving new domain ACL IP rule [00000000-0000-0000-0000-000000000000]: test error", err.Error())
	})
}

func Test_ddosxDomainACLIPRuleUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainACLIPRuleUpdateCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainACLIPRuleUpdateCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingRule_Error", func(t *testing.T) {
		cmd := ddosxDomainACLIPRuleUpdateCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing rule", err.Error())
	})
}

func Test_ddosxDomainACLIPRuleUpdate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainACLIPRuleUpdateCmd()
		cmd.Flags().Set("mode", "deny")
		cmd.Flags().Set("uri", "testuri")
		cmd.Flags().Set("ip", "1.2.3.4")

		expectedRequest := ddosx.PatchACLIPRuleRequest{
			Mode: ddosx.ACLIPModeDeny,
			URI:  ptr.String("testuri"),
			IP:   "1.2.3.4",
		}

		gomock.InOrder(
			service.EXPECT().PatchDomainACLIPRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Eq(expectedRequest)).Return(nil),
			service.EXPECT().GetDomainACLIPRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.ACLIPRule{}, nil),
		)

		ddosxDomainACLIPRuleUpdate(service, cmd, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("InvalidMode_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainACLIPRuleUpdateCmd()
		cmd.Flags().Set("mode", "invalidmode")

		test_output.AssertFatalOutputFunc(t, func(stdErr string) {
			assert.Contains(t, stdErr, "Invalid ddosx.ACLIPMode")
		}, func() {
			ddosxDomainACLIPRuleUpdate(service, cmd, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})

	t.Run("UpdateDomainACLIPRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().PatchDomainACLIPRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating domain ACL IP rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainACLIPRuleUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})

	t.Run("GetDomainACLIPRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchDomainACLIPRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(nil),
			service.EXPECT().GetDomainACLIPRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.ACLIPRule{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated domain ACL IP rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainACLIPRuleUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_ddosxDomainACLIPRuleDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainACLIPRuleDeleteCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainACLIPRuleDeleteCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingRule_Error", func(t *testing.T) {
		cmd := ddosxDomainACLIPRuleDeleteCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing rule", err.Error())
	})
}

func Test_ddosxDomainACLIPRuleDelete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomainACLIPRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(nil).Times(1)

		ddosxDomainACLIPRuleDelete(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("DeleteDomainACLIPRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomainACLIPRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error removing domain ACL IP rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainACLIPRuleDelete(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}
