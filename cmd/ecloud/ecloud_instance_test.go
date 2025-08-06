package ecloud

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/ptr"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
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

		assert.Equal(t, "error retrieving instances: test error", err.Error())
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
		assert.Equal(t, "missing instance", err.Error())
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
			Name:               "testinstance",
			ImageID:            "img-abcdef12",
			VCPUSockets:        1,
			VCPUCoresPerSocket: 1,
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
			Name:               "testinstance",
			ImageID:            "img-abcdef12",
			VCPUSockets:        1,
			VCPUCoresPerSocket: 1,
		}

		service.EXPECT().GetImages(connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				{
					Property: "name",
					Operator: connection.EQOperator,
					Value:    []string{"test"},
				},
			}}).Return([]ecloud.Image{{Name: "test", ID: "img-abcdef12"}}, nil)
		service.EXPECT().CreateInstance(req).Return("i-abcdef12", nil)
		service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{}, nil)

		err := ecloudInstanceCreate(service, cmd, []string{})
		assert.Nil(t, err)
	})

	t.Run("ImageRetrievalError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudInstanceCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testinstance", "--image=unknown"})

		service.EXPECT().GetImages(connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				{
					Property: "name",
					Operator: connection.EQOperator,
					Value:    []string{"unknown"},
				},
			},
		}).Return([]ecloud.Image{}, errors.New("test error"))

		err := ecloudInstanceCreate(service, cmd, []string{})
		assert.NotNil(t, err)
		assert.Equal(t, "error retrieving images: test error", err.Error())
	})

	t.Run("MultipleImagesRetrieved_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudInstanceCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testinstance", "--image=unknown"})

		service.EXPECT().GetImages(connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				{
					Property: "name",
					Operator: connection.EQOperator,
					Value:    []string{"unknown"},
				},
			},
		}).Return([]ecloud.Image{{
			ID: "img-1",
		}, {
			ID: "img-2",
		}}, nil)

		err := ecloudInstanceCreate(service, cmd, []string{})
		assert.NotNil(t, err)
		assert.Equal(t, "expected 1 image, got 2 images", err.Error())
	})

	t.Run("ImageNotFound_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudInstanceCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testinstance", "--image=unknown"})

		service.EXPECT().GetImages(connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				{
					Property: "name",
					Operator: connection.EQOperator,
					Value:    []string{"unknown"},
				},
			},
		}).Return([]ecloud.Image{}, nil)

		err := ecloudInstanceCreate(service, cmd, []string{})
		assert.NotNil(t, err)
		assert.Equal(t, "image not found with name 'unknown'", err.Error())
	})

	t.Run("CreateWithNetworkID_NoError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudInstanceCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testinstance", "--image=img-abcdef12", "--network=net-abcdef12"})

		req := ecloud.CreateInstanceRequest{
			Name:               "testinstance",
			ImageID:            "img-abcdef12",
			NetworkID:          "net-abcdef12",
			VCPUSockets:        1,
			VCPUCoresPerSocket: 1,
		}

		service.EXPECT().CreateInstance(req).Return("i-abcdef12", nil)
		service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{}, nil)

		err := ecloudInstanceCreate(service, cmd, []string{})
		assert.Nil(t, err)
	})

	t.Run("CreateInstanceError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudInstanceCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testinstance", "--image=img-abcdef12"})

		service.EXPECT().CreateInstance(gomock.Any()).Return("", errors.New("test error"))

		err := ecloudInstanceCreate(service, cmd, []string{})

		assert.Equal(t, "error creating instance: test error", err.Error())
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

		assert.Equal(t, "error retrieving new instance: test error", err.Error())
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
		assert.Equal(t, "missing instance", err.Error())
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

	t.Run("SingleInstanceWithVolumeGroup", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudInstanceUpdateCmd(nil)
		cmd.ParseFlags([]string{"--volume-group=volgroup-abcdef12"})

		req := ecloud.PatchInstanceRequest{
			VolumeGroupID: ptr.String("volgroup-abcdef12"),
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
		assert.Equal(t, "missing instance", err.Error())
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
		assert.Equal(t, "missing instance", err.Error())
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
		assert.Equal(t, "missing instance", err.Error())
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
		assert.Equal(t, "missing instance", err.Error())
	})
}

