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

func Test_ecloudFirewallPolicyTaskListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFirewallPolicyTaskListCmd(nil).Args(nil, []string{"i-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFirewallPolicyTaskListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing firewall policy", err.Error())
	})
}

func Test_ecloudFirewallPolicyTaskList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewallPolicyTasks("i-abcdef12", gomock.Any()).Return([]ecloud.Task{}, nil).Times(1)

		ecloudFirewallPolicyTaskList(service, &cobra.Command{}, []string{"i-abcdef12"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudFirewallPolicyTaskList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetFirewallPolicysError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewallPolicyTasks("i-abcdef12", gomock.Any()).Return([]ecloud.Task{}, errors.New("test error")).Times(1)

		err := ecloudFirewallPolicyTaskList(service, &cobra.Command{}, []string{"i-abcdef12"})

		assert.Equal(t, "Error retrieving firewall policy tasks: test error", err.Error())
	})
}
