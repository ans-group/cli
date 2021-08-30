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

func Test_ecloudVPNServiceList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworkPolicies(gomock.Any()).Return([]ecloud.VPNService{}, nil).Times(1)

		ecloudVPNServiceList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudVPNServiceList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetNetworkPoliciesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworkPolicies(gomock.Any()).Return([]ecloud.VPNService{}, errors.New("test error")).Times(1)

		err := ecloudVPNServiceList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving network policies: test error", err.Error())
	})
}

func Test_ecloudVPNServiceShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPNServiceShowCmd(nil).Args(nil, []string{"vpns-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPNServiceShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing network policy", err.Error())
	})
}

func Test_ecloudVPNServiceShow(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPNService("vpns-abcdef12").Return(ecloud.VPNService{}, nil).Times(1)

		ecloudVPNServiceShow(service, &cobra.Command{}, []string{"vpns-abcdef12"})
	})

	t.Run("MultipleNetworkPolicies", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetVPNService("vpns-abcdef12").Return(ecloud.VPNService{}, nil),
			service.EXPECT().GetVPNService("vpns-abcdef23").Return(ecloud.VPNService{}, nil),
		)

		ecloudVPNServiceShow(service, &cobra.Command{}, []string{"vpns-abcdef12", "vpns-abcdef23"})
	})

	t.Run("GetVPNServiceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPNService("vpns-abcdef12").Return(ecloud.VPNService{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving network policy [vpns-abcdef12]: test error\n", func() {
			ecloudVPNServiceShow(service, &cobra.Command{}, []string{"vpns-abcdef12"})
		})
	})
}

func Test_ecloudVPNServiceCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNServiceCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		req := ecloud.CreateVPNServiceRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpns-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateVPNService(req).Return(resp, nil),
			service.EXPECT().GetVPNService("vpns-abcdef12").Return(ecloud.VPNService{}, nil),
		)

		ecloudVPNServiceCreate(service, cmd, []string{})
	})

	t.Run("CreateWithWaitFlagNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNServiceCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.CreateVPNServiceRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpns-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateVPNService(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetVPNService("vpns-abcdef12").Return(ecloud.VPNService{}, nil),
		)

		ecloudVPNServiceCreate(service, cmd, []string{})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNServiceCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.CreateVPNServiceRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpns-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateVPNService(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		err := ecloudVPNServiceCreate(service, cmd, []string{})
		assert.Equal(t, "Error waiting for network policy task to complete: Error waiting for command: Failed to retrieve task status: test error", err.Error())
	})

	t.Run("CreateVPNServiceError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNServiceCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		service.EXPECT().CreateVPNService(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error")).Times(1)

		err := ecloudVPNServiceCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating network policy: test error", err.Error())
	})

	t.Run("GetVPNServiceError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudVPNServiceCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		gomock.InOrder(
			service.EXPECT().CreateVPNService(gomock.Any()).Return(ecloud.TaskReference{ResourceID: "vpns-abcdef12"}, nil),
			service.EXPECT().GetVPNService("vpns-abcdef12").Return(ecloud.VPNService{}, errors.New("test error")),
		)

		err := ecloudVPNServiceCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new network policy: test error", err.Error())
	})
}

func Test_ecloudVPNServiceUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPNServiceUpdateCmd(nil).Args(nil, []string{"vpns-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPNServiceUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing network policy", err.Error())
	})
}

func Test_ecloudVPNServiceUpdate(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNServiceUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy"})

		req := ecloud.PatchVPNServiceRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpns-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchVPNService("vpns-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetVPNService("vpns-abcdef12").Return(ecloud.VPNService{}, nil),
		)

		ecloudVPNServiceUpdate(service, cmd, []string{"vpns-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNServiceUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.PatchVPNServiceRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpns-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchVPNService("vpns-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetVPNService("vpns-abcdef12").Return(ecloud.VPNService{}, nil),
		)

		ecloudVPNServiceUpdate(service, cmd, []string{"vpns-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNServiceUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testpolicy", "--wait"})

		req := ecloud.PatchVPNServiceRequest{
			Name: "testpolicy",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpns-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchVPNService("vpns-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for network policy [vpns-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudVPNServiceUpdate(service, cmd, []string{"vpns-abcdef12"})
		})
	})

	t.Run("PatchVPNServiceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchVPNService("vpns-abcdef12", gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating network policy [vpns-abcdef12]: test error\n", func() {
			ecloudVPNServiceUpdate(service, &cobra.Command{}, []string{"vpns-abcdef12"})
		})
	})

	t.Run("GetVPNServiceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "vpns-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchVPNService("vpns-abcdef12", gomock.Any()).Return(resp, nil),
			service.EXPECT().GetVPNService("vpns-abcdef12").Return(ecloud.VPNService{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated network policy [vpns-abcdef12]: test error\n", func() {
			ecloudVPNServiceUpdate(service, &cobra.Command{}, []string{"vpns-abcdef12"})
		})
	})
}

func Test_ecloudVPNServiceDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPNServiceDeleteCmd(nil).Args(nil, []string{"vpns-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPNServiceDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing network policy", err.Error())
	})
}

func Test_ecloudVPNServiceDelete(t *testing.T) {
	t.Run("SinglePolicy", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteVPNService("vpns-abcdef12").Return("task-abcdef12", nil)

		ecloudVPNServiceDelete(service, &cobra.Command{}, []string{"vpns-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudVPNServiceDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteVPNService("vpns-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudVPNServiceDelete(service, cmd, []string{"vpns-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudVPNServiceDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		gomock.InOrder(
			service.EXPECT().DeleteVPNService("vpns-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for network policy [vpns-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudVPNServiceDelete(service, cmd, []string{"vpns-abcdef12"})
		})
	})

	t.Run("DeleteVPNServiceError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteVPNService("vpns-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing network policy [vpns-abcdef12]: test error\n", func() {
			ecloudVPNServiceDelete(service, &cobra.Command{}, []string{"vpns-abcdef12"})
		})
	})
}
