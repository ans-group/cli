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

func Test_ecloudNetworkRuleList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworkRules(gomock.Any()).Return([]ecloud.NetworkRule{}, nil).Times(1)

		ecloudNetworkRuleList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudNetworkRuleList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetNetworkRulesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworkRules(gomock.Any()).Return([]ecloud.NetworkRule{}, errors.New("test error")).Times(1)

		err := ecloudNetworkRuleList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "error retrieving network rules: test error", err.Error())
	})
}

func Test_ecloudNetworkRuleShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudNetworkRuleShowCmd(nil).Args(nil, []string{"nr-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudNetworkRuleShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing network rule", err.Error())
	})
}

func Test_ecloudNetworkRuleShow(t *testing.T) {
	t.Run("SingleNetworkRule", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworkRule("nr-abcdef12").Return(ecloud.NetworkRule{}, nil).Times(1)

		ecloudNetworkRuleShow(service, &cobra.Command{}, []string{"nr-abcdef12"})
	})

	t.Run("MultipleNetworkRules", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetNetworkRule("nr-abcdef12").Return(ecloud.NetworkRule{}, nil),
			service.EXPECT().GetNetworkRule("nr-abcdef23").Return(ecloud.NetworkRule{}, nil),
		)

		ecloudNetworkRuleShow(service, &cobra.Command{}, []string{"nr-abcdef12", "nr-abcdef23"})
	})

	t.Run("GetNetworkRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworkRule("nr-abcdef12").Return(ecloud.NetworkRule{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving network rule [nr-abcdef12]: test error\n", func() {
			ecloudNetworkRuleShow(service, &cobra.Command{}, []string{"nr-abcdef12"})
		})
	})
}

func Test_ecloudNetworkRuleCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNetworkRuleCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testrule", "--direction=IN", "--action=DROP"})

		req := ecloud.CreateNetworkRuleRequest{
			Name:      "testrule",
			Direction: ecloud.NetworkRuleDirectionIn,
			Action:    ecloud.NetworkRuleActionDrop,
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "nr-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateNetworkRule(req).Return(resp, nil),
			service.EXPECT().GetNetworkRule("nr-abcdef12").Return(ecloud.NetworkRule{}, nil),
		)

		ecloudNetworkRuleCreate(service, cmd, []string{})
	})

	t.Run("CreateWithWaitFlagNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNetworkRuleCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testrule", "--direction=IN", "--action=DROP", "--wait"})

		req := ecloud.CreateNetworkRuleRequest{
			Name:      "testrule",
			Direction: ecloud.NetworkRuleDirectionIn,
			Action:    ecloud.NetworkRuleActionDrop,
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "nr-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateNetworkRule(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetNetworkRule("nr-abcdef12").Return(ecloud.NetworkRule{}, nil),
		)

		ecloudNetworkRuleCreate(service, cmd, []string{})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNetworkRuleCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testrule", "--direction=IN", "--action=DROP", "--wait"})

		req := ecloud.CreateNetworkRuleRequest{
			Name:      "testrule",
			Direction: ecloud.NetworkRuleDirectionIn,
			Action:    ecloud.NetworkRuleActionDrop,
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "nr-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateNetworkRule(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		err := ecloudNetworkRuleCreate(service, cmd, []string{})
		assert.Equal(t, "error waiting for network rule task to complete: error waiting for command: failed to retrieve task status: test error", err.Error())
	})

	t.Run("CreateNetworkRuleError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNetworkRuleCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testrule", "--direction=IN", "--action=DROP"})

		service.EXPECT().CreateNetworkRule(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error")).Times(1)

		err := ecloudNetworkRuleCreate(service, cmd, []string{})

		assert.Equal(t, "error creating network rule: test error", err.Error())
	})

	t.Run("GetNetworkRuleError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNetworkRuleCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testrule", "--direction=IN", "--action=DROP"})

		gomock.InOrder(
			service.EXPECT().CreateNetworkRule(gomock.Any()).Return(ecloud.TaskReference{ResourceID: "nr-abcdef12"}, nil),
			service.EXPECT().GetNetworkRule("nr-abcdef12").Return(ecloud.NetworkRule{}, errors.New("test error")),
		)

		err := ecloudNetworkRuleCreate(service, cmd, []string{})

		assert.Equal(t, "error retrieving new network rule: test error", err.Error())
	})
}

