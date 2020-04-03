package ddosx

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

func Test_ddosxDomainHSTSRuleListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainHSTSRuleListCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainHSTSRuleListCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainHSTSRuleList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainHSTSRules("testdomain1.co.uk", gomock.Any()).Return([]ddosx.HSTSRule{}, nil).Times(1)

		ddosxDomainHSTSRuleList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ddosxDomainHSTSRuleList(service, cmd, []string{"testdomain1.co.uk"})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetDomainsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainHSTSRules("testdomain1.co.uk", gomock.Any()).Return([]ddosx.HSTSRule{}, errors.New("test error")).Times(1)

		err := ddosxDomainHSTSRuleList(service, &cobra.Command{}, []string{"testdomain1.co.uk"})

		assert.Equal(t, "Error retrieving HSTS rules: test error", err.Error())
	})
}

func Test_ddosxDomainHSTSRuleShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainHSTSRuleShowCmd(nil).Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		err := ddosxDomainHSTSRuleShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingHSTSRule_Error", func(t *testing.T) {
		err := ddosxDomainHSTSRuleShowCmd(nil).Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing rule", err.Error())
	})
}

func Test_ddosxDomainHSTSRuleShow(t *testing.T) {
	t.Run("SingleDomain", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainHSTSRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.HSTSRule{}, nil).Times(1)

		ddosxDomainHSTSRuleShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("MultipleDomains", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetDomainHSTSRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.HSTSRule{}, nil),
			service.EXPECT().GetDomainHSTSRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000001").Return(ddosx.HSTSRule{}, nil),
		)

		ddosxDomainHSTSRuleShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
	})

	t.Run("GetDomainHSTSRuleError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomainHSTSRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.HSTSRule{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving HSTS rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainHSTSRuleShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_ddosxDomainHSTSRuleCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainHSTSRuleCreateCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainHSTSRuleCreateCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainHSTSRuleCreate(t *testing.T) {
	t.Run("Valid_CreatesRule", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainHSTSRuleCreateCmd(nil)
		cmd.Flags().Set("type", "record")
		cmd.Flags().Set("record-name", "example.com")

		gomock.InOrder(
			service.EXPECT().CreateDomainHSTSRule("testdomain1.co.uk", gomock.Any()).Do(func(id string, req ddosx.CreateHSTSRuleRequest) {
				if req.Type != ddosx.HSTSRuleTypeRecord {
					t.Fatalf("Expected Type 'HSTSRuleTypeRecord', got '%s'", req.Type)
				}
				if req.RecordName == nil || *req.RecordName != "example.com" {
					t.Fatal("Expected RecordName 'example.com'")
				}
			}).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetDomainHSTSRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.HSTSRule{}, nil),
		)

		ddosxDomainHSTSRuleCreate(service, cmd, []string{"testdomain1.co.uk"})
	})

	t.Run("ParseHSTSRuleTypeError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainHSTSRuleCreateCmd(nil)
		cmd.Flags().Set("type", "invalid")

		err := ddosxDomainHSTSRuleCreate(service, cmd, []string{"testdomain1.co.uk"})
		assert.NotNil(t, err)
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("CreateDomainHSTSRuleError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainHSTSRuleCreateCmd(nil)
		cmd.Flags().Set("type", "record")

		service.EXPECT().CreateDomainHSTSRule("testdomain1.co.uk", gomock.Any()).Return("", errors.New("test error"))

		err := ddosxDomainHSTSRuleCreate(service, cmd, []string{"testdomain1.co.uk"})
		assert.NotNil(t, err)
		assert.Equal(t, "Error creating HSTS rule: test error", err.Error())
	})

	t.Run("GetDomainHSTSRuleError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainHSTSRuleCreateCmd(nil)
		cmd.Flags().Set("type", "record")

		gomock.InOrder(
			service.EXPECT().CreateDomainHSTSRule("testdomain1.co.uk", gomock.Any()).Return("00000000-0000-0000-0000-000000000000", nil),
			service.EXPECT().GetDomainHSTSRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.HSTSRule{}, errors.New("test error")),
		)

		err := ddosxDomainHSTSRuleCreate(service, cmd, []string{"testdomain1.co.uk"})
		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving new HSTS rule [00000000-0000-0000-0000-000000000000]: test error", err.Error())
	})
}

func Test_ddosxDomainHSTSRuleUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainHSTSRuleUpdateCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainHSTSRuleUpdateCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingRule_Error", func(t *testing.T) {
		cmd := ddosxDomainHSTSRuleUpdateCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing rule", err.Error())
	})
}

func Test_ddosxDomainHSTSRuleUpdate(t *testing.T) {
	t.Run("Valid_UpdatesRule", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainHSTSRuleUpdateCmd(nil)
		cmd.Flags().Set("max-age", "300")
		cmd.Flags().Set("preload", "true")
		cmd.Flags().Set("include-subdomains", "true")

		expectedRequest := ddosx.PatchHSTSRuleRequest{
			MaxAge:            ptr.Int(300),
			Preload:           ptr.Bool(true),
			IncludeSubdomains: ptr.Bool(true),
		}

		gomock.InOrder(
			service.EXPECT().PatchDomainHSTSRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Eq(expectedRequest)).Return(nil),
			service.EXPECT().GetDomainHSTSRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.HSTSRule{}, nil),
		)

		ddosxDomainHSTSRuleUpdate(service, cmd, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("PatchDomainHSTSRuleError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().PatchDomainHSTSRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating domain HSTS rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainHSTSRuleUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})

	t.Run("GetDomainHSTSRuleError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchDomainHSTSRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000", gomock.Any()).Return(nil),
			service.EXPECT().GetDomainHSTSRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(ddosx.HSTSRule{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated HSTS rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainHSTSRuleUpdate(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}

func Test_ddosxDomainHSTSRuleDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := ddosxDomainHSTSRuleDeleteCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingDomain_Error", func(t *testing.T) {
		cmd := ddosxDomainHSTSRuleDeleteCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})

	t.Run("MissingRule_Error", func(t *testing.T) {
		cmd := ddosxDomainHSTSRuleDeleteCmd(nil)
		err := cmd.Args(nil, []string{"testdomain1.co.uk"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing rule", err.Error())
	})
}

func Test_ddosxDomainHSTSRuleDelete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomainHSTSRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(nil).Times(1)

		ddosxDomainHSTSRuleDelete(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("DeleteDomainHSTSRule_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomainHSTSRule("testdomain1.co.uk", "00000000-0000-0000-0000-000000000000").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error removing domain HSTS rule [00000000-0000-0000-0000-000000000000]: test error\n", func() {
			ddosxDomainHSTSRuleDelete(service, &cobra.Command{}, []string{"testdomain1.co.uk", "00000000-0000-0000-0000-000000000000"})
		})
	})
}
