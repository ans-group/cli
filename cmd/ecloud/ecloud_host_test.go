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

func Test_ecloudHostList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetHosts(gomock.Any()).Return([]ecloud.Host{}, nil).Times(1)

		ecloudHostList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudHostList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetHostsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetHosts(gomock.Any()).Return([]ecloud.Host{}, errors.New("test error")).Times(1)

		err := ecloudHostList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "error retrieving hosts: test error", err.Error())
	})
}

func Test_ecloudHostShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudHostShowCmd(nil).Args(nil, []string{"h-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudHostShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing host", err.Error())
	})
}

func Test_ecloudHostShow(t *testing.T) {
	t.Run("SingleHost", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetHost("h-abcdef12").Return(ecloud.Host{}, nil).Times(1)

		ecloudHostShow(service, &cobra.Command{}, []string{"h-abcdef12"})
	})

	t.Run("MultipleHosts", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetHost("h-abcdef12").Return(ecloud.Host{}, nil),
			service.EXPECT().GetHost("h-abcdef23").Return(ecloud.Host{}, nil),
		)

		ecloudHostShow(service, &cobra.Command{}, []string{"h-abcdef12", "h-abcdef23"})
	})

	t.Run("GetHostError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetHost("h-abcdef12").Return(ecloud.Host{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving host [h-abcdef12]: test error\n", func() {
			ecloudHostShow(service, &cobra.Command{}, []string{"h-abcdef12"})
		})
	})
}

func Test_ecloudHostCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudHostCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgroup", "--host-group=hg-abcdef12"})

		req := ecloud.CreateHostRequest{
			Name:        "testgroup",
			HostGroupID: "hg-abcdef12",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "h-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateHost(req).Return(resp, nil),
			service.EXPECT().GetHost("h-abcdef12").Return(ecloud.Host{}, nil),
		)

		ecloudHostCreate(service, cmd, []string{})
	})

	t.Run("CreateWithWaitFlagNoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudHostCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgroup", "--host-group=hg-abcdef12", "--wait"})

		req := ecloud.CreateHostRequest{
			Name:        "testgroup",
			HostGroupID: "hg-abcdef12",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "h-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateHost(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetHost("h-abcdef12").Return(ecloud.Host{}, nil),
		)

		ecloudHostCreate(service, cmd, []string{})
	})

	t.Run("WithWaitFlag_GetTaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudHostCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgroup", "--host-group=hg-abcdef12", "--wait"})

		req := ecloud.CreateHostRequest{
			Name:        "testgroup",
			HostGroupID: "hg-abcdef12",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "h-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateHost(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		err := ecloudHostCreate(service, cmd, []string{})
		assert.Equal(t, "error waiting for host task to complete: error waiting for command: failed to retrieve task status: test error", err.Error())
	})

	t.Run("CreateHostError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudHostCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgroup"})

		service.EXPECT().CreateHost(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error")).Times(1)

		err := ecloudHostCreate(service, cmd, []string{})

		assert.Equal(t, "error creating host: test error", err.Error())
	})

	t.Run("GetHostError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudHostCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgroup"})

		gomock.InOrder(
			service.EXPECT().CreateHost(gomock.Any()).Return(ecloud.TaskReference{ResourceID: "h-abcdef12"}, nil),
			service.EXPECT().GetHost("h-abcdef12").Return(ecloud.Host{}, errors.New("test error")),
		)

		err := ecloudHostCreate(service, cmd, []string{})

		assert.Equal(t, "error retrieving new host: test error", err.Error())
	})
}

func Test_ecloudHostUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudHostUpdateCmd(nil).Args(nil, []string{"h-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudHostUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing host", err.Error())
	})
}

func Test_ecloudHostUpdate(t *testing.T) {
	t.Run("SingleGroup", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudHostUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgroup"})

		req := ecloud.PatchHostRequest{
			Name: "testgroup",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "h-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchHost("h-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetHost("h-abcdef12").Return(ecloud.Host{}, nil),
		)

		ecloudHostUpdate(service, cmd, []string{"h-abcdef12"})
	})

	t.Run("MultipleHosts", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp1 := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "h-abcdef12",
		}

		resp2 := ecloud.TaskReference{
			TaskID:     "task-abcdef23",
			ResourceID: "h-12abcdef",
		}

		gomock.InOrder(
			service.EXPECT().PatchHost("h-abcdef12", gomock.Any()).Return(resp1, nil),
			service.EXPECT().GetHost("h-abcdef12").Return(ecloud.Host{}, nil),
			service.EXPECT().PatchHost("h-12abcdef", gomock.Any()).Return(resp2, nil),
			service.EXPECT().GetHost("h-12abcdef").Return(ecloud.Host{}, nil),
		)

		ecloudHostUpdate(service, &cobra.Command{}, []string{"h-abcdef12", "h-12abcdef"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudHostUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgroup", "--wait"})

		req := ecloud.PatchHostRequest{
			Name: "testgroup",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "h-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchHost("h-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetHost("h-abcdef12").Return(ecloud.Host{}, nil),
		)

		ecloudHostUpdate(service, cmd, []string{"h-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudHostUpdateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgroup", "--wait"})

		req := ecloud.PatchHostRequest{
			Name: "testgroup",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "h-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchHost("h-abcdef12", req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for host [h-abcdef12]: error waiting for command: failed to retrieve task status: test error\n", func() {
			ecloudHostUpdate(service, cmd, []string{"h-abcdef12"})
		})
	})

	t.Run("PatchHostError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().PatchHost("h-abcdef12", gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error updating host [h-abcdef12]: test error\n", func() {
			ecloudHostUpdate(service, &cobra.Command{}, []string{"h-abcdef12"})
		})
	})

	t.Run("GetHostError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "h-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().PatchHost("h-abcdef12", gomock.Any()).Return(resp, nil),
			service.EXPECT().GetHost("h-abcdef12").Return(ecloud.Host{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error retrieving updated host [h-abcdef12]: test error\n", func() {
			ecloudHostUpdate(service, &cobra.Command{}, []string{"h-abcdef12"})
		})
	})
}

func Test_ecloudHostDeleteCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudHostDeleteCmd(nil).Args(nil, []string{"h-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudHostDeleteCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing host", err.Error())
	})
}

func Test_ecloudHostDelete(t *testing.T) {
	t.Run("SingleGroup", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteHost("h-abcdef12").Return("task-abcdef12", nil).Times(1)

		ecloudHostDelete(service, &cobra.Command{}, []string{"h-abcdef12"})
	})

	t.Run("MultipleHosts", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteHost("h-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().DeleteHost("h-12abcdef").Return("task-abcdef23", nil),
		)

		ecloudHostDelete(service, &cobra.Command{}, []string{"h-abcdef12", "h-12abcdef"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudHostDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteHost("h-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudHostDelete(service, cmd, []string{"h-abcdef12"})
	})

	t.Run("WithWaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		cmd := ecloudHostDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		gomock.InOrder(
			service.EXPECT().DeleteHost("h-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for host [h-abcdef12]: error waiting for command: failed to retrieve task status: test error\n", func() {
			ecloudHostDelete(service, cmd, []string{"h-abcdef12"})
		})
	})

	t.Run("DeleteHostError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteHost("h-abcdef12").Return("", errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error removing host [h-abcdef12]: test error\n", func() {
			ecloudHostDelete(service, &cobra.Command{}, []string{"h-abcdef12"})
		})
	})
}
