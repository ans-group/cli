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

func Test_ecloudIPAddressList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		session.EXPECT().GetIPAddresses(gomock.Any()).Return([]ecloud.IPAddress{}, nil).Times(1)

		ecloudIPAddressList(session, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudIPAddressList(session, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetIPAddressesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		session.EXPECT().GetIPAddresses(gomock.Any()).Return([]ecloud.IPAddress{}, errors.New("test error")).Times(1)

		err := ecloudIPAddressList(session, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving IP addresses: test error", err.Error())
	})
}

func Test_ecloudIPAddressShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudIPAddressShowCmd(nil).Args(nil, []string{"ip-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudIPAddressShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing IP address", err.Error())
	})
}

func Test_ecloudIPAddressShow(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		session.EXPECT().GetIPAddress("ip-abcdef12").Return(ecloud.IPAddress{}, nil).Times(1)

		ecloudIPAddressShow(session, &cobra.Command{}, []string{"ip-abcdef12"})
	})

	t.Run("MultipleIPAddresses", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			session.EXPECT().GetIPAddress("ip-abcdef12").Return(ecloud.IPAddress{}, nil),
			session.EXPECT().GetIPAddress("ip-abcdef23").Return(ecloud.IPAddress{}, nil),
		)

		ecloudIPAddressShow(session, &cobra.Command{}, []string{"ip-abcdef12", "ip-abcdef23"})
	})

	t.Run("GetIPAddressError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		session.EXPECT().GetIPAddress("ip-abcdef12").Return(ecloud.IPAddress{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving IP address [ip-abcdef12]: test error\n", func() {
			ecloudIPAddressShow(session, &cobra.Command{}, []string{"ip-abcdef12"})
		})
	})
}

func Test_ecloudIPAddressCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudIPAddressCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		req := ecloud.CreateIPAddressRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "ip-abcdef12",
		}

		gomock.InOrder(
			session.EXPECT().CreateIPAddress(req).Return(resp, nil),
			session.EXPECT().GetIPAddress("ip-abcdef12").Return(ecloud.IPAddress{}, nil),
		)

		ecloudIPAddressCreate(session, cmd, []string{})
	})

	t.Run("CreateWithWaitFlagNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudIPAddressCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.CreateIPAddressRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "ip-abcdef12",
		}

		gomock.InOrder(
			session.EXPECT().CreateIPAddress(req).Return(resp, nil),
			session.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			session.EXPECT().GetIPAddress("ip-abcdef12").Return(ecloud.IPAddress{}, nil),
		)

		ecloudIPAddressCreate(session, cmd, []string{})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudIPAddressCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.CreateIPAddressRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "ip-abcdef12",
		}

		gomock.InOrder(
			session.EXPECT().CreateIPAddress(req).Return(resp, nil),
			session.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		err := ecloudIPAddressCreate(session, cmd, []string{})
		assert.Equal(t, "Error waiting for IP address task to complete: Error waiting for command: Failed to retrieve task status: test error", err.Error())
	})

	t.Run("CreateIPAddressError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudIPAddressCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		session.EXPECT().CreateIPAddress(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error")).Times(1)

		err := ecloudIPAddressCreate(session, cmd, []string{})

		assert.Equal(t, "Error creating IP address: test error", err.Error())
	})

	t.Run("GetIPAddressError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudIPAddressCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		gomock.InOrder(
			session.EXPECT().CreateIPAddress(gomock.Any()).Return(ecloud.TaskReference{ResourceID: "ip-abcdef12"}, nil),
			session.EXPECT().GetIPAddress("ip-abcdef12").Return(ecloud.IPAddress{}, errors.New("test error")),
		)

		err := ecloudIPAddressCreate(session, cmd, []string{})

		assert.Equal(t, "Error retrieving new IP address: test error", err.Error())
	})
}

func Test_ecloudIPAddressUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudIPAddressUpdateCmd(nil).Args(nil, []string{"ip-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudIPAddressUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing IP address", err.Error())
	})
}

func Test_ecloudIPAddressUpdate(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudIPAddressUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		req := ecloud.PatchIPAddressRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "ip-abcdef12",
		}

		gomock.InOrder(
			session.EXPECT().PatchIPAddress("ip-abcdef12", req).Return(resp, nil),
			session.EXPECT().GetIPAddress("ip-abcdef12").Return(ecloud.IPAddress{}, nil),
		)

		ecloudIPAddressUpdate(session, cmd, []string{"ip-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudIPAddressUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.PatchIPAddressRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "ip-abcdef12",
		}

		gomock.InOrder(
			session.EXPECT().PatchIPAddress("ip-abcdef12", req).Return(resp, nil),
			session.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			session.EXPECT().GetIPAddress("ip-abcdef12").Return(ecloud.IPAddress{}, nil),
		)

		ecloudIPAddressUpdate(session, cmd, []string{"ip-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudIPAddressUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.PatchIPAddressRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "ip-abcdef12",
		}

		gomock.InOrder(
			session.EXPECT().PatchIPAddress("ip-abcdef12", req).Return(resp, nil),
			session.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for IP address [ip-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudIPAddressUpdate(session, cmd, []string{"ip-abcdef12"})
		})
	})

	t.Run("PatchIPAddressError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		session.EXPECT().PatchIPAddress("ip-abcdef12", gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating IP address [ip-abcdef12]: test error\n", func() {
			ecloudIPAddressUpdate(session, &cobra.Command{}, []string{"ip-abcdef12"})
		})
	})

	t.Run("GetIPAddressError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "ip-abcdef12",
		}

		gomock.InOrder(
			session.EXPECT().PatchIPAddress("ip-abcdef12", gomock.Any()).Return(resp, nil),
			session.EXPECT().GetIPAddress("ip-abcdef12").Return(ecloud.IPAddress{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated IP address [ip-abcdef12]: test error\n", func() {
			ecloudIPAddressUpdate(session, &cobra.Command{}, []string{"ip-abcdef12"})
		})
	})
}

func Test_ecloudIPAddressDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudIPAddressDeleteCmd(nil).Args(nil, []string{"ip-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudIPAddressDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing IP address", err.Error())
	})
}

func Test_ecloudIPAddressDelete(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		session.EXPECT().DeleteIPAddress("ip-abcdef12").Return("task-abcdef12", nil)

		ecloudIPAddressDelete(session, &cobra.Command{}, []string{"ip-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudIPAddressDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		session := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			session.EXPECT().DeleteIPAddress("ip-abcdef12").Return("task-abcdef12", nil),
			session.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudIPAddressDelete(session, cmd, []string{"ip-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudIPAddressDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		gomock.InOrder(
			session.EXPECT().DeleteIPAddress("ip-abcdef12").Return("task-abcdef12", nil),
			session.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for IP address [ip-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudIPAddressDelete(session, cmd, []string{"ip-abcdef12"})
		})
	})

	t.Run("DeleteIPAddressError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		session.EXPECT().DeleteIPAddress("ip-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing IP address [ip-abcdef12]: test error\n", func() {
			ecloudIPAddressDelete(session, &cobra.Command{}, []string{"ip-abcdef12"})
		})
	})
}
