package cmd

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

func Test_ddosxDomainWAFAdvancedRuleListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainWAFAdvancedRuleListCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFAdvancedRuleListCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainWAFAdvancedRuleList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainWAFAdvancedRules("testdomain1.co.uk", gomock.Any()).Return([]ddosx.WAFAdvancedRule{}, nil).Times(1)

		ddosxDomainWAFAdvancedRuleList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("MalformedFlag_OutputsFatal", func(t *testing.T) {
		defer func() { flagFilter = nil }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		test_output.AssertFatalOutput(t, "Missing value for filtering\n", func() {
			ddosxDomainWAFAdvancedRuleList(service, cmd, []string{"testdomain1.co.uk"})
		})
	})

	t.Run("GetDomainsError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainWAFAdvancedRules("testdomain1.co.uk", gomock.Any()).Return([]ddosx.WAFAdvancedRule{}, errors.New("test error")).Times(1)

		test_output.AssertFatalOutput(t, "Error retrieving domain WAF advanced rules: test error\n", func() {
			ddosxDomainWAFAdvancedRuleList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}

func Test_ddosxDomainWAFAdvancedRuleShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainWAFAdvancedRuleShowCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFAdvancedRuleShowCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingRule_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFAdvancedRuleShowCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing rule", err.Error())
	})
}

func Test_ddosxDomainWAFAdvancedRuleShow(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainWAFAdvancedRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.WAFAdvancedRule{}, nil)

		ddosxDomainWAFAdvancedRuleShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("GetDomainWAFAdvancedRule_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainWAFAdvancedRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.WAFAdvancedRule{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving domain WAF advanced rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainWAFAdvancedRuleShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_ddosxDomainWAFAdvancedRuleCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainWAFAdvancedRuleCreateCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFAdvancedRuleCreateCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainWAFAdvancedRuleCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainWAFAdvancedRuleCreateCmd()
		cmd.Flags().Set("section", "REQUEST_URI")
		cmd.Flags().Set("modifier", "contains")
		cmd.Flags().Set("phrase", "testphrase")
		cmd.Flags().Set("ip", "1.2.3.4")

		expectedRequest := ddosx.CreateWAFAdvancedRuleRequest{
			Section:  "REQUEST_URI",
			Modifier: ddosx.WAFAdvancedRuleModifierContains,
			Phrase:   "testphrase",
			IP:       "1.2.3.4",
		}

		gomock.InOrder(
			service.EXPECT().CreateDomainWAFAdvancedRule("testdomain1.co.uk", gomock.Eq(expectedRequest)).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetDomainWAFAdvancedRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.WAFAdvancedRule{}, nil),
		)

		ddosxDomainWAFAdvancedRuleCreate(service, cmd, []string{"testdomain1.co.uk"})
	})

	t.Run("InvalidModifier_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainWAFAdvancedRuleCreateCmd()
		cmd.Flags().Set("section", "REQUEST_URI")
		cmd.Flags().Set("modifier", "invalidmodifier")
		cmd.Flags().Set("phrase", "testphrase")
		cmd.Flags().Set("ip", "1.2.3.4")

		test_output.AssertFatalOutputFunc(t, func(stdErr string) {
			assert.Contains(t, stdErr, "Invalid ddosx.WAFAdvancedRuleModifier")
		}, func() {
			ddosxDomainWAFAdvancedRuleCreate(service, cmd, []string{"testdomain1.co.uk"})
		})
	})

	t.Run("CreateDomainWAFAdvancedRuleError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainWAFAdvancedRuleCreateCmd()
		cmd.Flags().Set("section", "REQUEST_URI")
		cmd.Flags().Set("modifier", "contains")
		cmd.Flags().Set("phrase", "testphrase")
		cmd.Flags().Set("ip", "1.2.3.4")

		service.EXPECT().CreateDomainWAFAdvancedRule("testdomain1.co.uk", gomock.Any()).Return("00000000-0000-0000-0000-000000000000", errors.New("test error")).Times(1)

		test_output.AssertFatalOutput(t, "Error creating domain WAF advanced rule: test error\n", func() {
			ddosxDomainWAFAdvancedRuleCreate(service, cmd, []string{"testdomain1.co.uk"})
		})
	})

	t.Run("GetDomainWAFAdvancedRuleError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainWAFAdvancedRuleCreateCmd()
		cmd.Flags().Set("section", "REQUEST_URI")
		cmd.Flags().Set("modifier", "contains")
		cmd.Flags().Set("phrase", "testphrase")
		cmd.Flags().Set("ip", "1.2.3.4")

		gomock.InOrder(
			service.EXPECT().CreateDomainWAFAdvancedRule("testdomain1.co.uk", gomock.Any()).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetDomainWAFAdvancedRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.WAFAdvancedRule{}, errors.New("test error")),
		)

		test_output.AssertFatalOutput(t, "Error retrieving new domain WAF advanced rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainWAFAdvancedRuleCreate(service, cmd, []string{"testdomain1.co.uk"})
		})
	})
}

