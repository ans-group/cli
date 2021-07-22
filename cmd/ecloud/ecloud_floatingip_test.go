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

func Test_ecloudFloatingIPList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFloatingIPs(gomock.Any()).Return([]ecloud.FloatingIP{}, nil).Times(1)

		ecloudFloatingIPList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudFloatingIPList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetFloatingIPsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFloatingIPs(gomock.Any()).Return([]ecloud.FloatingIP{}, errors.New("test error")).Times(1)

		err := ecloudFloatingIPList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving floating IPs: test error", err.Error())
	})
}

func Test_ecloudFloatingIPShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFloatingIPShowCmd(nil).Args(nil, []string{"fip-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFloatingIPShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing floating IP", err.Error())
	})
}

func Test_ecloudFloatingIPShow(t *testing.T) {
	t.Run("SingleFloatingIP", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFloatingIP("fip-abcdef12").Return(ecloud.FloatingIP{}, nil).Times(1)

		ecloudFloatingIPShow(service, &cobra.Command{}, []string{"fip-abcdef12"})
	})

	t.Run("MultipleFloatingIPs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetFloatingIP("fip-abcdef12").Return(ecloud.FloatingIP{}, nil),
			service.EXPECT().GetFloatingIP("fip-abcdef23").Return(ecloud.FloatingIP{}, nil),
		)

		ecloudFloatingIPShow(service, &cobra.Command{}, []string{"fip-abcdef12", "fip-abcdef23"})
	})

	t.Run("GetFloatingIPError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFloatingIP("fip-abcdef12").Return(ecloud.FloatingIP{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving floating IP [fip-abcdef12]: test error\n", func() {
			ecloudFloatingIPShow(service, &cobra.Command{}, []string{"fip-abcdef12"})
		})
	})
}

func Test_ecloudFloatingIPCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFloatingIPCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testfip"})

		req := ecloud.CreateFloatingIPRequest{
			Name: "testfip",
		}

		resp := ecloud.TaskReference{
			TaskID: "task-abcdef12",
			ResourceID: "fip-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateFloatingIP(req).Return(resp, nil),
			service.EXPECT().GetFloatingIP("fip-abcdef12").Return(ecloud.FloatingIP{}, nil),
		)

		ecloudFloatingIPCreate(service, cmd, []string{})
	})

	t.Run("CreateWithWaitFlagNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFloatingIPCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testfip", "--wait"})

		req := ecloud.CreateFloatingIPRequest{
			Name:  "testfip",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fip-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateFloatingIP(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetFloatingIP("fip-abcdef12").Return(ecloud.FloatingIP{}, nil),
		)

		ecloudFloatingIPCreate(service, cmd, []string{})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFloatingIPCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testfip", "--wait"})

		req := ecloud.CreateFloatingIPRequest{
			Name:  "testfip",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fip-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateFloatingIP(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		err := ecloudFloatingIPCreate(service, cmd, []string{})
		assert.Equal(t, "Error waiting for floating IP task to complete: Error waiting for command: Failed to retrieve task status: test error", err.Error())
	})

	t.Run("CreateFloatingIPError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFloatingIPCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testfip"})

		service.EXPECT().CreateFloatingIP(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error")).Times(1)

		err := ecloudFloatingIPCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating floating IP: test error", err.Error())
	})

	t.Run("GetFloatingIPError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFloatingIPCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testfip"})

		gomock.InOrder(
			service.EXPECT().CreateFloatingIP(gomock.Any()).Return(ecloud.TaskReference{ResourceID: "fip-abcdef12"}, nil),
			service.EXPECT().GetFloatingIP("fip-abcdef12").Return(ecloud.FloatingIP{}, errors.New("test error")),
		)

		err := ecloudFloatingIPCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new floating IP: test error", err.Error())
	})
}

func Test_ecloudFloatingIPUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFloatingIPUpdateCmd(nil).Args(nil, []string{"fip-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFloatingIPUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing floating IP", err.Error())
	})
}

