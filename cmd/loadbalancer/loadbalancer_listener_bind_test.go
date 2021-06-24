package loadbalancer

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/loadbalancer"
)

func Test_loadbalancerListenerBindList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetListenerBinds(123, gomock.Any()).Return([]loadbalancer.Bind{}, nil).Times(1)

		loadbalancerListenerBindList(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerListenerBindListCmd(nil)
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := loadbalancerListenerBindList(service, cmd, []string{"123"})

		assert.NotNil(t, err)
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetListenerBindsError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetListenerBinds(123, gomock.Any()).Return([]loadbalancer.Bind{}, errors.New("test error")).Times(1)

		err := loadbalancerListenerBindList(service, &cobra.Command{}, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving binds: test error", err.Error())
	})
}

func Test_loadbalancerListenerBindShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerListenerBindShowCmd(nil).Args(nil, []string{"123", "456"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerListenerBindShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing listener", err.Error())
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerListenerBindShowCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing bind", err.Error())
	})
}

func Test_loadbalancerListenerBindShow(t *testing.T) {
	t.Run("SingleListenerBind", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetListenerBind(123, 456).Return(loadbalancer.Bind{}, nil).Times(1)

		loadbalancerListenerBindShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("MultipleListenerBinds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetListenerBind(123, 456).Return(loadbalancer.Bind{}, nil),
			service.EXPECT().GetListenerBind(123, 789).Return(loadbalancer.Bind{}, nil),
		)

		loadbalancerListenerBindShow(service, &cobra.Command{}, []string{"123", "456", "789"})
	})

	t.Run("GetListenerID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		err := loadbalancerListenerBindShow(service, &cobra.Command{}, []string{"abc", "456"})

		assert.Equal(t, "Invalid listener ID", err.Error())
	})

	t.Run("GetBindID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid bind ID [abc]\n", func() {
			loadbalancerListenerBindShow(service, &cobra.Command{}, []string{"123", "abc"})
		})
	})

	t.Run("GetListenerBindError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetListenerBind(123, 456).Return(loadbalancer.Bind{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving bind [456]: test error\n", func() {
			loadbalancerListenerBindShow(service, &cobra.Command{}, []string{"123", "456"})
		})
	})
}

func Test_loadbalancerListenerBindCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerListenerBindCreateCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerListenerBindCreateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing listener", err.Error())
	})
}

func Test_loadbalancerListenerBindCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerListenerBindCreateCmd(nil)
		cmd.ParseFlags([]string{"--vip=1", "--port=443"})

		req := loadbalancer.CreateBindRequest{
			VIPID: 1,
			Port:  443,
		}

		gomock.InOrder(
			service.EXPECT().CreateListenerBind(123, req).Return(456, nil),
			service.EXPECT().GetListenerBind(123, 456).Return(loadbalancer.Bind{}, nil),
		)

		loadbalancerListenerBindCreate(service, cmd, []string{"123"})
	})

	t.Run("CreateListenerError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerListenerBindCreateCmd(nil)
		cmd.ParseFlags([]string{"--vip=1", "--port=443"})

		service.EXPECT().CreateListenerBind(123, gomock.Any()).Return(0, errors.New("test error"))

		err := loadbalancerListenerBindCreate(service, cmd, []string{"123"})

		assert.Equal(t, "Error creating bind: test error", err.Error())
	})

	t.Run("GetListenerError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerListenerBindCreateCmd(nil)
		cmd.ParseFlags([]string{"--vip=1", "--port=443"})

		gomock.InOrder(
			service.EXPECT().CreateListenerBind(123, gomock.Any()).Return(456, nil),
			service.EXPECT().GetListenerBind(123, 456).Return(loadbalancer.Bind{}, errors.New("test error")),
		)

		err := loadbalancerListenerBindCreate(service, cmd, []string{"123"})

		assert.Equal(t, "Error retrieving new bind: test error", err.Error())
	})
}

func Test_loadbalancerListenerBindUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerListenerBindUpdateCmd(nil).Args(nil, []string{"123", "456"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerListenerBindUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing listener", err.Error())
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerListenerBindUpdateCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing bind", err.Error())
	})
}

func Test_loadbalancerListenerBindUpdate(t *testing.T) {
	t.Run("SingleBind", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		cmd := loadbalancerListenerBindUpdateCmd(nil)
		cmd.ParseFlags([]string{"--port=80"})

		req := loadbalancer.PatchBindRequest{
			Port: 80,
		}

		gomock.InOrder(
			service.EXPECT().PatchListenerBind(123, 456, req).Return(nil),
			service.EXPECT().GetListenerBind(123, 456).Return(loadbalancer.Bind{}, nil),
		)

		loadbalancerListenerBindUpdate(service, cmd, []string{"123", "456"})
	})

	t.Run("PatchListenerError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().PatchListenerBind(123, 456, gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating bind [456]: test error\n", func() {
			loadbalancerListenerBindUpdate(service, &cobra.Command{}, []string{"123", "456"})
		})
	})

	t.Run("GetListenerError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchListenerBind(123, 456, gomock.Any()).Return(nil),
			service.EXPECT().GetListenerBind(123, 456).Return(loadbalancer.Bind{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated bind [456]: test error\n", func() {
			loadbalancerListenerBindUpdate(service, &cobra.Command{}, []string{"123", "456"})
		})
	})
}

func Test_loadbalancerListenerBindDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerListenerBindDeleteCmd(nil).Args(nil, []string{"123", "456"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerListenerBindUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing listener", err.Error())
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerListenerBindUpdateCmd(nil).Args(nil, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing bind", err.Error())
	})
}

func Test_loadbalancerListenerBindDelete(t *testing.T) {
	t.Run("SingleBind", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().DeleteListenerBind(123, 456).Return(nil).Times(1)

		loadbalancerListenerBindDelete(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("MultipleBinds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteListenerBind(1, 12).Return(nil),
			service.EXPECT().DeleteListenerBind(1, 123).Return(nil),
		)

		loadbalancerListenerBindDelete(service, &cobra.Command{}, []string{"1", "12", "123"})
	})

	t.Run("DeleteListenerError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().DeleteListenerBind(123, 456).Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing bind [456]: test error\n", func() {
			loadbalancerListenerBindDelete(service, &cobra.Command{}, []string{"123", "456"})
		})
	})
}
