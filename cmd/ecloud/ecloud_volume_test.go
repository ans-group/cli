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

func Test_ecloudVolumeList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVolumes(gomock.Any()).Return([]ecloud.Volume{}, nil).Times(1)

		ecloudVolumeList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudVolumeList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetVolumesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVolumes(gomock.Any()).Return([]ecloud.Volume{}, errors.New("test error")).Times(1)

		err := ecloudVolumeList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving volumes: test error", err.Error())
	})
}

func Test_ecloudVolumeShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVolumeShowCmd(nil).Args(nil, []string{"vol-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVolumeShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing volume", err.Error())
	})
}

func Test_ecloudVolumeShow(t *testing.T) {
	t.Run("SingleVolume", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVolume("vol-abcdef12").Return(ecloud.Volume{}, nil).Times(1)

		ecloudVolumeShow(service, &cobra.Command{}, []string{"vol-abcdef12"})
	})

	t.Run("MultipleVolumes", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetVolume("vol-abcdef12").Return(ecloud.Volume{}, nil),
			service.EXPECT().GetVolume("vol-abcdef23").Return(ecloud.Volume{}, nil),
		)

		ecloudVolumeShow(service, &cobra.Command{}, []string{"vol-abcdef12", "vol-abcdef23"})
	})

	t.Run("GetVolumeError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVolume("vol-abcdef12").Return(ecloud.Volume{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving volume [vol-abcdef12]: test error\n", func() {
			ecloudVolumeShow(service, &cobra.Command{}, []string{"vol-abcdef12"})
		})
	})
}

func Test_ecloudVolumeCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVolumeCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvolume", "--vpc=vpc-abcdef12", "--capacity=20"})

		req := ecloud.CreateVolumeRequest{
			Name:     "testvolume",
			VPCID:    "vpc-abcdef12",
			Capacity: 20,
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vol-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateVolume(req).Return(resp, nil),
			service.EXPECT().GetVolume("vol-abcdef12").Return(ecloud.Volume{}, nil),
		)

		ecloudVolumeCreate(service, cmd, []string{})
	})

	t.Run("CreateWithWaitFlagNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVolumeCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvolume", "--vpc=vpc-abcdef12", "--capacity=20", "--wait"})

		req := ecloud.CreateVolumeRequest{
			Name:     "testvolume",
			VPCID:    "vpc-abcdef12",
			Capacity: 20,
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vol-abcdef12",
		}

		taskReqParameters := connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				{
					Property: "id",
					Operator: connection.EQOperator,
					Value:    []string{"task-abcdef12"},
				},
			},
		}

		taskResp := []ecloud.Task{
			{
				Status: ecloud.TaskStatusComplete,
			},
		}

		gomock.InOrder(
			service.EXPECT().CreateVolume(req).Return(resp, nil),
			service.EXPECT().GetVolumeTasks("vol-abcdef12", taskReqParameters).Return(taskResp, nil),
			service.EXPECT().GetVolume("vol-abcdef12").Return(ecloud.Volume{}, nil),
		)

		ecloudVolumeCreate(service, cmd, []string{})
	})

	t.Run("GetVolumeTasksError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVolumeCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvolume", "--vpc=vpc-abcdef12", "--capacity=20", "--wait"})

		req := ecloud.CreateVolumeRequest{
			Name:     "testvolume",
			VPCID:    "vpc-abcdef12",
			Capacity: 20,
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vol-abcdef12",
		}

		taskReqParameters := connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				{
					Property: "id",
					Operator: connection.EQOperator,
					Value:    []string{"task-abcdef12"},
				},
			},
		}

		gomock.InOrder(
			service.EXPECT().CreateVolume(req).Return(resp, nil),
			service.EXPECT().GetVolumeTasks("vol-abcdef12", taskReqParameters).Return([]ecloud.Task{}, errors.New("test error")),
		)

		err := ecloudVolumeCreate(service, cmd, []string{})

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "test error")
	})

	t.Run("CreateVolumeError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVolumeCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvolume"})

		service.EXPECT().CreateVolume(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error")).Times(1)

		err := ecloudVolumeCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating volume: test error", err.Error())
	})

	t.Run("GetVolumeError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVolumeCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvolume"})

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vol-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateVolume(gomock.Any()).Return(resp, nil),
			service.EXPECT().GetVolume("vol-abcdef12").Return(ecloud.Volume{}, errors.New("test error")),
		)

		err := ecloudVolumeCreate(service, cmd, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving new volume: test error", err.Error())
	})
}

