package ecloud

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func Test_ecloudInstanceList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetInstances(gomock.Any()).Return([]ecloud.Instance{}, nil).Times(1)

		ecloudInstanceList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudInstanceList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetInstancesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetInstances(gomock.Any()).Return([]ecloud.Instance{}, errors.New("test error")).Times(1)

		err := ecloudInstanceList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving instances: test error", err.Error())
	})
}

func Test_ecloudInstanceShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudInstanceShowCmd(nil).Args(nil, []string{"i-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudInstanceShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing instance", err.Error())
	})
}

func Test_ecloudInstanceShow(t *testing.T) {
	t.Run("SingleInstance", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{}, nil).Times(1)

		ecloudInstanceShow(service, &cobra.Command{}, []string{"i-abcdef12"})
	})

	t.Run("MultipleInstances", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{}, nil),
			service.EXPECT().GetInstance("i-abcdef23").Return(ecloud.Instance{}, nil),
		)

		ecloudInstanceShow(service, &cobra.Command{}, []string{"i-abcdef12", "i-abcdef23"})
	})

	t.Run("GetInstanceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving instance [i-abcdef12]: test error\n", func() {
			ecloudInstanceShow(service, &cobra.Command{}, []string{"i-abcdef12"})
		})
	})
}

func Test_ecloudInstanceCreate(t *testing.T) {
	t.Run("CreateWithImageID_NoError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudInstanceCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testinstance", "--image=img-abcdef12"})

		req := ecloud.CreateInstanceRequest{
			Name:    "testinstance",
			ImageID: "img-abcdef12",
		}

		service.EXPECT().CreateInstance(req).Return("i-abcdef12", nil)
		service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{}, nil)

		err := ecloudInstanceCreate(service, cmd, []string{})
		assert.Nil(t, err)
	})

	t.Run("CreateWithImageName_NoError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudInstanceCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testinstance", "--image=test"})

		req := ecloud.CreateInstanceRequest{
			Name:    "testinstance",
			ImageID: "img-abcdef12",
		}

		service.EXPECT().GetImages(connection.APIRequestParameters{}).Return([]ecloud.Image{{Name: "test", ID: "img-abcdef12"}}, nil)
		service.EXPECT().CreateInstance(req).Return("i-abcdef12", nil)
		service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{}, nil)

		err := ecloudInstanceCreate(service, cmd, []string{})
		assert.Nil(t, err)
	})

	t.Run("ImageNotFound_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudInstanceCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testinstance", "--image=unknown"})

		service.EXPECT().GetImages(connection.APIRequestParameters{}).Return([]ecloud.Image{{Name: "test", ID: "img-abcdef12"}}, nil)

		err := ecloudInstanceCreate(service, cmd, []string{})
		assert.NotNil(t, err)
		assert.Equal(t, "Image not found with name 'unknown'", err.Error())
	})

	t.Run("CreateInstanceError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudInstanceCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testinstance", "--image=img-abcdef12"})

		service.EXPECT().CreateInstance(gomock.Any()).Return("", errors.New("test error"))

		err := ecloudInstanceCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating instance: test error", err.Error())
	})

	t.Run("GetInstanceError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudInstanceCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testinstance", "--image=img-abcdef12"})

		gomock.InOrder(
			service.EXPECT().CreateInstance(gomock.Any()).Return("i-abcdef12", nil),
			service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{}, errors.New("test error")),
		)

		err := ecloudInstanceCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new instance: test error", err.Error())
	})
}

func Test_ecloudInstanceUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudInstanceUpdateCmd(nil).Args(nil, []string{"i-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudInstanceUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing instance", err.Error())
	})
}

