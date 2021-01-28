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

func Test_ecloudDHCPList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetDHCPs(gomock.Any()).Return([]ecloud.DHCP{}, nil).Times(1)

		ecloudDHCPList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudDHCPList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetDHCPsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetDHCPs(gomock.Any()).Return([]ecloud.DHCP{}, errors.New("test error")).Times(1)

		err := ecloudDHCPList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving DHCPs: test error", err.Error())
	})
}

func Test_ecloudDHCPShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudDHCPShowCmd(nil).Args(nil, []string{"dhcp-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudDHCPShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing dhcp", err.Error())
	})
}

func Test_ecloudDHCPShow(t *testing.T) {
	t.Run("SingleDHCP", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetDHCP("dhcp-abcdef12").Return(ecloud.DHCP{}, nil).Times(1)

		ecloudDHCPShow(service, &cobra.Command{}, []string{"dhcp-abcdef12"})
	})

	t.Run("MultipleDHCPs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetDHCP("dhcp-abcdef12").Return(ecloud.DHCP{}, nil),
			service.EXPECT().GetDHCP("dhcp-abcdef23").Return(ecloud.DHCP{}, nil),
		)

		ecloudDHCPShow(service, &cobra.Command{}, []string{"dhcp-abcdef12", "dhcp-abcdef23"})
	})

	t.Run("GetDHCPError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetDHCP("dhcp-abcdef12").Return(ecloud.DHCP{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving DHCP [dhcp-abcdef12]: test error\n", func() {
			ecloudDHCPShow(service, &cobra.Command{}, []string{"dhcp-abcdef12"})
		})
	})
}
