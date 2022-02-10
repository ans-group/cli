package ecloud

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func Test_ecloudLoadBalancerLoadBalancerNetworkListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudLoadBalancerLoadBalancerNetworkListCmd(nil).Args(nil, []string{"lb-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudLoadBalancerLoadBalancerNetworkListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing load balancer", err.Error())
	})
}

func Test_ecloudLoadBalancerLoadBalancerNetworkList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetLoadBalancerLoadBalancerNetworks("lb-abcdef12", gomock.Any()).Return([]ecloud.LoadBalancerNetwork{}, nil).Times(1)

		ecloudLoadBalancerLoadBalancerNetworkList(service, &cobra.Command{}, []string{"lb-abcdef12"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudLoadBalancerLoadBalancerNetworkList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetNetworkRulesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetLoadBalancerLoadBalancerNetworks("lb-abcdef12", gomock.Any()).Return([]ecloud.LoadBalancerNetwork{}, errors.New("test error")).Times(1)

		err := ecloudLoadBalancerLoadBalancerNetworkList(service, &cobra.Command{}, []string{"lb-abcdef12"})

		assert.Equal(t, "Error retrieving load balancer networks: test error", err.Error())
	})
}
