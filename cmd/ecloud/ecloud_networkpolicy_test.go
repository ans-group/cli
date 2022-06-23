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

func Test_ecloudNetworkPolicyList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworkPolicies(gomock.Any()).Return([]ecloud.NetworkPolicy{}, nil).Times(1)

		ecloudNetworkPolicyList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudNetworkPolicyList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetNetworkPoliciesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworkPolicies(gomock.Any()).Return([]ecloud.NetworkPolicy{}, errors.New("test error")).Times(1)

		err := ecloudNetworkPolicyList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving network policies: test error", err.Error())
	})
}

func Test_ecloudNetworkPolicyShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudNetworkPolicyShowCmd(nil).Args(nil, []string{"np-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudNetworkPolicyShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing network policy", err.Error())
	})
}

func Test_ecloudNetworkPolicyShow(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworkPolicy("np-abcdef12").Return(ecloud.NetworkPolicy{}, nil).Times(1)

		ecloudNetworkPolicyShow(service, &cobra.Command{}, []string{"np-abcdef12"})
	})

	t.Run("MultipleNetworkPolicies", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetNetworkPolicy("np-abcdef12").Return(ecloud.NetworkPolicy{}, nil),
			service.EXPECT().GetNetworkPolicy("np-abcdef23").Return(ecloud.NetworkPolicy{}, nil),
		)

		ecloudNetworkPolicyShow(service, &cobra.Command{}, []string{"np-abcdef12", "np-abcdef23"})
	})

	t.Run("GetNetworkPolicyError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworkPolicy("np-abcdef12").Return(ecloud.NetworkPolicy{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving network policy [np-abcdef12]: test error\n", func() {
			ecloudNetworkPolicyShow(service, &cobra.Command{}, []string{"np-abcdef12"})
		})
	})
}

func Test_ecloudNetworkPolicyCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNetworkPolicyCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		req := ecloud.CreateNetworkPolicyRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "np-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateNetworkPolicy(req).Return(resp, nil),
			service.EXPECT().GetNetworkPolicy("np-abcdef12").Return(ecloud.NetworkPolicy{}, nil),
		)

		ecloudNetworkPolicyCreate(service, cmd, []string{})
	})

	t.Run("CreateWithWaitFlagNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNetworkPolicyCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.CreateNetworkPolicyRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "np-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateNetworkPolicy(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetNetworkPolicy("np-abcdef12").Return(ecloud.NetworkPolicy{}, nil),
		)

		ecloudNetworkPolicyCreate(service, cmd, []string{})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNetworkPolicyCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.CreateNetworkPolicyRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "np-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateNetworkPolicy(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		err := ecloudNetworkPolicyCreate(service, cmd, []string{})
		assert.Equal(t, "Error waiting for network policy task to complete: Error waiting for command: Failed to retrieve task status: test error", err.Error())
	})

	t.Run("CreateNetworkPolicyError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNetworkPolicyCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		service.EXPECT().CreateNetworkPolicy(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error")).Times(1)

		err := ecloudNetworkPolicyCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating network policy: test error", err.Error())
	})

	t.Run("GetNetworkPolicyError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNetworkPolicyCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		gomock.InOrder(
			service.EXPECT().CreateNetworkPolicy(gomock.Any()).Return(ecloud.TaskReference{ResourceID: "np-abcdef12"}, nil),
			service.EXPECT().GetNetworkPolicy("np-abcdef12").Return(ecloud.NetworkPolicy{}, errors.New("test error")),
		)

		err := ecloudNetworkPolicyCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new network policy: test error", err.Error())
	})
}

func Test_ecloudNetworkPolicyUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudNetworkPolicyUpdateCmd(nil).Args(nil, []string{"np-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudNetworkPolicyUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing network policy", err.Error())
	})
}

func Test_ecloudNetworkPolicyUpdate(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNetworkPolicyUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		req := ecloud.PatchNetworkPolicyRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "np-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchNetworkPolicy("np-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetNetworkPolicy("np-abcdef12").Return(ecloud.NetworkPolicy{}, nil),
		)

		ecloudNetworkPolicyUpdate(service, cmd, []string{"np-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNetworkPolicyUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.PatchNetworkPolicyRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "np-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchNetworkPolicy("np-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetNetworkPolicy("np-abcdef12").Return(ecloud.NetworkPolicy{}, nil),
		)

		ecloudNetworkPolicyUpdate(service, cmd, []string{"np-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNetworkPolicyUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.PatchNetworkPolicyRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "np-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchNetworkPolicy("np-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for network policy [np-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudNetworkPolicyUpdate(service, cmd, []string{"np-abcdef12"})
		})
	})

	t.Run("PatchNetworkPolicyError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchNetworkPolicy("np-abcdef12", gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating network policy [np-abcdef12]: test error\n", func() {
			ecloudNetworkPolicyUpdate(service, &cobra.Command{}, []string{"np-abcdef12"})
		})
	})

	t.Run("GetNetworkPolicyError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "np-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchNetworkPolicy("np-abcdef12", gomock.Any()).Return(resp, nil),
			service.EXPECT().GetNetworkPolicy("np-abcdef12").Return(ecloud.NetworkPolicy{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated network policy [np-abcdef12]: test error\n", func() {
			ecloudNetworkPolicyUpdate(service, &cobra.Command{}, []string{"np-abcdef12"})
		})
	})
}

func Test_ecloudNetworkPolicyDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudNetworkPolicyDeleteCmd(nil).Args(nil, []string{"np-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudNetworkPolicyDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing network policy", err.Error())
	})
}

func Test_ecloudNetworkPolicyDelete(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteNetworkPolicy("np-abcdef12").Return("task-abcdef12", nil)

		ecloudNetworkPolicyDelete(service, &cobra.Command{}, []string{"np-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudNetworkPolicyDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteNetworkPolicy("np-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudNetworkPolicyDelete(service, cmd, []string{"np-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNetworkPolicyDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		gomock.InOrder(
			service.EXPECT().DeleteNetworkPolicy("np-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for network policy [np-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudNetworkPolicyDelete(service, cmd, []string{"np-abcdef12"})
		})
	})

	t.Run("DeleteNetworkPolicyError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteNetworkPolicy("np-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing network policy [np-abcdef12]: test error\n", func() {
			ecloudNetworkPolicyDelete(service, &cobra.Command{}, []string{"np-abcdef12"})
		})
	})
}
