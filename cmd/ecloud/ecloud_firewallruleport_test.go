package ecloud

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ecloudFirewallRulePortList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewallRulePorts(gomock.Any()).Return([]ecloud.FirewallRulePort{}, nil).Times(1)

		ecloudFirewallRulePortList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudFirewallRulePortList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetFirewallRulePortsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewallRulePorts(gomock.Any()).Return([]ecloud.FirewallRulePort{}, errors.New("test error")).Times(1)

		err := ecloudFirewallRulePortList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving firewall rule ports: test error", err.Error())
	})
}

func Test_ecloudFirewallRulePortShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFirewallRulePortShowCmd(nil).Args(nil, []string{"fwrp-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFirewallRulePortShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing firewall rule port", err.Error())
	})
}

func Test_ecloudFirewallRulePortShow(t *testing.T) {
	t.Run("SingleFirewallRulePort", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewallRulePort("fwrp-abcdef12").Return(ecloud.FirewallRulePort{}, nil).Times(1)

		ecloudFirewallRulePortShow(service, &cobra.Command{}, []string{"fwrp-abcdef12"})
	})

	t.Run("MultipleFirewallRulePorts", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetFirewallRulePort("fwrp-abcdef12").Return(ecloud.FirewallRulePort{}, nil),
			service.EXPECT().GetFirewallRulePort("fwrp-abcdef23").Return(ecloud.FirewallRulePort{}, nil),
		)

		ecloudFirewallRulePortShow(service, &cobra.Command{}, []string{"fwrp-abcdef12", "fwrp-abcdef23"})
	})

	t.Run("GetFirewallRulePortError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewallRulePort("fwrp-abcdef12").Return(ecloud.FirewallRulePort{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving firewall rule port [fwrp-abcdef12]: test error\n", func() {
			ecloudFirewallRulePortShow(service, &cobra.Command{}, []string{"fwrp-abcdef12"})
		})
	})
}

func Test_ecloudFirewallRulePortCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFirewallRulePortCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testport", "--protocol=TCP"})

		req := ecloud.CreateFirewallRulePortRequest{
			Name:     "testport",
			Protocol: "TCP",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fwrp-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateFirewallRulePort(req).Return(resp, nil),
			service.EXPECT().GetFirewallRulePort("fwrp-abcdef12").Return(ecloud.FirewallRulePort{}, nil),
		)

		ecloudFirewallRulePortCreate(service, cmd, []string{})
	})

	t.Run("CreateWithWaitFlagNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFirewallRulePortCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testport", "--protocol=TCP", "--wait"})

		req := ecloud.CreateFirewallRulePortRequest{
			Name:     "testport",
			Protocol: "TCP",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fwrp-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateFirewallRulePort(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetFirewallRulePort("fwrp-abcdef12").Return(ecloud.FirewallRulePort{}, nil),
		)

		ecloudFirewallRulePortCreate(service, cmd, []string{})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFirewallRulePortCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testport", "--protocol=TCP", "--wait"})

		req := ecloud.CreateFirewallRulePortRequest{
			Name:     "testport",
			Protocol: "TCP",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fwrp-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateFirewallRulePort(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		err := ecloudFirewallRulePortCreate(service, cmd, []string{})
		assert.Equal(t, "Error waiting for firewall rule port task to complete: Error waiting for command: Failed to retrieve task status: test error", err.Error())
	})

	t.Run("CreateFirewallRulePortError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFirewallRulePortCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testport", "--protocol=TCP"})

		service.EXPECT().CreateFirewallRulePort(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error")).Times(1)

		err := ecloudFirewallRulePortCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating firewall rule port: test error", err.Error())
	})

	t.Run("GetFirewallRulePortError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFirewallRulePortCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testport", "--protocol=TCP"})

		gomock.InOrder(
			service.EXPECT().CreateFirewallRulePort(gomock.Any()).Return(ecloud.TaskReference{ResourceID: "fwrp-abcdef12"}, nil),
			service.EXPECT().GetFirewallRulePort("fwrp-abcdef12").Return(ecloud.FirewallRulePort{}, errors.New("test error")),
		)

		err := ecloudFirewallRulePortCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new firewall rule port: test error", err.Error())
	})
}

func Test_ecloudFirewallRulePortUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFirewallRulePortUpdateCmd(nil).Args(nil, []string{"fwrp-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFirewallRulePortUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing firewall rule port", err.Error())
	})
}

