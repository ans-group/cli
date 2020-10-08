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

func Test_loadbalancerGroupList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetGroups(gomock.Any()).Return([]loadbalancer.Group{}, nil).Times(1)

		loadbalancerGroupList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := loadbalancerGroupList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetGroupsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetGroups(gomock.Any()).Return([]loadbalancer.Group{}, errors.New("test error")).Times(1)

		err := loadbalancerGroupList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving groups: test error", err.Error())
	})
}

func Test_loadbalancerGroupShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerGroupShowCmd(nil).Args(nil, []string{"rtr-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerGroupShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing group", err.Error())
	})
}

func Test_loadbalancerGroupShow(t *testing.T) {
	t.Run("SingleGroup", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetGroup("rtr-abcdef12").Return(loadbalancer.Group{}, nil).Times(1)

		loadbalancerGroupShow(service, &cobra.Command{}, []string{"rtr-abcdef12"})
	})

	t.Run("MultipleGroups", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetGroup("rtr-abcdef12").Return(loadbalancer.Group{}, nil),
			service.EXPECT().GetGroup("rtr-abcdef23").Return(loadbalancer.Group{}, nil),
		)

		loadbalancerGroupShow(service, &cobra.Command{}, []string{"rtr-abcdef12", "rtr-abcdef23"})
	})

	t.Run("GetGroupError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetGroup("rtr-abcdef12").Return(loadbalancer.Group{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving group [rtr-abcdef12]: test error\n", func() {
			loadbalancerGroupShow(service, &cobra.Command{}, []string{"rtr-abcdef12"})
		})
	})
}