func Test_ecloudInstanceStart(t *testing.T) {
	t.Run("SingleInstance", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerOnInstance("i-abcdef12").Return("task-abcdef12", nil)

		ecloudInstanceStart(service, &cobra.Command{}, []string{"i-abcdef12"})
	})

	t.Run("MultipleInstances", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		gomock.InOrder(
			service.EXPECT().PowerOnInstance("i-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().PowerOnInstance("i-12abcdef").Return("task-abcdef23", nil),
		)

		ecloudInstanceStart(service, &cobra.Command{}, []string{"i-abcdef12", "i-12abcdef"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudInstanceStartCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerOnInstance("i-abcdef12").Return("task-abcdef12", nil)
		service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil)

		ecloudInstanceStart(service, cmd, []string{"i-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudInstanceStartCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerOnInstance("i-abcdef12").Return("task-abcdef12", nil)
		service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for instance [i-abcdef12]: error waiting for command: failed to retrieve task status: test error\n", func() {
			ecloudInstanceStart(service, cmd, []string{"i-abcdef12"})
		})
	})

	t.Run("PowerOnInstanceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerOnInstance("i-abcdef12").Return("", errors.New("test error"))

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
		assert.Equal(t, "missing instance", err.Error())
	})
}

func Test_ecloudInstanceStop(t *testing.T) {
	t.Run("WithoutForceFlag_CallsPowerShutdownInstance", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerShutdownInstance("i-abcdef12").Return("task-abcdef12", nil)

		ecloudInstanceStop(service, &cobra.Command{}, []string{"i-abcdef12"})
	})

	t.Run("WithForceFlag_CallsPowerOffInstance", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerOffInstance("i-abcdef12").Return("task-abcdef12", nil)

		cmd := ecloudInstanceStopCmd(nil)
		cmd.ParseFlags([]string{"--force"})

		ecloudInstanceStop(service, cmd, []string{"i-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudInstanceStopCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerShutdownInstance("i-abcdef12").Return("task-abcdef12", nil)
		service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil)

		ecloudInstanceStop(service, cmd, []string{"i-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudInstanceStopCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerShutdownInstance("i-abcdef12").Return("task-abcdef12", nil)
		service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for instance [i-abcdef12]: error waiting for command: failed to retrieve task status: test error\n", func() {
			ecloudInstanceStop(service, cmd, []string{"i-abcdef12"})
		})
	})

	t.Run("PowerShutdownInstanceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerShutdownInstance("i-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error stopping instance [i-abcdef12]: test error\n", func() {
			ecloudInstanceStop(service, &cobra.Command{}, []string{"i-abcdef12"})
		})
	})

	t.Run("PowerOffInstanceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerOffInstance("i-abcdef12").Return("", errors.New("test error"))

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
		assert.Equal(t, "missing instance", err.Error())
	})
}

func Test_ecloudInstanceRestart(t *testing.T) {
	t.Run("WithoutForceFlag_CallsPowerRestartInstance", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerRestartInstance("i-abcdef12").Return("task-abcdef12", nil)

		ecloudInstanceRestart(service, &cobra.Command{}, []string{"i-abcdef12"})
	})

	t.Run("WithForceFlag_CallsPowerResetInstance", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerResetInstance("i-abcdef12").Return("task-abcdef12", nil)

		cmd := ecloudInstanceRestartCmd(nil)
		cmd.ParseFlags([]string{"--force"})

		ecloudInstanceRestart(service, cmd, []string{"i-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudInstanceRestartCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerRestartInstance("i-abcdef12").Return("task-abcdef12", nil)
		service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil)

		ecloudInstanceRestart(service, cmd, []string{"i-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudInstanceRestartCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerRestartInstance("i-abcdef12").Return("task-abcdef12", nil)
		service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for instance [i-abcdef12]: error waiting for command: failed to retrieve task status: test error\n", func() {
			ecloudInstanceRestart(service, cmd, []string{"i-abcdef12"})
		})
	})

	t.Run("PowerRestartInstanceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerRestartInstance("i-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error restarting instance [i-abcdef12]: test error\n", func() {
			ecloudInstanceRestart(service, &cobra.Command{}, []string{"i-abcdef12"})
		})
	})

	t.Run("PowerResetInstanceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().PowerResetInstance("i-abcdef12").Return("", errors.New("test error"))

		cmd := ecloudInstanceRestartCmd(nil)
		cmd.ParseFlags([]string{"--force"})

		test_output.AssertErrorOutput(t, "Error restarting instance [i-abcdef12] (forced): test error\n", func() {
			ecloudInstanceRestart(service, cmd, []string{"i-abcdef12"})
		})
	})
}

