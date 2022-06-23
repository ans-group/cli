package loadbalancer

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_loadbalancerTargetGroupTargetListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerTargetGroupTargetListCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerTargetGroupTargetListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing target group", err.Error())
	})
}

func Test_loadbalancerTargetGroupTargetList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetTargetGroupTargets(123, gomock.Any()).Return([]loadbalancer.Target{}, nil).Times(1)

		loadbalancerTargetGroupTargetList(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerTargetGroupTargetListCmd(nil)
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := loadbalancerTargetGroupTargetList(service, cmd, []string{"123"})

		assert.NotNil(t, err)
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetTargetGroupTargetsError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetTargetGroupTargets(123, gomock.Any()).Return([]loadbalancer.Target{}, errors.New("test error")).Times(1)

		err := loadbalancerTargetGroupTargetList(service, &cobra.Command{}, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving targets: test error", err.Error())
	})
}

func Test_loadbalancerTargetGroupTargetShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerTargetGroupTargetShowCmd(nil).Args(nil, []string{"123", "456"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerTargetGroupTargetShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing target group", err.Error())
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerTargetGroupTargetShowCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing target", err.Error())
	})
}

func Test_loadbalancerTargetGroupTargetShow(t *testing.T) {
	t.Run("SingleTargetGroupTarget", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetTargetGroupTarget(123, 456).Return(loadbalancer.Target{}, nil).Times(1)

		loadbalancerTargetGroupTargetShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("MultipleTargetGroupTargets", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetTargetGroupTarget(123, 456).Return(loadbalancer.Target{}, nil),
			service.EXPECT().GetTargetGroupTarget(123, 789).Return(loadbalancer.Target{}, nil),
		)

		loadbalancerTargetGroupTargetShow(service, &cobra.Command{}, []string{"123", "456", "789"})
	})

	t.Run("GetTargetGroupID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		err := loadbalancerTargetGroupTargetShow(service, &cobra.Command{}, []string{"abc", "456"})

		assert.Equal(t, "Invalid target group ID", err.Error())
	})

	t.Run("GetTargetID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid target ID [abc]\n", func() {
			loadbalancerTargetGroupTargetShow(service, &cobra.Command{}, []string{"123", "abc"})
		})
	})

	t.Run("GetTargetGroupTargetError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetTargetGroupTarget(123, 456).Return(loadbalancer.Target{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving target [456]: test error\n", func() {
			loadbalancerTargetGroupTargetShow(service, &cobra.Command{}, []string{"123", "456"})
		})
	})
}

func Test_loadbalancerTargetGroupTargetCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerTargetGroupTargetCreateCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerTargetGroupTargetCreateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing target group", err.Error())
	})
}

func Test_loadbalancerTargetGroupTargetCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerTargetGroupTargetCreateCmd(nil)
		cmd.ParseFlags([]string{"--ip=1.2.3.4", "--port=443"})

		req := loadbalancer.CreateTargetRequest{
			IP:     connection.IPAddress("1.2.3.4"),
			Port:   443,
			Active: true,
		}

		gomock.InOrder(
			service.EXPECT().CreateTargetGroupTarget(123, req).Return(456, nil),
			service.EXPECT().GetTargetGroupTarget(123, 456).Return(loadbalancer.Target{}, nil),
		)

		loadbalancerTargetGroupTargetCreate(service, cmd, []string{"123"})
	})

	t.Run("CreateTargetGroupError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerTargetGroupTargetCreateCmd(nil)
		cmd.ParseFlags([]string{"--ip=1.2.3.4", "--port=443"})

		service.EXPECT().CreateTargetGroupTarget(123, gomock.Any()).Return(0, errors.New("test error"))

		err := loadbalancerTargetGroupTargetCreate(service, cmd, []string{"123"})

		assert.Equal(t, "Error creating target: test error", err.Error())
	})

	t.Run("GetTargetGroupError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerTargetGroupTargetCreateCmd(nil)
		cmd.ParseFlags([]string{"--ip=1.2.3.4", "--port=443"})

		gomock.InOrder(
			service.EXPECT().CreateTargetGroupTarget(123, gomock.Any()).Return(456, nil),
			service.EXPECT().GetTargetGroupTarget(123, 456).Return(loadbalancer.Target{}, errors.New("test error")),
		)

		err := loadbalancerTargetGroupTargetCreate(service, cmd, []string{"123"})

		assert.Equal(t, "Error retrieving new target: test error", err.Error())
	})
}

func Test_loadbalancerTargetGroupTargetUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerTargetGroupTargetUpdateCmd(nil).Args(nil, []string{"123", "456"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerTargetGroupTargetUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing target group", err.Error())
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerTargetGroupTargetUpdateCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing target", err.Error())
	})
}

func Test_loadbalancerTargetGroupTargetUpdate(t *testing.T) {
	t.Run("SingleTarget", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		cmd := loadbalancerTargetGroupTargetUpdateCmd(nil)
		cmd.ParseFlags([]string{"--port=443"})

		req := loadbalancer.PatchTargetRequest{
			Port: 443,
		}

		gomock.InOrder(
			service.EXPECT().PatchTargetGroupTarget(123, 456, req).Return(nil),
			service.EXPECT().GetTargetGroupTarget(123, 456).Return(loadbalancer.Target{}, nil),
		)

		loadbalancerTargetGroupTargetUpdate(service, cmd, []string{"123", "456"})
	})

	t.Run("PatchTargetGroupError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().PatchTargetGroupTarget(123, 456, gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating target [456]: test error\n", func() {
			loadbalancerTargetGroupTargetUpdate(service, &cobra.Command{}, []string{"123", "456"})
		})
	})

	t.Run("GetTargetGroupError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchTargetGroupTarget(123, 456, gomock.Any()).Return(nil),
			service.EXPECT().GetTargetGroupTarget(123, 456).Return(loadbalancer.Target{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated target [456]: test error\n", func() {
			loadbalancerTargetGroupTargetUpdate(service, &cobra.Command{}, []string{"123", "456"})
		})
	})
}

func Test_loadbalancerTargetGroupTargetDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerTargetGroupTargetDeleteCmd(nil).Args(nil, []string{"123", "456"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerTargetGroupTargetUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing target group", err.Error())
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerTargetGroupTargetUpdateCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing target", err.Error())
	})
}

func Test_loadbalancerTargetGroupTargetDelete(t *testing.T) {
	t.Run("SingleTarget", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().DeleteTargetGroupTarget(123, 456).Return(nil).Times(1)

		loadbalancerTargetGroupTargetDelete(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("MultipleTargets", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteTargetGroupTarget(1, 12).Return(nil),
			service.EXPECT().DeleteTargetGroupTarget(1, 123).Return(nil),
		)

		loadbalancerTargetGroupTargetDelete(service, &cobra.Command{}, []string{"1", "12", "123"})
	})

	t.Run("DeleteTargetGroupError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().DeleteTargetGroupTarget(123, 456).Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing target [456]: test error\n", func() {
			loadbalancerTargetGroupTargetDelete(service, &cobra.Command{}, []string{"123", "456"})
		})
	})
}
