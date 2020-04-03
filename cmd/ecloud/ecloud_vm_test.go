package ecloud

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
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func Test_ecloudVirtualMachineList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVirtualMachines(gomock.Any()).Return([]ecloud.VirtualMachine{}, nil).Times(1)

		ecloudVirtualMachineList(service, &cobra.Command{}, []string{})
	})

	t.Run("ExpectedFilterFromFlags", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineListCmd(nil)
		cmd.Flags().Set("name", "test vm 1")

		expectedParameters := connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				connection.APIRequestFiltering{
					Property: "name",
					Operator: connection.EQOperator,
					Value:    []string{"test vm 1"},
				},
			},
		}

		service.EXPECT().GetVirtualMachines(gomock.Eq(expectedParameters)).Return([]ecloud.VirtualMachine{}, nil).Times(1)

		ecloudVirtualMachineList(service, cmd, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudVirtualMachineList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetVirtualMachinesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVirtualMachines(gomock.Any()).Return([]ecloud.VirtualMachine{}, errors.New("test error")).Times(1)

		err := ecloudVirtualMachineList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving virtual machines: test error", err.Error())
	})
}

func Test_ecloudVirtualMachineShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVirtualMachineShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVirtualMachineShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing virtual machine", err.Error())
	})
}

func Test_ecloudVirtualMachineShow(t *testing.T) {
	t.Run("SingleVirtualMachine", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, nil).Times(1)

		ecloudVirtualMachineShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleVirtualMachines", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, nil),
			service.EXPECT().GetVirtualMachine(456).Return(ecloud.VirtualMachine{}, nil),
		)

		ecloudVirtualMachineShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("InvalidVirtualMachineID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid virtual machine ID [abc]\n", func() {
			ecloudVirtualMachineShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetVirtualMachineError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving virtual machine [123]: test error\n", func() {
			ecloudVirtualMachineShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_ecloudVirtualMachineCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineCreateCmd(nil)
		cmd.Flags().Set("name", "test vm 1")

		gomock.InOrder(
			service.EXPECT().CreateVirtualMachine(gomock.Any()).Do(func(req ecloud.CreateVirtualMachineRequest) {
				if req.Name != "test vm 1" {
					t.Fatalf("expected VM name 'test vm 1', got '%s", req.Name)
				}
			}).Return(123, nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, nil),
		)

		ecloudVirtualMachineCreate(service, cmd, []string{})
	})

	t.Run("WithTag_SetsTag", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineCreateCmd(nil)
		cmd.Flags().Set("tag", "abc=123")

		gomock.InOrder(
			service.EXPECT().CreateVirtualMachine(gomock.Any()).Do(func(req ecloud.CreateVirtualMachineRequest) {
				assert.NotNil(t, req.Tags)
			}).Return(123, nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, nil),
		)

		ecloudVirtualMachineCreate(service, cmd, []string{})
	})

	t.Run("WithParameter_SetsParameter", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineCreateCmd(nil)
		cmd.Flags().Set("parameter", "abc=123")

		gomock.InOrder(
			service.EXPECT().CreateVirtualMachine(gomock.Any()).Do(func(req ecloud.CreateVirtualMachineRequest) {
				assert.NotNil(t, req.Parameters)
			}).Return(123, nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, nil),
		)

		ecloudVirtualMachineCreate(service, cmd, []string{})
	})

	t.Run("WithSSHKey_SetsSSHKey", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineCreateCmd(nil)
		cmd.Flags().Set("ssh-key", "testkey1")

		gomock.InOrder(
			service.EXPECT().CreateVirtualMachine(gomock.Any()).Do(func(req ecloud.CreateVirtualMachineRequest) {
				if req.SSHKeys == nil || len(req.SSHKeys) < 1 || req.SSHKeys[0] != "testkey1" {
					t.Fatal("Expected SSHKeys to contain key [testkey1]")
				}
			}).Return(123, nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, nil),
		)

		ecloudVirtualMachineCreate(service, cmd, []string{})
	})

	t.Run("WithEncrypt_SetsEncrypt", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineCreateCmd(nil)
		cmd.Flags().Set("encrypt", "true")

		gomock.InOrder(
			service.EXPECT().CreateVirtualMachine(gomock.Any()).Do(func(req ecloud.CreateVirtualMachineRequest) {
				assert.NotNil(t, req.Encrypt)
				assert.True(t, *req.Encrypt)
			}).Return(123, nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, nil),
		)

		ecloudVirtualMachineCreate(service, cmd, []string{})
	})

	t.Run("CreateWithWait_WaitsForComplete", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineCreateCmd(nil)
		cmd.Flags().Set("wait", "true")

		gomock.InOrder(
			service.EXPECT().CreateVirtualMachine(gomock.Any()).Return(123, nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{Status: ecloud.VirtualMachineStatusComplete}, nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, nil),
		)

		ecloudVirtualMachineCreate(service, cmd, []string{"123"})
	})

	t.Run("CreateWithWaitFailedStatus_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineCreateCmd(nil)
		cmd.Flags().Set("wait", "true")

		gomock.InOrder(
			service.EXPECT().CreateVirtualMachine(gomock.Any()).Return(123, nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{Status: ecloud.VirtualMachineStatusFailed}, nil),
		)

		err := ecloudVirtualMachineCreate(service, cmd, []string{"123"})

		assert.Equal(t, "Error waiting for command: Virtual machine [123] in [Failed] state", err.Error())
	})

	t.Run("CreateWithWaitRetrieveStatusFailure_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineCreateCmd(nil)
		cmd.Flags().Set("wait", "true")

		gomock.InOrder(
			service.EXPECT().CreateVirtualMachine(gomock.Any()).Return(123, nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, errors.New("test error 1")),
		)

		err := ecloudVirtualMachineCreate(service, cmd, []string{"123"})

		assert.Equal(t, "Error waiting for command: Failed to retrieve virtual machine [123]: test error 1", err.Error())
	})

	t.Run("CreateVirtualMachineError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().CreateVirtualMachine(gomock.Any()).Return(0, errors.New("test error")).Times(1)

		err := ecloudVirtualMachineCreate(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error creating virtual machine: test error", err.Error())
	})

	t.Run("GetVirtualMachineError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().CreateVirtualMachine(gomock.Any()).Return(123, nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, errors.New("test error")),
		)

		err := ecloudVirtualMachineCreate(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving new virtual machine: test error", err.Error())
	})
}

