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

func Test_ecloudFirewallPolicyList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewallPolicies(gomock.Any()).Return([]ecloud.FirewallPolicy{}, nil).Times(1)

		ecloudFirewallPolicyList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudFirewallPolicyList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetFirewallPoliciesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewallPolicies(gomock.Any()).Return([]ecloud.FirewallPolicy{}, errors.New("test error")).Times(1)

		err := ecloudFirewallPolicyList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving firewall policies: test error", err.Error())
	})
}

func Test_ecloudFirewallPolicyShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFirewallPolicyShowCmd(nil).Args(nil, []string{"fwp-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFirewallPolicyShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing firewall policy", err.Error())
	})
}

func Test_ecloudFirewallPolicyShow(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewallPolicy("fwp-abcdef12").Return(ecloud.FirewallPolicy{}, nil).Times(1)

		ecloudFirewallPolicyShow(service, &cobra.Command{}, []string{"fwp-abcdef12"})
	})

	t.Run("MultipleFirewallPolicies", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetFirewallPolicy("fwp-abcdef12").Return(ecloud.FirewallPolicy{}, nil),
			service.EXPECT().GetFirewallPolicy("fwp-abcdef23").Return(ecloud.FirewallPolicy{}, nil),
		)

		ecloudFirewallPolicyShow(service, &cobra.Command{}, []string{"fwp-abcdef12", "fwp-abcdef23"})
	})

	t.Run("GetFirewallPolicyError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewallPolicy("fwp-abcdef12").Return(ecloud.FirewallPolicy{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving firewall policy [fwp-abcdef12]: test error\n", func() {
			ecloudFirewallPolicyShow(service, &cobra.Command{}, []string{"fwp-abcdef12"})
		})
	})
}

func Test_ecloudFirewallPolicyCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFirewallPolicyCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		req := ecloud.CreateFirewallPolicyRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fwp-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateFirewallPolicy(req).Return(resp, nil),
			service.EXPECT().GetFirewallPolicy("fwp-abcdef12").Return(ecloud.FirewallPolicy{}, nil),
		)

		ecloudFirewallPolicyCreate(service, cmd, []string{})
	})

	t.Run("CreateWithWaitFlagNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFirewallPolicyCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.CreateFirewallPolicyRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fwp-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateFirewallPolicy(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetFirewallPolicy("fwp-abcdef12").Return(ecloud.FirewallPolicy{}, nil),
		)

		ecloudFirewallPolicyCreate(service, cmd, []string{})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFirewallPolicyCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.CreateFirewallPolicyRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fwp-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateFirewallPolicy(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		err := ecloudFirewallPolicyCreate(service, cmd, []string{})
		assert.Equal(t, "Error waiting for firewall policy task to complete: Error waiting for command: Failed to retrieve task status: test error", err.Error())
	})

	t.Run("CreateFirewallPolicyError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFirewallPolicyCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		service.EXPECT().CreateFirewallPolicy(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error")).Times(1)

		err := ecloudFirewallPolicyCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating firewall policy: test error", err.Error())
	})

	t.Run("GetFirewallPolicyError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudFirewallPolicyCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		gomock.InOrder(
			service.EXPECT().CreateFirewallPolicy(gomock.Any()).Return(ecloud.TaskReference{ResourceID: "fwp-abcdef12"}, nil),
			service.EXPECT().GetFirewallPolicy("fwp-abcdef12").Return(ecloud.FirewallPolicy{}, errors.New("test error")),
		)

		err := ecloudFirewallPolicyCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new firewall policy: test error", err.Error())
	})
}

func Test_ecloudFirewallPolicyUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFirewallPolicyUpdateCmd(nil).Args(nil, []string{"fwp-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFirewallPolicyUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing firewall policy", err.Error())
	})
}

func Test_ecloudFirewallPolicyUpdate(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudFirewallPolicyUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		req := ecloud.PatchFirewallPolicyRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fwp-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchFirewallPolicy("fwp-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetFirewallPolicy("fwp-abcdef12").Return(ecloud.FirewallPolicy{}, nil),
		)

		ecloudFirewallPolicyUpdate(service, cmd, []string{"fwp-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudFirewallPolicyUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.PatchFirewallPolicyRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fwp-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchFirewallPolicy("fwp-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetFirewallPolicy("fwp-abcdef12").Return(ecloud.FirewallPolicy{}, nil),
		)

		ecloudFirewallPolicyUpdate(service, cmd, []string{"fwp-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudFirewallPolicyUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.PatchFirewallPolicyRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fwp-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchFirewallPolicy("fwp-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for firewall policy [fwp-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudFirewallPolicyUpdate(service, cmd, []string{"fwp-abcdef12"})
		})
	})

	t.Run("PatchFirewallPolicyError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchFirewallPolicy("fwp-abcdef12", gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating firewall policy [fwp-abcdef12]: test error\n", func() {
			ecloudFirewallPolicyUpdate(service, &cobra.Command{}, []string{"fwp-abcdef12"})
		})
	})

	t.Run("GetFirewallPolicyError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "fwp-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchFirewallPolicy("fwp-abcdef12", gomock.Any()).Return(resp, nil),
			service.EXPECT().GetFirewallPolicy("fwp-abcdef12").Return(ecloud.FirewallPolicy{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated firewall policy [fwp-abcdef12]: test error\n", func() {
			ecloudFirewallPolicyUpdate(service, &cobra.Command{}, []string{"fwp-abcdef12"})
		})
	})
}

func Test_ecloudFirewallPolicyDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFirewallPolicyDeleteCmd(nil).Args(nil, []string{"fwp-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFirewallPolicyDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing firewall policy", err.Error())
	})
}

func Test_ecloudFirewallPolicyDelete(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteFirewallPolicy("fwp-abcdef12").Return("task-abcdef12", nil)

		ecloudFirewallPolicyDelete(service, &cobra.Command{}, []string{"fwp-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudFirewallPolicyDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteFirewallPolicy("fwp-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudFirewallPolicyDelete(service, cmd, []string{"fwp-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudFirewallPolicyDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		gomock.InOrder(
			service.EXPECT().DeleteFirewallPolicy("fwp-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for firewall policy [fwp-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudFirewallPolicyDelete(service, cmd, []string{"fwp-abcdef12"})
		})
	})

	t.Run("DeleteFirewallPolicyError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteFirewallPolicy("fwp-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing firewall policy [fwp-abcdef12]: test error\n", func() {
			ecloudFirewallPolicyDelete(service, &cobra.Command{}, []string{"fwp-abcdef12"})
		})
	})
}
