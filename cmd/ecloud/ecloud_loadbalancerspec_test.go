package ecloud

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

func Test_ecloudLoadBalancerSpecList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetLoadBalancerSpecs(gomock.Any()).Return([]ecloud.LoadBalancerSpec{}, nil).Times(1)

		ecloudLoadBalancerSpecList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudLoadBalancerSpecList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetLoadBalancerSpecsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetLoadBalancerSpecs(gomock.Any()).Return([]ecloud.LoadBalancerSpec{}, errors.New("test error")).Times(1)

		err := ecloudLoadBalancerSpecList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving load balancer specs: test error", err.Error())
	})
}

func Test_ecloudLoadBalancerSpecShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudLoadBalancerSpecShowCmd(nil).Args(nil, []string{"lbn-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudLoadBalancerSpecShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing load balancer spec", err.Error())
	})
}

func Test_ecloudLoadBalancerSpecShow(t *testing.T) {
	t.Run("SingleLoadBalancerSpec", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetLoadBalancerSpec("lbn-abcdef12").Return(ecloud.LoadBalancerSpec{}, nil).Times(1)

		ecloudLoadBalancerSpecShow(service, &cobra.Command{}, []string{"lbn-abcdef12"})
	})

	t.Run("MultipleLoadBalancerSpecs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetLoadBalancerSpec("lbn-abcdef12").Return(ecloud.LoadBalancerSpec{}, nil),
			service.EXPECT().GetLoadBalancerSpec("lbn-abcdef23").Return(ecloud.LoadBalancerSpec{}, nil),
		)

		ecloudLoadBalancerSpecShow(service, &cobra.Command{}, []string{"lbn-abcdef12", "lbn-abcdef23"})
	})

	t.Run("GetLoadBalancerSpecError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetLoadBalancerSpec("lbn-abcdef12").Return(ecloud.LoadBalancerSpec{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving load balancer spec [lbn-abcdef12]: test error\n", func() {
			ecloudLoadBalancerSpecShow(service, &cobra.Command{}, []string{"lbn-abcdef12"})
		})
	})
}