func Test_ecloudNetworkRuleUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudNetworkRuleUpdateCmd(nil).Args(nil, []string{"nr-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudNetworkRuleUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing network rule", err.Error())
	})
}

func Test_ecloudNetworkRuleUpdate(t *testing.T) {
	t.Run("SingleRule", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNetworkRuleUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testrule"})

		req := ecloud.PatchNetworkRuleRequest{
			Name: "testrule",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "nr-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchNetworkRule("nr-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetNetworkRule("nr-abcdef12").Return(ecloud.NetworkRule{}, nil),
		)

		ecloudNetworkRuleUpdate(service, cmd, []string{"nr-abcdef12"})
	})

	t.Run("MultipleNetworkRules", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp1 := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "nr-abcdef12",
		}

		resp2 := ecloud.TaskReference{
			TaskID:     "task-abcdef23",
			ResourceID: "nr-12abcdef",
		}

		gomock.InOrder(
			service.EXPECT().PatchNetworkRule("nr-abcdef12", gomock.Any()).Return(resp1, nil),
			service.EXPECT().GetNetworkRule("nr-abcdef12").Return(ecloud.NetworkRule{}, nil),
			service.EXPECT().PatchNetworkRule("nr-12abcdef", gomock.Any()).Return(resp2, nil),
			service.EXPECT().GetNetworkRule("nr-12abcdef").Return(ecloud.NetworkRule{}, nil),
		)

		ecloudNetworkRuleUpdate(service, &cobra.Command{}, []string{"nr-abcdef12", "nr-12abcdef"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNetworkRuleUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testrule", "--wait"})

		req := ecloud.PatchNetworkRuleRequest{
			Name: "testrule",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "nr-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchNetworkRule("nr-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetNetworkRule("nr-abcdef12").Return(ecloud.NetworkRule{}, nil),
		)

		ecloudNetworkRuleUpdate(service, cmd, []string{"nr-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNetworkRuleUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testrule", "--wait"})

		req := ecloud.PatchNetworkRuleRequest{
			Name: "testrule",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "nr-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchNetworkRule("nr-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for network rule [nr-abcdef12]: error waiting for command: failed to retrieve task status: test error\n", func() {
			ecloudNetworkRuleUpdate(service, cmd, []string{"nr-abcdef12"})
		})
	})

	t.Run("PatchNetworkRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchNetworkRule("nr-abcdef12", gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating network rule [nr-abcdef12]: test error\n", func() {
			ecloudNetworkRuleUpdate(service, &cobra.Command{}, []string{"nr-abcdef12"})
		})
	})

	t.Run("GetNetworkRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "nr-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchNetworkRule("nr-abcdef12", gomock.Any()).Return(resp, nil),
			service.EXPECT().GetNetworkRule("nr-abcdef12").Return(ecloud.NetworkRule{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated network rule [nr-abcdef12]: test error\n", func() {
			ecloudNetworkRuleUpdate(service, &cobra.Command{}, []string{"nr-abcdef12"})
		})
	})
}

func Test_ecloudNetworkRuleDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudNetworkRuleDeleteCmd(nil).Args(nil, []string{"nr-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudNetworkRuleDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing network rule", err.Error())
	})
}

func Test_ecloudNetworkRuleDelete(t *testing.T) {
	t.Run("SingleRule", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteNetworkRule("nr-abcdef12").Return("task-abcdef12", nil).Times(1)

		ecloudNetworkRuleDelete(service, &cobra.Command{}, []string{"nr-abcdef12"})
	})

	t.Run("MultipleNetworkRules", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteNetworkRule("nr-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().DeleteNetworkRule("nr-12abcdef").Return("task-abcdef23", nil),
		)

		ecloudNetworkRuleDelete(service, &cobra.Command{}, []string{"nr-abcdef12", "nr-12abcdef"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudNetworkRuleDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteNetworkRule("nr-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudNetworkRuleDelete(service, cmd, []string{"nr-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNetworkRuleDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		gomock.InOrder(
			service.EXPECT().DeleteNetworkRule("nr-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for network rule [nr-abcdef12]: error waiting for command: failed to retrieve task status: test error\n", func() {
			ecloudNetworkRuleDelete(service, cmd, []string{"nr-abcdef12"})
		})
	})

	t.Run("DeleteNetworkRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteNetworkRule("nr-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing network rule [nr-abcdef12]: test error\n", func() {
			ecloudNetworkRuleDelete(service, &cobra.Command{}, []string{"nr-abcdef12"})
		})
	})
}
