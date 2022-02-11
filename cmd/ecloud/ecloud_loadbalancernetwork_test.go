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

func Test_ecloudLoadBalancerNetworkList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetLoadBalancerNetworks(gomock.Any()).Return([]ecloud.LoadBalancerNetwork{}, nil).Times(1)

		ecloudLoadBalancerNetworkList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudLoadBalancerNetworkList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetLoadBalancerNetworksError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetLoadBalancerNetworks(gomock.Any()).Return([]ecloud.LoadBalancerNetwork{}, errors.New("test error")).Times(1)

		err := ecloudLoadBalancerNetworkList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving load balancer networks: test error", err.Error())
	})
}

func Test_ecloudLoadBalancerNetworkShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudLoadBalancerNetworkShowCmd(nil).Args(nil, []string{"lbn-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudLoadBalancerNetworkShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing load balancer network", err.Error())
	})
}

func Test_ecloudLoadBalancerNetworkShow(t *testing.T) {
	t.Run("SingleLoadBalancerNetwork", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetLoadBalancerNetwork("lbn-abcdef12").Return(ecloud.LoadBalancerNetwork{}, nil).Times(1)

		ecloudLoadBalancerNetworkShow(service, &cobra.Command{}, []string{"lbn-abcdef12"})
	})

	t.Run("MultipleLoadBalancerNetworks", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetLoadBalancerNetwork("lbn-abcdef12").Return(ecloud.LoadBalancerNetwork{}, nil),
			service.EXPECT().GetLoadBalancerNetwork("lbn-abcdef23").Return(ecloud.LoadBalancerNetwork{}, nil),
		)

		ecloudLoadBalancerNetworkShow(service, &cobra.Command{}, []string{"lbn-abcdef12", "lbn-abcdef23"})
	})

	t.Run("GetLoadBalancerNetworkError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetLoadBalancerNetwork("lbn-abcdef12").Return(ecloud.LoadBalancerNetwork{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving load balancer network [lbn-abcdef12]: test error\n", func() {
			ecloudLoadBalancerNetworkShow(service, &cobra.Command{}, []string{"lbn-abcdef12"})
		})
	})
}

func Test_ecloudLoadBalancerNetworkCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudLoadBalancerNetworkCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testlb"})

		req := ecloud.CreateLoadBalancerNetworkRequest{
			Name: "testlb",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "lbn-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateLoadBalancerNetwork(req).Return(resp, nil),
			service.EXPECT().GetLoadBalancerNetwork("lbn-abcdef12").Return(ecloud.LoadBalancerNetwork{}, nil),
		)

		ecloudLoadBalancerNetworkCreate(service, cmd, []string{})
	})

	t.Run("CreateWithWaitFlagNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudLoadBalancerNetworkCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testlb", "--wait"})

		req := ecloud.CreateLoadBalancerNetworkRequest{
			Name: "testlb",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "lbn-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateLoadBalancerNetwork(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetLoadBalancerNetwork("lbn-abcdef12").Return(ecloud.LoadBalancerNetwork{}, nil),
		)

		ecloudLoadBalancerNetworkCreate(service, cmd, []string{})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudLoadBalancerNetworkCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testlb", "--wait"})

		req := ecloud.CreateLoadBalancerNetworkRequest{
			Name: "testlb",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "lbn-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateLoadBalancerNetwork(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		err := ecloudLoadBalancerNetworkCreate(service, cmd, []string{})
		assert.Equal(t, "Error waiting for load balancer network task to complete: Error waiting for command: Failed to retrieve task status: test error", err.Error())
	})

	t.Run("CreateLoadBalancerNetworkError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudLoadBalancerNetworkCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testlb"})

		service.EXPECT().CreateLoadBalancerNetwork(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error")).Times(1)

		err := ecloudLoadBalancerNetworkCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating load balancer network: test error", err.Error())
	})

	t.Run("GetLoadBalancerNetworkError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudLoadBalancerNetworkCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testlb"})

		gomock.InOrder(
			service.EXPECT().CreateLoadBalancerNetwork(gomock.Any()).Return(ecloud.TaskReference{ResourceID: "lbn-abcdef12"}, nil),
			service.EXPECT().GetLoadBalancerNetwork("lbn-abcdef12").Return(ecloud.LoadBalancerNetwork{}, errors.New("test error")),
		)

		err := ecloudLoadBalancerNetworkCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new load balancer network: test error", err.Error())
	})
}

func Test_ecloudLoadBalancerNetworkUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudLoadBalancerNetworkUpdateCmd(nil).Args(nil, []string{"lbn-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudLoadBalancerNetworkUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing load balancer network", err.Error())
	})
}

