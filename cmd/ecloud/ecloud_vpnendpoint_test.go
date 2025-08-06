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

func Test_ecloudVPNEndpointList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		endpoint := mocks.NewMockECloudService(mockCtrl)

		endpoint.EXPECT().GetVPNEndpoints(gomock.Any()).Return([]ecloud.VPNEndpoint{}, nil).Times(1)

		ecloudVPNEndpointList(endpoint, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		endpoint := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudVPNEndpointList(endpoint, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetVPNEndpointsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		endpoint := mocks.NewMockECloudService(mockCtrl)

		endpoint.EXPECT().GetVPNEndpoints(gomock.Any()).Return([]ecloud.VPNEndpoint{}, errors.New("test error")).Times(1)

		err := ecloudVPNEndpointList(endpoint, &cobra.Command{}, []string{})

		assert.Equal(t, "error retrieving VPN endpoints: test error", err.Error())
	})
}

func Test_ecloudVPNEndpointShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPNEndpointShowCmd(nil).Args(nil, []string{"vpne-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPNEndpointShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing VPN endpoint", err.Error())
	})
}

func Test_ecloudVPNEndpointShow(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		endpoint := mocks.NewMockECloudService(mockCtrl)

		endpoint.EXPECT().GetVPNEndpoint("vpne-abcdef12").Return(ecloud.VPNEndpoint{}, nil).Times(1)

		ecloudVPNEndpointShow(endpoint, &cobra.Command{}, []string{"vpne-abcdef12"})
	})

	t.Run("MultipleVPNEndpoints", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		endpoint := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			endpoint.EXPECT().GetVPNEndpoint("vpne-abcdef12").Return(ecloud.VPNEndpoint{}, nil),
			endpoint.EXPECT().GetVPNEndpoint("vpne-abcdef23").Return(ecloud.VPNEndpoint{}, nil),
		)

		ecloudVPNEndpointShow(endpoint, &cobra.Command{}, []string{"vpne-abcdef12", "vpne-abcdef23"})
	})

	t.Run("GetVPNEndpointError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		endpoint := mocks.NewMockECloudService(mockCtrl)

		endpoint.EXPECT().GetVPNEndpoint("vpne-abcdef12").Return(ecloud.VPNEndpoint{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving VPN endpoint [vpne-abcdef12]: test error\n", func() {
			ecloudVPNEndpointShow(endpoint, &cobra.Command{}, []string{"vpne-abcdef12"})
		})
	})
}

func Test_ecloudVPNEndpointCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		endpoint := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNEndpointCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		req := ecloud.CreateVPNEndpointRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpne-abcdef12",
		}

		gomock.InOrder(
			endpoint.EXPECT().CreateVPNEndpoint(req).Return(resp, nil),
			endpoint.EXPECT().GetVPNEndpoint("vpne-abcdef12").Return(ecloud.VPNEndpoint{}, nil),
		)

		ecloudVPNEndpointCreate(endpoint, cmd, []string{})
	})

	t.Run("CreateWithWaitFlagNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		endpoint := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNEndpointCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.CreateVPNEndpointRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpne-abcdef12",
		}

		gomock.InOrder(
			endpoint.EXPECT().CreateVPNEndpoint(req).Return(resp, nil),
			endpoint.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			endpoint.EXPECT().GetVPNEndpoint("vpne-abcdef12").Return(ecloud.VPNEndpoint{}, nil),
		)

		ecloudVPNEndpointCreate(endpoint, cmd, []string{})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		endpoint := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNEndpointCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.CreateVPNEndpointRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpne-abcdef12",
		}

		gomock.InOrder(
			endpoint.EXPECT().CreateVPNEndpoint(req).Return(resp, nil),
			endpoint.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		err := ecloudVPNEndpointCreate(endpoint, cmd, []string{})
		assert.Equal(t, "error waiting for VPN endpoint task to complete: error waiting for command: failed to retrieve task status: test error", err.Error())
	})

	t.Run("CreateVPNEndpointError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		endpoint := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNEndpointCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		endpoint.EXPECT().CreateVPNEndpoint(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error")).Times(1)

		err := ecloudVPNEndpointCreate(endpoint, cmd, []string{})

		assert.Equal(t, "error creating VPN endpoint: test error", err.Error())
	})

	t.Run("GetVPNEndpointError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		endpoint := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNEndpointCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		gomock.InOrder(
			endpoint.EXPECT().CreateVPNEndpoint(gomock.Any()).Return(ecloud.TaskReference{ResourceID: "vpne-abcdef12"}, nil),
			endpoint.EXPECT().GetVPNEndpoint("vpne-abcdef12").Return(ecloud.VPNEndpoint{}, errors.New("test error")),
		)

		err := ecloudVPNEndpointCreate(endpoint, cmd, []string{})

		assert.Equal(t, "error retrieving new VPN endpoint: test error", err.Error())
	})
}