func Test_ecloudVirtualMachineUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVirtualMachineUpdateCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVirtualMachineUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing virtual machine", err.Error())
	})
}

func Test_ecloudVirtualMachineUpdate(t *testing.T) {
	t.Run("SingleVirtualMachine_SetsCPU", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineUpdateCmd(nil)
		cmd.Flags().Set("cpu", "2")
		cmd.Flags().Set("ram", "4")
		cmd.Flags().Set("name", "test vm name 1")

		gomock.InOrder(
			service.EXPECT().PatchVirtualMachine(123, gomock.Any()).Do(func(vmID int, patch ecloud.PatchVirtualMachineRequest) {
				if patch.CPU != 2 {
					t.Fatal("Unexpected CPU count")
				}
				if patch.RAM != 4 {
					t.Fatal("Unexpected RAM count")
				}
				if patch.Name == nil || *patch.Name != "test vm name 1" {
					t.Fatal("Unexpected name")
				}
			}).Return(nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{Status: ecloud.VirtualMachineStatusComplete}, nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, nil),
		)

		ecloudVirtualMachineUpdate(service, cmd, []string{"123"})
	})

	t.Run("MultipleVirtualMachines", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineUpdateCmd(nil)
		cmd.Flags().Set("cpu", "2")

		expectedPatch := ecloud.PatchVirtualMachineRequest{
			CPU: 2,
		}

		gomock.InOrder(
			service.EXPECT().PatchVirtualMachine(123, gomock.Eq(expectedPatch)).Return(nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{Status: ecloud.VirtualMachineStatusComplete}, nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, nil),
			service.EXPECT().PatchVirtualMachine(456, gomock.Eq(expectedPatch)).Return(nil),
			service.EXPECT().GetVirtualMachine(456).Return(ecloud.VirtualMachine{Status: ecloud.VirtualMachineStatusComplete}, nil),
			service.EXPECT().GetVirtualMachine(456).Return(ecloud.VirtualMachine{}, nil),
		)

		ecloudVirtualMachineUpdate(service, cmd, []string{"123", "456"})
	})

	t.Run("InvalidVirtualMachineID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid virtual machine ID [abc]\n", func() {
			ecloudVirtualMachineUpdate(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("PatchVirtualMachineError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchVirtualMachine(123, gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating virtual machine [123]: test error\n", func() {
			ecloudVirtualMachineUpdate(service, &cobra.Command{}, []string{"123"})
		})
	})

	t.Run("WaitGetVirtualMachineError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchVirtualMachine(123, gomock.Any()).Return(nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error updating virtual machine [123]: Error waiting for command: Failed to retrieve virtual machine [123]: test error\n", func() {
			ecloudVirtualMachineUpdate(service, &cobra.Command{}, []string{"123"})
		})
	})

	t.Run("GetVirtualMachineError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchVirtualMachine(123, gomock.Any()).Return(nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{Status: ecloud.VirtualMachineStatusComplete}, nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated virtual machine [123]: test error\n", func() {
			ecloudVirtualMachineUpdate(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_ecloudVirtualMachineStartCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVirtualMachineStartCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVirtualMachineStartCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing virtual machine", err.Error())
	})
}

