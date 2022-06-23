package ddosx

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/ptr"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ddosxDomainWAFRuleSetListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleSetListCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleSetListCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainWAFRuleSetList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainWAFRuleSets("testdomain1.co.uk", gomock.Any()).Return([]ddosx.WAFRuleSet{}, nil).Times(1)

		ddosxDomainWAFRuleSetList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ddosxDomainWAFRuleSetList(service, cmd, []string{"testdomain1.co.uk"})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetDomainsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainWAFRuleSets("testdomain1.co.uk", gomock.Any()).Return([]ddosx.WAFRuleSet{}, errors.New("test error")).Times(1)

		err := ddosxDomainWAFRuleSetList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})

		assert.Equal(t, "Error retrieving domain waf rule sets: test error", err.Error())
	})
}

func Test_ddosxDomainWAFRuleSetShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleSetShowCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleSetShowCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingRuleSet_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleSetShowCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing rule set", err.Error())
	})
}

func Test_ddosxDomainWAFRuleSetShow(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainWAFRuleSet("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.WAFRuleSet{}, nil)

		ddosxDomainWAFRuleSetShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("GetDomainWAFRuleSet_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainWAFRuleSet("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.WAFRuleSet{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving domain WAF rule set [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainWAFRuleSetShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_ddosxDomainWAFRuleSetUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleSetUpdateCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleSetUpdateCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingRuleSet_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleSetUpdateCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing rule set", err.Error())
	})
}

func Test_ddosxDomainWAFRuleSetUpdate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainWAFRuleSetUpdateCmd(nil)
		cmd.Flags().Set("active", "true")

		expectedRequest := ddosx.PatchWAFRuleSetRequest{
			Active: ptr.Bool(true),
		}

		gomock.InOrder(
			service.EXPECT().PatchDomainWAFRuleSet("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Eq(expectedRequest)).Return(nil),
			service.EXPECT().GetDomainWAFRuleSet("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.WAFRuleSet{}, nil),
		)

		ddosxDomainWAFRuleSetUpdate(service, cmd, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("PatchDomainWAFRuleSet_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().PatchDomainWAFRuleSet("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating domain WAF rule set [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainWAFRuleSetUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})

	t.Run("GetDomainWAFRuleSet_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchDomainWAFRuleSet("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(nil),
			service.EXPECT().GetDomainWAFRuleSet("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.WAFRuleSet{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated domain WAF rule set [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainWAFRuleSetUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}