func Test_ecloudInstanceMigrateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudInstanceMigrateCmd(nil).Args(nil, []string{"i-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudInstanceMigrateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing instance", err.Error())
	})
}

func Test_ecloudInstanceMigrate(t *testing.T) {
	t.Run("WithResourceTier_CallsMigrate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudInstanceMigrateCmd(nil)
		cmd.ParseFlags([]string{"--resource-tier=rt-abcdef12"})

		req := ecloud.MigrateInstanceRequest{
			ResourceTierID: "rt-abcdef12",
		}

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().MigrateInstance("i-abcdef12", req).Return("task-abcdef12", nil)

		ecloudInstanceMigrate(service, cmd, []string{"i-abcdef12"})
	})

	t.Run("WithHostGroup_CallsMigrate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudInstanceMigrateCmd(nil)
		cmd.ParseFlags([]string{"--host-group=hg-abcdef12"})

		req := ecloud.MigrateInstanceRequest{
			HostGroupID: "hg-abcdef12",
		}

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().MigrateInstance("i-abcdef12", req).Return("task-abcdef12", nil)

		ecloudInstanceMigrate(service, cmd, []string{"i-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudInstanceMigrateCmd(nil)
		cmd.ParseFlags([]string{"--resource-tier=rt-abcdef12", "--wait"})

		req := ecloud.MigrateInstanceRequest{
			ResourceTierID: "rt-abcdef12",
		}

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().MigrateInstance("i-abcdef12", req).Return("task-abcdef12", nil)
		service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil)

		ecloudInstanceMigrate(service, cmd, []string{"i-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudInstanceMigrateCmd(nil)
		cmd.ParseFlags([]string{"--resource-tier=rt-abcdef12", "--wait"})

		req := ecloud.MigrateInstanceRequest{
			ResourceTierID: "rt-abcdef12",
		}

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().MigrateInstance("i-abcdef12", req).Return("task-abcdef12", nil)
		service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for instance [i-abcdef12]: error waiting for command: failed to retrieve task status: test error\n", func() {
			ecloudInstanceMigrate(service, cmd, []string{"i-abcdef12"})
		})
	})

	t.Run("MigrateInstanceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().MigrateInstance("i-abcdef12", ecloud.MigrateInstanceRequest{}).Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error migrating instance [i-abcdef12]: test error\n", func() {
			ecloudInstanceMigrate(service, &cobra.Command{}, []string{"i-abcdef12"})
		})
	})
}

func Test_ecloudInstanceEncryptCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudInstanceEncryptCmd(nil).Args(nil, []string{"i-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudInstanceEncryptCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing instance", err.Error())
	})
}

func Test_ecloudInstanceEncrypt(t *testing.T) {
	t.Run("WithNoWaitFlag_ReturnsTaskID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudInstanceEncryptCmd(nil)

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().EncryptInstance("i-abcdef12").Return("task-abcdef12", nil)

		ecloudInstanceEncrypt(service, cmd, []string{"i-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudInstanceEncryptCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().EncryptInstance("i-abcdef12").Return("task-abcdef12", nil)
		service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil)

		ecloudInstanceEncrypt(service, cmd, []string{"i-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudInstanceEncryptCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().EncryptInstance("i-abcdef12").Return("task-abcdef12", nil)
		service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for instance [i-abcdef12]: error waiting for command: failed to retrieve task status: test error\n", func() {
			ecloudInstanceEncrypt(service, cmd, []string{"i-abcdef12"})
		})
	})

	t.Run("EncryptInstanceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().EncryptInstance("i-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error encrypting instance [i-abcdef12]: test error\n", func() {
			ecloudInstanceEncrypt(service, &cobra.Command{}, []string{"i-abcdef12"})
		})
	})
}

func Test_ecloudInstanceDecryptCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudInstanceDecryptCmd(nil).Args(nil, []string{"i-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudInstanceDecryptCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing instance", err.Error())
	})
}

