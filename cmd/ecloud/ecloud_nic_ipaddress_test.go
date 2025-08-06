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

func Test_ecloudNICIPAddressListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudNICIPAddressListCmd(nil).Args(nil, []string{"ip-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudNICIPAddressListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing NIC", err.Error())
	})
}

func Test_ecloudNICIPAddressList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNICIPAddresses("ip-abcdef12", gomock.Any()).Return([]ecloud.IPAddress{}, nil).Times(1)

		ecloudNICIPAddressList(service, &cobra.Command{}, []string{"ip-abcdef12"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudNICIPAddressList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetNICIPAddressesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNICIPAddresses("ip-abcdef12", gomock.Any()).Return([]ecloud.IPAddress{}, errors.New("test error")).Times(1)

		err := ecloudNICIPAddressList(service, &cobra.Command{}, []string{"ip-abcdef12"})

		assert.Equal(t, "error retrieving NIC IP addresses: test error", err.Error())
	})
}

func Test_ecloudNICIPAddressAssignCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudNICIPAddressAssignCmd(nil).Args(nil, []string{"ip-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudNICIPAddressAssignCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing NIC", err.Error())
	})
}

func Test_ecloudNICIPAddressAssign(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNICIPAddressAssignCmd(nil)
		cmd.Flags().Set("ip-address", "ip-abcdef12")

		req := ecloud.AssignIPAddressRequest{
			IPAddressID: "ip-abcdef12",
		}

		session.EXPECT().AssignNICIPAddress("nic-abcdef12", req).Return("task-abcdef12", nil)

		ecloudNICIPAddressAssign(session, cmd, []string{"nic-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNICIPAddressAssignCmd(nil)
		cmd.Flags().Set("ip-address", "ip-abcdef12")
		cmd.Flags().Set("wait", "true")

		req := ecloud.AssignIPAddressRequest{
			IPAddressID: "ip-abcdef12",
		}

		gomock.InOrder(
			session.EXPECT().AssignNICIPAddress("nic-abcdef12", req).Return("task-abcdef12", nil),
			session.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudNICIPAddressAssign(session, cmd, []string{"nic-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNICIPAddressAssignCmd(nil)
		cmd.Flags().Set("ip-address", "ip-abcdef12")
		cmd.Flags().Set("wait", "true")

		req := ecloud.AssignIPAddressRequest{
			IPAddressID: "ip-abcdef12",
		}

		gomock.InOrder(
			session.EXPECT().AssignNICIPAddress("nic-abcdef12", req).Return("task-abcdef12", nil),
			session.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for NIC [nic-abcdef12]: error waiting for command: failed to retrieve task status: test error\n", func() {
			ecloudNICIPAddressAssign(session, cmd, []string{"nic-abcdef12"})
		})
	})

	t.Run("AssignNICIPAddressError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		session.EXPECT().AssignNICIPAddress("nic-abcdef12", gomock.Any()).Return("", errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error assigning IP address to NIC [nic-abcdef12]: test error\n", func() {
			ecloudNICIPAddressAssign(session, &cobra.Command{}, []string{"nic-abcdef12"})
		})
	})
}

func Test_ecloudNICIPAddressUnassignCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudNICIPAddressUnassignCmd(nil).Args(nil, []string{"ip-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudNICIPAddressUnassignCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing NIC", err.Error())
	})
}

func Test_ecloudNICIPAddressUnassign(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNICIPAddressUnassignCmd(nil)
		cmd.Flags().Set("ip-address", "ip-abcdef12")

		session.EXPECT().UnassignNICIPAddress("nic-abcdef12", "ip-abcdef12").Return("task-abcdef12", nil)

		ecloudNICIPAddressUnassign(session, cmd, []string{"nic-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNICIPAddressUnassignCmd(nil)
		cmd.Flags().Set("ip-address", "ip-abcdef12")
		cmd.Flags().Set("wait", "true")

		gomock.InOrder(
			session.EXPECT().UnassignNICIPAddress("nic-abcdef12", "ip-abcdef12").Return("task-abcdef12", nil),
			session.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudNICIPAddressUnassign(session, cmd, []string{"nic-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudNICIPAddressUnassignCmd(nil)
		cmd.Flags().Set("ip-address", "ip-abcdef12")
		cmd.Flags().Set("wait", "true")

		gomock.InOrder(
			session.EXPECT().UnassignNICIPAddress("nic-abcdef12", "ip-abcdef12").Return("task-abcdef12", nil),
			session.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for NIC [nic-abcdef12]: error waiting for command: failed to retrieve task status: test error\n", func() {
			ecloudNICIPAddressUnassign(session, cmd, []string{"nic-abcdef12"})
		})
	})

	t.Run("UnassignNICIPAddressError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		session := mocks.NewMockECloudService(mockCtrl)

		session.EXPECT().UnassignNICIPAddress("nic-abcdef12", gomock.Any()).Return("", errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error unassigning IP address from NIC [nic-abcdef12]: test error\n", func() {
			ecloudNICIPAddressUnassign(session, &cobra.Command{}, []string{"nic-abcdef12"})
		})
	})
}
