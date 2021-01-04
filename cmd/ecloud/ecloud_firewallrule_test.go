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

func Test_ecloudFirewallRuleList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewallRules(gomock.Any()).Return([]ecloud.FirewallRule{}, nil).Times(1)

		ecloudFirewallRuleList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudFirewallRuleList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetFirewallRulesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewallRules(gomock.Any()).Return([]ecloud.FirewallRule{}, errors.New("test error")).Times(1)

		err := ecloudFirewallRuleList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving firewall rules: test error", err.Error())
	})
}

func Test_ecloudFirewallRuleShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFirewallRuleShowCmd(nil).Args(nil, []string{"fwr-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFirewallRuleShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing firewall rule", err.Error())
	})
}

func Test_ecloudFirewallRuleShow(t *testing.T) {
	t.Run("SingleFirewallRule", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewallRule("fwr-abcdef12").Return(ecloud.FirewallRule{}, nil).Times(1)

		ecloudFirewallRuleShow(service, &cobra.Command{}, []string{"fwr-abcdef12"})
	})

	t.Run("MultipleFirewallRules", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetFirewallRule("fwr-abcdef12").Return(ecloud.FirewallRule{}, nil),
			service.EXPECT().GetFirewallRule("fwr-abcdef23").Return(ecloud.FirewallRule{}, nil),
		)

		ecloudFirewallRuleShow(service, &cobra.Command{}, []string{"fwr-abcdef12", "fwr-abcdef23"})
	})

	t.Run("GetFirewallRuleError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewallRule("fwr-abcdef12").Return(ecloud.FirewallRule{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving firewall rule [fwr-abcdef12]: test error\n", func() {
			ecloudFirewallRuleShow(service, &cobra.Command{}, []string{"fwr-abcdef12"})
		})
	})
}