func Test_ecloudFloatingIPUpdate(t *testing.T) {
	t.Run("SingleFloatingIP", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudFloatingIPUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testfip"})

		req := ecloud.PatchFloatingIPRequest{
			Name: "testfip",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fip-abcef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchFloatingIP("fip-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetFloatingIP("fip-abcdef12").Return(ecloud.FloatingIP{}, nil),
		)

		ecloudFloatingIPUpdate(service, cmd, []string{"fip-abcdef12"})
	})

	t.Run("MultipleFloatingIPs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp1 := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fip-abcdef12",
		}

		resp2 := ecloud.TaskReference{
			TaskID:     "task-abcdef13",
			ResourceID: "fip-12abcdef",
		}

		gomock.InOrder(
			service.EXPECT().PatchFloatingIP("fip-abcdef12", gomock.Any()).Return(resp1, nil),
			service.EXPECT().GetFloatingIP("fip-abcdef12").Return(ecloud.FloatingIP{}, nil),
			service.EXPECT().PatchFloatingIP("fip-12abcdef", gomock.Any()).Return(resp2, nil),
			service.EXPECT().GetFloatingIP("fip-12abcdef").Return(ecloud.FloatingIP{}, nil),
		)

		ecloudFloatingIPUpdate(service, &cobra.Command{}, []string{"fip-abcdef12", "fip-12abcdef"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFloatingIPCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testfip", "--wait"})

		req := ecloud.PatchFloatingIPRequest{
			Name:  "testfip",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fip-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchFloatingIP("fip-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetFloatingIP("fip-abcdef12").Return(ecloud.FloatingIP{}, nil),
		)

		ecloudFloatingIPUpdate(service, cmd, []string{"fip-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFloatingIPCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testfip", "--wait"})

		req := ecloud.PatchFloatingIPRequest{
			Name:  "testfip",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fip-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchFloatingIP("fip-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for floating IP [fip-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudFloatingIPUpdate(service, cmd, []string{"fip-abcdef12"})
		})
	})

	t.Run("PatchFloatingIPError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchFloatingIP("fip-abcdef12", gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating floating IP [fip-abcdef12]: test error\n", func() {
			ecloudFloatingIPUpdate(service, &cobra.Command{}, []string{"fip-abcdef12"})
		})
	})

	t.Run("GetFloatingIPError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fip-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchFloatingIP("fip-abcdef12", gomock.Any()).Return(resp, nil),
			service.EXPECT().GetFloatingIP("fip-abcdef12").Return(ecloud.FloatingIP{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated floating IP [fip-abcdef12]: test error\n", func() {
			ecloudFloatingIPUpdate(service, &cobra.Command{}, []string{"fip-abcdef12"})
		})
	})
}

func Test_ecloudFloatingIPDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFloatingIPDeleteCmd(nil).Args(nil, []string{"fip-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFloatingIPDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing floating IP", err.Error())
	})
}

func Test_ecloudFloatingIPDelete(t *testing.T) {
	t.Run("SingleFloatingIP", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteFloatingIP("fip-abcdef12").Return("task-abcdef12", nil).Times(1)

		ecloudFloatingIPDelete(service, &cobra.Command{}, []string{"fip-abcdef12"})
	})

	t.Run("MultipleFloatingIPs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteFloatingIP("fip-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().DeleteFloatingIP("fip-12abcdef").Return("task-abcdef23", nil),
		)

		ecloudFloatingIPDelete(service, &cobra.Command{}, []string{"fip-abcdef12", "fip-12abcdef"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudFloatingIPDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteFloatingIP("fip-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudFloatingIPDelete(service, cmd, []string{"fip-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudFloatingIPDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		gomock.InOrder(
			service.EXPECT().DeleteFloatingIP("fip-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for floating IP [fip-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudFloatingIPDelete(service, cmd, []string{"fip-abcdef12"})
		})
	})

	t.Run("DeleteFloatingIPError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteFloatingIP("fip-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing floating IP [fip-abcdef12]: test error\n", func() {
			ecloudFloatingIPDelete(service, &cobra.Command{}, []string{"fip-abcdef12"})
		})
	})
}

func Test_ecloudFloatingIPAssignCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFloatingIPAssignCmd(nil).Args(nil, []string{"fip-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFloatingIPAssignCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing floating IP", err.Error())
	})
}

