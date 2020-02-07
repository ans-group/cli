package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func Test_ecloudVirtualMachineTemplateCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVirtualMachineTemplateCreateCmd().Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVirtualMachineTemplateCreateCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing virtual machine", err.Error())
	})
}

func Test_ecloudVirtualMachineTemplateCreate(t *testing.T) {
	t.Run("InvalidVirtualMachineID_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudVirtualMachineTemplateCreate(service, &cobra.Command{}, []string{"abc"})
		assert.Equal(t, "Invalid virtual machine ID [abc]", err.Error())
	})

	t.Run("InvalidType_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineTemplateCreateCmd()
		cmd.Flags().Set("type", "invalid")

		err := ecloudVirtualMachineTemplateCreate(service, cmd, []string{"123"})
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("CreateVirtualMachineTemplateError_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineTemplateCreateCmd()
		cmd.Flags().Set("type", "solution")

		service.EXPECT().CreateVirtualMachineTemplate(123, gomock.Any()).Return(errors.New("test error 1"))

		err := ecloudVirtualMachineTemplateCreate(service, cmd, []string{"123"})
		assert.Equal(t, "Error creating virtual machine template: test error 1", err.Error())
	})

	t.Run("WaitForCommandError_OutputsFatal", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer testResetViper()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineTemplateCreateCmd()
		cmd.Flags().Set("type", "solution")
		cmd.Flags().Set("wait", "true")

		gomock.InOrder(
			service.EXPECT().CreateVirtualMachineTemplate(123, gomock.Any()).Return(nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, errors.New("test error 1")),
		)

		err := ecloudVirtualMachineTemplateCreate(service, cmd, []string{"123"})
		assert.Equal(t, "Error waiting for command: Failed to retrieve virtual machine [123]: test error 1", err.Error())
	})
}
