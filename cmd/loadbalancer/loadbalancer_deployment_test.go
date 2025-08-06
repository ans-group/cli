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

func Test_loadbalancerDeploymentList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetDeployments(gomock.Any()).Return([]loadbalancer.Deployment{}, nil).Times(1)

		loadbalancerDeploymentList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerDeploymentListCmd(nil)
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := loadbalancerDeploymentList(service, cmd, []string{})

		assert.NotNil(t, err)
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetDeploymentsError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetDeployments(gomock.Any()).Return([]loadbalancer.Deployment{}, errors.New("test error")).Times(1)

		err := loadbalancerDeploymentList(service, &cobra.Command{}, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "error retrieving deployments: test error", err.Error())
	})
}

func Test_loadbalancerDeploymentShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerDeploymentShowCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerDeploymentShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing deployment", err.Error())
	})
}

func Test_loadbalancerDeploymentShow(t *testing.T) {
	t.Run("SingleDeployment", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetDeployment(123).Return(loadbalancer.Deployment{}, nil).Times(1)

		loadbalancerDeploymentShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleDeployments", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetDeployment(123).Return(loadbalancer.Deployment{}, nil),
			service.EXPECT().GetDeployment(456).Return(loadbalancer.Deployment{}, nil),
		)

		loadbalancerDeploymentShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetDeploymentID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid deployment ID [abc]\n", func() {
			loadbalancerDeploymentShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetDeploymentError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetDeployment(123).Return(loadbalancer.Deployment{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving deployment [123]: test error\n", func() {
			loadbalancerDeploymentShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}
