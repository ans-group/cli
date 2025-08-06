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

func Test_ecloudVPNGatewaySpecificationList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPNGatewaySpecifications(gomock.Any()).Return([]ecloud.VPNGatewaySpecification{}, nil).Times(1)

		ecloudVPNGatewaySpecificationList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudVPNGatewaySpecificationList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetVPNGatewaySpecificationsError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPNGatewaySpecifications(gomock.Any()).Return([]ecloud.VPNGatewaySpecification{}, errors.New("test error")).Times(1)

		err := ecloudVPNGatewaySpecificationList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "error retrieving VPN gateway specifications: test error", err.Error())
	})
}

func Test_ecloudVPNGatewaySpecificationShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPNGatewaySpecificationShowCmd(nil).Args(nil, []string{"vpngs-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPNGatewaySpecificationShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing VPN gateway specification", err.Error())
	})
}

func Test_ecloudVPNGatewaySpecificationShow(t *testing.T) {
	t.Run("SingleSpecification", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPNGatewaySpecification("vpngs-abcdef12").Return(ecloud.VPNGatewaySpecification{}, nil).Times(1)

		ecloudVPNGatewaySpecificationShow(service, &cobra.Command{}, []string{"vpngs-abcdef12"})
	})

	t.Run("MultipleSpecifications", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetVPNGatewaySpecification("vpngs-abcdef12").Return(ecloud.VPNGatewaySpecification{}, nil),
			service.EXPECT().GetVPNGatewaySpecification("vpngs-abcdef23").Return(ecloud.VPNGatewaySpecification{}, nil),
		)

		ecloudVPNGatewaySpecificationShow(service, &cobra.Command{}, []string{"vpngs-abcdef12", "vpngs-abcdef23"})
	})

	t.Run("GetVPNGatewaySpecificationError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPNGatewaySpecification("vpngs-abcdef12").Return(ecloud.VPNGatewaySpecification{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving VPN gateway specification [vpngs-abcdef12]: test error\n", func() {
			ecloudVPNGatewaySpecificationShow(service, &cobra.Command{}, []string{"vpngs-abcdef12"})
		})
	})
}
