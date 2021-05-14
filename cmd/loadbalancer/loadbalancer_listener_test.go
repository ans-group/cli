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
		assert.Equal(t, "Error retrieving listeners: test error", err.Error())
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
		assert.Equal(t, "Missing listener", err.Error())
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