func Test_ecloudInstanceUpdate(t *testing.T) {
	t.Run("SingleInstance", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudInstanceUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testinstance", "--vcpu=2", "--ram=2"})

		req := ecloud.PatchInstanceRequest{
			Name:        "testinstance",
			VCPUCores:   2,
			RAMCapacity: 2,
		}

		gomock.InOrder(
			service.EXPECT().PatchInstance("i-abcdef12", req).Return(nil),
			service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{}, nil),
		)

		ecloudInstanceUpdate(service, cmd, []string{"i-abcdef12"})
	})

	t.Run("MultipleInstances", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchInstance("i-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{}, nil),
			service.EXPECT().PatchInstance("i-12abcdef", gomock.Any()).Return(nil),
			service.EXPECT().GetInstance("i-12abcdef").Return(ecloud.Instance{}, nil),
		)

		ecloudInstanceUpdate(service, &cobra.Command{}, []string{"i-abcdef12", "i-12abcdef"})
	})

	t.Run("PatchInstanceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchInstance("i-abcdef12", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating instance [i-abcdef12]: test error\n", func() {
			ecloudInstanceUpdate(service, &cobra.Command{}, []string{"i-abcdef12"})
		})
	})

	t.Run("GetInstanceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchInstance("i-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated instance [i-abcdef12]: test error\n", func() {
			ecloudInstanceUpdate(service, &cobra.Command{}, []string{"i-abcdef12"})
		})
	})
}

func Test_ecloudInstanceDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudInstanceDeleteCmd(nil).Args(nil, []string{"i-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudInstanceDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing instance", err.Error())
	})
}

func Test_ecloudInstanceDelete(t *testing.T) {
	t.Run("SingleInstance", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteInstance("i-abcdef12").Return(nil).Times(1)

		ecloudInstanceDelete(service, &cobra.Command{}, []string{"i-abcdef12"})
	})

	t.Run("MultipleInstances", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteInstance("i-abcdef12").Return(nil),
			service.EXPECT().DeleteInstance("i-12abcdef").Return(nil),
		)

		ecloudInstanceDelete(service, &cobra.Command{}, []string{"i-abcdef12", "i-12abcdef"})
	})

	t.Run("DeleteInstanceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteInstance("i-abcdef12").Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing instance [i-abcdef12]: test error\n", func() {
			ecloudInstanceDelete(service, &cobra.Command{}, []string{"i-abcdef12"})
		})
	})
}

func Test_ecloudInstanceLockCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudInstanceLockCmd(nil).Args(nil, []string{"i-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudInstanceLockCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing instance", err.Error())
	})
}

func Test_ecloudInstanceLock(t *testing.T) {
	t.Run("SingleInstance", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().LockInstance("i-abcdef12").Return(nil).Times(1)

		ecloudInstanceLock(service, &cobra.Command{}, []string{"i-abcdef12"})
	})

	t.Run("MultipleInstances", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().LockInstance("i-abcdef12").Return(nil),
			service.EXPECT().LockInstance("i-12abcdef").Return(nil),
		)

		ecloudInstanceLock(service, &cobra.Command{}, []string{"i-abcdef12", "i-12abcdef"})
	})

	t.Run("LockInstanceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().LockInstance("i-abcdef12").Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error locking instance [i-abcdef12]: test error\n", func() {
			ecloudInstanceLock(service, &cobra.Command{}, []string{"i-abcdef12"})
		})
	})
}

func Test_ecloudInstanceUnlockCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudInstanceUnlockCmd(nil).Args(nil, []string{"i-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudInstanceUnlockCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing instance", err.Error())
	})
}

func Test_ecloudInstanceUnlock(t *testing.T) {
	t.Run("SingleInstance", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().UnlockInstance("i-abcdef12").Return(nil)

		ecloudInstanceUnlock(service, &cobra.Command{}, []string{"i-abcdef12"})
	})

	t.Run("MultipleInstances", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		gomock.InOrder(
			service.EXPECT().UnlockInstance("i-abcdef12").Return(nil),
			service.EXPECT().UnlockInstance("i-12abcdef").Return(nil),
		)

		ecloudInstanceUnlock(service, &cobra.Command{}, []string{"i-abcdef12", "i-12abcdef"})
	})

	t.Run("UnlockInstanceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().UnlockInstance("i-abcdef12").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error unlocking instance [i-abcdef12]: test error\n", func() {
			ecloudInstanceUnlock(service, &cobra.Command{}, []string{"i-abcdef12"})
		})
	})
}

