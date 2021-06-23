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

func Test_ecloudHostGroupList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetHostGroups(gomock.Any()).Return([]ecloud.HostGroup{}, nil).Times(1)

		ecloudHostGroupList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudHostGroupList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetHostGroupsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetHostGroups(gomock.Any()).Return([]ecloud.HostGroup{}, errors.New("test error")).Times(1)

		err := ecloudHostGroupList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving host groups: test error", err.Error())
	})
}

func Test_ecloudHostGroupShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudHostGroupShowCmd(nil).Args(nil, []string{"hg-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudHostGroupShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing host group", err.Error())
	})
}

func Test_ecloudHostGroupShow(t *testing.T) {
	t.Run("SingleHostGroup", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetHostGroup("hg-abcdef12").Return(ecloud.HostGroup{}, nil).Times(1)

		ecloudHostGroupShow(service, &cobra.Command{}, []string{"hg-abcdef12"})
	})

	t.Run("MultipleHostGroups", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetHostGroup("hg-abcdef12").Return(ecloud.HostGroup{}, nil),
			service.EXPECT().GetHostGroup("hg-abcdef23").Return(ecloud.HostGroup{}, nil),
		)

		ecloudHostGroupShow(service, &cobra.Command{}, []string{"hg-abcdef12", "hg-abcdef23"})
	})

	t.Run("GetHostGroupError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetHostGroup("hg-abcdef12").Return(ecloud.HostGroup{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving host group [hg-abcdef12]: test error\n", func() {
			ecloudHostGroupShow(service, &cobra.Command{}, []string{"hg-abcdef12"})
		})
	})
}

func Test_ecloudHostGroupCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudHostGroupCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgroup", "--vpc=vpc-abcdef12"})

		req := ecloud.CreateHostGroupRequest{
			Name:  "testgroup",
			VPCID: "vpc-abcdef12",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "hg-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateHostGroup(req).Return(resp, nil),
			service.EXPECT().GetHostGroup("hg-abcdef12").Return(ecloud.HostGroup{}, nil),
		)

		ecloudHostGroupCreate(service, cmd, []string{})
	})

	t.Run("CreateWithWaitFlagNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudHostGroupCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgroup", "--vpc=vpc-abcdef12", "--wait"})

		req := ecloud.CreateHostGroupRequest{
			Name:  "testgroup",
			VPCID: "vpc-abcdef12",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "hg-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateHostGroup(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetHostGroup("hg-abcdef12").Return(ecloud.HostGroup{}, nil),
		)

		ecloudHostGroupCreate(service, cmd, []string{})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudHostGroupCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgroup", "--vpc=vpc-abcdef12", "--wait"})

		req := ecloud.CreateHostGroupRequest{
			Name:  "testgroup",
			VPCID: "vpc-abcdef12",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "hg-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateHostGroup(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		err := ecloudHostGroupCreate(service, cmd, []string{})
		assert.Equal(t, "Error waiting for host group task to complete: Error waiting for command: Failed to retrieve task status: test error", err.Error())
	})

	t.Run("CreateHostGroupError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudHostGroupCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgroup", "--direction=IN", "--action=DROP"})

		service.EXPECT().CreateHostGroup(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error")).Times(1)

		err := ecloudHostGroupCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating host group: test error", err.Error())
	})

	t.Run("GetHostGroupError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudHostGroupCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgroup", "--direction=IN", "--action=DROP"})

		gomock.InOrder(
			service.EXPECT().CreateHostGroup(gomock.Any()).Return(ecloud.TaskReference{ResourceID: "hg-abcdef12"}, nil),
			service.EXPECT().GetHostGroup("hg-abcdef12").Return(ecloud.HostGroup{}, errors.New("test error")),
		)

		err := ecloudHostGroupCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new host group: test error", err.Error())
	})
}

func Test_ecloudHostGroupUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudHostGroupUpdateCmd(nil).Args(nil, []string{"hg-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudHostGroupUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing host group", err.Error())
	})
}

