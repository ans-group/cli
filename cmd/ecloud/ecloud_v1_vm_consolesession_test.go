package ecloud

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ecloudVirtualMachineConsoleSessionCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVirtualMachineConsoleSessionCreateCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("MissingVirtualMachine_Error", func(t *testing.T) {
		err := ecloudVirtualMachineConsoleSessionCreateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing virtual machine", err.Error())
	})
}

func Test_ecloudVirtualMachineConsoleSessionCreate(t *testing.T) {
	t.Run("CreateSuccess_NoError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().CreateVirtualMachineConsoleSession(123).Return(ecloud.ConsoleSession{}, nil).Times(1)

		err := ecloudVirtualMachineConsoleSessionCreate(service, &cobra.Command{}, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidVirtualMachineID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudVirtualMachineConsoleSessionCreate(service, &cobra.Command{}, []string{"abc"})

		assert.Equal(t, "Invalid virtual machine ID [abc]", err.Error())
	})

	t.Run("GetVirtualMachineConsoleSessionError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().CreateVirtualMachineConsoleSession(123).Return(ecloud.ConsoleSession{}, errors.New("test error 1")).Times(1)

		err := ecloudVirtualMachineConsoleSessionCreate(service, &cobra.Command{}, []string{"123"})

		assert.Equal(t, "Error creating virtual machine console session: test error 1", err.Error())
	})
}
