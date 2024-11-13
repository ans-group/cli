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

func Test_ecloudVPNGatewayUserList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPNGatewayUsers(gomock.Any()).Return([]ecloud.VPNGatewayUser{}, nil).Times(1)

		ecloudVPNGatewayUserList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudVPNGatewayUserList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetVPNGatewayUsersError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPNGatewayUsers(gomock.Any()).Return([]ecloud.VPNGatewayUser{}, errors.New("test error")).Times(1)

		err := ecloudVPNGatewayUserList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving VPN gateway users: test error", err.Error())
	})
}

func Test_ecloudVPNGatewayUserShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPNGatewayUserShowCmd(nil).Args(nil, []string{"vpngu-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPNGatewayUserShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing VPN gateway user", err.Error())
	})
}

func Test_ecloudVPNGatewayUserShow(t *testing.T) {
	t.Run("SingleUser", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPNGatewayUser("vpngu-abcdef12").Return(ecloud.VPNGatewayUser{}, nil).Times(1)

		ecloudVPNGatewayUserShow(service, &cobra.Command{}, []string{"vpngu-abcdef12"})
	})

	t.Run("MultipleUsers", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetVPNGatewayUser("vpngu-abcdef12").Return(ecloud.VPNGatewayUser{}, nil),
			service.EXPECT().GetVPNGatewayUser("vpngu-abcdef23").Return(ecloud.VPNGatewayUser{}, nil),
		)

		ecloudVPNGatewayUserShow(service, &cobra.Command{}, []string{"vpngu-abcdef12", "vpngu-abcdef23"})
	})

	t.Run("GetVPNGatewayUserError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPNGatewayUser("vpngu-abcdef12").Return(ecloud.VPNGatewayUser{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving VPN gateway user [vpngu-abcdef12]: test error\n", func() {
			ecloudVPNGatewayUserShow(service, &cobra.Command{}, []string{"vpngu-abcdef12"})
		})
	})
}

func Test_ecloudVPNGatewayUserCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNGatewayUserCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testuser", "--vpngateway=vpng-abcdef12", "--username=user1", "--password=pass123"})

		req := ecloud.CreateVPNGatewayUserRequest{
			Name:         "testuser",
			VPNGatewayID: "vpng-abcdef12",
			Username:     "user1",
			Password:     "pass123",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpngu-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateVPNGatewayUser(req).Return(resp, nil),
			service.EXPECT().GetVPNGatewayUser("vpngu-abcdef12").Return(ecloud.VPNGatewayUser{}, nil),
		)

		ecloudVPNGatewayUserCreate(service, cmd, []string{})
	})

	t.Run("CreateWithWaitFlagNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNGatewayUserCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testuser", "--vpngateway=vpng-abcdef12", "--username=user1", "--password=pass123", "--wait"})

		req := ecloud.CreateVPNGatewayUserRequest{
			Name:         "testuser",
			VPNGatewayID: "vpng-abcdef12",
			Username:     "user1",
			Password:     "pass123",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpngu-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateVPNGatewayUser(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetVPNGatewayUser("vpngu-abcdef12").Return(ecloud.VPNGatewayUser{}, nil),
		)

		ecloudVPNGatewayUserCreate(service, cmd, []string{})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNGatewayUserCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testuser", "--vpngateway=vpng-abcdef12", "--username=user1", "--password=pass123", "--wait"})

		req := ecloud.CreateVPNGatewayUserRequest{
			Name:         "testuser",
			VPNGatewayID: "vpng-abcdef12",
			Username:     "user1",
			Password:     "pass123",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpngu-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateVPNGatewayUser(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		err := ecloudVPNGatewayUserCreate(service, cmd, []string{})
		assert.Equal(t, "Error waiting for VPN gateway user task to complete: Error waiting for command: Failed to retrieve task status: test error", err.Error())
	})

	t.Run("CreateVPNGatewayUserError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNGatewayUserCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testuser", "--vpngateway=vpng-abcdef12", "--username=user1", "--password=pass123"})

		service.EXPECT().CreateVPNGatewayUser(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error")).Times(1)

		err := ecloudVPNGatewayUserCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating VPN gateway user: test error", err.Error())
	})

	t.Run("GetVPNGatewayUserError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNGatewayUserCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testuser", "--vpngateway=vpng-abcdef12", "--username=user1", "--password=pass123"})

		gomock.InOrder(
			service.EXPECT().CreateVPNGatewayUser(gomock.Any()).Return(ecloud.TaskReference{ResourceID: "vpngu-abcdef12"}, nil),
			service.EXPECT().GetVPNGatewayUser("vpngu-abcdef12").Return(ecloud.VPNGatewayUser{}, errors.New("test error")),
		)

		err := ecloudVPNGatewayUserCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new VPN gateway user: test error", err.Error())
	})
}

func Test_ecloudVPNGatewayUserUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPNGatewayUserUpdateCmd(nil).Args(nil, []string{"vpngu-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPNGatewayUserUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing VPN gateway user", err.Error())
	})
}

func Test_ecloudVPNGatewayUserUpdate(t *testing.T) {
	t.Run("SingleUser", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNGatewayUserUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testuser", "--username=user1"})

		req := ecloud.PatchVPNGatewayUserRequest{
			Name: "testuser",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpngu-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchVPNGatewayUser("vpngu-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetVPNGatewayUser("vpngu-abcdef12").Return(ecloud.VPNGatewayUser{}, nil),
		)

		ecloudVPNGatewayUserUpdate(service, cmd, []string{"vpngu-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNGatewayUserUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testuser", "--wait"})

		req := ecloud.PatchVPNGatewayUserRequest{
			Name: "testuser",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpngu-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchVPNGatewayUser("vpngu-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetVPNGatewayUser("vpngu-abcdef12").Return(ecloud.VPNGatewayUser{}, nil),
		)

		ecloudVPNGatewayUserUpdate(service, cmd, []string{"vpngu-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNGatewayUserUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testuser", "--wait"})

		req := ecloud.PatchVPNGatewayUserRequest{
			Name: "testuser",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpngu-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchVPNGatewayUser("vpngu-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for VPN gateway user [vpngu-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudVPNGatewayUserUpdate(service, cmd, []string{"vpngu-abcdef12"})
		})
	})

	t.Run("PatchVPNGatewayUserError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchVPNGatewayUser("vpngu-abcdef12", gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating VPN gateway user [vpngu-abcdef12]: test error\n", func() {
			ecloudVPNGatewayUserUpdate(service, &cobra.Command{}, []string{"vpngu-abcdef12"})
		})
	})

	t.Run("GetVPNGatewayUserError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpngu-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchVPNGatewayUser("vpngu-abcdef12", gomock.Any()).Return(resp, nil),
			service.EXPECT().GetVPNGatewayUser("vpngu-abcdef12").Return(ecloud.VPNGatewayUser{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated VPN gateway user [vpngu-abcdef12]: test error\n", func() {
			ecloudVPNGatewayUserUpdate(service, &cobra.Command{}, []string{"vpngu-abcdef12"})
		})
	})
}

func Test_ecloudVPNGatewayUserDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPNGatewayUserDeleteCmd(nil).Args(nil, []string{"vpngu-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPNGatewayUserDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing VPN gateway user", err.Error())
	})
}

func Test_ecloudVPNGatewayUserDelete(t *testing.T) {
	t.Run("SingleUser", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteVPNGatewayUser("vpngu-abcdef12").Return("task-abcdef12", nil)

		ecloudVPNGatewayUserDelete(service, &cobra.Command{}, []string{"vpngu-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudVPNGatewayUserDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteVPNGatewayUser("vpngu-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudVPNGatewayUserDelete(service, cmd, []string{"vpngu-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNGatewayUserDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		gomock.InOrder(
			service.EXPECT().DeleteVPNGatewayUser("vpngu-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for VPN gateway user [vpngu-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudVPNGatewayUserDelete(service, cmd, []string{"vpngu-abcdef12"})
		})
	})

	t.Run("DeleteVPNGatewayUserError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteVPNGatewayUser("vpngu-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing VPN gateway user [vpngu-abcdef12]: test error\n", func() {
			ecloudVPNGatewayUserDelete(service, &cobra.Command{}, []string{"vpngu-abcdef12"})
		})
	})
}
