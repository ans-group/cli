package ecloudv2

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

func Test_ecloudVPNList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPNs(gomock.Any()).Return([]ecloud.VPN{}, nil).Times(1)

		ecloudVPNList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudVPNList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetVPNsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPNs(gomock.Any()).Return([]ecloud.VPN{}, errors.New("test error")).Times(1)

		err := ecloudVPNList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving VPNs: test error", err.Error())
	})
}

func Test_ecloudVPNShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPNShowCmd(nil).Args(nil, []string{"vpn-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPNShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing VPN", err.Error())
	})
}

func Test_ecloudVPNShow(t *testing.T) {
	t.Run("SingleVPN", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPN("vpn-abcdef12").Return(ecloud.VPN{}, nil).Times(1)

		ecloudVPNShow(service, &cobra.Command{}, []string{"vpn-abcdef12"})
	})

	t.Run("MultipleVPNs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetVPN("vpn-abcdef12").Return(ecloud.VPN{}, nil),
			service.EXPECT().GetVPN("vpn-abcdef23").Return(ecloud.VPN{}, nil),
		)

		ecloudVPNShow(service, &cobra.Command{}, []string{"vpn-abcdef12", "vpn-abcdef23"})
	})

	t.Run("GetVPNError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPN("vpn-abcdef12").Return(ecloud.VPN{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving VPN [vpn-abcdef12]: test error\n", func() {
			ecloudVPNShow(service, &cobra.Command{}, []string{"vpn-abcdef12"})
		})
	})
}

func Test_ecloudVPNCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNCreateCmd(nil)
		cmd.ParseFlags([]string{"--router=rtr-abcdef12"})

		req := ecloud.CreateVPNRequest{
			RouterID: "rtr-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateVPN(req).Return("vpn-abcdef12", nil),
			service.EXPECT().GetVPN("vpn-abcdef12").Return(ecloud.VPN{}, nil),
		)

		ecloudVPNCreate(service, cmd, []string{})
	})

	t.Run("CreateVPNError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNCreateCmd(nil)
		cmd.ParseFlags([]string{"--router=rtr-abcdef12"})

		service.EXPECT().CreateVPN(gomock.Any()).Return("", errors.New("test error")).Times(1)

		err := ecloudVPNCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating VPN: test error", err.Error())
	})

	t.Run("GetVPNError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNCreateCmd(nil)
		cmd.ParseFlags([]string{"--router=rtr-abcdef12"})

		gomock.InOrder(
			service.EXPECT().CreateVPN(gomock.Any()).Return("vpn-abcdef12", nil),
			service.EXPECT().GetVPN("vpn-abcdef12").Return(ecloud.VPN{}, errors.New("test error")),
		)

		err := ecloudVPNCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new VPN: test error", err.Error())
	})
}

func Test_ecloudVPNDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPNDeleteCmd(nil).Args(nil, []string{"vpn-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPNDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing vpn", err.Error())
	})
}

func Test_ecloudVPNDelete(t *testing.T) {
	t.Run("SingleVPN", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteVPN("vpn-abcdef12").Return(nil).Times(1)

		ecloudVPNDelete(service, &cobra.Command{}, []string{"vpn-abcdef12"})
	})

	t.Run("MultipleVPNs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteVPN("vpn-abcdef12").Return(nil),
			service.EXPECT().DeleteVPN("vpn-12abcdef").Return(nil),
		)

		ecloudVPNDelete(service, &cobra.Command{}, []string{"vpn-abcdef12", "vpn-12abcdef"})
	})

	t.Run("DeleteVPNError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteVPN("vpn-abcdef12").Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing VPN [vpn-abcdef12]: test error\n", func() {
			ecloudVPNDelete(service, &cobra.Command{}, []string{"vpn-abcdef12"})
		})
	})
}
