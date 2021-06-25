package loadbalancer

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/loadbalancer"
)

func Test_loadbalancerListenerAccessIPList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetListenerAccessIPs(123, gomock.Any()).Return([]loadbalancer.AccessIP{}, nil).Times(1)

		loadbalancerListenerAccessIPList(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerListenerAccessIPListCmd(nil)
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := loadbalancerListenerAccessIPList(service, cmd, []string{"123"})

		assert.NotNil(t, err)
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetListenerAccessIPsError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetListenerAccessIPs(123, gomock.Any()).Return([]loadbalancer.AccessIP{}, errors.New("test error")).Times(1)

		err := loadbalancerListenerAccessIPList(service, &cobra.Command{}, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving access IPs: test error", err.Error())
	})
}

func Test_loadbalancerListenerAccessIPCreateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerListenerAccessIPCreateCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerListenerAccessIPCreateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing listener", err.Error())
	})
}

func Test_loadbalancerListenerAccessIPCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerListenerAccessIPCreateCmd(nil)
		cmd.ParseFlags([]string{"--ip=1.2.3.4"})

		req := loadbalancer.CreateAccessIPRequest{
			IP: "1.2.3.4",
		}

		gomock.InOrder(
			service.EXPECT().CreateListenerAccessIP(123, req).Return(456, nil),
			service.EXPECT().GetAccessIP(456).Return(loadbalancer.AccessIP{}, nil),
		)

		loadbalancerListenerAccessIPCreate(service, cmd, []string{"123"})
	})

	t.Run("CreateListenerError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerListenerAccessIPCreateCmd(nil)
		cmd.ParseFlags([]string{"--ip=1.2.3.4"})

		service.EXPECT().CreateListenerAccessIP(123, gomock.Any()).Return(0, errors.New("test error"))

		err := loadbalancerListenerAccessIPCreate(service, cmd, []string{"123"})

		assert.Equal(t, "Error creating access IP: test error", err.Error())
	})

	t.Run("GetListenerError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerListenerAccessIPCreateCmd(nil)
		cmd.ParseFlags([]string{"--ip=1.2.3.4"})

		gomock.InOrder(
			service.EXPECT().CreateListenerAccessIP(123, gomock.Any()).Return(456, nil),
			service.EXPECT().GetAccessIP(456).Return(loadbalancer.AccessIP{}, errors.New("test error")),
		)

		err := loadbalancerListenerAccessIPCreate(service, cmd, []string{"123"})

		assert.Equal(t, "Error retrieving new access IP: test error", err.Error())
	})
}
