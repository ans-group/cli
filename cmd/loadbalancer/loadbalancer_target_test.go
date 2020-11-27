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

func Test_loadbalancerTargetList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetTargets(gomock.Any()).Return([]loadbalancer.Target{}, nil).Times(1)

		loadbalancerTargetList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := loadbalancerTargetList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetTargetsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetTargets(gomock.Any()).Return([]loadbalancer.Target{}, errors.New("test error")).Times(1)

		err := loadbalancerTargetList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving targets: test error", err.Error())
	})
}

func Test_loadbalancerTargetShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerTargetShowCmd(nil).Args(nil, []string{"rtr-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerTargetShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing target", err.Error())
	})
}

func Test_loadbalancerTargetShow(t *testing.T) {
	t.Run("SingleTarget", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetTarget("rtr-abcdef12").Return(loadbalancer.Target{}, nil).Times(1)

		loadbalancerTargetShow(service, &cobra.Command{}, []string{"rtr-abcdef12"})
	})

	t.Run("MultipleTargets", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetTarget("rtr-abcdef12").Return(loadbalancer.Target{}, nil),
			service.EXPECT().GetTarget("rtr-abcdef23").Return(loadbalancer.Target{}, nil),
		)

		loadbalancerTargetShow(service, &cobra.Command{}, []string{"rtr-abcdef12", "rtr-abcdef23"})
	})

	t.Run("GetTargetError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetTarget("rtr-abcdef12").Return(loadbalancer.Target{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving target [rtr-abcdef12]: test error\n", func() {
			loadbalancerTargetShow(service, &cobra.Command{}, []string{"rtr-abcdef12"})
		})
	})
}
