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

func Test_loadbalancerListenerList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetListeners(gomock.Any()).Return([]loadbalancer.Listener{}, nil).Times(1)

		loadbalancerListenerList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerListenerListCmd(nil)
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := loadbalancerListenerList(service, cmd, []string{})

		assert.NotNil(t, err)
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetListenersError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetListeners(gomock.Any()).Return([]loadbalancer.Listener{}, errors.New("test error")).Times(1)

		err := loadbalancerListenerList(service, &cobra.Command{}, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "error retrieving listeners: test error", err.Error())
	})
}

func Test_loadbalancerListenerShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerListenerShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerListenerShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing listener", err.Error())
	})
}

func Test_loadbalancerListenerShow(t *testing.T) {
	t.Run("SingleListener", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetListener(123).Return(loadbalancer.Listener{}, nil).Times(1)

		loadbalancerListenerShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleListeners", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetListener(123).Return(loadbalancer.Listener{}, nil),
			service.EXPECT().GetListener(456).Return(loadbalancer.Listener{}, nil),
		)

		loadbalancerListenerShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetListenerID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid listener ID [abc]\n", func() {
			loadbalancerListenerShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetListenerError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetListener(123).Return(loadbalancer.Listener{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving listener [123]: test error\n", func() {
			loadbalancerListenerShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_loadbalancerListenerCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerListenerCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testlistener", "--mode=http"})

		req := loadbalancer.CreateListenerRequest{
			Name: "testlistener",
			Mode: loadbalancer.ModeHTTP,
		}

		gomock.InOrder(
			service.EXPECT().CreateListener(req).Return(123, nil),
			service.EXPECT().GetListener(123).Return(loadbalancer.Listener{}, nil),
		)

		loadbalancerListenerCreate(service, cmd, []string{})
	})

	t.Run("CreateListenerError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerListenerCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testlistener", "--mode=http"})

		service.EXPECT().CreateListener(gomock.Any()).Return(0, errors.New("test error"))

		err := loadbalancerListenerCreate(service, cmd, []string{})

		assert.Equal(t, "error creating listener: test error", err.Error())
	})

	t.Run("GetListenerError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerListenerCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testlistener", "--mode=http"})

		gomock.InOrder(
			service.EXPECT().CreateListener(gomock.Any()).Return(123, nil),
			service.EXPECT().GetListener(123).Return(loadbalancer.Listener{}, errors.New("test error")),
		)

		err := loadbalancerListenerCreate(service, cmd, []string{})

		assert.Equal(t, "error retrieving new listener: test error", err.Error())
	})
}

func Test_loadbalancerListenerUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerListenerUpdateCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerListenerUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing listener", err.Error())
	})
}

func Test_loadbalancerListenerUpdate(t *testing.T) {
	t.Run("SingleGroup", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		cmd := loadbalancerListenerUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testlistener"})

		req := loadbalancer.PatchListenerRequest{
			Name: "testlistener",
		}

		gomock.InOrder(
			service.EXPECT().PatchListener(123, req).Return(nil),
			service.EXPECT().GetListener(123).Return(loadbalancer.Listener{}, nil),
		)

		loadbalancerListenerUpdate(service, cmd, []string{"123"})
	})

	t.Run("MultipleListeners", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchListener(123, gomock.Any()).Return(nil),
			service.EXPECT().GetListener(123).Return(loadbalancer.Listener{}, nil),
			service.EXPECT().PatchListener(456, gomock.Any()).Return(nil),
			service.EXPECT().GetListener(456).Return(loadbalancer.Listener{}, nil),
		)

		loadbalancerListenerUpdate(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("PatchListenerError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().PatchListener(123, gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating listener [123]: test error\n", func() {
			loadbalancerListenerUpdate(service, &cobra.Command{}, []string{"123"})
		})
	})

	t.Run("GetListenerError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchListener(123, gomock.Any()).Return(nil),
			service.EXPECT().GetListener(123).Return(loadbalancer.Listener{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated listener [123]: test error\n", func() {
			loadbalancerListenerUpdate(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_loadbalancerListenerDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerListenerDeleteCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerListenerDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing listener", err.Error())
	})
}

func Test_loadbalancerListenerDelete(t *testing.T) {
	t.Run("SingleGroup", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().DeleteListener(123).Return(nil).Times(1)

		loadbalancerListenerDelete(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleListeners", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteListener(123).Return(nil),
			service.EXPECT().DeleteListener(456).Return(nil),
		)

		loadbalancerListenerDelete(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("DeleteListenerError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().DeleteListener(123).Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing listener [123]: test error\n", func() {
			loadbalancerListenerDelete(service, &cobra.Command{}, []string{"123"})
		})
	})
}
