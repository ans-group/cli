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

func Test_ecloudLoadBalancerList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetLoadBalancers(gomock.Any()).Return([]ecloud.LoadBalancer{}, nil).Times(1)

		ecloudLoadBalancerList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudLoadBalancerList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetLoadBalancersError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetLoadBalancers(gomock.Any()).Return([]ecloud.LoadBalancer{}, errors.New("test error")).Times(1)

		err := ecloudLoadBalancerList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving load balancers: test error", err.Error())
	})
}

func Test_ecloudLoadBalancerShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudLoadBalancerShowCmd(nil).Args(nil, []string{"lb-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudLoadBalancerShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing load balancer", err.Error())
	})
}

func Test_ecloudLoadBalancerShow(t *testing.T) {
	t.Run("SingleLoadBalancer", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetLoadBalancer("lb-abcdef12").Return(ecloud.LoadBalancer{}, nil).Times(1)

		ecloudLoadBalancerShow(service, &cobra.Command{}, []string{"lb-abcdef12"})
	})

	t.Run("MultipleLoadBalancers", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetLoadBalancer("lb-abcdef12").Return(ecloud.LoadBalancer{}, nil),
			service.EXPECT().GetLoadBalancer("lb-abcdef23").Return(ecloud.LoadBalancer{}, nil),
		)

		ecloudLoadBalancerShow(service, &cobra.Command{}, []string{"lb-abcdef12", "lb-abcdef23"})
	})

	t.Run("GetLoadBalancerError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetLoadBalancer("lb-abcdef12").Return(ecloud.LoadBalancer{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving load balancer [lb-abcdef12]: test error\n", func() {
			ecloudLoadBalancerShow(service, &cobra.Command{}, []string{"lb-abcdef12"})
		})
	})
}

func Test_ecloudLoadBalancerCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudLoadBalancerCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testlb"})

		req := ecloud.CreateLoadBalancerRequest{
			Name: "testlb",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "lb-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateLoadBalancer(req).Return(resp, nil),
			service.EXPECT().GetLoadBalancer("lb-abcdef12").Return(ecloud.LoadBalancer{}, nil),
		)

		ecloudLoadBalancerCreate(service, cmd, []string{})
	})

	t.Run("CreateWithWaitFlagNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudLoadBalancerCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testlb", "--wait"})

		req := ecloud.CreateLoadBalancerRequest{
			Name: "testlb",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "lb-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateLoadBalancer(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetLoadBalancer("lb-abcdef12").Return(ecloud.LoadBalancer{}, nil),
		)

		ecloudLoadBalancerCreate(service, cmd, []string{})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudLoadBalancerCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testlb", "--wait"})

		req := ecloud.CreateLoadBalancerRequest{
			Name: "testlb",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "lb-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateLoadBalancer(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		err := ecloudLoadBalancerCreate(service, cmd, []string{})
		assert.Equal(t, "Error waiting for load balancer task to complete: Error waiting for command: Failed to retrieve task status: test error", err.Error())
	})

	t.Run("CreateLoadBalancerError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudLoadBalancerCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testlb"})

		service.EXPECT().CreateLoadBalancer(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error")).Times(1)

		err := ecloudLoadBalancerCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating load balancer: test error", err.Error())
	})

	t.Run("GetLoadBalancerError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudLoadBalancerCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testlb"})

		gomock.InOrder(
			service.EXPECT().CreateLoadBalancer(gomock.Any()).Return(ecloud.TaskReference{ResourceID: "lb-abcdef12"}, nil),
			service.EXPECT().GetLoadBalancer("lb-abcdef12").Return(ecloud.LoadBalancer{}, errors.New("test error")),
		)

		err := ecloudLoadBalancerCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new load balancer: test error", err.Error())
	})
}

func Test_ecloudLoadBalancerUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudLoadBalancerUpdateCmd(nil).Args(nil, []string{"lb-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudLoadBalancerUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing load balancer", err.Error())
	})
}

