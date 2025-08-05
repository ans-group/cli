package ecloud

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ecloudNetworkPolicyNetworkRuleListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudNetworkPolicyNetworkRuleListCmd(nil).Args(nil, []string{"np-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudNetworkPolicyNetworkRuleListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing network policy", err.Error())
	})
}

func Test_ecloudNetworkPolicyNetworkRuleList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworkPolicyNetworkRules("np-abcdef12", gomock.Any()).Return([]ecloud.NetworkRule{}, nil).Times(1)

		ecloudNetworkPolicyNetworkRuleList(service, &cobra.Command{}, []string{"np-abcdef12"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudNetworkPolicyNetworkRuleList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetNetworkPolicysError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworkPolicyNetworkRules("np-abcdef12", gomock.Any()).Return([]ecloud.NetworkRule{}, errors.New("test error")).Times(1)

		err := ecloudNetworkPolicyNetworkRuleList(service, &cobra.Command{}, []string{"np-abcdef12"})

		assert.Equal(t, "error retrieving network policy network rules: test error", err.Error())
	})
}
