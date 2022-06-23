package loadbalancer

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_loadbalancerTargetGroupList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetTargetGroups(gomock.Any()).Return([]loadbalancer.TargetGroup{}, nil).Times(1)

		loadbalancerTargetGroupList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerTargetGroupListCmd(nil)
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := loadbalancerTargetGroupList(service, cmd, []string{})

		assert.NotNil(t, err)
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetTargetGroupsError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetTargetGroups(gomock.Any()).Return([]loadbalancer.TargetGroup{}, errors.New("test error")).Times(1)

		err := loadbalancerTargetGroupList(service, &cobra.Command{}, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving target groups: test error", err.Error())
	})
}

func Test_loadbalancerTargetGroupShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerTargetGroupShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerTargetGroupShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing target group", err.Error())
	})
}

func Test_loadbalancerTargetGroupShow(t *testing.T) {
	t.Run("SingleTargetGroup", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetTargetGroup(123).Return(loadbalancer.TargetGroup{}, nil).Times(1)

		loadbalancerTargetGroupShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleTargetGroups", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetTargetGroup(123).Return(loadbalancer.TargetGroup{}, nil),
			service.EXPECT().GetTargetGroup(456).Return(loadbalancer.TargetGroup{}, nil),
		)

		loadbalancerTargetGroupShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetTargetGroupID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid target group ID [abc]\n", func() {
			loadbalancerTargetGroupShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetTargetGroupError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetTargetGroup(123).Return(loadbalancer.TargetGroup{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving target group [123]: test error\n", func() {
			loadbalancerTargetGroupShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_loadbalancerTargetGroupCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerTargetGroupCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgroup", "--balance=roundrobin", "--mode=http"})

		req := loadbalancer.CreateTargetGroupRequest{
			Name:    "testgroup",
			Balance: loadbalancer.TargetGroupBalanceRoundRobin,
			Mode:    loadbalancer.ModeHTTP,
		}

		gomock.InOrder(
			service.EXPECT().CreateTargetGroup(req).Return(123, nil),
			service.EXPECT().GetTargetGroup(123).Return(loadbalancer.TargetGroup{}, nil),
		)

		loadbalancerTargetGroupCreate(service, cmd, []string{})
	})

	t.Run("CreateTargetGroupError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerTargetGroupCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgroup", "--balance=roundrobin", "--mode=http"})

		service.EXPECT().CreateTargetGroup(gomock.Any()).Return(0, errors.New("test error"))

		err := loadbalancerTargetGroupCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating target group: test error", err.Error())
	})

	t.Run("GetTargetGroupError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerTargetGroupCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgroup", "--balance=roundrobin", "--mode=http"})

		gomock.InOrder(
			service.EXPECT().CreateTargetGroup(gomock.Any()).Return(123, nil),
			service.EXPECT().GetTargetGroup(123).Return(loadbalancer.TargetGroup{}, errors.New("test error")),
		)

		err := loadbalancerTargetGroupCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new target group: test error", err.Error())
	})
}

func Test_loadbalancerTargetGroupUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerTargetGroupUpdateCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerTargetGroupUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing target group", err.Error())
	})
}

func Test_loadbalancerTargetGroupUpdate(t *testing.T) {
	t.Run("SingleGroup", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		cmd := loadbalancerTargetGroupUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgroup"})

		req := loadbalancer.PatchTargetGroupRequest{
			Name: "testgroup",
		}

		gomock.InOrder(
			service.EXPECT().PatchTargetGroup(123, req).Return(nil),
			service.EXPECT().GetTargetGroup(123).Return(loadbalancer.TargetGroup{}, nil),
		)

		loadbalancerTargetGroupUpdate(service, cmd, []string{"123"})
	})

	t.Run("MultipleTargetGroups", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchTargetGroup(123, gomock.Any()).Return(nil),
			service.EXPECT().GetTargetGroup(123).Return(loadbalancer.TargetGroup{}, nil),
			service.EXPECT().PatchTargetGroup(456, gomock.Any()).Return(nil),
			service.EXPECT().GetTargetGroup(456).Return(loadbalancer.TargetGroup{}, nil),
		)

		loadbalancerTargetGroupUpdate(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("PatchTargetGroupError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().PatchTargetGroup(123, gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating target group [123]: test error\n", func() {
			loadbalancerTargetGroupUpdate(service, &cobra.Command{}, []string{"123"})
		})
	})

	t.Run("GetTargetGroupError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchTargetGroup(123, gomock.Any()).Return(nil),
			service.EXPECT().GetTargetGroup(123).Return(loadbalancer.TargetGroup{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated target group [123]: test error\n", func() {
			loadbalancerTargetGroupUpdate(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_loadbalancerTargetGroupDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerTargetGroupDeleteCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerTargetGroupDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing target group", err.Error())
	})
}

func Test_loadbalancerTargetGroupDelete(t *testing.T) {
	t.Run("SingleGroup", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().DeleteTargetGroup(123).Return(nil).Times(1)

		loadbalancerTargetGroupDelete(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleTargetGroups", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteTargetGroup(123).Return(nil),
			service.EXPECT().DeleteTargetGroup(456).Return(nil),
		)

		loadbalancerTargetGroupDelete(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("DeleteTargetGroupError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().DeleteTargetGroup(123).Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing target group [123]: test error\n", func() {
			loadbalancerTargetGroupDelete(service, &cobra.Command{}, []string{"123"})
		})
	})
}
