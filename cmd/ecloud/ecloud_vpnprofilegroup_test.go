package ecloud

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ecloudVPNProfileGroupList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		session.EXPECT().GetVPNProfileGroups(gomock.Any()).Return([]ecloud.VPNProfileGroup{}, nil).Times(1)

		ecloudVPNProfileGroupList(session, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudVPNProfileGroupList(session, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetVPNProfileGroupsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		session.EXPECT().GetVPNProfileGroups(gomock.Any()).Return([]ecloud.VPNProfileGroup{}, errors.New("test error")).Times(1)

		err := ecloudVPNProfileGroupList(session, &cobra.Command{}, []string{})

		assert.Equal(t, "error retrieving VPN sessions: test error", err.Error())
	})
}

func Test_ecloudVPNProfileGroupShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPNProfileGroupShowCmd(nil).Args(nil, []string{"vpnpg-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPNProfileGroupShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing VPN session", err.Error())
	})
}

func Test_ecloudVPNProfileGroupShow(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		session.EXPECT().GetVPNProfileGroup("vpnpg-abcdef12").Return(ecloud.VPNProfileGroup{}, nil).Times(1)

		ecloudVPNProfileGroupShow(session, &cobra.Command{}, []string{"vpnpg-abcdef12"})
	})

	t.Run("MultipleVPNProfileGroups", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			session.EXPECT().GetVPNProfileGroup("vpnpg-abcdef12").Return(ecloud.VPNProfileGroup{}, nil),
			session.EXPECT().GetVPNProfileGroup("vpnpg-abcdef23").Return(ecloud.VPNProfileGroup{}, nil),
		)

		ecloudVPNProfileGroupShow(session, &cobra.Command{}, []string{"vpnpg-abcdef12", "vpnpg-abcdef23"})
	})

	t.Run("GetVPNProfileGroupError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		session.EXPECT().GetVPNProfileGroup("vpnpg-abcdef12").Return(ecloud.VPNProfileGroup{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving VPN session [vpnpg-abcdef12]: test error\n", func() {
			ecloudVPNProfileGroupShow(session, &cobra.Command{}, []string{"vpnpg-abcdef12"})
		})
	})
}
