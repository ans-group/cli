package ddosx

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func Test_ddosxDomainList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomains(gomock.Any()).Return([]ddosx.Domain{}, nil).Times(1)

		ddosxDomainList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ddosxDomainList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetDomainsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomains(gomock.Any()).Return([]ddosx.Domain{}, errors.New("test error")).Times(1)

		err := ddosxDomainList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving domains: test error", err.Error())
	})
}

func Test_ddosxDomainShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainShowCmd(nil).Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ddosxDomainShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainShow(t *testing.T) {
	t.Run("SingleDomain", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{}, nil).Times(1)

		ddosxDomainShow(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("MultipleDomains", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{}, nil),
			service.EXPECT().GetDomain("testdomain2.co.uk").Return(ddosx.Domain{}, nil),
		)

		ddosxDomainShow(service, &cobra.Command{}, []string{"testdomain1.co.uk", "testdomain2.co.uk"})
	})

	t.Run("GetDomainError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving domain [testdomain1.co.uk]: test error\n", func() {
			ddosxDomainShow(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}

func Test_ddosxDomainCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainCreateCmd(nil)
		cmd.Flags().Set("name", "testdomain1.co.uk")

		expectedRequest := ddosx.CreateDomainRequest{
			Name: "testdomain1.co.uk",
		}

		gomock.InOrder(
			service.EXPECT().CreateDomain(gomock.Eq(expectedRequest)).Return(nil),
			service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{Name: "testdomain1.co.uk"}, nil),
		)

		ddosxDomainCreate(service, cmd, []string{})
	})

	t.Run("CreateDomain_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().CreateDomain(gomock.Any()).Return(errors.New("test error")).Times(1)

		err := ddosxDomainCreate(service, &cobra.Command{}, []string{"testdomain1.co.uk"})

		assert.Equal(t, "Error creating domain: test error", err.Error())
	})

	t.Run("CreateDomain_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainCreateCmd(nil)
		cmd.Flags().Set("name", "testdomain1.co.uk")

		gomock.InOrder(
			service.EXPECT().CreateDomain(gomock.Any()).Return(nil),
			service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{}, errors.New("test error")),
		)

		err := ddosxDomainCreate(service, cmd, []string{"testdomain1.co.uk"})

		assert.Equal(t, "Error retrieving new domain: test error", err.Error())
	})
}

func Test_ddosxDomainDelete(t *testing.T) {
	t.Run("SingleDomain", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		cmd := ddosxDomainDeleteCmd(nil)
		cmd.ParseFlags([]string{"--summary=testsummary", "--description=testdescription"})

		req := ddosx.DeleteDomainRequest{
			Summary:     "testsummary",
			Description: "testdescription",
		}

		service.EXPECT().DeleteDomain("testdomain1.co.uk", gomock.Eq(req)).Return(nil)

		ddosxDomainDelete(service, cmd, []string{"testdomain1.co.uk"})
	})

	t.Run("MultipleDomains", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		cmd := ddosxDomainDeleteCmd(nil)
		cmd.ParseFlags([]string{"--summary=testsummary", "--description=testdescription"})

		req := ddosx.DeleteDomainRequest{
			Summary:     "testsummary",
			Description: "testdescription",
		}

		gomock.InOrder(
			service.EXPECT().DeleteDomain("testdomain1.co.uk", gomock.Eq(req)).Return(nil),
			service.EXPECT().DeleteDomain("testdomain2.co.uk", gomock.Eq(req)).Return(nil),
		)

		ddosxDomainDelete(service, cmd, []string{"testdomain1.co.uk", "testdomain2.co.uk"})
	})

	t.Run("DeleteDomainError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeleteDomain("testdomain1.co.uk", ddosx.DeleteDomainRequest{}).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error removing domain [testdomain1.co.uk]: test error\n", func() {
			ddosxDomainDelete(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}

func Test_ddosxDomainDeployCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainDeployCmd(nil).Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ddosxDomainDeployCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing domain", err.Error())
	})
}

func Test_ddosxDomainDeploy(t *testing.T) {
	t.Run("SingleDomain", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeployDomain("testdomain1.co.uk").Return(nil),
			service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{}, nil),
		)

		ddosxDomainDeploy(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
	})

	t.Run("MultipleDomains", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeployDomain("testdomain1.co.uk").Return(nil),
			service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{}, nil),
			service.EXPECT().DeployDomain("testdomain2.co.uk").Return(nil),
			service.EXPECT().GetDomain("testdomain2.co.uk").Return(ddosx.Domain{}, nil),
		)

		ddosxDomainDeploy(service, &cobra.Command{}, []string{"testdomain1.co.uk", "testdomain2.co.uk"})
	})

	t.Run("WithWait", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainDeployCmd(nil)
		cmd.Flags().Set("wait", "true")

		gomock.InOrder(
			service.EXPECT().DeployDomain("testdomain1.co.uk").Return(nil),
			service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{Status: ddosx.DomainStatusConfigured}, nil),
			service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{}, nil),
		)

		ddosxDomainDeploy(service, cmd, []string{"testdomain1.co.uk"})
	})

	t.Run("WithWaitFailedStatus_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainDeployCmd(nil)
		cmd.Flags().Set("wait", "true")

		gomock.InOrder(
			service.EXPECT().DeployDomain("testdomain1.co.uk").Return(nil),
			service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{Status: ddosx.DomainStatusFailed}, nil),
		)

		test_output.AssertErrorOutput(t, "Error deploying domain [testdomain1.co.uk]: Error waiting for command: Domain [testdomain1.co.uk] in [Failed] state\n", func() {
			ddosxDomainDeploy(service, cmd, []string{"testdomain1.co.uk"})
		})
	})

	t.Run("DeployDomainError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeployDomain("testdomain1.co.uk").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error deploying domain [testdomain1.co.uk]: test error\n", func() {
			ddosxDomainDeploy(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})

	t.Run("GetDomainError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeployDomain("testdomain1.co.uk").Return(nil),
			service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving domain [testdomain1.co.uk]: test error\n", func() {
			ddosxDomainDeploy(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})
	})
}

func TestDomainStatusWaitFunc(t *testing.T) {
	t.Run("GetDomain_Error_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{}, errors.New("test error 1"))

		finished, err := DomainStatusWaitFunc(service, "testdomain1.co.uk", ddosx.DomainStatusConfigured)()

		assert.NotNil(t, err)
		assert.Equal(t, "Failed to retrieve domain [testdomain1.co.uk]: test error 1", err.Error())
		assert.False(t, finished)
	})

	t.Run("GetDomain_FailedStatus_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{Status: ddosx.DomainStatusFailed}, nil)

		finished, err := DomainStatusWaitFunc(service, "testdomain1.co.uk", ddosx.DomainStatusConfigured)()

		assert.NotNil(t, err)
		assert.Equal(t, "Domain [testdomain1.co.uk] in [Failed] state", err.Error())
		assert.False(t, finished)
	})

	t.Run("GetDomain_ExpectedStatus_ReturnsExpected", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{Status: ddosx.DomainStatusConfigured}, nil)

		finished, err := DomainStatusWaitFunc(service, "testdomain1.co.uk", ddosx.DomainStatusConfigured)()

		assert.Nil(t, err)
		assert.True(t, finished)
	})

	t.Run("GetDomain_UnexpectedStatus_ReturnsExpected", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{Status: ddosx.DomainStatusNotConfigured}, nil)

		finished, err := DomainStatusWaitFunc(service, "testdomain1.co.uk", ddosx.DomainStatusConfigured)()

		assert.Nil(t, err)
		assert.False(t, finished)
	})
}
