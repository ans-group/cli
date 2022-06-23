package ecloud

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ecloudTaskList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetTasksPaginated(gomock.Any()).Return(connection.NewPaginated(&connection.APIResponseBodyData[[]ecloud.Task]{}, connection.APIRequestParameters{}, nil), nil).Times(1)

		ecloudTaskList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudTaskList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetTasksError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetTasksPaginated(gomock.Any()).Return(connection.NewPaginated(&connection.APIResponseBodyData[[]ecloud.Task]{}, connection.APIRequestParameters{}, nil), errors.New("test error")).Times(1)

		err := ecloudTaskList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving tasks: test error", err.Error())
	})
}

func Test_ecloudTaskShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudTaskShowCmd(nil).Args(nil, []string{"task-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudTaskShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing task", err.Error())
	})
}

func Test_ecloudTaskShow(t *testing.T) {
	t.Run("SingleTask", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{}, nil).Times(1)

		ecloudTaskShow(service, &cobra.Command{}, []string{"task-abcdef12"})
	})

	t.Run("MultipleTasks", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{}, nil),
			service.EXPECT().GetTask("task-abcdef23").Return(ecloud.Task{}, nil),
		)

		ecloudTaskShow(service, &cobra.Command{}, []string{"task-abcdef12", "task-abcdef23"})
	})

	t.Run("GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving task [task-abcdef12]: test error\n", func() {
			ecloudTaskShow(service, &cobra.Command{}, []string{"task-abcdef12"})
		})
	})
}

func Test_ecloudTaskWaitCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudTaskWaitCmd(nil).Args(nil, []string{"task-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudTaskWaitCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing task", err.Error())
	})
}

func Test_ecloudTaskWait(t *testing.T) {
	t.Run("ValidFlags_ExpectedCalls", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil)

		ecloudTaskWait(service, &cobra.Command{}, []string{"task-abcdef12"})
	})

	t.Run("InvalidStatus_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudTaskWaitCmd(nil)
		cmd.ParseFlags([]string{"--status=invalid"})

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudTaskWait(service, cmd, []string{"task-abcdef12"})

		assert.Equal(t, "Failed to parse status: Invalid ecloud.TaskStatus. Valid values: complete, failed, in-progress", err.Error())
	})

	t.Run("GetTaskError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{}, errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error waiting for task [task-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudTaskWait(service, &cobra.Command{}, []string{"task-abcdef12"})
		})
	})
}
