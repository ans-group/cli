package ecloud

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ecloudMonitoringGatewayList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetMonitoringGateways(gomock.Any()).Return([]ecloud.MonitoringGateway{}, nil).Times(1)

		ecloudMonitoringGatewayList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudMonitoringGatewayList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetMonitoringGatewaysError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetMonitoringGateways(gomock.Any()).Return([]ecloud.MonitoringGateway{}, errors.New("test error")).Times(1)

		err := ecloudMonitoringGatewayList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "error retrieving monitoring gateways: test error", err.Error())
	})
}

func Test_ecloudMonitoringGatewayShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudMonitoringGatewayShowCmd(nil).Args(nil, []string{"mgw-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudMonitoringGatewayShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing monitoring gateway ID", err.Error())
	})
}

func Test_ecloudMonitoringGatewayShow(t *testing.T) {
	t.Run("SingleMonitoringGateway", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetMonitoringGateway("mgw-abcdef12").Return(ecloud.MonitoringGateway{}, nil).Times(1)

		ecloudMonitoringGatewayShow(service, &cobra.Command{}, []string{"mgw-abcdef12"})
	})

	t.Run("MultipleMonitoringGateways", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetMonitoringGateway("mgw-abcdef12").Return(ecloud.MonitoringGateway{}, nil),
			service.EXPECT().GetMonitoringGateway("mgw-abcdef23").Return(ecloud.MonitoringGateway{}, nil),
		)

		ecloudMonitoringGatewayShow(service, &cobra.Command{}, []string{"mgw-abcdef12", "mgw-abcdef23"})
	})

	t.Run("GetMonitoringGatewayError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetMonitoringGateway("mgw-abcdef12").Return(ecloud.MonitoringGateway{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving monitoring gateway [mgw-abcdef12]: test error\n", func() {
			ecloudMonitoringGatewayShow(service, &cobra.Command{}, []string{"mgw-abcdef12"})
		})
	})
}