func Test_ecloudFirewallRulePortUpdate(t *testing.T) {
	t.Run("SingleRule", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudFirewallRulePortCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testport"})

		req := ecloud.PatchFirewallRulePortRequest{
			Name: "testport",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fwrp-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchFirewallRulePort("fwrp-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetFirewallRulePort("fwrp-abcdef12").Return(ecloud.FirewallRulePort{}, nil),
		)

		ecloudFirewallRulePortUpdate(service, cmd, []string{"fwrp-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudFirewallRulePortUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testport", "--wait"})

		req := ecloud.PatchFirewallRulePortRequest{
			Name: "testport",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fwrp-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchFirewallRulePort("fwrp-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetFirewallRulePort("fwrp-abcdef12").Return(ecloud.FirewallRulePort{}, nil),
		)

		ecloudFirewallRulePortUpdate(service, cmd, []string{"fwrp-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudFirewallRulePortUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testrule", "--wait"})

		req := ecloud.PatchFirewallRulePortRequest{
			Name: "testrule",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fwrp-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchFirewallRulePort("fwrp-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for firewall rule port [fwrp-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudFirewallRulePortUpdate(service, cmd, []string{"fwrp-abcdef12"})
		})
	})

	t.Run("MultipleFirewallRulePorts", func(t *testing.T) {
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
			service.EXPECT().PatchFirewallRulePort("fwrp-abcdef12", gomock.Any()).Return(resp1, nil),
			service.EXPECT().GetFirewallRulePort("fwrp-abcdef12").Return(ecloud.FirewallRulePort{}, nil),
			service.EXPECT().PatchFirewallRulePort("fwrp-12abcdef", gomock.Any()).Return(resp2, nil),
			service.EXPECT().GetFirewallRulePort("fwrp-12abcdef").Return(ecloud.FirewallRulePort{}, nil),
		)

		ecloudFirewallRulePortUpdate(service, &cobra.Command{}, []string{"fwrp-abcdef12", "fwrp-12abcdef"})
	})

	t.Run("PatchFirewallRulePortError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchFirewallRulePort("fwrp-abcdef12", gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating firewall rule port [fwrp-abcdef12]: test error\n", func() {
			ecloudFirewallRulePortUpdate(service, &cobra.Command{}, []string{"fwrp-abcdef12"})
		})
	})

	t.Run("GetFirewallRulePortError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fwrp-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchFirewallRulePort("fwrp-abcdef12", gomock.Any()).Return(resp, nil),
			service.EXPECT().GetFirewallRulePort("fwrp-abcdef12").Return(ecloud.FirewallRulePort{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated firewall rule port [fwrp-abcdef12]: test error\n", func() {
			ecloudFirewallRulePortUpdate(service, &cobra.Command{}, []string{"fwrp-abcdef12"})
		})
	})
}

func Test_ecloudFirewallRulePortDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFirewallRulePortDeleteCmd(nil).Args(nil, []string{"fwrp-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFirewallRulePortDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing firewall rule port", err.Error())
	})
}

func Test_ecloudFirewallRulePortDelete(t *testing.T) {
	t.Run("SingleRule", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteFirewallRulePort("fwrp-abcdef12").Return("task-abcdef12", nil).Times(1)

		ecloudFirewallRulePortDelete(service, &cobra.Command{}, []string{"fwrp-abcdef12"})
	})

	t.Run("MultipleFirewallRulePorts", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteFirewallRulePort("fwrp-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().DeleteFirewallRulePort("fwrp-12abcdef").Return("task-abcdef23", nil),
		)

		ecloudFirewallRulePortDelete(service, &cobra.Command{}, []string{"fwrp-abcdef12", "fwrp-12abcdef"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudFirewallRulePortDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteFirewallRulePort("fwrp-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudFirewallRulePortDelete(service, cmd, []string{"fwrp-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudFirewallRulePortDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		gomock.InOrder(
			service.EXPECT().DeleteFirewallRulePort("fwrp-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for firewall rule port [fwrp-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudFirewallRulePortDelete(service, cmd, []string{"fwrp-abcdef12"})
		})
	})

	t.Run("DeleteFirewallRulePortError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteFirewallRulePort("fwrp-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing firewall rule port [fwrp-abcdef12]: test error\n", func() {
			ecloudFirewallRulePortDelete(service, &cobra.Command{}, []string{"fwrp-abcdef12"})
		})
	})
}
