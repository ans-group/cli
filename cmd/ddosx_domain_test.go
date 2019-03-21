package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/cli/test"
	"github.com/ukfast/cli/test/mocks"
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
			ddosxDomainList(service, &cobra.Command{}, []string{})
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

		service.EXPECT().GetDomains(gomock.Any()).Return([]ddosx.Domain{}, errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			ddosxDomainList(service, &cobra.Command{}, []string{})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error retrieving domains: test error\n", output)
	})
}

func Test_ddosxDomainShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainShowCmd().Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ddosxDomainShowCmd().Args(nil, []string{})

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

		output := test.CatchStdErr(t, func() {
			ddosxDomainShow(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})

		assert.Equal(t, "Error retrieving domain [testdomain1.co.uk]: test error\n", output)
	})
}

func Test_ddosxDomainCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainCreateCmd()
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

	t.Run("CreateDomain_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().CreateDomain(gomock.Any()).Return(errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			ddosxDomainCreate(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error creating domain: test error\n", output)
	})

	t.Run("CreateDomain_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainCreateCmd()
		cmd.Flags().Set("name", "testdomain1.co.uk")

		gomock.InOrder(
			service.EXPECT().CreateDomain(gomock.Any()).Return(nil),
			service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{}, errors.New("test error")),
		)

		output := test.CatchStdErr(t, func() {
			ddosxDomainCreate(service, cmd, []string{"testdomain1.co.uk"})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error retrieving new domain: test error\n", output)
	})
}

func Test_ddosxDomainDeployCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ddosxDomainDeployCmd().Args(nil, []string{"testdomain1.co.uk"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ddosxDomainDeployCmd().Args(nil, []string{})

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
		defer testResetViper()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainDeployCmd()
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
		defer testResetViper()

		service := mocks.NewMockDDoSXService(mockCtrl)
		cmd := ddosxDomainDeployCmd()
		cmd.Flags().Set("wait", "true")

		gomock.InOrder(
			service.EXPECT().DeployDomain("testdomain1.co.uk").Return(nil),
			service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{Status: ddosx.DomainStatusFailed}, nil),
		)

		output := test.CatchStdErr(t, func() {
			ddosxDomainDeploy(service, cmd, []string{"testdomain1.co.uk"})
		})

		assert.Equal(t, "Error deploying domain [testdomain1.co.uk]: Error waiting for command: Domain [testdomain1.co.uk] in [Failed] state\n", output)
	})

	t.Run("DeployDomainError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().DeployDomain("testdomain1.co.uk").Return(errors.New("test error"))

		output := test.CatchStdErr(t, func() {
			ddosxDomainDeploy(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})

		assert.Equal(t, "Error deploying domain [testdomain1.co.uk]: test error\n", output)
	})

	t.Run("GetDomainError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDDoSXService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeployDomain("testdomain1.co.uk").Return(nil),
			service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{}, errors.New("test error")),
		)

		output := test.CatchStdErr(t, func() {
			ddosxDomainDeploy(service, &cobra.Command{}, []string{"testdomain1.co.uk"})
		})

		assert.Equal(t, "Error retrieving domain [testdomain1.co.uk]: test error\n", output)
	})
}

func TestDomainStatusWaitFunc(t *testing.T) {
	t.Run("GetDomain_Error_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer testResetViper()

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
		defer testResetViper()

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
		defer testResetViper()

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
		defer testResetViper()

		service := mocks.NewMockDDoSXService(mockCtrl)

		service.EXPECT().GetDomain("testdomain1.co.uk").Return(ddosx.Domain{Status: ddosx.DomainStatusNotConfigured}, nil)

		finished, err := DomainStatusWaitFunc(service, "testdomain1.co.uk", ddosx.DomainStatusConfigured)()

		assert.Nil(t, err)
		assert.False(t, finished)
	})
}