func Test_ecloudVolumeUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVolumeUpdateCmd(nil).Args(nil, []string{"vol-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVolumeUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing volume", err.Error())
	})
}

func Test_ecloudVolumeUpdate(t *testing.T) {
	t.Run("Default_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVolumeUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvolume"})

		req := ecloud.PatchVolumeRequest{
			Name: "testvolume",
		}

		gomock.InOrder(
			service.EXPECT().PatchVolume("vol-abcdef12", req).Return("task-abcdef12", nil),
			service.EXPECT().GetVolume("vol-abcdef12").Return(ecloud.Volume{}, nil),
		)

		ecloudVolumeUpdate(service, cmd, []string{"vol-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVolumeUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvolume", "--wait"})

		req := ecloud.PatchVolumeRequest{
			Name: "testvolume",
		}

		taskReqParameters := connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				{
					Property: "id",
					Operator: connection.EQOperator,
					Value:    []string{"task-abcdef12"},
				},
			},
		}

		taskResp := []ecloud.Task{
			{
				Status: ecloud.TaskStatusComplete,
			},
		}

		gomock.InOrder(
			service.EXPECT().PatchVolume("vol-abcdef12", req).Return("task-abcdef12", nil),
			service.EXPECT().GetVolumeTasks("vol-abcdef12", taskReqParameters).Return(taskResp, nil),
			service.EXPECT().GetVolume("vol-abcdef12").Return(ecloud.Volume{}, nil),
		)

		ecloudVolumeUpdate(service, cmd, []string{"vol-abcdef12"})
	})

	t.Run("GetVolumeTasksError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVolumeUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvolume", "--wait"})

		req := ecloud.PatchVolumeRequest{
			Name: "testvolume",
		}

		taskReqParameters := connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				{
					Property: "id",
					Operator: connection.EQOperator,
					Value:    []string{"task-abcdef12"},
				},
			},
		}

		gomock.InOrder(
			service.EXPECT().PatchVolume("vol-abcdef12", req).Return("task-abcdef12", nil),
			service.EXPECT().GetVolumeTasks("vol-abcdef12", taskReqParameters).Return([]ecloud.Task{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for volume task to complete for volume [vol-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudVolumeUpdate(service, cmd, []string{"vol-abcdef12"})
		})
	})

	t.Run("PatchVolumeError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchVolume("vol-abcdef12", gomock.Any()).Return("task-abcdef12", errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating volume [vol-abcdef12]: test error\n", func() {
			ecloudVolumeUpdate(service, &cobra.Command{}, []string{"vol-abcdef12"})
		})
	})

	t.Run("GetVolumeError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchVolume("vol-abcdef12", gomock.Any()).Return("task-abcdef12", nil),
			service.EXPECT().GetVolume("vol-abcdef12").Return(ecloud.Volume{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated volume [vol-abcdef12]: test error\n", func() {
			ecloudVolumeUpdate(service, &cobra.Command{}, []string{"vol-abcdef12"})
		})
	})
}

func Test_ecloudVolumeDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVolumeDeleteCmd(nil).Args(nil, []string{"vol-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVolumeDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing volume", err.Error())
	})
}

func Test_ecloudVolumeDelete(t *testing.T) {
	t.Run("SingleVolume", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteVolume("vol-abcdef12").Return("task-abcdef12", nil).Times(1)

		ecloudVolumeDelete(service, &cobra.Command{}, []string{"vol-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudVolumeDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().DeleteVolume("vol-abcdef12").Return("task-abcdef12", nil)
		service.EXPECT().GetVolume("vol-abcdef12").Return(ecloud.Volume{}, &ecloud.VolumeNotFoundError{})

		ecloudVolumeDelete(service, cmd, []string{"vol-abcdef12"})
	})

	t.Run("WithWaitFlag_GetVolumeError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudVolumeDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().DeleteVolume("vol-abcdef12").Return("task-abcdef12", nil)
		service.EXPECT().GetVolume("vol-abcdef12").Return(ecloud.Volume{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error waiting for volume [vol-abcdef12] to be removed: Error waiting for command: Failed to retrieve volume [vol-abcdef12]: test error\n", func() {
			ecloudVolumeDelete(service, cmd, []string{"vol-abcdef12"})
		})
	})

	t.Run("DeleteVolumeError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteVolume("vol-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing volume [vol-abcdef12]: test error\n", func() {
			ecloudVolumeDelete(service, &cobra.Command{}, []string{"vol-abcdef12"})
		})
	})
}
