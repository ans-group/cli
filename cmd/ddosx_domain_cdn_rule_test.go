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

func Test_ddosxDomainCDNRuleListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainCDNRuleListCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainCDNRuleListCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainCDNRuleList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainCDNRules("testdomain1.co.uk", gomock.Any()).Return([]ddosx.CDNRule{}, nil).Times(1)

		ddosxDomainCDNRuleList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("MalformedFlag_OutputsFatal", func(t *testing.T) {
		defer func() { flagFilter = nil }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		flagFilter = []string{"invalidfilter"}

		test_output.AssertFatalOutput(t, "Missing value for filtering\n", func() {
			ddosxDomainCDNRuleList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})

	t.Run("GetDomainsError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainCDNRules("testdomain1.co.uk", gomock.Any()).Return([]ddosx.CDNRule{}, errors.New("test error")).Times(1)

		test_output.AssertFatalOutput(t, "Error retrieving CDN rules: test error\n", func() {
			ddosxDomainCDNRuleList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}

func Test_ddosxDomainCDNRuleShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainCDNRuleShowCmd().Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		err := ddosxDomainCDNRuleShowCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingCDNRule_Error", func(t *testing.T) {
		err := ddosxDomainCDNRuleShowCmd().Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing rule", err.Error())
	})
}

func Test_ddosxDomainCDNRuleShow(t *testing.T) {
	t.Run("SingleDomain", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainCDNRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.CDNRule{}, nil).Times(1)

		ddosxDomainCDNRuleShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("MultipleDomains", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetDomainCDNRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.CDNRule{}, nil),
			service.EXPECT().GetDomainCDNRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000001").Return(ddosx.CDNRule{}, nil),
		)

		ddosxDomainCDNRuleShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
	})

	t.Run("GetDomainCDNRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainCDNRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.CDNRule{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving CDN rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainCDNRuleShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_ddosxDomainCDNRuleCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainCDNRuleCreateCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainCDNRuleCreateCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainCDNRuleCreate(t *testing.T) {
	t.Run("Valid_CreatesRule", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainCDNRuleCreateCmd()
		cmd.Flags().Set("uri", "test.html")
		cmd.Flags().Set("cache-control", "custom")
		cmd.Flags().Set("mime-type", "test/*")
		cmd.Flags().Set("type", "global")

		expectedRequest := ddosx.CreateCDNRuleRequest{
			URI:          "test.html",
			CacheControl: ddosx.CDNRuleCacheControlCustom,
			MimeTypes:    []string{"test/*"},
			Type:         ddosx.CDNRuleTypeGlobal,
		}

		gomock.InOrder(
			service.EXPECT().CreateDomainCDNRule("testdomain1.co.uk", gomock.Eq(expectedRequest)).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetDomainCDNRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.CDNRule{}, nil),
		)

		ddosxDomainCDNRuleCreate(service, cmd, []string{"testdomain1.co.uk"})
	})

	t.Run("ParseCDNRuleCacheControlError_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainCDNRuleCreateCmd()
		cmd.Flags().Set("cache-control", "invalid")

		test_output.AssertFatalOutputFunc(t, func() {
			ddosxDomainCDNRuleCreate(service, cmd, []string{"testdomain1.co.uk"})
		}, func(stdErr string) {
			assert.Contains(t, stdErr, "Invalid value 'invalid' provided for 'cache-control'")
		})
	})

	t.Run("ParseCDNRuleTypeError_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainCDNRuleCreateCmd()
		cmd.Flags().Set("cache-control", "custom")
		cmd.Flags().Set("type", "invalid")

		test_output.AssertFatalOutputFunc(t, func() {
			ddosxDomainCDNRuleCreate(service, cmd, []string{"testdomain1.co.uk"})
		}, func(stdErr string) {
			assert.Contains(t, stdErr, "Invalid value 'invalid' provided for 'type'")
		})
	})

	t.Run("ParseDurationError_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainCDNRuleCreateCmd()
		cmd.Flags().Set("cache-control", "custom")
		cmd.Flags().Set("type", "per-uri")
		cmd.Flags().Set("cache-control-duration", "invalid")

		test_output.AssertFatalOutputFunc(t, func() {
			ddosxDomainCDNRuleCreate(service, cmd, []string{"testdomain1.co.uk"})
		}, func(stdErr string) {
			assert.Contains(t, stdErr, "Invalid value 'invalid' provided for 'cache-control-duration'")
		})
	})

	t.Run("CreateDomainCDNRuleError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainCDNRuleCreateCmd()
		cmd.Flags().Set("uri", "test.html")
		cmd.Flags().Set("cache-control", "custom")
		cmd.Flags().Set("mime-type", "test/*")
		cmd.Flags().Set("type", "global")

		service.EXPECT().CreateDomainCDNRule("testdomain1.co.uk", gomock.Any()).Return("", errors.New("test error"))

		test_output.AssertFatalOutput(t, "Error creating CDN rule: test error\n", func() {
			ddosxDomainCDNRuleCreate(service, cmd, []string{"testdomain1.co.uk"})
		})
	})

	t.Run("GetDomainCDNRuleError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainCDNRuleCreateCmd()
		cmd.Flags().Set("uri", "test.html")
		cmd.Flags().Set("cache-control", "custom")
		cmd.Flags().Set("mime-type", "test/*")
		cmd.Flags().Set("type", "global")

		gomock.InOrder(
			service.EXPECT().CreateDomainCDNRule("testdomain1.co.uk", gomock.Any()).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetDomainCDNRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.CDNRule{}, errors.New("test error")),
		)

		test_output.AssertFatalOutput(t, "Error retrieving new CDN rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainCDNRuleCreate(service, cmd, []string{"testdomain1.co.uk"})
		})
	})
}

func Test_ddosxDomainCDNRuleUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainCDNRuleUpdateCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainCDNRuleUpdateCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingRule_Error", func(t *testing.T) {
		cmd := ddosxDomainCDNRuleUpdateCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing rule", err.Error())
	})
}

func Test_ddosxDomainCDNRuleUpdate(t *testing.T) {
	t.Run("Valid_UpdatesRule", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainCDNRuleUpdateCmd()
		cmd.Flags().Set("uri", "test.html")
		cmd.Flags().Set("cache-control", "custom")
		cmd.Flags().Set("mime-type", "test/*")
		cmd.Flags().Set("type", "global")

		expectedRequest := ddosx.PatchCDNRuleRequest{
			URI:          "test.html",
			CacheControl: ddosx.CDNRuleCacheControlCustom,
			MimeTypes:    []string{"test/*"},
			Type:         ddosx.CDNRuleTypeGlobal,
		}

		gomock.InOrder(
			service.EXPECT().PatchDomainCDNRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Eq(expectedRequest)).Return(nil),
			service.EXPECT().GetDomainCDNRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.CDNRule{}, nil),
		)

		ddosxDomainCDNRuleUpdate(service, cmd, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("ParseCDNRuleCacheControlError_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainCDNRuleUpdateCmd()
		cmd.Flags().Set("cache-control", "invalid")

		test_output.AssertFatalOutputFunc(t, func() {
			ddosxDomainCDNRuleUpdate(service, cmd, []string{"testdomain1.co.uk"})
		}, func(stdErr string) {
			assert.Contains(t, stdErr, "Invalid value 'invalid' provided for 'cache-control'")
		})
	})

	t.Run("ParseCDNRuleTypeError_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainCDNRuleUpdateCmd()
		cmd.Flags().Set("type", "invalid")

		test_output.AssertFatalOutputFunc(t, func() {
			ddosxDomainCDNRuleUpdate(service, cmd, []string{"testdomain1.co.uk"})
		}, func(stdErr string) {
			assert.Contains(t, stdErr, "Invalid value 'invalid' provided for 'type'")
		})
	})

	t.Run("ParseDurationError_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainCDNRuleUpdateCmd()
		cmd.Flags().Set("cache-control-duration", "invalid")

		test_output.AssertFatalOutputFunc(t, func() {
			ddosxDomainCDNRuleUpdate(service, cmd, []string{"testdomain1.co.uk"})
		}, func(stdErr string) {
			assert.Contains(t, stdErr, "Invalid value 'invalid' provided for 'cache-control-duration'")
		})
	})

	t.Run("PatchDomainCDNRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().PatchDomainCDNRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating domain CDN rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainCDNRuleUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})

	t.Run("GetDomainCDNRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchDomainCDNRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(nil),
			service.EXPECT().GetDomainCDNRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.CDNRule{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated CDN rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainCDNRuleUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_ddosxDomainCDNRuleDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainCDNRuleDeleteCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainCDNRuleDeleteCmd()
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingRule_Error", func(t *testing.T) {
		cmd := ddosxDomainCDNRuleDeleteCmd()
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing rule", err.Error())
	})
}

func Test_ddosxDomainCDNRuleDelete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomainCDNRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(nil).Times(1)

		ddosxDomainCDNRuleDelete(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("DeleteDomainCDNRule_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomainCDNRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error removing domain CDN rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainCDNRuleDelete(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}
