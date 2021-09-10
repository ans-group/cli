package ecloud

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

func Test_ecloudVPNSessionPreSharedKeyShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPNSessionPreSharedKeyShowCmd(nil).Args(nil, []string{"vpn-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPNSessionPreSharedKeyShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing VPN session", err.Error())
	})
}

func Test_ecloudVPNSessionPreSharedKeyShow(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPNSessionPreSharedKey("vpn-abcdef12").Return(ecloud.VPNSessionPreSharedKey{}, nil).Times(1)

		ecloudVPNSessionPreSharedKeyShow(service, &cobra.Command{}, []string{"vpn-abcdef12"})
	})

	t.Run("MultipleVPNSessionPreSharedKeys", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetVPNSessionPreSharedKey("vpn-abcdef12").Return(ecloud.VPNSessionPreSharedKey{}, nil),
			service.EXPECT().GetVPNSessionPreSharedKey("vpn-abcdef23").Return(ecloud.VPNSessionPreSharedKey{}, nil),
		)

		ecloudVPNSessionPreSharedKeyShow(service, &cobra.Command{}, []string{"vpn-abcdef12", "vpn-abcdef23"})
	})

	t.Run("GetVPNSessionPreSharedKeyError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPNSessionPreSharedKey("vpn-abcdef12").Return(ecloud.VPNSessionPreSharedKey{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving VPN session [vpn-abcdef12] pre-shared key: test error\n", func() {
			ecloudVPNSessionPreSharedKeyShow(service, &cobra.Command{}, []string{"vpn-abcdef12"})
		})
	})
}
