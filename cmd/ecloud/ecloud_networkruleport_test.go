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

func Test_ecloudNetworkRulePortList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworkRulePorts(gomock.Any()).Return([]ecloud.NetworkRulePort{}, nil).Times(1)

		ecloudNetworkRulePortList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudNetworkRulePortList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetNetworkRulePortsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworkRulePorts(gomock.Any()).Return([]ecloud.NetworkRulePort{}, errors.New("test error")).Times(1)

		err := ecloudNetworkRulePortList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving network rule ports: test error", err.Error())
	})
}

func Test_ecloudNetworkRulePortShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudNetworkRulePortShowCmd(nil).Args(nil, []string{"fwrp-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudNetworkRulePortShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing network rule port", err.Error())
	})
}

func Test_ecloudNetworkRulePortShow(t *testing.T) {
	t.Run("SingleNetworkRulePort", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworkRulePort("fwrp-abcdef12").Return(ecloud.NetworkRulePort{}, nil).Times(1)

		ecloudNetworkRulePortShow(service, &cobra.Command{}, []string{"fwrp-abcdef12"})
	})

	t.Run("MultipleNetworkRulePorts", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetNetworkRulePort("fwrp-abcdef12").Return(ecloud.NetworkRulePort{}, nil),
			service.EXPECT().GetNetworkRulePort("fwrp-abcdef23").Return(ecloud.NetworkRulePort{}, nil),
		)

		ecloudNetworkRulePortShow(service, &cobra.Command{}, []string{"fwrp-abcdef12", "fwrp-abcdef23"})
	})

	t.Run("GetNetworkRulePortError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworkRulePort("fwrp-abcdef12").Return(ecloud.NetworkRulePort{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving network rule port [fwrp-abcdef12]: test error\n", func() {
			ecloudNetworkRulePortShow(service, &cobra.Command{}, []string{"fwrp-abcdef12"})
		})
	})
}

func Test_ecloudNetworkRulePortCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNetworkRulePortCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testport", "--protocol=TCP"})

		req := ecloud.CreateNetworkRulePortRequest{
			Name:     "testport",
			Protocol: "TCP",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fwrp-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateNetworkRulePort(req).Return(resp, nil),
			service.EXPECT().GetNetworkRulePort("fwrp-abcdef12").Return(ecloud.NetworkRulePort{}, nil),
		)

		ecloudNetworkRulePortCreate(service, cmd, []string{})
	})

	t.Run("CreateWithWaitFlagNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNetworkRulePortCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testport", "--protocol=TCP", "--wait"})

		req := ecloud.CreateNetworkRulePortRequest{
			Name:     "testport",
			Protocol: "TCP",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fwrp-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateNetworkRulePort(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetNetworkRulePort("fwrp-abcdef12").Return(ecloud.NetworkRulePort{}, nil),
		)

		ecloudNetworkRulePortCreate(service, cmd, []string{})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNetworkRulePortCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testport", "--protocol=TCP", "--wait"})

		req := ecloud.CreateNetworkRulePortRequest{
			Name:     "testport",
			Protocol: "TCP",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fwrp-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateNetworkRulePort(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		err := ecloudNetworkRulePortCreate(service, cmd, []string{})
		assert.Equal(t, "Error waiting for network rule port task to complete: Error waiting for command: Failed to retrieve task status: test error", err.Error())
	})

	t.Run("CreateNetworkRulePortError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNetworkRulePortCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testport", "--protocol=TCP"})

		service.EXPECT().CreateNetworkRulePort(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error")).Times(1)

		err := ecloudNetworkRulePortCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating network rule port: test error", err.Error())
	})

	t.Run("GetNetworkRulePortError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNetworkRulePortCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testport", "--protocol=TCP"})

		gomock.InOrder(
			service.EXPECT().CreateNetworkRulePort(gomock.Any()).Return(ecloud.TaskReference{ResourceID: "fwrp-abcdef12"}, nil),
			service.EXPECT().GetNetworkRulePort("fwrp-abcdef12").Return(ecloud.NetworkRulePort{}, errors.New("test error")),
		)

		err := ecloudNetworkRulePortCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new network rule port: test error", err.Error())
	})
}

func Test_ecloudNetworkRulePortUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudNetworkRulePortUpdateCmd(nil).Args(nil, []string{"fwrp-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudNetworkRulePortUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing network rule port", err.Error())
	})
}