func Test_ecloudVirtualMachineStart(t *testing.T) {
	t.Run("SingleVirtualMachine", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PowerOnVirtualMachine(123).Return(nil).Times(1)

		ecloudVirtualMachineStart(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleVirtualMachines", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PowerOnVirtualMachine(123).Return(nil),
			service.EXPECT().PowerOnVirtualMachine(456).Return(nil),
		)

		ecloudVirtualMachineStart(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("InvalidVirtualMachineID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid virtual machine ID [abc]\n", func() {
			ecloudVirtualMachineStart(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("PowerOnVirtualMachineError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PowerOnVirtualMachine(123).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error powering on virtual machine [123]: test error\n", func() {
			ecloudVirtualMachineStart(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_ecloudVirtualMachineStopCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVirtualMachineStopCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVirtualMachineStopCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing virtual machine", err.Error())
	})
}

func Test_ecloudVirtualMachineStop(t *testing.T) {
	t.Run("SingleVirtualMachine", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PowerShutdownVirtualMachine(123).Return(nil).Times(1)

		ecloudVirtualMachineStop(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleVirtualMachines", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PowerShutdownVirtualMachine(123).Return(nil),
			service.EXPECT().PowerShutdownVirtualMachine(456).Return(nil),
		)

		ecloudVirtualMachineStop(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("SingleVirtualMachine_Force", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineStopCmd(nil)
		cmd.Flags().Set("force", "true")

		service.EXPECT().PowerOffVirtualMachine(123).Return(nil).Times(1)

		ecloudVirtualMachineStop(service, cmd, []string{"123"})
	})

	t.Run("MultipleVirtualMachines_Force", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineStopCmd(nil)
		cmd.Flags().Set("force", "true")

		gomock.InOrder(
			service.EXPECT().PowerOffVirtualMachine(123).Return(nil),
			service.EXPECT().PowerOffVirtualMachine(456).Return(nil),
		)

		ecloudVirtualMachineStop(service, cmd, []string{"123", "456"})
	})

	t.Run("InvalidVirtualMachineID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid virtual machine ID [abc]\n", func() {
			ecloudVirtualMachineStop(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("PowerShutdownVirtualMachineError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PowerShutdownVirtualMachine(123).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error powering off virtual machine [123]: test error\n", func() {
			ecloudVirtualMachineStop(service, &cobra.Command{}, []string{"123"})
		})
	})

	t.Run("PowerOffVirtualMachineError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineStopCmd(nil)
		cmd.Flags().Set("force", "true")

		service.EXPECT().PowerOffVirtualMachine(123).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error powering off (forced) virtual machine [123]: test error\n", func() {
			ecloudVirtualMachineStop(service, cmd, []string{"123"})
		})
	})
}

func Test_ecloudVirtualMachineRestartCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVirtualMachineRestartCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVirtualMachineRestartCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing virtual machine", err.Error())
	})
}

func Test_ecloudVirtualMachineRestart(t *testing.T) {
	t.Run("SingleVirtualMachine", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PowerRestartVirtualMachine(123).Return(nil).Times(1)

		ecloudVirtualMachineRestart(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleVirtualMachines", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PowerRestartVirtualMachine(123).Return(nil),
			service.EXPECT().PowerRestartVirtualMachine(456).Return(nil),
		)

		ecloudVirtualMachineRestart(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("SingleVirtualMachine_Force", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineRestartCmd(nil)
		cmd.Flags().Set("force", "true")

		service.EXPECT().PowerResetVirtualMachine(123).Return(nil).Times(1)

		ecloudVirtualMachineRestart(service, cmd, []string{"123"})
	})

	t.Run("MultipleVirtualMachines_Force", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineRestartCmd(nil)
		cmd.Flags().Set("force", "true")

		gomock.InOrder(
			service.EXPECT().PowerResetVirtualMachine(123).Return(nil),
			service.EXPECT().PowerResetVirtualMachine(456).Return(nil),
		)

		ecloudVirtualMachineRestart(service, cmd, []string{"123", "456"})
	})

	t.Run("InvalidVirtualMachineID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid virtual machine ID [abc]\n", func() {
			ecloudVirtualMachineRestart(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("PowerRestartVirtualMachineError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PowerRestartVirtualMachine(123).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error restarting virtual machine [123]: test error\n", func() {
			ecloudVirtualMachineRestart(service, &cobra.Command{}, []string{"123"})
		})
	})

	t.Run("PowerResetVirtualMachineError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineRestartCmd(nil)
		cmd.Flags().Set("force", "true")

		service.EXPECT().PowerResetVirtualMachine(123).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error restarting (forced) virtual machine [123]: test error\n", func() {
			ecloudVirtualMachineRestart(service, cmd, []string{"123"})
		})
	})
}

