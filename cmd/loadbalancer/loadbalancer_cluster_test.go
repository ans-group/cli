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
		cmd := loadbalancerClusterListCmd(nil)
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := loadbalancerClusterList(service, cmd, []string{})

		assert.NotNil(t, err)
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetClusterError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetClusters(gomock.Any()).Return([]loadbalancer.Cluster{}, errors.New("test error")).Times(1)

		err := loadbalancerClusterList(service, &cobra.Command{}, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving clusters: test error", err.Error())
	})
}

func Test_loadbalancerClusterShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerClusterShowCmd(nil).Args(nil, []string{"123"})

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

		service.EXPECT().GetCluster(123).Return(loadbalancer.Cluster{}, nil).Times(1)

		loadbalancerClusterShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleClusters", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetCluster(123).Return(loadbalancer.Cluster{}, nil),
			service.EXPECT().GetCluster(456).Return(loadbalancer.Cluster{}, nil),
		)

		loadbalancerClusterShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetClusterID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid cluster ID [abc]\n", func() {
			loadbalancerClusterShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetClusterError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetCluster(123).Return(loadbalancer.Cluster{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving cluster [123]: test error\n", func() {
			loadbalancerClusterShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_loadbalancerClusterUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerClusterUpdateCmd(nil).Args(nil, []string{"123"})

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
		cmd.ParseFlags([]string{"--name=test"})

		req := loadbalancer.PatchClusterRequest{
			Name: "test",
		}

		service.EXPECT().PatchCluster(123, req).Return(nil).Times(1)
		service.EXPECT().GetCluster(123).Return(loadbalancer.Cluster{}, nil).Times(1)

		loadbalancerClusterUpdate(service, cmd, []string{"123"})
	})

	t.Run("GetClusterID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid cluster ID [abc]\n", func() {
			loadbalancerClusterUpdate(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("PatchCluster_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().PatchCluster(123, gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating cluster [123]: test error\n", func() {
			loadbalancerClusterUpdate(service, &cobra.Command{}, []string{"123"})
		})
	})

	t.Run("GetClusterError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().PatchCluster(123, gomock.Any()).Return(nil)
		service.EXPECT().GetCluster(123).Return(loadbalancer.Cluster{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving updated cluster [123]: test error\n", func() {
			loadbalancerClusterUpdate(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_loadbalancerClusterDeployCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerClusterDeployCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerClusterDeployCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing cluster", err.Error())
	})
}

func Test_loadbalancerClusterDeploy(t *testing.T) {
	t.Run("SingleCluster", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().DeployCluster(123).Return(nil).Times(1)

		loadbalancerClusterDeploy(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("GetClusterID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid cluster ID [abc]\n", func() {
			loadbalancerClusterDeploy(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetClusterError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().DeployCluster(123).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error deploying cluster [123]: test error\n", func() {
			loadbalancerClusterDeploy(service, &cobra.Command{}, []string{"123"})
		})
	})
}

func Test_loadbalancerClusterValidateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerClusterValidateCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerClusterValidateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing cluster", err.Error())
	})
}

func Test_loadbalancerClusterValidate(t *testing.T) {
	t.Run("SingleCluster", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().ValidateCluster(123).Return(nil).Times(1)

		loadbalancerClusterValidate(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("GetClusterID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid cluster ID [abc]\n", func() {
			loadbalancerClusterValidate(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetClusterError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().ValidateCluster(123).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error validating cluster [123]: test error\n", func() {
			loadbalancerClusterValidate(service, &cobra.Command{}, []string{"123"})
		})
	})
}
