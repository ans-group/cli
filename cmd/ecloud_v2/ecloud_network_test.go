package ecloud_v2

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

func Test_ecloudNetworkList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworks(gomock.Any()).Return([]ecloud.Network{}, nil).Times(1)

		ecloudNetworkList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudNetworkList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetNetworksError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworks(gomock.Any()).Return([]ecloud.Network{}, errors.New("test error")).Times(1)

		err := ecloudNetworkList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving networks: test error", err.Error())
	})
}

func Test_ecloudNetworkShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudNetworkShowCmd(nil).Args(nil, []string{"network-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudNetworkShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing network", err.Error())
	})
}

func Test_ecloudNetworkShow(t *testing.T) {
	t.Run("SingleNetwork", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetwork("network-abcdef12").Return(ecloud.Network{}, nil).Times(1)

		ecloudNetworkShow(service, &cobra.Command{}, []string{"network-abcdef12"})
	})

	t.Run("MultipleNetworks", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetNetwork("network-abcdef12").Return(ecloud.Network{}, nil),
			service.EXPECT().GetNetwork("network-abcdef23").Return(ecloud.Network{}, nil),
		)

		ecloudNetworkShow(service, &cobra.Command{}, []string{"network-abcdef12", "network-abcdef23"})
	})

	t.Run("GetNetworkError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetwork("network-abcdef12").Return(ecloud.Network{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving network [network-abcdef12]: test error\n", func() {
			ecloudNetworkShow(service, &cobra.Command{}, []string{"network-abcdef12"})
		})
	})
}
