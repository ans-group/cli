package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/cli/test"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/ptr"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func Test_ddosxDomainWAFRuleSetListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleSetListCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleSetListCmd()
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

	t.Run("MalformedFlag_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()
		defer func() { flagFilter = nil }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		flagFilter = []string{"invalidfilter"}

		output := test.CatchStdErr(t, func() {
			ddosxDomainWAFRuleSetList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Missing value for filtering\n", output)
	})

	t.Run("GetDomainsError_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainWAFRuleSets("testdomain1.co.uk", gomock.Any()).Return([]ddosx.WAFRuleSet{}, errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			ddosxDomainWAFRuleSetList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error retrieving domain waf rule sets: test error\n", output)
	})
}

func Test_ddosxDomainWAFRuleSetShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleSetShowCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleSetShowCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingRuleSet_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleSetShowCmd()
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

		service.EXPECT().GetDomainWAFRuleSet("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(ddosx.WAFRuleSet{}, nil)

		ddosxDomainWAFRuleSetShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("GetDomainWAFRuleSet_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainWAFRuleSet("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(ddosx.WAFRuleSet{}, errors.New("test error"))

		output := test.CatchStdErr(t, func() {
			ddosxDomainWAFRuleSetShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})

		assert.Equal(t, "Error retrieving domain WAF rule set [00000000-0000-0000-0000-000000000000]: test error\n", output)
	})
}

func Test_ddosxDomainWAFRuleSetUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleSetUpdateCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleSetUpdateCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingRuleSet_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFRuleSetUpdateCmd()
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
		cmd := ddosxDomainWAFRuleSetUpdateCmd()
		cmd.Flags().Set("active", "true")

		expectedRequest := ddosx.PatchWAFRuleSetRequest{
			Active: ptr.Bool(true),
		}

		gomock.InOrder(
			service.EXPECT().PatchDomainWAFRuleSet("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Eq(expectedRequest)).Return(nil),
			service.EXPECT().GetDomainWAFRuleSet("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(ddosx.WAFRuleSet{}, nil),
		)

		ddosxDomainWAFRuleSetUpdate(service, cmd, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("PatchDomainWAFRuleSet_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().PatchDomainWAFRuleSet("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(errors.New("test error"))

		output := test.CatchStdErr(t, func() {
			ddosxDomainWAFRuleSetUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})

		assert.Equal(t, "Error updating domain WAF rule set [00000000-0000-0000-0000-000000000000]: test error\n", output)
	})

	t.Run("GetDomainWAFRuleSet_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchDomainWAFRuleSet("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(nil),
			service.EXPECT().GetDomainWAFRuleSet("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(ddosx.WAFRuleSet{}, errors.New("test error")),
		)

		output := test.CatchStdErr(t, func() {
			ddosxDomainWAFRuleSetUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})

		assert.Equal(t, "Error retrieving updated domain WAF rule set [00000000-0000-0000-0000-000000000000]: test error\n", output)
	})
}