func Test_ddosxDomainWAFAdvancedRuleUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainWAFAdvancedRuleUpdateCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFAdvancedRuleUpdateCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingAdvancedRule_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFAdvancedRuleUpdateCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing advanced rule", err.Error())
	})
}

func Test_ddosxDomainWAFAdvancedRuleUpdate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainWAFAdvancedRuleUpdateCmd()
		cmd.Flags().Set("section", "REQUEST_URI")
		cmd.Flags().Set("modifier", "contains")
		cmd.Flags().Set("phrase", "testphrase")
		cmd.Flags().Set("ip", "1.2.3.4")

		expectedRequest := ddosx.PatchWAFAdvancedRuleRequest{
			Section:  "REQUEST_URI",
			Modifier: ddosx.WAFAdvancedRuleModifierContains,
			Phrase:   "testphrase",
			IP:       "1.2.3.4",
		}

		gomock.InOrder(
			service.EXPECT().PatchDomainWAFAdvancedRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Eq(expectedRequest)).Return(nil),
			service.EXPECT().GetDomainWAFAdvancedRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.WAFAdvancedRule{}, nil),
		)

		ddosxDomainWAFAdvancedRuleUpdate(service, cmd, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("InvalidModifier_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainWAFAdvancedRuleUpdateCmd()
		cmd.Flags().Set("section", "REQUEST_URI")
		cmd.Flags().Set("modifier", "invalidmodifier")
		cmd.Flags().Set("phrase", "testphrase")
		cmd.Flags().Set("ip", "1.2.3.4")

		test_output.AssertFatalOutputFunc(t, func(stdErr string) {
			assert.Contains(t, stdErr, "Invalid ddosx.WAFAdvancedRuleModifier")
		}, func() {
			ddosxDomainWAFAdvancedRuleUpdate(service, cmd, []string{"testdomain1.co.uk"})
		})
	})

	t.Run("PatchDomainWAFAdvancedRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().PatchDomainWAFAdvancedRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating domain WAF advanced rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainWAFAdvancedRuleUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})

	t.Run("GetDomainWAFAdvancedRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchDomainWAFAdvancedRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(nil),
			service.EXPECT().GetDomainWAFAdvancedRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.WAFAdvancedRule{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated domain WAF advanced rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainWAFAdvancedRuleUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_ddosxDomainWAFAdvancedRuleDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainWAFAdvancedRuleDeleteCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFAdvancedRuleDeleteCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingAdvancedRule_Error", func(t *testing.T) {
		cmd := ddosxDomainWAFAdvancedRuleDeleteCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing advanced rule", err.Error())
	})
}

func Test_ddosxDomainWAFAdvancedRuleDelete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomainWAFAdvancedRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(nil).Times(1)

		ddosxDomainWAFAdvancedRuleDelete(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("DeleteDomainWAFAdvancedRule_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomainWAFAdvancedRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error removing domain WAF advanced rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainWAFAdvancedRuleDelete(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}