func Test_ecloudInstanceDecrypt(t *testing.T) {
	t.Run("WithNoWaitFlag_ReturnsTaskID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudInstanceDecryptCmd(nil)

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().DecryptInstance("i-abcdef12").Return("task-abcdef12", nil)

		ecloudInstanceDecrypt(service, cmd, []string{"i-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudInstanceDecryptCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().DecryptInstance("i-abcdef12").Return("task-abcdef12", nil)
		service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil)

		ecloudInstanceDecrypt(service, cmd, []string{"i-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudInstanceDecryptCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().DecryptInstance("i-abcdef12").Return("task-abcdef12", nil)
		service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for instance [i-abcdef12]: error waiting for command: failed to retrieve task status: test error\n", func() {
			ecloudInstanceDecrypt(service, cmd, []string{"i-abcdef12"})
		})
	})

	t.Run("DecryptInstanceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().DecryptInstance("i-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error decrypting instance [i-abcdef12]: test error\n", func() {
			ecloudInstanceDecrypt(service, &cobra.Command{}, []string{"i-abcdef12"})
		})
	})
}

func Test_tagLookup(t *testing.T) {
	t.Run("EmptyTag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		_, err := tagLookup(service, "")

		assert.NotNil(t, err)
		assert.Equal(t, "cannot lookup tag with empty value", err.Error())
	})

	t.Run("TagIDWithPrefix_ReturnsID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		tagID, err := tagLookup(service, "tag-abcdef12")

		assert.Nil(t, err)
		assert.Equal(t, "tag-abcdef12", tagID)
	})

	t.Run("ScopedTagName_ReturnsID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		expectedParams := connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				{
					Property: "scope",
					Operator: connection.EQOperator,
					Value:    []string{"test-scope"},
				},
				{
					Property: "name",
					Operator: connection.EQOperator,
					Value:    []string{"test-tag"},
				},
			},
		}

		service.EXPECT().GetTags(expectedParams).Return([]ecloud.Tag{{ID: "tag-abcdef12", Scope: "test-scope", Name: "test-tag"}}, nil)

		tagID, err := tagLookup(service, "test-scope:test-tag")

		assert.Nil(t, err)
		assert.Equal(t, "tag-abcdef12", tagID)
	})

	t.Run("TagNameOnly_ReturnsID", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		expectedParams := connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				{
					Property: "name",
					Operator: connection.EQOperator,
					Value:    []string{"test-tag"},
				},
			},
		}

		service.EXPECT().GetTags(expectedParams).Return([]ecloud.Tag{{ID: "tag-abcdef12", Name: "test-tag"}}, nil)

		tagID, err := tagLookup(service, "test-tag")

		assert.Nil(t, err)
		assert.Equal(t, "tag-abcdef12", tagID)
	})

	t.Run("InvalidScopedFormat_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		_, err := tagLookup(service, "invalid:scope:format")

		assert.NotNil(t, err)
		assert.Equal(t, "invalid tag format 'invalid:scope:format', expected '<scope>:<name>'", err.Error())
	})

	t.Run("GetTagsError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetTags(gomock.Any()).Return([]ecloud.Tag{}, errors.New("test error"))

		_, err := tagLookup(service, "test-tag")

		assert.NotNil(t, err)
		assert.Equal(t, "error retrieving tags: test error", err.Error())
	})

	t.Run("TagNotFound_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetTags(gomock.Any()).Return([]ecloud.Tag{}, nil)

		_, err := tagLookup(service, "nonexistent-tag")

		assert.NotNil(t, err)
		assert.Equal(t, "tag 'nonexistent-tag' not found, create the tag first", err.Error())
	})

	t.Run("MultipleTags_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetTags(gomock.Any()).Return([]ecloud.Tag{
			{ID: "tag-abcdef12", Scope: "scope1", Name: "test-tag"},
			{ID: "tag-12abcdef", Scope: "scope2", Name: "test-tag"},
		}, nil)

		_, err := tagLookup(service, "test-tag")

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "expected 1 tag, got 2 tags")
	})
}

