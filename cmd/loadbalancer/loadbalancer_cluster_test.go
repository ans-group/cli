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
	"github.com/ukfast/sdk-go/pkg/ptr"
	"github.com/ukfast/sdk-go/pkg/service/loadbalancer"
)

func Test_loadbalancerClusterList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetClusters(gomock.Any()).Return([]loadbalancer.Cluster{}, nil).Times(1)

		loadbalancerClusterList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := loadbalancerClusterList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetClustersError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetClusters(gomock.Any()).Return([]loadbalancer.Cluster{}, errors.New("test error")).Times(1)

		err := loadbalancerClusterList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving clusters: test error", err.Error())
	})
}

func Test_loadbalancerClusterShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerClusterShowCmd(nil).Args(nil, []string{"rtr-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerClusterShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing cluster", err.Error())
	})
}

func Test_loadbalancerClusterShow(t *testing.T) {
	t.Run("SingleCluster", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetCluster("rtr-abcdef12").Return(loadbalancer.Cluster{}, nil).Times(1)

		loadbalancerClusterShow(service, &cobra.Command{}, []string{"rtr-abcdef12"})
	})

	t.Run("MultipleClusters", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetCluster("rtr-abcdef12").Return(loadbalancer.Cluster{}, nil),
			service.EXPECT().GetCluster("rtr-abcdef23").Return(loadbalancer.Cluster{}, nil),
		)

		loadbalancerClusterShow(service, &cobra.Command{}, []string{"rtr-abcdef12", "rtr-abcdef23"})
	})

	t.Run("GetClusterError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetCluster("rtr-abcdef12").Return(loadbalancer.Cluster{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving cluster [rtr-abcdef12]: test error\n", func() {
			loadbalancerClusterShow(service, &cobra.Command{}, []string{"rtr-abcdef12"})
		})
	})
}

func Test_loadbalancerClusterUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerClusterUpdateCmd(nil).Args(nil, []string{"rtr-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerClusterUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing cluster", err.Error())
	})
}

func Test_loadbalancerClusterUpdate(t *testing.T) {
	t.Run("SingleCluster", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		cmd := loadbalancerClusterUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testcluster"})

		req := loadbalancer.PatchClusterRequest{
			Name: ptr.String("testcluster"),
		}

		gomock.InOrder(
			service.EXPECT().PatchCluster("rtr-abcdef12", req).Return(nil),
			service.EXPECT().GetCluster("rtr-abcdef12").Return(loadbalancer.Cluster{}, nil),
		)

		loadbalancerClusterUpdate(service, cmd, []string{"rtr-abcdef12"})
	})

	t.Run("MultipleClusters", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchCluster("rtr-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetCluster("rtr-abcdef12").Return(loadbalancer.Cluster{}, nil),
			service.EXPECT().PatchCluster("rtr-12abcdef", gomock.Any()).Return(nil),
			service.EXPECT().GetCluster("rtr-12abcdef").Return(loadbalancer.Cluster{}, nil),
		)

		loadbalancerClusterUpdate(service, &cobra.Command{}, []string{"rtr-abcdef12", "rtr-12abcdef"})
	})

	t.Run("PatchClusterError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().PatchCluster("rtr-abcdef12", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating cluster [rtr-abcdef12]: test error\n", func() {
			loadbalancerClusterUpdate(service, &cobra.Command{}, []string{"rtr-abcdef12"})
		})
	})

	t.Run("GetClusterError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchCluster("rtr-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetCluster("rtr-abcdef12").Return(loadbalancer.Cluster{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated cluster [rtr-abcdef12]: test error\n", func() {
			loadbalancerClusterUpdate(service, &cobra.Command{}, []string{"rtr-abcdef12"})
		})
	})
}

func Test_loadbalancerClusterDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerClusterDeleteCmd(nil).Args(nil, []string{"rtr-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerClusterDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing cluster", err.Error())
	})
}

func Test_loadbalancerClusterDelete(t *testing.T) {
	t.Run("SingleCluster", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().DeleteCluster("rtr-abcdef12").Return(nil).Times(1)

		loadbalancerClusterDelete(service, &cobra.Command{}, []string{"rtr-abcdef12"})
	})

	t.Run("MultipleClusters", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteCluster("rtr-abcdef12").Return(nil),
			service.EXPECT().DeleteCluster("rtr-12abcdef").Return(nil),
		)

		loadbalancerClusterDelete(service, &cobra.Command{}, []string{"rtr-abcdef12", "rtr-12abcdef"})
	})

	t.Run("DeleteClusterError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().DeleteCluster("rtr-abcdef12").Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing cluster [rtr-abcdef12]: test error\n", func() {
			loadbalancerClusterDelete(service, &cobra.Command{}, []string{"rtr-abcdef12"})
		})
	})
}