func Test_ecloudLoadBalancerUpdate(t *testing.T) {
	t.Run("SingleRule", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudLoadBalancerUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testlb"})

		req := ecloud.PatchLoadBalancerRequest{
			Name: "testlb",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "lb-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchLoadBalancer("lb-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetLoadBalancer("lb-abcdef12").Return(ecloud.LoadBalancer{}, nil),
		)

		ecloudLoadBalancerUpdate(service, cmd, []string{"lb-abcdef12"})
	})

	t.Run("MultipleLoadBalancers", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp1 := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "lb-abcdef12",
		}

		resp2 := ecloud.TaskReference{
			TaskID:     "task-abcdef23",
			ResourceID: "lb-12abcdef",
		}

		gomock.InOrder(
			service.EXPECT().PatchLoadBalancer("lb-abcdef12", gomock.Any()).Return(resp1, nil),
			service.EXPECT().GetLoadBalancer("lb-abcdef12").Return(ecloud.LoadBalancer{}, nil),
			service.EXPECT().PatchLoadBalancer("lb-12abcdef", gomock.Any()).Return(resp2, nil),
			service.EXPECT().GetLoadBalancer("lb-12abcdef").Return(ecloud.LoadBalancer{}, nil),
		)

		ecloudLoadBalancerUpdate(service, &cobra.Command{}, []string{"lb-abcdef12", "lb-12abcdef"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudLoadBalancerUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testlb", "--wait"})

		req := ecloud.PatchLoadBalancerRequest{
			Name: "testlb",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "lb-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchLoadBalancer("lb-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetLoadBalancer("lb-abcdef12").Return(ecloud.LoadBalancer{}, nil),
		)

		ecloudLoadBalancerUpdate(service, cmd, []string{"lb-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudLoadBalancerUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testlb", "--wait"})

		req := ecloud.PatchLoadBalancerRequest{
			Name: "testlb",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "lb-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchLoadBalancer("lb-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for load balancer [lb-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudLoadBalancerUpdate(service, cmd, []string{"lb-abcdef12"})
		})
	})

	t.Run("PatchLoadBalancerError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchLoadBalancer("lb-abcdef12", gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating load balancer [lb-abcdef12]: test error\n", func() {
			ecloudLoadBalancerUpdate(service, &cobra.Command{}, []string{"lb-abcdef12"})
		})
	})

	t.Run("GetLoadBalancerError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "lb-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchLoadBalancer("lb-abcdef12", gomock.Any()).Return(resp, nil),
			service.EXPECT().GetLoadBalancer("lb-abcdef12").Return(ecloud.LoadBalancer{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated load balancer [lb-abcdef12]: test error\n", func() {
			ecloudLoadBalancerUpdate(service, &cobra.Command{}, []string{"lb-abcdef12"})
		})
	})
}

func Test_ecloudLoadBalancerDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudLoadBalancerDeleteCmd(nil).Args(nil, []string{"lb-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudLoadBalancerDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing load balancer", err.Error())
	})
}

func Test_ecloudLoadBalancerDelete(t *testing.T) {
	t.Run("SingleRule", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteLoadBalancer("lb-abcdef12").Return("task-abcdef12", nil).Times(1)

		ecloudLoadBalancerDelete(service, &cobra.Command{}, []string{"lb-abcdef12"})
	})

	t.Run("MultipleLoadBalancers", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteLoadBalancer("lb-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().DeleteLoadBalancer("lb-12abcdef").Return("task-abcdef23", nil),
		)

		ecloudLoadBalancerDelete(service, &cobra.Command{}, []string{"lb-abcdef12", "lb-12abcdef"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudLoadBalancerDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteLoadBalancer("lb-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudLoadBalancerDelete(service, cmd, []string{"lb-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudLoadBalancerDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		gomock.InOrder(
			service.EXPECT().DeleteLoadBalancer("lb-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for load balancer [lb-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudLoadBalancerDelete(service, cmd, []string{"lb-abcdef12"})
		})
	})

	t.Run("DeleteLoadBalancerError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteLoadBalancer("lb-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing load balancer [lb-abcdef12]: test error\n", func() {
			ecloudLoadBalancerDelete(service, &cobra.Command{}, []string{"lb-abcdef12"})
		})
	})
}
