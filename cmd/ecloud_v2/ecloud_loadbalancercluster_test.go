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
	"github.com/ukfast/sdk-go/pkg/ptr"
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

func Test_ecloudLoadBalancerClusterCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudLoadBalancerClusterCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testlbc", "--router=rtr-abcdef12"})

		req := ecloud.CreateLoadBalancerClusterRequest{
			Name: ptr.String("testlbc"),
		}

		gomock.InOrder(
			service.EXPECT().CreateLoadBalancerCluster(req).Return("lbc-abcdef12", nil),
			service.EXPECT().GetLoadBalancerCluster("lbc-abcdef12").Return(ecloud.LoadBalancerCluster{}, nil),
		)

		ecloudLoadBalancerClusterCreate(service, cmd, []string{})
	})

	t.Run("CreateLoadBalancerClusterError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudLoadBalancerClusterCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testnetwork", "--router=rtr-abcdef12"})

		service.EXPECT().CreateLoadBalancerCluster(gomock.Any()).Return("", errors.New("test error")).Times(1)

		err := ecloudLoadBalancerClusterCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating load balancer cluster: test error", err.Error())
	})

	t.Run("GetLoadBalancerClusterError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudLoadBalancerClusterCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testnetwork", "--router=rtr-abcdef12"})

		gomock.InOrder(
			service.EXPECT().CreateLoadBalancerCluster(gomock.Any()).Return("lbc-abcdef12", nil),
			service.EXPECT().GetLoadBalancerCluster("lbc-abcdef12").Return(ecloud.LoadBalancerCluster{}, errors.New("test error")),
		)

		err := ecloudLoadBalancerClusterCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new load balancer cluster: test error", err.Error())
	})
}

func Test_ecloudLoadBalancerClusterUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudLoadBalancerClusterUpdateCmd(nil).Args(nil, []string{"lbc-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudLoadBalancerClusterUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing load balancer cluster", err.Error())
	})
}

func Test_ecloudLoadBalancerClusterUpdate(t *testing.T) {
	t.Run("SingleLoadBalancerCluster", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudLoadBalancerClusterCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testnetwork"})

		req := ecloud.PatchLoadBalancerClusterRequest{
			Name: ptr.String("testnetwork"),
		}

		gomock.InOrder(
			service.EXPECT().PatchLoadBalancerCluster("lbc-abcdef12", req).Return(nil),
			service.EXPECT().GetLoadBalancerCluster("lbc-abcdef12").Return(ecloud.LoadBalancerCluster{}, nil),
		)

		ecloudLoadBalancerClusterUpdate(service, cmd, []string{"lbc-abcdef12"})
	})

	t.Run("MultipleLoadBalancerClusters", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchLoadBalancerCluster("lbc-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetLoadBalancerCluster("lbc-abcdef12").Return(ecloud.LoadBalancerCluster{}, nil),
			service.EXPECT().PatchLoadBalancerCluster("lbc-12abcdef", gomock.Any()).Return(nil),
			service.EXPECT().GetLoadBalancerCluster("lbc-12abcdef").Return(ecloud.LoadBalancerCluster{}, nil),
		)

		ecloudLoadBalancerClusterUpdate(service, &cobra.Command{}, []string{"lbc-abcdef12", "lbc-12abcdef"})
	})

	t.Run("PatchLoadBalancerClusterError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchLoadBalancerCluster("lbc-abcdef12", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating load balancer cluster [lbc-abcdef12]: test error\n", func() {
			ecloudLoadBalancerClusterUpdate(service, &cobra.Command{}, []string{"lbc-abcdef12"})
		})
	})

	t.Run("GetLoadBalancerClusterError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchLoadBalancerCluster("lbc-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetLoadBalancerCluster("lbc-abcdef12").Return(ecloud.LoadBalancerCluster{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated load balancer cluster [lbc-abcdef12]: test error\n", func() {
			ecloudLoadBalancerClusterUpdate(service, &cobra.Command{}, []string{"lbc-abcdef12"})
		})
	})
}

func Test_ecloudLoadBalancerClusterDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudLoadBalancerClusterDeleteCmd(nil).Args(nil, []string{"lbc-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudLoadBalancerClusterDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing load balancer cluster", err.Error())
	})
}

func Test_ecloudLoadBalancerClusterDelete(t *testing.T) {
	t.Run("SingleLoadBalancerCluster", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteLoadBalancerCluster("lbc-abcdef12").Return(nil).Times(1)

		ecloudLoadBalancerClusterDelete(service, &cobra.Command{}, []string{"lbc-abcdef12"})
	})

	t.Run("MultipleLoadBalancerClusters", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteLoadBalancerCluster("lbc-abcdef12").Return(nil),
			service.EXPECT().DeleteLoadBalancerCluster("lbc-12abcdef").Return(nil),
		)

		ecloudLoadBalancerClusterDelete(service, &cobra.Command{}, []string{"lbc-abcdef12", "lbc-12abcdef"})
	})

	t.Run("DeleteLoadBalancerClusterError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteLoadBalancerCluster("lbc-abcdef12").Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing load balancer cluster [lbc-abcdef12]: test error\n", func() {
			ecloudLoadBalancerClusterDelete(service, &cobra.Command{}, []string{"lbc-abcdef12"})
		})
	})
}
