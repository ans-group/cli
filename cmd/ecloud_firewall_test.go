package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/cli/test/test_output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func Test_ecloudFirewallList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewalls(gomock.Any()).Return([]ecloud.Firewall{}, nil).Times(1)

		ecloudFirewallList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_OutputsFatal", func(t *testing.T) {
		defer func() { flagFilter = nil }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		flagFilter = []string{"invalidfilter"}

		test_output.AssertFatalOutput(t, "Missing value for filtering\n", func() {
			ecloudFirewallList(service, &cobra.Command{}, []string{})
		})
	})

	t.Run("GetFirewallsError_OutputsFatal", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewalls(gomock.Any()).Return([]ecloud.Firewall{}, errors.New("test error")).Times(1)

		test_output.AssertFatalOutput(t, "Error retrieving firewalls: test error\n", func() {
			ecloudFirewallList(service, &cobra.Command{}, []string{})
		})
	})
}

func Test_ecloudFirewallShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudFirewallShowCmd().Args(nil, []string{"123"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudFirewallShowCmd().Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing firewall", err.Error())
	})
}

func Test_ecloudFirewallShow(t *testing.T) {
	t.Run("SingleFirewall", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewall(123).Return(ecloud.Firewall{}, nil).Times(1)

		ecloudFirewallShow(service, &cobra.Command{}, []string{"123"})
	})

	t.Run("MultipleFirewalls", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetFirewall(123).Return(ecloud.Firewall{}, nil),
			service.EXPECT().GetFirewall(456).Return(ecloud.Firewall{}, nil),
		)

		ecloudFirewallShow(service, &cobra.Command{}, []string{"123", "456"})
	})

	t.Run("GetFirewallID_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		test_output.AssertErrorOutput(t, "Invalid firewall ID [abc]\n", func() {
			ecloudFirewallShow(service, &cobra.Command{}, []string{"abc"})
		})
	})

	t.Run("GetFirewallError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewall(123).Return(ecloud.Firewall{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving firewall [123]: test error\n", func() {
			ecloudFirewallShow(service, &cobra.Command{}, []string{"123"})
		})
	})
}