func Test_ecloudNetworkRulePortUpdate(t *testing.T) {
	t.Run("SingleRule", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNetworkRulePortCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testport"})

		req := ecloud.PatchNetworkRulePortRequest{
			Name: "testport",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fwrp-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchNetworkRulePort("fwrp-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetNetworkRulePort("fwrp-abcdef12").Return(ecloud.NetworkRulePort{}, nil),
		)

		ecloudNetworkRulePortUpdate(service, cmd, []string{"fwrp-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNetworkRulePortUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testport", "--wait"})

		req := ecloud.PatchNetworkRulePortRequest{
			Name: "testport",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fwrp-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchNetworkRulePort("fwrp-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetNetworkRulePort("fwrp-abcdef12").Return(ecloud.NetworkRulePort{}, nil),
		)

		ecloudNetworkRulePortUpdate(service, cmd, []string{"fwrp-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNetworkRulePortUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testrule", "--wait"})

		req := ecloud.PatchNetworkRulePortRequest{
			Name: "testrule",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fwrp-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchNetworkRulePort("fwrp-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for network rule port [fwrp-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudNetworkRulePortUpdate(service, cmd, []string{"fwrp-abcdef12"})
		})
	})

	t.Run("MultipleNetworkRulePorts", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp1 := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fwrp-abcdef12",
		}

		resp2 := ecloud.TaskReference{
			TaskID:     "task-abcdef23",
			ResourceID: "fwrp-12abcdef",
		}

		gomock.InOrder(
			service.EXPECT().PatchNetworkRulePort("fwrp-abcdef12", gomock.Any()).Return(resp1, nil),
			service.EXPECT().GetNetworkRulePort("fwrp-abcdef12").Return(ecloud.NetworkRulePort{}, nil),
			service.EXPECT().PatchNetworkRulePort("fwrp-12abcdef", gomock.Any()).Return(resp2, nil),
			service.EXPECT().GetNetworkRulePort("fwrp-12abcdef").Return(ecloud.NetworkRulePort{}, nil),
		)

		ecloudNetworkRulePortUpdate(service, &cobra.Command{}, []string{"fwrp-abcdef12", "fwrp-12abcdef"})
	})

	t.Run("PatchNetworkRulePortError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchNetworkRulePort("fwrp-abcdef12", gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating network rule port [fwrp-abcdef12]: test error\n", func() {
			ecloudNetworkRulePortUpdate(service, &cobra.Command{}, []string{"fwrp-abcdef12"})
		})
	})

	t.Run("GetNetworkRulePortError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fwrp-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchNetworkRulePort("fwrp-abcdef12", gomock.Any()).Return(resp, nil),
			service.EXPECT().GetNetworkRulePort("fwrp-abcdef12").Return(ecloud.NetworkRulePort{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated network rule port [fwrp-abcdef12]: test error\n", func() {
			ecloudNetworkRulePortUpdate(service, &cobra.Command{}, []string{"fwrp-abcdef12"})
		})
	})
}

func Test_ecloudNetworkRulePortDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudNetworkRulePortDeleteCmd(nil).Args(nil, []string{"fwrp-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudNetworkRulePortDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing network rule port", err.Error())
	})
}

func Test_ecloudNetworkRulePortDelete(t *testing.T) {
	t.Run("SingleRule", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteNetworkRulePort("fwrp-abcdef12").Return("task-abcdef12", nil).Times(1)

		ecloudNetworkRulePortDelete(service, &cobra.Command{}, []string{"fwrp-abcdef12"})
	})

	t.Run("MultipleNetworkRulePorts", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteNetworkRulePort("fwrp-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().DeleteNetworkRulePort("fwrp-12abcdef").Return("task-abcdef23", nil),
		)

		ecloudNetworkRulePortDelete(service, &cobra.Command{}, []string{"fwrp-abcdef12", "fwrp-12abcdef"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudNetworkRulePortDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteNetworkRulePort("fwrp-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudNetworkRulePortDelete(service, cmd, []string{"fwrp-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNetworkRulePortDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		gomock.InOrder(
			service.EXPECT().DeleteNetworkRulePort("fwrp-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for network rule port [fwrp-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudNetworkRulePortDelete(service, cmd, []string{"fwrp-abcdef12"})
		})
	})

	t.Run("DeleteNetworkRulePortError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteNetworkRulePort("fwrp-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing network rule port [fwrp-abcdef12]: test error\n", func() {
			ecloudNetworkRulePortDelete(service, &cobra.Command{}, []string{"fwrp-abcdef12"})
		})
	})
}
