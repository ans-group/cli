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

func Test_ecloudRouterFirewallPolicyListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudRouterFirewallPolicyListCmd(nil).Args(nil, []string{"rtr-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudRouterFirewallPolicyListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing router", err.Error())
	})
}

func Test_ecloudRouterFirewallPolicyList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetRouterFirewallPolicies("rtr-abcdef12", gomock.Any()).Return([]ecloud.FirewallPolicy{}, nil).Times(1)

		ecloudRouterFirewallPolicyList(service, &cobra.Command{}, []string{"rtr-abcdef12"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudRouterFirewallPolicyList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetRoutersError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetRouterFirewallPolicies("rtr-abcdef12", gomock.Any()).Return([]ecloud.FirewallPolicy{}, errors.New("test error")).Times(1)

		err := ecloudRouterFirewallPolicyList(service, &cobra.Command{}, []string{"rtr-abcdef12"})

		assert.Equal(t, "Error retrieving router firewall policies: test error", err.Error())
	})
}
