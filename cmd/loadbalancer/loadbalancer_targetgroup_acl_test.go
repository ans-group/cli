package loadbalancer

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/service/loadbalancer"
)

func Test_loadbalancerTargetGroupACLListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := loadbalancerTargetGroupACLListCmd(nil).Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := loadbalancerTargetGroupACLListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing target group", err.Error())
	})

	t.Run("InvalidListenerID_Error", func(t *testing.T) {
		err := loadbalancerTargetGroupACLListCmd(nil).Args(nil, []string{"abc"})

		assert.NotNil(t, err)
		assert.Equal(t, "Invalid target group ID", err.Error())
	})
}

func Test_loadbalancerTargetGroupACLList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		params := connection.NewAPIRequestParameters().WithFilter(connection.APIRequestFiltering{
			Property: "target_group_id",
			Operator: connection.EQOperator,
			Value:    []string{"123"},
		})

		service.EXPECT().GetACLs(*params).Return([]loadbalancer.ACL{}, nil).Times(1)

		loadbalancerTargetGroupACLList(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)
		cmd := loadbalancerTargetGroupACLListCmd(nil)
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := loadbalancerTargetGroupACLList(service, cmd, []string{"123"})

		assert.NotNil(t, err)
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetACLsError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockLoadBalancerService(mockCtrl)

		service.EXPECT().GetACLs(gomock.Any()).Return([]loadbalancer.ACL{}, errors.New("test error")).Times(1)

		err := loadbalancerTargetGroupACLList(service, &cobra.Command{}, []string{"123"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving ACLs: test error", err.Error())
	})
}
