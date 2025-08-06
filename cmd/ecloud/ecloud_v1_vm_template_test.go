package ecloud

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/sdk-go/pkg/config"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ecloudVirtualMachineTemplateCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVirtualMachineTemplateCreateCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVirtualMachineTemplateCreateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing virtual machine", err.Error())
	})
}

func Test_ecloudVirtualMachineTemplateCreate(t *testing.T) {
	t.Run("InvalidVirtualMachineID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudVirtualMachineTemplateCreate(service, &cobra.Command{}, []string{"abc"})
		assert.Equal(t, "invalid virtual machine ID [abc]", err.Error())
	})

	t.Run("InvalidType_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineTemplateCreateCmd(nil)
		cmd.Flags().Set("type", "invalid")

		err := ecloudVirtualMachineTemplateCreate(service, cmd, []string{"123"})
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("CreateVirtualMachineTemplateError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineTemplateCreateCmd(nil)
		cmd.Flags().Set("type", "solution")

		service.EXPECT().CreateVirtualMachineTemplate(123, gomock.Any()).Return(errors.New("test error 1"))

		err := ecloudVirtualMachineTemplateCreate(service, cmd, []string{"123"})
		assert.Equal(t, "error creating virtual machine template: test error 1", err.Error())
	})

	t.Run("WaitForCommandError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		config.Set("test", "command_wait_timeout_seconds", 1200)
		config.Set("test", "command_wait_sleep_seconds", 1)
		defer config.Reset()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineTemplateCreateCmd(nil)
		cmd.Flags().Set("type", "solution")
		cmd.Flags().Set("wait", "true")

		gomock.InOrder(
			service.EXPECT().CreateVirtualMachineTemplate(123, gomock.Any()).Return(nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, errors.New("test error 1")),
		)

		err := ecloudVirtualMachineTemplateCreate(service, cmd, []string{"123"})
		assert.Equal(t, "error waiting for command: failed to retrieve virtual machine [123]: test error 1", err.Error())
	})
}
