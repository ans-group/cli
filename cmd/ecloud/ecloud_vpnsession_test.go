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

func Test_ecloudVPNSessionList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		session.EXPECT().GetVPNSessions(gomock.Any()).Return([]ecloud.VPNSession{}, nil).Times(1)

		ecloudVPNSessionList(session, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudVPNSessionList(session, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetVPNSessionsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		session.EXPECT().GetVPNSessions(gomock.Any()).Return([]ecloud.VPNSession{}, errors.New("test error")).Times(1)

		err := ecloudVPNSessionList(session, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving VPN sessions: test error", err.Error())
	})
}

func Test_ecloudVPNSessionShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPNSessionShowCmd(nil).Args(nil, []string{"vpns-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPNSessionShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing VPN session", err.Error())
	})
}

func Test_ecloudVPNSessionShow(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		session.EXPECT().GetVPNSession("vpns-abcdef12").Return(ecloud.VPNSession{}, nil).Times(1)

		ecloudVPNSessionShow(session, &cobra.Command{}, []string{"vpns-abcdef12"})
	})

	t.Run("MultipleVPNSessions", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			session.EXPECT().GetVPNSession("vpns-abcdef12").Return(ecloud.VPNSession{}, nil),
			session.EXPECT().GetVPNSession("vpns-abcdef23").Return(ecloud.VPNSession{}, nil),
		)

		ecloudVPNSessionShow(session, &cobra.Command{}, []string{"vpns-abcdef12", "vpns-abcdef23"})
	})

	t.Run("GetVPNSessionError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		session.EXPECT().GetVPNSession("vpns-abcdef12").Return(ecloud.VPNSession{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving VPN session [vpns-abcdef12]: test error\n", func() {
			ecloudVPNSessionShow(session, &cobra.Command{}, []string{"vpns-abcdef12"})
		})
	})
}

func Test_ecloudVPNSessionCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNSessionCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		req := ecloud.CreateVPNSessionRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpns-abcdef12",
		}

		gomock.InOrder(
			session.EXPECT().CreateVPNSession(req).Return(resp, nil),
			session.EXPECT().GetVPNSession("vpns-abcdef12").Return(ecloud.VPNSession{}, nil),
		)

		ecloudVPNSessionCreate(session, cmd, []string{})
	})

	t.Run("CreateWithWaitFlagNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNSessionCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.CreateVPNSessionRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpns-abcdef12",
		}

		gomock.InOrder(
			session.EXPECT().CreateVPNSession(req).Return(resp, nil),
			session.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			session.EXPECT().GetVPNSession("vpns-abcdef12").Return(ecloud.VPNSession{}, nil),
		)

		ecloudVPNSessionCreate(session, cmd, []string{})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNSessionCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.CreateVPNSessionRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpns-abcdef12",
		}

		gomock.InOrder(
			session.EXPECT().CreateVPNSession(req).Return(resp, nil),
			session.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		err := ecloudVPNSessionCreate(session, cmd, []string{})
		assert.Equal(t, "Error waiting for VPN session task to complete: Error waiting for command: Failed to retrieve task status: test error", err.Error())
	})

	t.Run("CreateVPNSessionError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNSessionCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		session.EXPECT().CreateVPNSession(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error")).Times(1)

		err := ecloudVPNSessionCreate(session, cmd, []string{})

		assert.Equal(t, "Error creating VPN session: test error", err.Error())
	})

	t.Run("GetVPNSessionError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNSessionCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		gomock.InOrder(
			session.EXPECT().CreateVPNSession(gomock.Any()).Return(ecloud.TaskReference{ResourceID: "vpns-abcdef12"}, nil),
			session.EXPECT().GetVPNSession("vpns-abcdef12").Return(ecloud.VPNSession{}, errors.New("test error")),
		)

		err := ecloudVPNSessionCreate(session, cmd, []string{})

		assert.Equal(t, "Error retrieving new VPN session: test error", err.Error())
	})
}

func Test_ecloudVPNSessionUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPNSessionUpdateCmd(nil).Args(nil, []string{"vpns-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPNSessionUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing VPN session", err.Error())
	})
}

func Test_ecloudVPNSessionUpdate(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNSessionUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		req := ecloud.PatchVPNSessionRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpns-abcdef12",
		}

		gomock.InOrder(
			session.EXPECT().PatchVPNSession("vpns-abcdef12", req).Return(resp, nil),
			session.EXPECT().GetVPNSession("vpns-abcdef12").Return(ecloud.VPNSession{}, nil),
		)

		ecloudVPNSessionUpdate(session, cmd, []string{"vpns-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNSessionUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.PatchVPNSessionRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpns-abcdef12",
		}

		gomock.InOrder(
			session.EXPECT().PatchVPNSession("vpns-abcdef12", req).Return(resp, nil),
			session.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			session.EXPECT().GetVPNSession("vpns-abcdef12").Return(ecloud.VPNSession{}, nil),
		)

		ecloudVPNSessionUpdate(session, cmd, []string{"vpns-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNSessionUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.PatchVPNSessionRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpns-abcdef12",
		}

		gomock.InOrder(
			session.EXPECT().PatchVPNSession("vpns-abcdef12", req).Return(resp, nil),
			session.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for VPN session [vpns-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudVPNSessionUpdate(session, cmd, []string{"vpns-abcdef12"})
		})
	})

	t.Run("PatchVPNSessionError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		session.EXPECT().PatchVPNSession("vpns-abcdef12", gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating VPN session [vpns-abcdef12]: test error\n", func() {
			ecloudVPNSessionUpdate(session, &cobra.Command{}, []string{"vpns-abcdef12"})
		})
	})

	t.Run("GetVPNSessionError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpns-abcdef12",
		}

		gomock.InOrder(
			session.EXPECT().PatchVPNSession("vpns-abcdef12", gomock.Any()).Return(resp, nil),
			session.EXPECT().GetVPNSession("vpns-abcdef12").Return(ecloud.VPNSession{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated VPN session [vpns-abcdef12]: test error\n", func() {
			ecloudVPNSessionUpdate(session, &cobra.Command{}, []string{"vpns-abcdef12"})
		})
	})
}

func Test_ecloudVPNSessionDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPNSessionDeleteCmd(nil).Args(nil, []string{"vpns-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPNSessionDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing VPN session", err.Error())
	})
}

func Test_ecloudVPNSessionDelete(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		session.EXPECT().DeleteVPNSession("vpns-abcdef12").Return("task-abcdef12", nil)

		ecloudVPNSessionDelete(session, &cobra.Command{}, []string{"vpns-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudVPNSessionDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		session := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			session.EXPECT().DeleteVPNSession("vpns-abcdef12").Return("task-abcdef12", nil),
			session.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudVPNSessionDelete(session, cmd, []string{"vpns-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNSessionDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		gomock.InOrder(
			session.EXPECT().DeleteVPNSession("vpns-abcdef12").Return("task-abcdef12", nil),
			session.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for VPN session [vpns-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudVPNSessionDelete(session, cmd, []string{"vpns-abcdef12"})
		})
	})

	t.Run("DeleteVPNSessionError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		session.EXPECT().DeleteVPNSession("vpns-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing VPN session [vpns-abcdef12]: test error\n", func() {
			ecloudVPNSessionDelete(session, &cobra.Command{}, []string{"vpns-abcdef12"})
		})
	})
}