func Test_ecloudFloatingIPAssign(t *testing.T) {
	t.Run("AssignFloatingIP_NoError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		req := ecloud.AssignFloatingIPRequest{
			ResourceID: "i-abcdef12",
		}

		cmd := ecloudFloatingIPAssignCmd(nil)
		cmd.ParseFlags([]string{"--resource=i-abcdef12"})

		service := mocks.NewMockECloudService(mockCtrl)
		service.EXPECT().AssignFloatingIP("fip-abcdef12", gomock.Eq(req)).Return("task-abcdef12", nil)
		service.EXPECT().GetFloatingIP("fip-abcdef12").Return(ecloud.FloatingIP{}, nil)

		err := ecloudFloatingIPAssign(service, cmd, []string{"fip-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("AssignFloatingIPError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().AssignFloatingIP("fip-abcdef12", gomock.Any()).Return("", errors.New("test error"))

		err := ecloudFloatingIPAssign(service, &cobra.Command{}, []string{"fip-abcdef12"})

		assert.NotNil(t, err)
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudFloatingIPAssignCmd(nil)
		cmd.ParseFlags([]string{"--resource=i-abcdef12","--wait"})

		service := mocks.NewMockECloudService(mockCtrl)

		req := ecloud.AssignFloatingIPRequest{
			ResourceID: "i-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().AssignFloatingIP("fip-abcdef12", gomock.Eq(req)).Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetFloatingIP("fip-abcdef12").Return(ecloud.FloatingIP{}, nil),
		)

		ecloudFloatingIPAssign(service, cmd, []string{"fip-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudFloatingIPAssignCmd(nil)
		cmd.ParseFlags([]string{"--resource=i-abcdef12","--wait"})

		gomock.InOrder(
			service.EXPECT().AssignFloatingIP("fip-abcdef12", gomock.Any()).Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{}, errors.New("test error")),
		)

		err := ecloudFloatingIPAssign(service, cmd, []string{"fip-abcdef12"})
		assert.Equal(t, "Error waiting for floating IP [fip-abcdef12] task: Error waiting for command: Failed to retrieve task status: test error", err.Error())
	})
}

func Test_ecloudFloatingIPUnassignCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFloatingIPUnassignCmd(nil).Args(nil, []string{"fip-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFloatingIPUnassignCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing floating IP", err.Error())
	})
}

func Test_ecloudFloatingIPUnassign(t *testing.T) {
	t.Run("SingleFloatingIP", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().UnassignFloatingIP("fip-abcdef12").Return("task-abcdef12", nil).Times(1)

		ecloudFloatingIPUnassign(service, &cobra.Command{}, []string{"fip-abcdef12"})
	})

	t.Run("MultipleFloatingIPs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().UnassignFloatingIP("fip-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().UnassignFloatingIP("fip-12abcdef").Return("task-abcdef23", nil),
		)

		ecloudFloatingIPUnassign(service, &cobra.Command{}, []string{"fip-abcdef12", "fip-12abcdef"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudFloatingIPUnassignCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().UnassignFloatingIP("fip-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudFloatingIPUnassign(service, cmd, []string{"fip-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudFloatingIPUnassignCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		gomock.InOrder(
			service.EXPECT().UnassignFloatingIP("fip-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{}, errors.New("test error")),
		)

		// err := ecloudFloatingIPUnassign(service, cmd, []string{"fip-abcdef12"})
		// assert.Equal(t, "Error waiting for floating IP task: Error waiting for command: Failed to retrieve task status: test error", err)

		test_output.AssertErrorOutput(t, "Error waiting for floating IP [fip-abcdef12] task: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudFloatingIPUnassign(service, cmd, []string{"fip-abcdef12"})
		})
	})

	t.Run("UnassignFloatingIPError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().UnassignFloatingIP("fip-abcdef12").Return("", errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error unassigning floating IP [fip-abcdef12]: test error\n", func() {
			ecloudFloatingIPUnassign(service, &cobra.Command{}, []string{"fip-abcdef12"})
		})
	})
}
