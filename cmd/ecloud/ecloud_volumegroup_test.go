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
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func Test_ecloudVolumeGroupList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVolumeGroups(gomock.Any()).Return([]ecloud.VolumeGroup{}, nil).Times(1)

		ecloudVolumeGroupList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudVolumeGroupList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetVolumeGroupsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVolumeGroups(gomock.Any()).Return([]ecloud.VolumeGroup{}, errors.New("test error")).Times(1)

		err := ecloudVolumeGroupList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving volume groups: test error", err.Error())
	})
}

func Test_ecloudVolumeGroupShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVolumeGroupShowCmd(nil).Args(nil, []string{"volgroup-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVolumeGroupShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing volume-group", err.Error())
	})
}

func Test_ecloudVolumeGroupShow(t *testing.T) {
	t.Run("SingleVolumeGroup", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVolumeGroup("volgroup-abcdef12").Return(ecloud.VolumeGroup{}, nil).Times(1)

		ecloudVolumeGroupShow(service, &cobra.Command{}, []string{"volgroup-abcdef12"})
	})

	t.Run("MultipleVolumeGroups", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetVolumeGroup("volgroup-abcdef12").Return(ecloud.VolumeGroup{}, nil),
			service.EXPECT().GetVolumeGroup("volgroup-abcdef23").Return(ecloud.VolumeGroup{}, nil),
		)

		ecloudVolumeGroupShow(service, &cobra.Command{}, []string{"volgroup-abcdef12", "volgroup-abcdef23"})
	})

	t.Run("GetVolumeGroupError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVolumeGroup("volgroup-abcdef12").Return(ecloud.VolumeGroup{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving volume group [volgroup-abcdef12]: test error\n", func() {
			ecloudVolumeGroupShow(service, &cobra.Command{}, []string{"volgroup-abcdef12"})
		})
	})
}

func Test_ecloudVolumeGroupCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVolumeGroupCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvolumegroup", "--vpc=vpc-abcdef12", "--availability-zone=az-abcdef12"})

		req := ecloud.CreateVolumeGroupRequest{
			Name:     "testvolumegroup",
			VPCID:    "vpc-abcdef12",
			AvailabilityZoneID: "az-abcdef12",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "volgroup-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateVolumeGroup(req).Return(resp, nil),
			service.EXPECT().GetVolumeGroup("volgroup-abcdef12").Return(ecloud.VolumeGroup{}, nil),
		)

		ecloudVolumeGroupCreate(service, cmd, []string{})
	})

	t.Run("CreateWithWaitFlagNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVolumeGroupCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvolumegroup", "--vpc=vpc-abcdef12", "--availability-zone=az-abcdef12", "--wait"})

		req := ecloud.CreateVolumeGroupRequest{
			Name:     "testvolumegroup",
			VPCID:    "vpc-abcdef12",
			AvailabilityZoneID: "az-abcdef12",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "volgroup-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateVolumeGroup(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetVolumeGroup("volgroup-abcdef12").Return(ecloud.VolumeGroup{}, nil),
		)

		ecloudVolumeGroupCreate(service, cmd, []string{})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVolumeGroupCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvolumegroup", "--vpc=vpc-abcdef12", "--availability-zone=az-abcdef12", "--wait"})

		req := ecloud.CreateVolumeGroupRequest{
			Name:     "testvolumegroup",
			VPCID:    "vpc-abcdef12",
			AvailabilityZoneID: "az-abcdef12",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "volgroup-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateVolumeGroup(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{}, errors.New("test error")),
		)

		err := ecloudVolumeGroupCreate(service, cmd, []string{})

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "test error")
	})

	t.Run("CreateVolumeGroupError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVolumeGroupCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvolumegroup"})

		service.EXPECT().CreateVolumeGroup(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error")).Times(1)

		err := ecloudVolumeGroupCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating volume group: test error", err.Error())
	})

	t.Run("GetVolumeGroupError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVolumeGroupCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvolumegroup"})

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "volgroup-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateVolumeGroup(gomock.Any()).Return(resp, nil),
			service.EXPECT().GetVolumeGroup("volgroup-abcdef12").Return(ecloud.VolumeGroup{}, errors.New("test error")),
		)

		err := ecloudVolumeGroupCreate(service, cmd, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving new volume group: test error", err.Error())
	})
}