func Test_ecloudVirtualMachineDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVirtualMachineDeleteCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVirtualMachineDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing virtual machine", err.Error())
	})
}

func Test_ecloudVirtualMachineDelete(t *testing.T) {
	t.Run("SingleVirtualMachine", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteVirtualMachine(123).Return(nil).Times(1)

		ecloudVirtualMachineDelete(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleVirtualMachines", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteVirtualMachine(123).Return(nil),
			service.EXPECT().DeleteVirtualMachine(456).Return(nil),
		)

		ecloudVirtualMachineDelete(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("WithWait", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineDeleteCmd(nil)
		cmd.Flags().Set("wait", "true")

		gomock.InOrder(
			service.EXPECT().DeleteVirtualMachine(123).Return(nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, &ecloud.VirtualMachineNotFoundError{}),
		)

		ecloudVirtualMachineDelete(service, cmd, []string{"123"})
	})

	t.Run("WithWaitFailedStatus_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineDeleteCmd(nil)
		cmd.Flags().Set("wait", "true")

		gomock.InOrder(
			service.EXPECT().DeleteVirtualMachine(123).Return(nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{Status: ecloud.VirtualMachineStatusFailed}, nil),
		)

		test_output.AssertErrorOutput(t, "Error removing virtual machine [123]: Error waiting for command: Virtual machine [123] in [Failed] state\n", func() {
			ecloudVirtualMachineDelete(service, cmd, []string{"123"})
		})
	})

	t.Run("WithWaitGetStatusError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVirtualMachineDeleteCmd(nil)
		cmd.Flags().Set("wait", "true")

		gomock.InOrder(
			service.EXPECT().DeleteVirtualMachine(123).Return(nil),
			service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, errors.New("test error 1")),
		)

		test_output.AssertErrorOutput(t, "Error removing virtual machine [123]: Error waiting for command: Failed to retrieve virtual machine [123]: test error 1\n", func() {
			ecloudVirtualMachineDelete(service, cmd, []string{"123"})
		})
	})

	t.Run("InvalidVirtualMachineID_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid virtual machine ID [abc]\n", func() {
			ecloudVirtualMachineDelete(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetVirtualMachineError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteVirtualMachine(123).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error removing virtual machine [123]: test error\n", func() {
			ecloudVirtualMachineDelete(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func TestVirtualMachineNotFoundWaitFunc(t *testing.T) {
	t.Run("GetVirtualMachine_VirtualMachineNotFoundError_ReturnsExpected", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, &ecloud.VirtualMachineNotFoundError{})

		finished, err := VirtualMachineNotFoundWaitFunc(service, 123)()

		assert.Nil(t, err)
		assert.True(t, finished)
	})

	t.Run("GetVirtualMachine_Error_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, errors.New("test error 1"))

		finished, err := VirtualMachineNotFoundWaitFunc(service, 123)()

		assert.NotNil(t, err)
		assert.Equal(t, "Failed to retrieve virtual machine [123]: test error 1", err.Error())
		assert.False(t, finished)
	})

	t.Run("GetVirtualMachine_FailedStatus_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{Status: ecloud.VirtualMachineStatusFailed}, nil)

		finished, err := VirtualMachineNotFoundWaitFunc(service, 123)()

		assert.NotNil(t, err)
		assert.Equal(t, "Virtual machine [123] in [Failed] state", err.Error())
		assert.False(t, finished)
	})

	t.Run("GetVirtualMachine_NonFailedStatus_ReturnsExpected", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{Status: ecloud.VirtualMachineStatusBeingBuilt}, nil)

		finished, err := VirtualMachineNotFoundWaitFunc(service, 123)()

		assert.Nil(t, err)
		assert.False(t, finished)
	})
}

