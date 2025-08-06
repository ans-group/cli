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

func Test_ecloudVPNGatewayList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPNGateways(gomock.Any()).Return([]ecloud.VPNGateway{}, nil).Times(1)

		ecloudVPNGatewayList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudVPNGatewayList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetVPNGatewaysError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPNGateways(gomock.Any()).Return([]ecloud.VPNGateway{}, errors.New("test error")).Times(1)

		err := ecloudVPNGatewayList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "error retrieving VPN gateways: test error", err.Error())
	})
}

func Test_ecloudVPNGatewayShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPNGatewayShowCmd(nil).Args(nil, []string{"vpng-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPNGatewayShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing VPN gateway", err.Error())
	})
}

func Test_ecloudVPNGatewayShow(t *testing.T) {
	t.Run("SingleGateway", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPNGateway("vpng-abcdef12").Return(ecloud.VPNGateway{}, nil).Times(1)

		ecloudVPNGatewayShow(service, &cobra.Command{}, []string{"vpng-abcdef12"})
	})

	t.Run("MultipleVPNGateways", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetVPNGateway("vpng-abcdef12").Return(ecloud.VPNGateway{}, nil),
			service.EXPECT().GetVPNGateway("vpng-abcdef23").Return(ecloud.VPNGateway{}, nil),
		)

		ecloudVPNGatewayShow(service, &cobra.Command{}, []string{"vpng-abcdef12", "vpng-abcdef23"})
	})

	t.Run("GetVPNGatewayError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPNGateway("vpng-abcdef12").Return(ecloud.VPNGateway{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving VPN gateway [vpng-abcdef12]: test error\n", func() {
			ecloudVPNGatewayShow(service, &cobra.Command{}, []string{"vpng-abcdef12"})
		})
	})
}

func Test_ecloudVPNGatewayCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNGatewayCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgateway", "--router=rtr-abcdef12", "--specification=vpngs-abcdef12"})

		req := ecloud.CreateVPNGatewayRequest{
			Name:            "testgateway",
			RouterID:        "rtr-abcdef12",
			SpecificationID: "vpngs-abcdef12",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpng-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateVPNGateway(req).Return(resp, nil),
			service.EXPECT().GetVPNGateway("vpng-abcdef12").Return(ecloud.VPNGateway{}, nil),
		)

		ecloudVPNGatewayCreate(service, cmd, []string{})
	})

	t.Run("CreateWithWaitFlagNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNGatewayCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgateway", "--router=rtr-abcdef12", "--specification=vpngs-abcdef12", "--wait"})

		req := ecloud.CreateVPNGatewayRequest{
			Name:            "testgateway",
			RouterID:        "rtr-abcdef12",
			SpecificationID: "vpngs-abcdef12",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpng-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateVPNGateway(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetVPNGateway("vpng-abcdef12").Return(ecloud.VPNGateway{}, nil),
		)

		ecloudVPNGatewayCreate(service, cmd, []string{})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNGatewayCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgateway", "--router=rtr-abcdef12", "--specification=vpngs-abcdef12", "--wait"})

		req := ecloud.CreateVPNGatewayRequest{
			Name:            "testgateway",
			RouterID:        "rtr-abcdef12",
			SpecificationID: "vpngs-abcdef12",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpng-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateVPNGateway(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		err := ecloudVPNGatewayCreate(service, cmd, []string{})
		assert.Equal(t, "error waiting for VPN gateway task to complete: error waiting for command: failed to retrieve task status: test error", err.Error())
	})

	t.Run("CreateVPNGatewayError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNGatewayCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgateway", "--router=rtr-abcdef12", "--specification=vpngs-abcdef12"})

		service.EXPECT().CreateVPNGateway(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error")).Times(1)

		err := ecloudVPNGatewayCreate(service, cmd, []string{})

		assert.Equal(t, "error creating VPN gateway: test error", err.Error())
	})

	t.Run("GetVPNGatewayError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNGatewayCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgateway", "--router=rtr-abcdef12", "--specification=vpngs-abcdef12"})

		gomock.InOrder(
			service.EXPECT().CreateVPNGateway(gomock.Any()).Return(ecloud.TaskReference{ResourceID: "vpng-abcdef12"}, nil),
			service.EXPECT().GetVPNGateway("vpng-abcdef12").Return(ecloud.VPNGateway{}, errors.New("test error")),
		)

		err := ecloudVPNGatewayCreate(service, cmd, []string{})

		assert.Equal(t, "error retrieving new VPN gateway: test error", err.Error())
	})
}

