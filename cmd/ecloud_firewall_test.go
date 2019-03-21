package cmd

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/cli/test"
	"github.com/ukfast/cli/test/mocks"
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
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()
		defer func() { flagFilter = nil }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		flagFilter = []string{"invalidfilter"}

		output := test.CatchStdErr(t, func() {
			ecloudFirewallList(service, &cobra.Command{}, []string{})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Missing value for filtering\n", output)
	})

	t.Run("GetFirewallsError_OutputsFatal", func(t *testing.T) {
		code := 0
		oldOutputExit := output.SetOutputExit(func(c int) {
			code = c
		})
		defer func() { output.SetOutputExit(oldOutputExit) }()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewalls(gomock.Any()).Return([]ecloud.Firewall{}, errors.New("test error")).Times(1)

		output := test.CatchStdErr(t, func() {
			ecloudFirewallList(service, &cobra.Command{}, []string{})
		})

		assert.Equal(t, 1, code)
		assert.Equal(t, "Error retrieving firewalls: test error\n", output)
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

		output := test.CatchStdErr(t, func() {
			ecloudFirewallShow(service, &cobra.Command{}, []string{"abc"})
		})

		assert.Equal(t, "Invalid firewall ID [abc]\n", output)
	})

	t.Run("GetFirewallError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetFirewall(123).Return(ecloud.Firewall{}, errors.New("test error"))

		output := test.CatchStdErr(t, func() {
			ecloudFirewallShow(service, &cobra.Command{}, []string{"123"})
		})

		assert.Equal(t, "Error retrieving firewall [123]: test error\n", output)
	})
}
