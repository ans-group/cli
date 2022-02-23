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

func Test_ecloudVIPList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVIPs(gomock.Any()).Return([]ecloud.VIP{}, nil).Times(1)

		ecloudVIPList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudVIPList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetVIPsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVIPs(gomock.Any()).Return([]ecloud.VIP{}, errors.New("test error")).Times(1)

		err := ecloudVIPList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving VIPs: test error", err.Error())
	})
}

func Test_ecloudVIPShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVIPShowCmd(nil).Args(nil, []string{"vip-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVIPShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing VIP", err.Error())
	})
}

func Test_ecloudVIPShow(t *testing.T) {
	t.Run("SingleVIP", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVIP("vip-abcdef12").Return(ecloud.VIP{}, nil).Times(1)

		ecloudVIPShow(service, &cobra.Command{}, []string{"vip-abcdef12"})
	})

	t.Run("MultipleVIPs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetVIP("vip-abcdef12").Return(ecloud.VIP{}, nil),
			service.EXPECT().GetVIP("vip-abcdef23").Return(ecloud.VIP{}, nil),
		)

		ecloudVIPShow(service, &cobra.Command{}, []string{"vip-abcdef12", "vip-abcdef23"})
	})

	t.Run("GetVIPError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVIP("vip-abcdef12").Return(ecloud.VIP{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving VIP [vip-abcdef12]: test error\n", func() {
			ecloudVIPShow(service, &cobra.Command{}, []string{"vip-abcdef12"})
		})
	})
}

func Test_ecloudVIPCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVIPCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvip"})

		req := ecloud.CreateVIPRequest{
			Name: "testvip",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vip-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateVIP(req).Return(resp, nil),
			service.EXPECT().GetVIP("vip-abcdef12").Return(ecloud.VIP{}, nil),
		)

		ecloudVIPCreate(service, cmd, []string{})
	})

	t.Run("CreateWithWaitFlagNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVIPCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvip", "--wait"})

		req := ecloud.CreateVIPRequest{
			Name: "testvip",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vip-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateVIP(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetVIP("vip-abcdef12").Return(ecloud.VIP{}, nil),
		)

		ecloudVIPCreate(service, cmd, []string{})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVIPCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvip", "--wait"})

		req := ecloud.CreateVIPRequest{
			Name: "testvip",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vip-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateVIP(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		err := ecloudVIPCreate(service, cmd, []string{})
		assert.Equal(t, "Error waiting for VIP task to complete: Error waiting for command: Failed to retrieve task status: test error", err.Error())
	})

	t.Run("CreateVIPError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVIPCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvip"})

		service.EXPECT().CreateVIP(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error")).Times(1)

		err := ecloudVIPCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating VIP: test error", err.Error())
	})

	t.Run("GetVIPError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVIPCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvip"})

		gomock.InOrder(
			service.EXPECT().CreateVIP(gomock.Any()).Return(ecloud.TaskReference{ResourceID: "vip-abcdef12"}, nil),
			service.EXPECT().GetVIP("vip-abcdef12").Return(ecloud.VIP{}, errors.New("test error")),
		)

		err := ecloudVIPCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new VIP: test error", err.Error())
	})
}

func Test_ecloudVIPUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVIPUpdateCmd(nil).Args(nil, []string{"vip-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVIPUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing VIP", err.Error())
	})
}

func Test_ecloudVIPUpdate(t *testing.T) {
	t.Run("SingleRule", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVIPUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvip"})

		req := ecloud.PatchVIPRequest{
			Name: "testvip",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vip-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchVIP("vip-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetVIP("vip-abcdef12").Return(ecloud.VIP{}, nil),
		)

		ecloudVIPUpdate(service, cmd, []string{"vip-abcdef12"})
	})

	t.Run("MultipleVIPs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp1 := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vip-abcdef12",
		}

		resp2 := ecloud.TaskReference{
			TaskID:     "task-abcdef23",
			ResourceID: "vip-12abcdef",
		}

		gomock.InOrder(
			service.EXPECT().PatchVIP("vip-abcdef12", gomock.Any()).Return(resp1, nil),
			service.EXPECT().GetVIP("vip-abcdef12").Return(ecloud.VIP{}, nil),
			service.EXPECT().PatchVIP("vip-12abcdef", gomock.Any()).Return(resp2, nil),
			service.EXPECT().GetVIP("vip-12abcdef").Return(ecloud.VIP{}, nil),
		)

		ecloudVIPUpdate(service, &cobra.Command{}, []string{"vip-abcdef12", "vip-12abcdef"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVIPUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvip", "--wait"})

		req := ecloud.PatchVIPRequest{
			Name: "testvip",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vip-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchVIP("vip-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetVIP("vip-abcdef12").Return(ecloud.VIP{}, nil),
		)

		ecloudVIPUpdate(service, cmd, []string{"vip-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVIPUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testvip", "--wait"})

		req := ecloud.PatchVIPRequest{
			Name: "testvip",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vip-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchVIP("vip-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for VIP [vip-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudVIPUpdate(service, cmd, []string{"vip-abcdef12"})
		})
	})

	t.Run("PatchVIPError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchVIP("vip-abcdef12", gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating VIP [vip-abcdef12]: test error\n", func() {
			ecloudVIPUpdate(service, &cobra.Command{}, []string{"vip-abcdef12"})
		})
	})

	t.Run("GetVIPError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vip-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchVIP("vip-abcdef12", gomock.Any()).Return(resp, nil),
			service.EXPECT().GetVIP("vip-abcdef12").Return(ecloud.VIP{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated VIP [vip-abcdef12]: test error\n", func() {
			ecloudVIPUpdate(service, &cobra.Command{}, []string{"vip-abcdef12"})
		})
	})
}

func Test_ecloudVIPDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVIPDeleteCmd(nil).Args(nil, []string{"vip-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVIPDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing VIP", err.Error())
	})
}

func Test_ecloudVIPDelete(t *testing.T) {
	t.Run("SingleRule", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteVIP("vip-abcdef12").Return("task-abcdef12", nil).Times(1)

		ecloudVIPDelete(service, &cobra.Command{}, []string{"vip-abcdef12"})
	})

	t.Run("MultipleVIPs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteVIP("vip-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().DeleteVIP("vip-12abcdef").Return("task-abcdef23", nil),
		)

		ecloudVIPDelete(service, &cobra.Command{}, []string{"vip-abcdef12", "vip-12abcdef"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudVIPDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteVIP("vip-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudVIPDelete(service, cmd, []string{"vip-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVIPDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		gomock.InOrder(
			service.EXPECT().DeleteVIP("vip-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for VIP [vip-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudVIPDelete(service, cmd, []string{"vip-abcdef12"})
		})
	})

	t.Run("DeleteVIPError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteVIP("vip-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing VIP [vip-abcdef12]: test error\n", func() {
			ecloudVIPDelete(service, &cobra.Command{}, []string{"vip-abcdef12"})
		})
	})
}