func Test_ecloudHostGroupUpdate(t *testing.T) {
	t.Run("SingleGroup", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudHostGroupUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgroup"})

		req := ecloud.PatchHostGroupRequest{
			Name: "testgroup",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "hg-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchHostGroup("hg-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetHostGroup("hg-abcdef12").Return(ecloud.HostGroup{}, nil),
		)

		ecloudHostGroupUpdate(service, cmd, []string{"hg-abcdef12"})
	})

	t.Run("MultipleHostGroups", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp1 := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "hg-abcdef12",
		}

		resp2 := ecloud.TaskReference{
			TaskID:     "task-abcdef23",
			ResourceID: "hg-12abcdef",
		}

		gomock.InOrder(
			service.EXPECT().PatchHostGroup("hg-abcdef12", gomock.Any()).Return(resp1, nil),
			service.EXPECT().GetHostGroup("hg-abcdef12").Return(ecloud.HostGroup{}, nil),
			service.EXPECT().PatchHostGroup("hg-12abcdef", gomock.Any()).Return(resp2, nil),
			service.EXPECT().GetHostGroup("hg-12abcdef").Return(ecloud.HostGroup{}, nil),
		)

		ecloudHostGroupUpdate(service, &cobra.Command{}, []string{"hg-abcdef12", "hg-12abcdef"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudHostGroupUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgroup", "--wait"})

		req := ecloud.PatchHostGroupRequest{
			Name: "testgroup",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "hg-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchHostGroup("hg-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetHostGroup("hg-abcdef12").Return(ecloud.HostGroup{}, nil),
		)

		ecloudHostGroupUpdate(service, cmd, []string{"hg-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudHostGroupUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgroup", "--wait"})

		req := ecloud.PatchHostGroupRequest{
			Name: "testgroup",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "hg-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchHostGroup("hg-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for host group [hg-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudHostGroupUpdate(service, cmd, []string{"hg-abcdef12"})
		})
	})

	t.Run("PatchHostGroupError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchHostGroup("hg-abcdef12", gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating host group [hg-abcdef12]: test error\n", func() {
			ecloudHostGroupUpdate(service, &cobra.Command{}, []string{"hg-abcdef12"})
		})
	})

	t.Run("GetHostGroupError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "hg-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchHostGroup("hg-abcdef12", gomock.Any()).Return(resp, nil),
			service.EXPECT().GetHostGroup("hg-abcdef12").Return(ecloud.HostGroup{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated host group [hg-abcdef12]: test error\n", func() {
			ecloudHostGroupUpdate(service, &cobra.Command{}, []string{"hg-abcdef12"})
		})
	})
}

func Test_ecloudHostGroupDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudHostGroupDeleteCmd(nil).Args(nil, []string{"hg-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudHostGroupDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing host group", err.Error())
	})
}

func Test_ecloudHostGroupDelete(t *testing.T) {
	t.Run("SingleGroup", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteHostGroup("hg-abcdef12").Return("task-abcdef12", nil).Times(1)

		ecloudHostGroupDelete(service, &cobra.Command{}, []string{"hg-abcdef12"})
	})

	t.Run("MultipleHostGroups", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteHostGroup("hg-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().DeleteHostGroup("hg-12abcdef").Return("task-abcdef23", nil),
		)

		ecloudHostGroupDelete(service, &cobra.Command{}, []string{"hg-abcdef12", "hg-12abcdef"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudHostGroupDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteHostGroup("hg-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudHostGroupDelete(service, cmd, []string{"hg-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudHostGroupDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		gomock.InOrder(
			service.EXPECT().DeleteHostGroup("hg-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for host group [hg-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudHostGroupDelete(service, cmd, []string{"hg-abcdef12"})
		})
	})

	t.Run("DeleteHostGroupError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteHostGroup("hg-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing host group [hg-abcdef12]: test error\n", func() {
			ecloudHostGroupDelete(service, &cobra.Command{}, []string{"hg-abcdef12"})
		})
	})
}
