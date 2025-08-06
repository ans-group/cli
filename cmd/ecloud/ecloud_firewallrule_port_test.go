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

func Test_ecloudFirewallRuleFirewallRulePortListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFirewallRuleFirewallRulePortListCmd(nil).Args(nil, []string{"fwp-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFirewallRuleFirewallRulePortListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing firewall rule", err.Error())
	})
}

func Test_ecloudFirewallRuleFirewallRulePortList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewallRuleFirewallRulePorts("fwp-abcdef12", gomock.Any()).Return([]ecloud.FirewallRulePort{}, nil).Times(1)

		ecloudFirewallRuleFirewallRulePortList(service, &cobra.Command{}, []string{"fwp-abcdef12"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudFirewallRuleFirewallRulePortList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetFirewallRulesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewallRuleFirewallRulePorts("fwp-abcdef12", gomock.Any()).Return([]ecloud.FirewallRulePort{}, errors.New("test error")).Times(1)

		err := ecloudFirewallRuleFirewallRulePortList(service, &cobra.Command{}, []string{"fwp-abcdef12"})

		assert.Equal(t, "error retrieving firewall rule ports: test error", err.Error())
	})
}
