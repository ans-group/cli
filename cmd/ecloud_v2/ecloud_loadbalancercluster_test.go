package ecloud_v2

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

func Test_ecloudLoadBalancerClusterList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetLoadBalancerClusters(gomock.Any()).Return([]ecloud.LoadBalancerCluster{}, nil).Times(1)

		ecloudLoadBalancerClusterList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudLoadBalancerClusterList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetLoadBalancerClustersError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetLoadBalancerClusters(gomock.Any()).Return([]ecloud.LoadBalancerCluster{}, errors.New("test error")).Times(1)

		err := ecloudLoadBalancerClusterList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving load balancer clusters: test error", err.Error())
	})
}

func Test_ecloudLoadBalancerClusterShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudLoadBalancerClusterShowCmd(nil).Args(nil, []string{"lbcs-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudLoadBalancerClusterShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing load balancer cluster", err.Error())
	})
}

func Test_ecloudLoadBalancerClusterShow(t *testing.T) {
	t.Run("SingleLoadBalancerCluster", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetLoadBalancerCluster("lbcs-abcdef12").Return(ecloud.LoadBalancerCluster{}, nil).Times(1)

		ecloudLoadBalancerClusterShow(service, &cobra.Command{}, []string{"lbcs-abcdef12"})
	})

	t.Run("MultipleLoadBalancerClusters", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetLoadBalancerCluster("lbcs-abcdef12").Return(ecloud.LoadBalancerCluster{}, nil),
			service.EXPECT().GetLoadBalancerCluster("lbcs-abcdef23").Return(ecloud.LoadBalancerCluster{}, nil),
		)

		ecloudLoadBalancerClusterShow(service, &cobra.Command{}, []string{"lbcs-abcdef12", "lbcs-abcdef23"})
	})

	t.Run("GetLoadBalancerClusterError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetLoadBalancerCluster("lbcs-abcdef12").Return(ecloud.LoadBalancerCluster{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving load balancer cluster [lbcs-abcdef12]: test error\n", func() {
			ecloudLoadBalancerClusterShow(service, &cobra.Command{}, []string{"lbcs-abcdef12"})
		})
	})
}