func Test_ecloudLoadBalancerNetworkUpdate(t *testing.T) {
	t.Run("SingleRule", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudLoadBalancerNetworkUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testlb"})

		req := ecloud.PatchLoadBalancerNetworkRequest{
			Name: "testlb",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "lbn-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchLoadBalancerNetwork("lbn-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetLoadBalancerNetwork("lbn-abcdef12").Return(ecloud.LoadBalancerNetwork{}, nil),
		)

		ecloudLoadBalancerNetworkUpdate(service, cmd, []string{"lbn-abcdef12"})
	})

	t.Run("MultipleLoadBalancerNetworks", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp1 := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "lbn-abcdef12",
		}

		resp2 := ecloud.TaskReference{
			TaskID:     "task-abcdef23",
			ResourceID: "lbn-12abcdef",
		}

		gomock.InOrder(
			service.EXPECT().PatchLoadBalancerNetwork("lbn-abcdef12", gomock.Any()).Return(resp1, nil),
			service.EXPECT().GetLoadBalancerNetwork("lbn-abcdef12").Return(ecloud.LoadBalancerNetwork{}, nil),
			service.EXPECT().PatchLoadBalancerNetwork("lbn-12abcdef", gomock.Any()).Return(resp2, nil),
			service.EXPECT().GetLoadBalancerNetwork("lbn-12abcdef").Return(ecloud.LoadBalancerNetwork{}, nil),
		)

		ecloudLoadBalancerNetworkUpdate(service, &cobra.Command{}, []string{"lbn-abcdef12", "lbn-12abcdef"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudLoadBalancerNetworkUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testlb", "--wait"})

		req := ecloud.PatchLoadBalancerNetworkRequest{
			Name: "testlb",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "lbn-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchLoadBalancerNetwork("lbn-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetLoadBalancerNetwork("lbn-abcdef12").Return(ecloud.LoadBalancerNetwork{}, nil),
		)

		ecloudLoadBalancerNetworkUpdate(service, cmd, []string{"lbn-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudLoadBalancerNetworkUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testlb", "--wait"})

		req := ecloud.PatchLoadBalancerNetworkRequest{
			Name: "testlb",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "lbn-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchLoadBalancerNetwork("lbn-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for load balancer network [lbn-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudLoadBalancerNetworkUpdate(service, cmd, []string{"lbn-abcdef12"})
		})
	})

	t.Run("PatchLoadBalancerNetworkError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchLoadBalancerNetwork("lbn-abcdef12", gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating load balancer network [lbn-abcdef12]: test error\n", func() {
			ecloudLoadBalancerNetworkUpdate(service, &cobra.Command{}, []string{"lbn-abcdef12"})
		})
	})

	t.Run("GetLoadBalancerNetworkError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "lbn-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchLoadBalancerNetwork("lbn-abcdef12", gomock.Any()).Return(resp, nil),
			service.EXPECT().GetLoadBalancerNetwork("lbn-abcdef12").Return(ecloud.LoadBalancerNetwork{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated load balancer network [lbn-abcdef12]: test error\n", func() {
			ecloudLoadBalancerNetworkUpdate(service, &cobra.Command{}, []string{"lbn-abcdef12"})
		})
	})
}

func Test_ecloudLoadBalancerNetworkDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudLoadBalancerNetworkDeleteCmd(nil).Args(nil, []string{"lbn-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudLoadBalancerNetworkDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing load balancer network", err.Error())
	})
}

func Test_ecloudLoadBalancerNetworkDelete(t *testing.T) {
	t.Run("SingleRule", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteLoadBalancerNetwork("lbn-abcdef12").Return("task-abcdef12", nil).Times(1)

		ecloudLoadBalancerNetworkDelete(service, &cobra.Command{}, []string{"lbn-abcdef12"})
	})

	t.Run("MultipleLoadBalancerNetworks", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteLoadBalancerNetwork("lbn-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().DeleteLoadBalancerNetwork("lbn-12abcdef").Return("task-abcdef23", nil),
		)

		ecloudLoadBalancerNetworkDelete(service, &cobra.Command{}, []string{"lbn-abcdef12", "lbn-12abcdef"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudLoadBalancerNetworkDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteLoadBalancerNetwork("lbn-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudLoadBalancerNetworkDelete(service, cmd, []string{"lbn-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudLoadBalancerNetworkDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		gomock.InOrder(
			service.EXPECT().DeleteLoadBalancerNetwork("lbn-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for load balancer network [lbn-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudLoadBalancerNetworkDelete(service, cmd, []string{"lbn-abcdef12"})
		})
	})

	t.Run("DeleteLoadBalancerNetworkError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteLoadBalancerNetwork("lbn-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing load balancer network [lbn-abcdef12]: test error\n", func() {
			ecloudLoadBalancerNetworkDelete(service, &cobra.Command{}, []string{"lbn-abcdef12"})
		})
	})
}
