package ecloud

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/sdk-go/pkg/config"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ecloudVirtualMachineDiskListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVirtualMachineDiskListCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVirtualMachineDiskListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing virtual machine", err.Error())
	})
}

func Test_ecloudVirtualMachineDiskList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, nil).Times(1)

		ecloudVirtualMachineDiskList(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("InvalidVirtualMachineID_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudVirtualMachineDiskList(service, &cobra.Command{}, []string{"abc"})

		assert.Equal(t, "Invalid virtual machine ID [abc]", err.Error())
	})

	t.Run("GetVirtualMachineError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, errors.New("test error 1")).Times(1)

		err := ecloudVirtualMachineDiskList(service, &cobra.Command{}, []string{"123"})

		assert.Equal(t, "Error retrieving virtual machine [123]: test error 1", err.Error())
	})
}

func Test_ecloudVirtualMachineDiskUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVirtualMachineDiskUpdateCmd(nil).Args(nil, []string{"123", "00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingVirtualMachine_Error", func(t *testing.T) {
		err := ecloudVirtualMachineDiskUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing virtual machine", err.Error())
	})

	t.Run("MissingDisk_Error", func(t *testing.T) {
		err := ecloudVirtualMachineDiskUpdateCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing disk", err.Error())
	})
}

func Test_ecloudVirtualMachineDiskUpdate(t *testing.T) {
	t.Run("SingleVirtualMachine_SetsCapacity", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		config.Set("test", "command_wait_timeout_seconds", 1200)
		config.Set("test", "command_wait_sleep_seconds", 1)
		config.SwitchCurrentContext("test")
		defer config.Reset()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineDiskUpdateCmd(nil)
		cmd.Flags().Set("capacity", "25")

		gomock.InOrder(
			service.EXPECT().PatchVirtualMachine(123, gomock.Any()).Do(func(vmID int, patch ecloud.PatchVirtualMachineRequest) {
				if patch.Disks == nil || len(patch.Disks) < 1 || patch.Disks[0].Capacity != 25 {
					t.Fatal("Unexpected disk patch request")
				}
			}).Return(nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{Status: ecloud.VirtualMachineStatusComplete}, nil),
		)

		ecloudVirtualMachineDiskUpdate(service, cmd, []string{"123", "00000000-0000-0000-0000-000000000000"})
	})

	t.Run("InvalidVirtualMachineID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudVirtualMachineDiskUpdate(service, &cobra.Command{}, []string{"abc", "00000000-0000-0000-0000-000000000000"})

		assert.Equal(t, "Invalid virtual machine ID [abc]", err.Error())
	})

	t.Run("PatchVirtualMachineError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchVirtualMachine(123, gomock.Any()).Return(errors.New("test error"))

		err := ecloudVirtualMachineDiskUpdate(service, &cobra.Command{}, []string{"123", "00000000-0000-0000-0000-000000000000"})

		assert.Equal(t, "Error updating virtual machine [123]: test error", err.Error())
	})

	t.Run("WaitGetVirtualMachineError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		config.Set("test", "command_wait_timeout_seconds", 1200)
		config.Set("test", "command_wait_sleep_seconds", 1)
		config.SwitchCurrentContext("test")
		defer config.Reset()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchVirtualMachine(123, gomock.Any()).Return(nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, errors.New("test error")),
		)

		err := ecloudVirtualMachineDiskUpdate(service, &cobra.Command{}, []string{"123", "00000000-0000-0000-0000-000000000000"})

		assert.Equal(t, "Error updating virtual machine [123]: Error waiting for command: Failed to retrieve virtual machine [123]: test error", err.Error())
	})
}