func Test_addRemoveTags(t *testing.T) {
	t.Run("AddTags_WithExistingTagsInRequest_Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		existingTags := []string{"tag-existing1", "tag-existing2"}
		request := &ecloud.PatchInstanceRequest{
			TagIDs: &existingTags,
		}

		service.EXPECT().GetTags(gomock.Any()).Return([]ecloud.Tag{{ID: "tag-new1", Name: "new-tag"}}, nil)

		err := addRemoveTags(service, request, "i-abcdef12", []string{"new-tag"}, true)

		assert.Nil(t, err)
		assert.Len(t, *request.TagIDs, 3)
		assert.Contains(t, *request.TagIDs, "tag-existing1")
		assert.Contains(t, *request.TagIDs, "tag-existing2")
		assert.Contains(t, *request.TagIDs, "tag-new1")
	})

	t.Run("AddTags_WithoutExistingTagsInRequest_Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		request := &ecloud.PatchInstanceRequest{}

		instance := ecloud.Instance{
			Tags: []ecloud.ResourceTag{
				{ID: "tag-existing1", Name: "existing-tag1"},
				{ID: "tag-existing2", Name: "existing-tag2"},
			},
		}

		service.EXPECT().GetInstance("i-abcdef12").Return(instance, nil)
		service.EXPECT().GetTags(gomock.Any()).Return([]ecloud.Tag{{ID: "tag-new1", Name: "new-tag"}}, nil)

		err := addRemoveTags(service, request, "i-abcdef12", []string{"new-tag"}, true)

		assert.Nil(t, err)
		assert.Len(t, *request.TagIDs, 3)
		assert.Contains(t, *request.TagIDs, "tag-existing1")
		assert.Contains(t, *request.TagIDs, "tag-existing2")
		assert.Contains(t, *request.TagIDs, "tag-new1")
	})

	t.Run("RemoveTags_WithExistingTagsInRequest_Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		existingTags := []string{"tag-keep1", "tag-remove1", "tag-keep2"}
		request := &ecloud.PatchInstanceRequest{
			TagIDs: &existingTags,
		}

		service.EXPECT().GetTags(gomock.Any()).Return([]ecloud.Tag{{ID: "tag-remove1", Name: "remove-tag"}}, nil)

		err := addRemoveTags(service, request, "i-abcdef12", []string{"remove-tag"}, false)

		assert.Nil(t, err)
		assert.Len(t, *request.TagIDs, 2)
		assert.Contains(t, *request.TagIDs, "tag-keep1")
		assert.Contains(t, *request.TagIDs, "tag-keep2")
		assert.NotContains(t, *request.TagIDs, "tag-remove1")
	})

	t.Run("RemoveTags_WithoutExistingTagsInRequest_Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		request := &ecloud.PatchInstanceRequest{}

		instance := ecloud.Instance{
			Tags: []ecloud.ResourceTag{
				{ID: "tag-keep1", Name: "keep-tag1"},
				{ID: "tag-remove1", Name: "remove-tag"},
				{ID: "tag-keep2", Name: "keep-tag2"},
			},
		}

		service.EXPECT().GetInstance("i-abcdef12").Return(instance, nil)
		service.EXPECT().GetTags(gomock.Any()).Return([]ecloud.Tag{{ID: "tag-remove1", Name: "remove-tag"}}, nil)

		err := addRemoveTags(service, request, "i-abcdef12", []string{"remove-tag"}, false)

		assert.Nil(t, err)
		assert.Len(t, *request.TagIDs, 2)
		assert.Contains(t, *request.TagIDs, "tag-keep1")
		assert.Contains(t, *request.TagIDs, "tag-keep2")
		assert.NotContains(t, *request.TagIDs, "tag-remove1")
	})

	t.Run("AddTags_EmptyAndWhitespaceTagsSkipped", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		request := &ecloud.PatchInstanceRequest{}

		instance := ecloud.Instance{
			Tags: []ecloud.ResourceTag{{ID: "tag-existing1", Name: "existing-tag1"}},
		}

		service.EXPECT().GetInstance("i-abcdef12").Return(instance, nil)
		service.EXPECT().GetTags(gomock.Any()).Return([]ecloud.Tag{{ID: "tag-new1", Name: "new-tag"}}, nil)

		err := addRemoveTags(service, request, "i-abcdef12", []string{"", "  ", "new-tag", "\t"}, true)

		assert.Nil(t, err)
		assert.Len(t, *request.TagIDs, 2)
		assert.Contains(t, *request.TagIDs, "tag-existing1")
		assert.Contains(t, *request.TagIDs, "tag-new1")
	})

	t.Run("AddTags_DuplicateTagsNotDuplicated", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		existingTags := []string{"tag-existing1"}
		request := &ecloud.PatchInstanceRequest{
			TagIDs: &existingTags,
		}

		service.EXPECT().GetTags(gomock.Any()).Return([]ecloud.Tag{{ID: "tag-existing1", Name: "existing-tag"}}, nil)

		err := addRemoveTags(service, request, "i-abcdef12", []string{"existing-tag"}, true)

		assert.Nil(t, err)
		assert.Len(t, *request.TagIDs, 1)
		assert.Contains(t, *request.TagIDs, "tag-existing1")
	})

	t.Run("GetInstanceError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		request := &ecloud.PatchInstanceRequest{}

		service.EXPECT().GetInstance("i-abcdef12").Return(ecloud.Instance{}, errors.New("test error"))

		err := addRemoveTags(service, request, "i-abcdef12", []string{"new-tag"}, true)

		assert.NotNil(t, err)
		assert.Equal(t, "error retrieving instance tags [i-abcdef12]: test error", err.Error())
	})

	t.Run("TagLookupError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		request := &ecloud.PatchInstanceRequest{}

		instance := ecloud.Instance{Tags: []ecloud.ResourceTag{}}
		service.EXPECT().GetInstance("i-abcdef12").Return(instance, nil)
		service.EXPECT().GetTags(gomock.Any()).Return([]ecloud.Tag{}, errors.New("tag lookup error"))

		err := addRemoveTags(service, request, "i-abcdef12", []string{"nonexistent-tag"}, true)

		assert.NotNil(t, err)
		assert.Equal(t, "error retrieving tags: tag lookup error", err.Error())
	})

	t.Run("FinalTagsSortedAlphabetically", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		request := &ecloud.PatchInstanceRequest{}

		instance := ecloud.Instance{
			Tags: []ecloud.ResourceTag{
				{ID: "tag-zebra", Name: "zebra"},
				{ID: "tag-alpha", Name: "alpha"},
			},
		}

		service.EXPECT().GetInstance("i-abcdef12").Return(instance, nil)
		service.EXPECT().GetTags(gomock.Any()).Return([]ecloud.Tag{{ID: "tag-beta", Name: "beta"}}, nil)

		err := addRemoveTags(service, request, "i-abcdef12", []string{"beta"}, true)

		assert.Nil(t, err)
		assert.Len(t, *request.TagIDs, 3)

		// Verify tags are sorted alphabetically
		expectedOrder := []string{"tag-alpha", "tag-beta", "tag-zebra"}
		assert.Equal(t, expectedOrder, *request.TagIDs)
	})

	t.Run("AddMultipleTags_Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		request := &ecloud.PatchInstanceRequest{}

		instance := ecloud.Instance{
			Tags: []ecloud.ResourceTag{{ID: "tag-existing1", Name: "existing-tag1"}},
		}

		service.EXPECT().GetInstance("i-abcdef12").Return(instance, nil)

		// Multiple calls to GetTags for each tag lookup
		service.EXPECT().GetTags(gomock.Any()).Return([]ecloud.Tag{{ID: "tag-new1", Name: "new-tag1"}}, nil)
		service.EXPECT().GetTags(gomock.Any()).Return([]ecloud.Tag{{ID: "tag-new2", Name: "new-tag2"}}, nil)

		err := addRemoveTags(service, request, "i-abcdef12", []string{"new-tag1", "new-tag2"}, true)

		assert.Nil(t, err)
		assert.Len(t, *request.TagIDs, 3)
		assert.Contains(t, *request.TagIDs, "tag-existing1")
		assert.Contains(t, *request.TagIDs, "tag-new1")
		assert.Contains(t, *request.TagIDs, "tag-new2")
	})

	t.Run("RemoveMultipleTags_Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		request := &ecloud.PatchInstanceRequest{}

		instance := ecloud.Instance{
			Tags: []ecloud.ResourceTag{
				{ID: "tag-keep1", Name: "keep-tag1"},
				{ID: "tag-remove1", Name: "remove-tag1"},
				{ID: "tag-remove2", Name: "remove-tag2"},
				{ID: "tag-keep2", Name: "keep-tag2"},
			},
		}

		service.EXPECT().GetInstance("i-abcdef12").Return(instance, nil)

		// Multiple calls to GetTags for each tag lookup
		service.EXPECT().GetTags(gomock.Any()).Return([]ecloud.Tag{{ID: "tag-remove1", Name: "remove-tag1"}}, nil)
		service.EXPECT().GetTags(gomock.Any()).Return([]ecloud.Tag{{ID: "tag-remove2", Name: "remove-tag2"}}, nil)

		err := addRemoveTags(service, request, "i-abcdef12", []string{"remove-tag1", "remove-tag2"}, false)

		assert.Nil(t, err)
		assert.Len(t, *request.TagIDs, 2)
		assert.Contains(t, *request.TagIDs, "tag-keep1")
		assert.Contains(t, *request.TagIDs, "tag-keep2")
		assert.NotContains(t, *request.TagIDs, "tag-remove1")
		assert.NotContains(t, *request.TagIDs, "tag-remove2")
	})
}
