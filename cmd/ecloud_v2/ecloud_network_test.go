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
	"github.com/ukfast/sdk-go/pkg/ptr"
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
		err := ecloudNetworkShowCmd(nil).Args(nil, []string{"net-abcdef12"})

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

		service.EXPECT().GetNetwork("net-abcdef12").Return(ecloud.Network{}, nil).Times(1)

		ecloudNetworkShow(service, &cobra.Command{}, []string{"net-abcdef12"})
	})

	t.Run("MultipleNetworks", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetNetwork("net-abcdef12").Return(ecloud.Network{}, nil),
			service.EXPECT().GetNetwork("net-abcdef23").Return(ecloud.Network{}, nil),
		)

		ecloudNetworkShow(service, &cobra.Command{}, []string{"net-abcdef12", "net-abcdef23"})
	})

	t.Run("GetNetworkError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetwork("net-abcdef12").Return(ecloud.Network{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving network [net-abcdef12]: test error\n", func() {
			ecloudNetworkShow(service, &cobra.Command{}, []string{"net-abcdef12"})
		})
	})
}

func Test_ecloudNetworkCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNetworkCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testnetwork", "--router=rtr-abcdef12"})

		req := ecloud.CreateNetworkRequest{
			Name:     ptr.String("testnetwork"),
			RouterID: "rtr-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateNetwork(req).Return("net-abcdef12", nil),
			service.EXPECT().GetNetwork("net-abcdef12").Return(ecloud.Network{}, nil),
		)

		ecloudNetworkCreate(service, cmd, []string{})
	})

	t.Run("CreateNetworkError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNetworkCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testnetwork", "--router=rtr-abcdef12"})

		service.EXPECT().CreateNetwork(gomock.Any()).Return("", errors.New("test error")).Times(1)

		err := ecloudNetworkCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating network: test error", err.Error())
	})

	t.Run("GetNetworkError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudNetworkCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testnetwork", "--router=rtr-abcdef12"})

		gomock.InOrder(
			service.EXPECT().CreateNetwork(gomock.Any()).Return("net-abcdef12", nil),
			service.EXPECT().GetNetwork("net-abcdef12").Return(ecloud.Network{}, errors.New("test error")),
		)

		err := ecloudNetworkCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new network: test error", err.Error())
	})
}

func Test_ecloudNetworkUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudNetworkUpdateCmd(nil).Args(nil, []string{"net-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudNetworkUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing network", err.Error())
	})
}

func Test_ecloudNetworkUpdate(t *testing.T) {
	t.Run("SingleNetwork", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNetworkCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testnetwork"})

		req := ecloud.PatchNetworkRequest{
			Name: ptr.String("testnetwork"),
		}

		gomock.InOrder(
			service.EXPECT().PatchNetwork("net-abcdef12", req).Return(nil),
			service.EXPECT().GetNetwork("net-abcdef12").Return(ecloud.Network{}, nil),
		)

		ecloudNetworkUpdate(service, cmd, []string{"net-abcdef12"})
	})

	t.Run("MultipleNetworks", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchNetwork("net-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetNetwork("net-abcdef12").Return(ecloud.Network{}, nil),
			service.EXPECT().PatchNetwork("net-12abcdef", gomock.Any()).Return(nil),
			service.EXPECT().GetNetwork("net-12abcdef").Return(ecloud.Network{}, nil),
		)

		ecloudNetworkUpdate(service, &cobra.Command{}, []string{"net-abcdef12", "net-12abcdef"})
	})

	t.Run("PatchNetworkError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchNetwork("net-abcdef12", gomock.Any()).Return(errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating network [net-abcdef12]: test error\n", func() {
			ecloudNetworkUpdate(service, &cobra.Command{}, []string{"net-abcdef12"})
		})
	})

	t.Run("GetNetworkError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().PatchNetwork("net-abcdef12", gomock.Any()).Return(nil),
			service.EXPECT().GetNetwork("net-abcdef12").Return(ecloud.Network{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated network [net-abcdef12]: test error\n", func() {
			ecloudNetworkUpdate(service, &cobra.Command{}, []string{"net-abcdef12"})
		})
	})
}

func Test_ecloudNetworkDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudNetworkDeleteCmd(nil).Args(nil, []string{"net-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudNetworkDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing network", err.Error())
	})
}

func Test_ecloudNetworkDelete(t *testing.T) {
	t.Run("SingleNetwork", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteNetwork("net-abcdef12").Return(nil).Times(1)

		ecloudNetworkDelete(service, &cobra.Command{}, []string{"net-abcdef12"})
	})

	t.Run("MultipleNetworks", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteNetwork("net-abcdef12").Return(nil),
			service.EXPECT().DeleteNetwork("net-12abcdef").Return(nil),
		)

		ecloudNetworkDelete(service, &cobra.Command{}, []string{"net-abcdef12", "net-12abcdef"})
	})

	t.Run("DeleteNetworkError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteNetwork("net-abcdef12").Return(errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing network [net-abcdef12]: test error\n", func() {
			ecloudNetworkDelete(service, &cobra.Command{}, []string{"net-abcdef12"})
		})
	})
}
