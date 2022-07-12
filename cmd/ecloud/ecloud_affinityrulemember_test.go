package ecloud

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ecloudAffinityRuleMemberShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudAffinityRuleMemberShowCmd(nil).Args(nil, []string{"arm-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudAffinityRuleMemberShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing affinity rule member", err.Error())
	})
}

func Test_ecloudAffinityRuleMemberShow(t *testing.T) {
	t.Run("SingleRuleMember", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetAffinityRuleMember("arm-abcdef12").Return(ecloud.AffinityRuleMember{}, nil).Times(1)

		ecloudAffinityRuleMemberShow(service, &cobra.Command{}, []string{"arm-abcdef12"})
	})

	t.Run("MultipleAffinityRuleMembers", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetAffinityRuleMember("arm-abcdef12").Return(ecloud.AffinityRuleMember{}, nil),
			service.EXPECT().GetAffinityRuleMember("arm-abcdef23").Return(ecloud.AffinityRuleMember{}, nil),
		)

		ecloudAffinityRuleMemberShow(service, &cobra.Command{}, []string{"arm-abcdef12", "arm-abcdef23"})
	})

	t.Run("GetAffinityRuleMemberError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetAffinityRuleMember("arm-abcdef12").Return(ecloud.AffinityRuleMember{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving affinity rule member [arm-abcdef12]: test error\n", func() {
			ecloudAffinityRuleMemberShow(service, &cobra.Command{}, []string{"arm-abcdef12"})
		})
	})
}

func Test_ecloudAffinityRuleMemberCreate(t *testing.T) {
	t.Run("CreateWithRequiredArgsNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudAffinityRuleMemberCreateCmd(nil)
		cmd.ParseFlags([]string{"--affinity-rule=ar-abcdef12", "--instance=i-abcdef12"})

		req := ecloud.CreateAffinityRuleMemberRequest{
			InstanceID: "i-abcdef12",
			AffinityRuleID: "ar-abcdef12",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "arm-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateAffinityRuleMember(req).Return(resp, nil),
			service.EXPECT().GetAffinityRuleMember("arm-abcdef12").Return(ecloud.AffinityRuleMember{}, nil),
		)

		ecloudAffinityRuleMemberCreate(service, cmd, []string{})
	})

	t.Run("CreateWithWaitFlagNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudAffinityRuleMemberCreateCmd(nil)
		cmd.ParseFlags([]string{"--affinity-rule=ar-abcdef12", "--instance=i-abcdef12", "--wait"})

		req := ecloud.CreateAffinityRuleMemberRequest{
			InstanceID: "i-abcdef12",
			AffinityRuleID: "ar-abcdef12",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "arm-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateAffinityRuleMember(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetAffinityRuleMember("arm-abcdef12").Return(ecloud.AffinityRuleMember{}, nil),
		)

		ecloudAffinityRuleMemberCreate(service, cmd, []string{})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudAffinityRuleMemberCreateCmd(nil)
		cmd.ParseFlags([]string{"--affinity-rule=ar-abcdef12", "--instance=i-abcdef12", "--wait"})

		req := ecloud.CreateAffinityRuleMemberRequest{
			InstanceID: "i-abcdef12",
			AffinityRuleID: "ar-abcdef12",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "arm-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateAffinityRuleMember(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		err := ecloudAffinityRuleMemberCreate(service, cmd, []string{})
		assert.Equal(t, "Error waiting for affinity rule member task to complete: Error waiting for command: Failed to retrieve task status: test error", err.Error())
	})

	t.Run("CreateAffinityRuleMemberError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudAffinityRuleMemberCreateCmd(nil)
		cmd.ParseFlags([]string{"--affinity-rule=ar-abcdef12", "--instance=i-abcdef12"})

		service.EXPECT().CreateAffinityRuleMember(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error")).Times(1)

		err := ecloudAffinityRuleMemberCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating affinity rule member: test error", err.Error())
	})

	t.Run("GetAffinityRuleMemberError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudAffinityRuleMemberCreateCmd(nil)
		cmd.ParseFlags([]string{"--affinity-rule=ar-abcdef12", "--instance=i-abcdef12"})

		gomock.InOrder(
			service.EXPECT().CreateAffinityRuleMember(gomock.Any()).Return(ecloud.TaskReference{ResourceID: "arm-abcdef12"}, nil),
			service.EXPECT().GetAffinityRuleMember("arm-abcdef12").Return(ecloud.AffinityRuleMember{}, errors.New("test error")),
		)

		err := ecloudAffinityRuleMemberCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new affinity rule member: test error", err.Error())
	})
}

func Test_ecloudAffinityRuleMemberDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudAffinityRuleMemberDeleteCmd(nil).Args(nil, []string{"arm-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudAffinityRuleMemberDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing affinity rule member", err.Error())
	})
}

func Test_ecloudAffinityRuleMemberDelete(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteAffinityRuleMember("arm-abcdef12").Return("task-abcdef12", nil)

		ecloudAffinityRuleMemberDelete(service, &cobra.Command{}, []string{"arm-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudAffinityRuleMemberDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteAffinityRuleMember("arm-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudAffinityRuleMemberDelete(service, cmd, []string{"arm-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudAffinityRuleMemberDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		gomock.InOrder(
			service.EXPECT().DeleteAffinityRuleMember("arm-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for affinity rule member [arm-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudAffinityRuleMemberDelete(service, cmd, []string{"arm-abcdef12"})
		})
	})

	t.Run("DeleteAffinityRuleMemberError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteAffinityRuleMember("arm-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing affinity rule member [arm-abcdef12]: test error\n", func() {
			ecloudAffinityRuleMemberDelete(service, &cobra.Command{}, []string{"arm-abcdef12"})
		})
	})
}