func TestVirtualMachineStatusWaitFunc(t *testing.T) {
	t.Run("GetVirtualMachine_Error_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{}, errors.New("test error 1"))

		finished, err := VirtualMachineStatusWaitFunc(service, 123, ecloud.VirtualMachineStatusComplete)()

		assert.NotNil(t, err)
		assert.Equal(t, "Failed to retrieve virtual machine [123]: test error 1", err.Error())
		assert.False(t, finished)
	})

	t.Run("GetVirtualMachine_FailedStatus_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{Status: ecloud.VirtualMachineStatusFailed}, nil)

		finished, err := VirtualMachineStatusWaitFunc(service, 123, ecloud.VirtualMachineStatusComplete)()

		assert.NotNil(t, err)
		assert.Equal(t, "Virtual machine [123] in [Failed] state", err.Error())
		assert.False(t, finished)
	})

	t.Run("GetVirtualMachine_ExpectedStatus_ReturnsExpected", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{Status: ecloud.VirtualMachineStatusComplete}, nil)

		finished, err := VirtualMachineStatusWaitFunc(service, 123, ecloud.VirtualMachineStatusComplete)()

		assert.Nil(t, err)
		assert.True(t, finished)
	})

	t.Run("GetVirtualMachine_UnexpectedStatus_ReturnsExpected", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		viper.SetDefault("command_wait_timeout_seconds", 1200)
		viper.SetDefault("command_wait_sleep_seconds", 1)
		defer test.TestResetViper()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVirtualMachine(123).Return(ecloud.VirtualMachine{Status: ecloud.VirtualMachineStatusBeingBuilt}, nil)

		finished, err := VirtualMachineStatusWaitFunc(service, 123, ecloud.VirtualMachineStatusComplete)()

		assert.Nil(t, err)
		assert.False(t, finished)
	})
}