func Test_ecloudInstanceStartCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudInstanceStartCmd(nil).Args(nil, []string{"i-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudInstanceStartCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing instance", err.Error())
	})
}

func Test_ecloudInstanceStart(t *testing.T) {
	t.Run("SingleInstance", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerOnInstance("i-abcdef12").Return(nil)

		ecloudInstanceStart(service, &cobra.Command{}, []string{"i-abcdef12"})
	})

	t.Run("MultipleInstances", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		gomock.InOrder(
			service.EXPECT().PowerOnInstance("i-abcdef12").Return(nil),
			service.EXPECT().PowerOnInstance("i-12abcdef").Return(nil),
		)

		ecloudInstanceStart(service, &cobra.Command{}, []string{"i-abcdef12", "i-12abcdef"})
	})

	t.Run("PowerOnInstanceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerOnInstance("i-abcdef12").Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error starting instance [i-abcdef12]: test error\n", func() {
			ecloudInstanceStart(service, &cobra.Command{}, []string{"i-abcdef12"})
		})
	})
}

func Test_ecloudInstanceStopCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudInstanceStopCmd(nil).Args(nil, []string{"i-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudInstanceStopCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing instance", err.Error())
	})
}

func Test_ecloudInstanceStop(t *testing.T) {
	t.Run("WithoutForceFlag_CallsPowerShutdownInstance", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerShutdownInstance("i-abcdef12").Return(nil)

		ecloudInstanceStop(service, &cobra.Command{}, []string{"i-abcdef12"})
	})

	t.Run("WithForceFlag_CallsPowerOffInstance", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerOffInstance("i-abcdef12").Return(nil)

		cmd := ecloudInstanceStopCmd(nil)
		cmd.ParseFlags([]string{"--force"})

		ecloudInstanceStop(service, cmd, []string{"i-abcdef12"})
	})

	t.Run("PowerShutdownInstanceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerShutdownInstance("i-abcdef12").Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error stopping instance [i-abcdef12]: test error\n", func() {
			ecloudInstanceStop(service, &cobra.Command{}, []string{"i-abcdef12"})
		})
	})

	t.Run("PowerOffInstanceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerOffInstance("i-abcdef12").Return(errors.New("test error"))

		cmd := ecloudInstanceStopCmd(nil)
		cmd.ParseFlags([]string{"--force"})

		test_output.AssertErrorOutput(t, "Error stopping instance [i-abcdef12] (forced): test error\n", func() {
			ecloudInstanceStop(service, cmd, []string{"i-abcdef12"})
		})
	})
}

func Test_ecloudInstanceRestartCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudInstanceRestartCmd(nil).Args(nil, []string{"i-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudInstanceRestartCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing instance", err.Error())
	})
}

func Test_ecloudInstanceRestart(t *testing.T) {
	t.Run("WithoutForceFlag_CallsPowerRestartInstance", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerRestartInstance("i-abcdef12").Return(nil)

		ecloudInstanceRestart(service, &cobra.Command{}, []string{"i-abcdef12"})
	})

	t.Run("WithForceFlag_CallsPowerResetInstance", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerResetInstance("i-abcdef12").Return(nil)

		cmd := ecloudInstanceRestartCmd(nil)
		cmd.ParseFlags([]string{"--force"})

		ecloudInstanceRestart(service, cmd, []string{"i-abcdef12"})
	})

	t.Run("PowerRestartInstanceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerRestartInstance("i-abcdef12").Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error restarting instance [i-abcdef12]: test error\n", func() {
			ecloudInstanceRestart(service, &cobra.Command{}, []string{"i-abcdef12"})
		})
	})

	t.Run("PowerResetInstanceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerResetInstance("i-abcdef12").Return(errors.New("test error"))

		cmd := ecloudInstanceRestartCmd(nil)
		cmd.ParseFlags([]string{"--force"})

		test_output.AssertErrorOutput(t, "Error restarting instance [i-abcdef12] (forced): test error\n", func() {
			ecloudInstanceRestart(service, cmd, []string{"i-abcdef12"})
		})
	})
}
