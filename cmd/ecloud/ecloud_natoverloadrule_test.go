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

func Test_ecloudNATOverloadRuleList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNATOverloadRules(gomock.Any()).Return([]ecloud.NATOverloadRule{}, nil).Times(1)

		ecloudNATOverloadRuleList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudNATOverloadRuleList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetNATOverloadRulesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNATOverloadRules(gomock.Any()).Return([]ecloud.NATOverloadRule{}, errors.New("test error")).Times(1)

		err := ecloudNATOverloadRuleList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "error retrieving NAT overload rules: test error", err.Error())
	})
}

func Test_ecloudNATOverloadRuleShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudNATOverloadRuleShowCmd(nil).Args(nil, []string{"nor-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudNATOverloadRuleShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing NAT overload rule", err.Error())
	})
}

func Test_ecloudNATOverloadRuleShow(t *testing.T) {
	t.Run("SingleNATOverloadRule", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNATOverloadRule("nor-abcdef12").Return(ecloud.NATOverloadRule{}, nil).Times(1)

		ecloudNATOverloadRuleShow(service, &cobra.Command{}, []string{"nor-abcdef12"})
	})

	t.Run("MultipleNATOverloadRules", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetNATOverloadRule("nor-abcdef12").Return(ecloud.NATOverloadRule{}, nil),
			service.EXPECT().GetNATOverloadRule("nor-abcdef23").Return(ecloud.NATOverloadRule{}, nil),
		)

		ecloudNATOverloadRuleShow(service, &cobra.Command{}, []string{"nor-abcdef12", "nor-abcdef23"})
	})

	t.Run("GetNATOverloadRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNATOverloadRule("nor-abcdef12").Return(ecloud.NATOverloadRule{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving NAT overload rule [nor-abcdef12]: test error\n", func() {
			ecloudNATOverloadRuleShow(service, &cobra.Command{}, []string{"nor-abcdef12"})
		})
	})
}

func Test_ecloudNATOverloadRuleCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNATOverloadRuleCreateCmd(nil)
		cmd.Flags().Set("name", "testrule")
		cmd.Flags().Set("network", "net-abcdef12")
		cmd.Flags().Set("subnet", "10.0.0.0/24")
		cmd.Flags().Set("floating-ip", "fip-abcdef12")
		cmd.Flags().Set("action", "ALLOW")

		req := ecloud.CreateNATOverloadRuleRequest{
			Name:         "testrule",
			NetworkID:    "net-abcdef12",
			Subnet:       "10.0.0.0/24",
			FloatingIPID: "fip-abcdef12",
			Action:       ecloud.NATOverloadRuleActionAllow,
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "nor-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateNATOverloadRule(req).Return(resp, nil),
			service.EXPECT().GetNATOverloadRule("nor-abcdef12").Return(ecloud.NATOverloadRule{}, nil),
		)

		ecloudNATOverloadRuleCreate(service, cmd, []string{})
	})

	t.Run("CreateWithWaitFlagNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNATOverloadRuleCreateCmd(nil)
		cmd.Flags().Set("name", "testrule")
		cmd.Flags().Set("network", "net-abcdef12")
		cmd.Flags().Set("subnet", "10.0.0.0/24")
		cmd.Flags().Set("floating-ip", "fip-abcdef12")
		cmd.Flags().Set("action", "ALLOW")
		cmd.Flags().Set("wait", "true")

		req := ecloud.CreateNATOverloadRuleRequest{
			Name:         "testrule",
			NetworkID:    "net-abcdef12",
			Subnet:       "10.0.0.0/24",
			FloatingIPID: "fip-abcdef12",
			Action:       ecloud.NATOverloadRuleActionAllow,
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "nor-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateNATOverloadRule(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetNATOverloadRule("nor-abcdef12").Return(ecloud.NATOverloadRule{}, nil),
		)

		ecloudNATOverloadRuleCreate(service, cmd, []string{})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNATOverloadRuleCreateCmd(nil)
		cmd.Flags().Set("name", "testrule")
		cmd.Flags().Set("network", "net-abcdef12")
		cmd.Flags().Set("subnet", "10.0.0.0/24")
		cmd.Flags().Set("floating-ip", "fip-abcdef12")
		cmd.Flags().Set("action", "ALLOW")
		cmd.Flags().Set("wait", "true")

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "nor-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateNATOverloadRule(gomock.Any()).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		err := ecloudNATOverloadRuleCreate(service, cmd, []string{})
		assert.Equal(t, "error waiting for NAT overload rule task to complete: error waiting for command: failed to retrieve task status: test error", err.Error())
	})

	t.Run("CreateNATOverloadRuleError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNATOverloadRuleCreateCmd(nil)
		cmd.Flags().Set("name", "testrule")
		cmd.Flags().Set("network", "net-abcdef12")
		cmd.Flags().Set("subnet", "10.0.0.0/24")
		cmd.Flags().Set("floating-ip", "fip-abcdef12")
		cmd.Flags().Set("action", "ALLOW")

		service.EXPECT().CreateNATOverloadRule(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error")).Times(1)

		err := ecloudNATOverloadRuleCreate(service, cmd, []string{})

		assert.Equal(t, "error creating NAT overload rule: test error", err.Error())
	})

	t.Run("GetNATOverloadRuleError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNATOverloadRuleCreateCmd(nil)
		cmd.Flags().Set("name", "testrule")
		cmd.Flags().Set("network", "net-abcdef12")
		cmd.Flags().Set("subnet", "10.0.0.0/24")
		cmd.Flags().Set("floating-ip", "fip-abcdef12")
		cmd.Flags().Set("action", "ALLOW")

		gomock.InOrder(
			service.EXPECT().CreateNATOverloadRule(gomock.Any()).Return(ecloud.TaskReference{ResourceID: "nor-abcdef12"}, nil),
			service.EXPECT().GetNATOverloadRule("nor-abcdef12").Return(ecloud.NATOverloadRule{}, errors.New("test error")),
		)

		err := ecloudNATOverloadRuleCreate(service, cmd, []string{})

		assert.Equal(t, "error retrieving new NAT overload rule: test error", err.Error())
	})
}

func Test_ecloudNATOverloadRuleUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudNATOverloadRuleUpdateCmd(nil).Args(nil, []string{"nor-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudNATOverloadRuleUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing NAT overload rule", err.Error())
	})
}

func Test_ecloudNATOverloadRuleUpdate(t *testing.T) {
	t.Run("SingleRule", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNATOverloadRuleCreateCmd(nil)
		cmd.Flags().Set("name", "testrule")

		req := ecloud.PatchNATOverloadRuleRequest{
			Name: "testrule",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "nor-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchNATOverloadRule("nor-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetNATOverloadRule("nor-abcdef12").Return(ecloud.NATOverloadRule{}, nil),
		)

		ecloudNATOverloadRuleUpdate(service, cmd, []string{"nor-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNATOverloadRuleUpdateCmd(nil)
		cmd.Flags().Set("name", "testrule")
		cmd.Flags().Set("wait", "true")

		req := ecloud.PatchNATOverloadRuleRequest{
			Name: "testrule",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "nor-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchNATOverloadRule("nor-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetNATOverloadRule("nor-abcdef12").Return(ecloud.NATOverloadRule{}, nil),
		)

		ecloudNATOverloadRuleUpdate(service, cmd, []string{"nor-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNATOverloadRuleUpdateCmd(nil)
		cmd.Flags().Set("name", "testrule")
		cmd.Flags().Set("wait", "true")

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "nor-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchNATOverloadRule("nor-abcdef12", gomock.Any()).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for NAT overload rule [nor-abcdef12]: error waiting for command: failed to retrieve task status: test error\n", func() {
			ecloudNATOverloadRuleUpdate(service, cmd, []string{"nor-abcdef12"})
		})
	})

	t.Run("MultipleNATOverloadRules", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNATOverloadRuleUpdateCmd(nil)
		cmd.Flags().Set("name", "testrule")

		resp1 := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "nor-abcdef12",
		}

		resp2 := ecloud.TaskReference{
			TaskID:     "task-abcdef23",
			ResourceID: "nor-12abcdef",
		}

		gomock.InOrder(
			service.EXPECT().PatchNATOverloadRule("nor-abcdef12", gomock.Any()).Return(resp1, nil),
			service.EXPECT().GetNATOverloadRule("nor-abcdef12").Return(ecloud.NATOverloadRule{}, nil),
			service.EXPECT().PatchNATOverloadRule("nor-12abcdef", gomock.Any()).Return(resp2, nil),
			service.EXPECT().GetNATOverloadRule("nor-12abcdef").Return(ecloud.NATOverloadRule{}, nil),
		)

		ecloudNATOverloadRuleUpdate(service, cmd, []string{"nor-abcdef12", "nor-12abcdef"})
	})

	t.Run("PatchNATOverloadRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchNATOverloadRule("nor-abcdef12", gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating NAT overload rule [nor-abcdef12]: test error\n", func() {
			ecloudNATOverloadRuleUpdate(service, &cobra.Command{}, []string{"nor-abcdef12"})
		})
	})

	t.Run("GetNATOverloadRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "nor-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchNATOverloadRule("nor-abcdef12", gomock.Any()).Return(resp, nil),
			service.EXPECT().GetNATOverloadRule("nor-abcdef12").Return(ecloud.NATOverloadRule{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated NAT overload rule [nor-abcdef12]: test error\n", func() {
			ecloudNATOverloadRuleUpdate(service, &cobra.Command{}, []string{"nor-abcdef12"})
		})
	})
}

func Test_ecloudNATOverloadRuleDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudNATOverloadRuleDeleteCmd(nil).Args(nil, []string{"nor-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudNATOverloadRuleDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing NAT overload rule", err.Error())
	})
}

func Test_ecloudNATOverloadRuleDelete(t *testing.T) {
	t.Run("SingleRule", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteNATOverloadRule("nor-abcdef12").Return("task-abcdef12", nil).Times(1)

		ecloudNATOverloadRuleDelete(service, &cobra.Command{}, []string{"nor-abcdef12"})
	})

	t.Run("MultipleNATOverloadRules", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteNATOverloadRule("nor-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().DeleteNATOverloadRule("nor-12abcdef").Return("task-abcdef23", nil),
		)

		ecloudNATOverloadRuleDelete(service, &cobra.Command{}, []string{"nor-abcdef12", "nor-12abcdef"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudNATOverloadRuleDeleteCmd(nil)
		cmd.Flags().Set("wait", "true")

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteNATOverloadRule("nor-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudNATOverloadRuleDelete(service, cmd, []string{"nor-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNATOverloadRuleDeleteCmd(nil)
		cmd.Flags().Set("wait", "true")

		gomock.InOrder(
			service.EXPECT().DeleteNATOverloadRule("nor-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for NAT overload rule [nor-abcdef12]: error waiting for command: failed to retrieve task status: test error\n", func() {
			ecloudNATOverloadRuleDelete(service, cmd, []string{"nor-abcdef12"})
		})
	})

	t.Run("DeleteNATOverloadRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteNATOverloadRule("nor-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing NAT overload rule [nor-abcdef12]: test error\n", func() {
			ecloudNATOverloadRuleDelete(service, &cobra.Command{}, []string{"nor-abcdef12"})
		})
	})
}