func Test_ecloudVolumeGroupUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVolumeGroupUpdateCmd(nil).Args(nil, []string{"volgroup-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVolumeGroupUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing volume-group", err.Error())
	})
}

func Test_ecloudVolumeGroupUpdate(t *testing.T) {
	t.Run("Default_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVolumeGroupUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvolumegroup"})

		req := ecloud.PatchVolumeGroupRequest{
			Name: "testvolumegroup",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "volgroup-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchVolumeGroup("volgroup-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetVolumeGroup("volgroup-abcdef12").Return(ecloud.VolumeGroup{}, nil),
		)

		ecloudVolumeGroupUpdate(service, cmd, []string{"volgroup-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVolumeGroupUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvolumegroup", "--wait"})

		req := ecloud.PatchVolumeGroupRequest{
			Name: "testvolumegroup",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "volgroup-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchVolumeGroup("volgroup-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetVolumeGroup("volgroup-abcdef12").Return(ecloud.VolumeGroup{}, nil),
		)

		ecloudVolumeGroupUpdate(service, cmd, []string{"volgroup-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVolumeGroupUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvolumegroup", "--wait"})

		req := ecloud.PatchVolumeGroupRequest{
			Name: "testvolumegroup",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "volgroup-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchVolumeGroup("volgroup-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for volume group [volgroup-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudVolumeGroupUpdate(service, cmd, []string{"volgroup-abcdef12"})
		})
	})

	t.Run("PatchVolumeGroupError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchVolumeGroup("volgroup-abcdef12", gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating volume group [volgroup-abcdef12]: test error\n", func() {
			ecloudVolumeGroupUpdate(service, &cobra.Command{}, []string{"volgroup-abcdef12"})
		})
	})

	t.Run("GetVolumeGroupError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "volgroup-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchVolumeGroup("volgroup-abcdef12", gomock.Any()).Return(resp, nil),
			service.EXPECT().GetVolumeGroup("volgroup-abcdef12").Return(ecloud.VolumeGroup{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated volume group [volgroup-abcdef12]: test error\n", func() {
			ecloudVolumeGroupUpdate(service, &cobra.Command{}, []string{"volgroup-abcdef12"})
		})
	})
}

func Test_ecloudVolumeGroupDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVolumeGroupDeleteCmd(nil).Args(nil, []string{"volgroup-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVolumeGroupDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing volume-group", err.Error())
	})
}

func Test_ecloudVolumeGroupDelete(t *testing.T) {
	t.Run("SingleVolumeGroup", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteVolumeGroup("volgroup-abcdef12").Return("task-abcdef12", nil).Times(1)

		ecloudVolumeGroupDelete(service, &cobra.Command{}, []string{"volgroup-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudVolumeGroupDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().DeleteVolumeGroup("volgroup-abcdef12").Return("task-abcdef12", nil)
		service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil)

		ecloudVolumeGroupDelete(service, cmd, []string{"volgroup-abcdef12"})
	})

	t.Run("WithWaitFlag_GetVolumeGroupError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudVolumeGroupDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().DeleteVolumeGroup("volgroup-abcdef12").Return("task-abcdef12", nil)
		service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for volume group [volgroup-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudVolumeGroupDelete(service, cmd, []string{"volgroup-abcdef12"})
		})
	})

	t.Run("DeleteVolumeGroupError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteVolumeGroup("volgroup-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing volume group [volgroup-abcdef12]: test error\n", func() {
			ecloudVolumeGroupDelete(service, &cobra.Command{}, []string{"volgroup-abcdef12"})
		})
	})
}
