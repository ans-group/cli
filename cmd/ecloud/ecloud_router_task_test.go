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

func Test_ecloudRouterTaskListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudRouterTaskListCmd(nil).Args(nil, []string{"rtr-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudRouterTaskListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing router", err.Error())
	})
}

func Test_ecloudRouterTaskList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetRouterTasks("rtr-abcdef12", gomock.Any()).Return([]ecloud.Task{}, nil).Times(1)

		ecloudRouterTaskList(service, &cobra.Command{}, []string{"rtr-abcdef12"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudRouterTaskList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetRoutersError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetRouterTasks("rtr-abcdef12", gomock.Any()).Return([]ecloud.Task{}, errors.New("test error")).Times(1)

		err := ecloudRouterTaskList(service, &cobra.Command{}, []string{"rtr-abcdef12"})

		assert.Equal(t, "Error retrieving router tasks: test error", err.Error())
	})
}

func Test_ecloudRouterTaskWaitCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudRouterTaskWaitCmd(nil).Args(nil, []string{"rtr-abcdef12", "task-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_MissingRouter_Error", func(t *testing.T) {
		err := ecloudRouterTaskWaitCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing router", err.Error())
	})

	t.Run("InvalidArgs_MissingTask_Error", func(t *testing.T) {
		err := ecloudRouterTaskWaitCmd(nil).Args(nil, []string{"rtr-abcdef12"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing task", err.Error())
	})
}

func Test_ecloudRouterTaskWait(t *testing.T) {
	t.Run("ValidFlags_ExpectedCalls", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetRouterTasks("rtr-abcdef12", gomock.Any()).Return([]ecloud.Task{}, nil).Times(1)

		ecloudRouterTaskWait(service, &cobra.Command{}, []string{"rtr-abcdef12", "task-abcdef12"})
	})

	t.Run("InvalidStatus_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudRouterTaskWaitCmd(nil)
		cmd.ParseFlags([]string{"--status=invalid"})

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudRouterTaskWait(service, cmd, []string{"rtr-abcdef12", "task-abcdef12"})

		assert.Equal(t, "Failed to parse status: Invalid ecloud.TaskStatus. Valid values: complete, failed, in-progress", err.Error())
	})

	t.Run("GetRouterTasksError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetRouterTasks("rtr-abcdef12", gomock.Any()).Return([]ecloud.Task{}, errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error waiting for router task [task-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudRouterTaskWait(service, &cobra.Command{}, []string{"rtr-abcdef12", "task-abcdef12"})
		})
	})
}
