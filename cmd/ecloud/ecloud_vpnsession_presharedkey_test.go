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
		err := ecloudVPNSessionPreSharedKeyShowCmd(nil).Args(nil, []string{"vpns-abcdef12"})

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

		service.EXPECT().GetVPNSessionPreSharedKey("vpns-abcdef12").Return(ecloud.VPNSessionPreSharedKey{}, nil).Times(1)

		ecloudVPNSessionPreSharedKeyShow(service, &cobra.Command{}, []string{"vpns-abcdef12"})
	})

	t.Run("MultipleVPNSessionPreSharedKeys", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetVPNSessionPreSharedKey("vpns-abcdef12").Return(ecloud.VPNSessionPreSharedKey{}, nil),
			service.EXPECT().GetVPNSessionPreSharedKey("vpns-abcdef23").Return(ecloud.VPNSessionPreSharedKey{}, nil),
		)

		ecloudVPNSessionPreSharedKeyShow(service, &cobra.Command{}, []string{"vpns-abcdef12", "vpns-abcdef23"})
	})

	t.Run("GetVPNSessionPreSharedKeyError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPNSessionPreSharedKey("vpns-abcdef12").Return(ecloud.VPNSessionPreSharedKey{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving VPN session [vpns-abcdef12] pre-shared key: test error\n", func() {
			ecloudVPNSessionPreSharedKeyShow(service, &cobra.Command{}, []string{"vpns-abcdef12"})
		})
	})
}

func Test_ecloudVPNSessionPreSharedKeyUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPNSessionPreSharedKeyUpdateCmd(nil).Args(nil, []string{"vpns-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPNSessionPreSharedKeyUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing VPN session", err.Error())
	})
}

func Test_ecloudVPNSessionPreSharedKeyUpdate(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNSessionPreSharedKeyUpdateCmd(nil)
		cmd.Flags().Set("psk", "testkey")

		req := ecloud.UpdateVPNSessionPreSharedKeyRequest{
			PSK: "testkey",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpns-abcdef12",
		}

		service.EXPECT().UpdateVPNSessionPreSharedKey("vpns-abcdef12", req).Return(resp, nil)

		ecloudVPNSessionPreSharedKeyUpdate(service, cmd, []string{"vpns-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNSessionPreSharedKeyUpdateCmd(nil)
		cmd.Flags().Set("psk", "testkey")
		cmd.Flags().Set("wait", "true")

		req := ecloud.UpdateVPNSessionPreSharedKeyRequest{
			PSK: "testkey",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpns-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().UpdateVPNSessionPreSharedKey("vpns-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		err := ecloudVPNSessionPreSharedKeyUpdate(service, cmd, []string{"vpns-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNSessionPreSharedKeyUpdateCmd(nil)
		cmd.Flags().Set("psk", "testkey")
		cmd.Flags().Set("wait", "true")

		req := ecloud.UpdateVPNSessionPreSharedKeyRequest{
			PSK: "testkey",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpns-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().UpdateVPNSessionPreSharedKey("vpns-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		err := ecloudVPNSessionPreSharedKeyUpdate(service, cmd, []string{"vpns-abcdef12"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error waiting for task to complete for VPN session: Error waiting for command: Failed to retrieve task status: test error", err.Error())
	})

	t.Run("UpdateVPNSessionPreSharedKeyError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().UpdateVPNSessionPreSharedKey("vpns-abcdef12", gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error"))

		err := ecloudVPNSessionPreSharedKeyUpdate(service, &cobra.Command{}, []string{"vpns-abcdef12"})

		assert.NotNil(t, err)
		assert.Equal(t, "Error updating VPN session pre-shared key: test error", err.Error())
	})
}