func Test_ecloudVPNGatewayUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPNGatewayUpdateCmd(nil).Args(nil, []string{"vpng-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPNGatewayUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing VPN gateway", err.Error())
	})
}

func Test_ecloudVPNGatewayUpdate(t *testing.T) {
	t.Run("SingleGateway", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNGatewayUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgateway"})

		req := ecloud.PatchVPNGatewayRequest{
			Name: "testgateway",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpng-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchVPNGateway("vpng-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetVPNGateway("vpng-abcdef12").Return(ecloud.VPNGateway{}, nil),
		)

		ecloudVPNGatewayUpdate(service, cmd, []string{"vpng-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNGatewayUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgateway", "--wait"})

		req := ecloud.PatchVPNGatewayRequest{
			Name: "testgateway",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpng-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchVPNGateway("vpng-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetVPNGateway("vpng-abcdef12").Return(ecloud.VPNGateway{}, nil),
		)

		ecloudVPNGatewayUpdate(service, cmd, []string{"vpng-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNGatewayUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgateway", "--wait"})

		req := ecloud.PatchVPNGatewayRequest{
			Name: "testgateway",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpng-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchVPNGateway("vpng-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for VPN gateway [vpng-abcdef12]: error waiting for command: failed to retrieve task status: test error\n", func() {
			ecloudVPNGatewayUpdate(service, cmd, []string{"vpng-abcdef12"})
		})
	})

	t.Run("PatchVPNGatewayError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchVPNGateway("vpng-abcdef12", gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating VPN gateway [vpng-abcdef12]: test error\n", func() {
			ecloudVPNGatewayUpdate(service, &cobra.Command{}, []string{"vpng-abcdef12"})
		})
	})

	t.Run("GetVPNGatewayError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpng-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchVPNGateway("vpng-abcdef12", gomock.Any()).Return(resp, nil),
			service.EXPECT().GetVPNGateway("vpng-abcdef12").Return(ecloud.VPNGateway{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated VPN gateway [vpng-abcdef12]: test error\n", func() {
			ecloudVPNGatewayUpdate(service, &cobra.Command{}, []string{"vpng-abcdef12"})
		})
	})
}

func Test_ecloudVPNGatewayDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPNGatewayDeleteCmd(nil).Args(nil, []string{"vpng-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPNGatewayDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing VPN gateway", err.Error())
	})
}

func Test_ecloudVPNGatewayDelete(t *testing.T) {
	t.Run("SingleGateway", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteVPNGateway("vpng-abcdef12").Return("task-abcdef12", nil)

		ecloudVPNGatewayDelete(service, &cobra.Command{}, []string{"vpng-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudVPNGatewayDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteVPNGateway("vpng-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudVPNGatewayDelete(service, cmd, []string{"vpng-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNGatewayDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		gomock.InOrder(
			service.EXPECT().DeleteVPNGateway("vpng-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for VPN gateway [vpng-abcdef12]: error waiting for command: failed to retrieve task status: test error\n", func() {
			ecloudVPNGatewayDelete(service, cmd, []string{"vpng-abcdef12"})
		})
	})

	t.Run("DeleteVPNGatewayError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteVPNGateway("vpng-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing VPN gateway [vpng-abcdef12]: test error\n", func() {
			ecloudVPNGatewayDelete(service, &cobra.Command{}, []string{"vpng-abcdef12"})
		})
	})
}