func Test_ecloudVPNEndpointUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPNEndpointUpdateCmd(nil).Args(nil, []string{"vpne-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPNEndpointUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing VPN endpoint", err.Error())
	})
}

func Test_ecloudVPNEndpointUpdate(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		endpoint := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNEndpointUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		req := ecloud.PatchVPNEndpointRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpne-abcdef12",
		}

		gomock.InOrder(
			endpoint.EXPECT().PatchVPNEndpoint("vpne-abcdef12", req).Return(resp, nil),
			endpoint.EXPECT().GetVPNEndpoint("vpne-abcdef12").Return(ecloud.VPNEndpoint{}, nil),
		)

		ecloudVPNEndpointUpdate(endpoint, cmd, []string{"vpne-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		endpoint := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNEndpointUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.PatchVPNEndpointRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpne-abcdef12",
		}

		gomock.InOrder(
			endpoint.EXPECT().PatchVPNEndpoint("vpne-abcdef12", req).Return(resp, nil),
			endpoint.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			endpoint.EXPECT().GetVPNEndpoint("vpne-abcdef12").Return(ecloud.VPNEndpoint{}, nil),
		)

		ecloudVPNEndpointUpdate(endpoint, cmd, []string{"vpne-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		endpoint := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNEndpointUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.PatchVPNEndpointRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpne-abcdef12",
		}

		gomock.InOrder(
			endpoint.EXPECT().PatchVPNEndpoint("vpne-abcdef12", req).Return(resp, nil),
			endpoint.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for VPN endpoint [vpne-abcdef12]: error waiting for command: failed to retrieve task status: test error\n", func() {
			ecloudVPNEndpointUpdate(endpoint, cmd, []string{"vpne-abcdef12"})
		})
	})

	t.Run("PatchVPNEndpointError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		endpoint := mocks.NewMockECloudService(mockCtrl)

		endpoint.EXPECT().PatchVPNEndpoint("vpne-abcdef12", gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating VPN endpoint [vpne-abcdef12]: test error\n", func() {
			ecloudVPNEndpointUpdate(endpoint, &cobra.Command{}, []string{"vpne-abcdef12"})
		})
	})

	t.Run("GetVPNEndpointError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		endpoint := mocks.NewMockECloudService(mockCtrl)

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpne-abcdef12",
		}

		gomock.InOrder(
			endpoint.EXPECT().PatchVPNEndpoint("vpne-abcdef12", gomock.Any()).Return(resp, nil),
			endpoint.EXPECT().GetVPNEndpoint("vpne-abcdef12").Return(ecloud.VPNEndpoint{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated VPN endpoint [vpne-abcdef12]: test error\n", func() {
			ecloudVPNEndpointUpdate(endpoint, &cobra.Command{}, []string{"vpne-abcdef12"})
		})
	})
}

func Test_ecloudVPNEndpointDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPNEndpointDeleteCmd(nil).Args(nil, []string{"vpne-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPNEndpointDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing VPN endpoint", err.Error())
	})
}

func Test_ecloudVPNEndpointDelete(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		endpoint := mocks.NewMockECloudService(mockCtrl)

		endpoint.EXPECT().DeleteVPNEndpoint("vpne-abcdef12").Return("task-abcdef12", nil)

		ecloudVPNEndpointDelete(endpoint, &cobra.Command{}, []string{"vpne-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudVPNEndpointDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		endpoint := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			endpoint.EXPECT().DeleteVPNEndpoint("vpne-abcdef12").Return("task-abcdef12", nil),
			endpoint.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudVPNEndpointDelete(endpoint, cmd, []string{"vpne-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		endpoint := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNEndpointDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		gomock.InOrder(
			endpoint.EXPECT().DeleteVPNEndpoint("vpne-abcdef12").Return("task-abcdef12", nil),
			endpoint.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for VPN endpoint [vpne-abcdef12]: error waiting for command: failed to retrieve task status: test error\n", func() {
			ecloudVPNEndpointDelete(endpoint, cmd, []string{"vpne-abcdef12"})
		})
	})

	t.Run("DeleteVPNEndpointError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		endpoint := mocks.NewMockECloudService(mockCtrl)

		endpoint.EXPECT().DeleteVPNEndpoint("vpne-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing VPN endpoint [vpne-abcdef12]: test error\n", func() {
			ecloudVPNEndpointDelete(endpoint, &cobra.Command{}, []string{"vpne-abcdef12"})
		})
	})
}
