package ecloud

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ecloudInstanceVolumeListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudInstanceVolumeListCmd(nil).Args(nil, []string{"i-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudInstanceVolumeListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing instance", err.Error())
	})
}

func Test_ecloudInstanceVolumeList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetInstanceVolumes("i-abcdef12", gomock.Any()).Return([]ecloud.Volume{}, nil).Times(1)

		ecloudInstanceVolumeList(service, &cobra.Command{}, []string{"i-abcdef12"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudInstanceVolumeList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetInstancesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetInstanceVolumes("i-abcdef12", gomock.Any()).Return([]ecloud.Volume{}, errors.New("test error")).Times(1)

		err := ecloudInstanceVolumeList(service, &cobra.Command{}, []string{"i-abcdef12"})

		assert.Equal(t, "Error retrieving instance volumes: test error", err.Error())
	})
}

func Test_ecloudInstanceVolumeAttach(t *testing.T) {
	t.Run("SingleInstance", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudInstanceVolumeAttachCmd(nil)
		cmd.ParseFlags([]string{"--volume=vol-abcdef12"})

		req := ecloud.AttachDetachInstanceVolumeRequest{
			VolumeID: "vol-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().AttachInstanceVolume("i-abcdef12", req).Return("task-abcdef12", nil),
		)

		ecloudInstanceVolumeAttach(service, cmd, []string{"i-abcdef12"})
	})

	t.Run("WithWaitFlag", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudInstanceVolumeAttachCmd(nil)
		cmd.ParseFlags([]string{"--volume=vol-abcdef12", "--wait"})

		req := ecloud.AttachDetachInstanceVolumeRequest{
			VolumeID: "vol-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().AttachInstanceVolume("i-abcdef12", req).Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudInstanceVolumeAttach(service, cmd, []string{"i-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudInstanceVolumeAttachCmd(nil)
		cmd.ParseFlags([]string{"--volume=vol-abcdef12", "--wait"})

		req := ecloud.AttachDetachInstanceVolumeRequest{
			VolumeID: "vol-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().AttachInstanceVolume("i-abcdef12", req).Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{}, errors.New("test error")),
		)

		err := ecloudInstanceVolumeAttach(service, cmd, []string{"i-abcdef12"})
		assert.Equal(t, "Error waiting for task: Error waiting for command: Failed to retrieve task status: test error", err.Error())
	})

	t.Run("AttachInstanceVolumeError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().AttachInstanceVolume("i-abcdef12", gomock.Any()).Return("", errors.New("test error"))

		err := ecloudInstanceVolumeAttach(service, &cobra.Command{}, []string{"i-abcdef12"})
		assert.Equal(t, "Error attaching instance volume: test error", err.Error())
	})
}

func Test_ecloudInstanceVolumeDetach(t *testing.T) {
	t.Run("SingleInstance", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudInstanceVolumeDetachCmd(nil)
		cmd.ParseFlags([]string{"--volume=vol-abcdef12"})

		req := ecloud.AttachDetachInstanceVolumeRequest{
			VolumeID: "vol-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().DetachInstanceVolume("i-abcdef12", req).Return("task-abcdef12", nil),
		)

		ecloudInstanceVolumeDetach(service, cmd, []string{"i-abcdef12"})
	})

	t.Run("WithWaitFlag", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudInstanceVolumeDetachCmd(nil)
		cmd.ParseFlags([]string{"--volume=vol-abcdef12", "--wait"})

		req := ecloud.AttachDetachInstanceVolumeRequest{
			VolumeID: "vol-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().DetachInstanceVolume("i-abcdef12", req).Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudInstanceVolumeDetach(service, cmd, []string{"i-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudInstanceVolumeDetachCmd(nil)
		cmd.ParseFlags([]string{"--volume=vol-abcdef12", "--wait"})

		req := ecloud.AttachDetachInstanceVolumeRequest{
			VolumeID: "vol-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().DetachInstanceVolume("i-abcdef12", req).Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{}, errors.New("test error")),
		)

		err := ecloudInstanceVolumeDetach(service, cmd, []string{"i-abcdef12"})
		assert.Equal(t, "Error waiting for task: Error waiting for command: Failed to retrieve task status: test error", err.Error())
	})

	t.Run("DetachInstanceVolumeError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DetachInstanceVolume("i-abcdef12", gomock.Any()).Return("", errors.New("test error"))

		err := ecloudInstanceVolumeDetach(service, &cobra.Command{}, []string{"i-abcdef12"})
		assert.Equal(t, "Error detaching instance volume: test error", err.Error())
	})
}
