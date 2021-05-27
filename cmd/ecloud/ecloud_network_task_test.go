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

func Test_ecloudNetworkTaskListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudNetworkTaskListCmd(nil).Args(nil, []string{"net-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudNetworkTaskListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing network", err.Error())
	})
}

func Test_ecloudNetworkTaskList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworkTasks("net-abcdef12", gomock.Any()).Return([]ecloud.Task{}, nil).Times(1)

		ecloudNetworkTaskList(service, &cobra.Command{}, []string{"net-abcdef12"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudNetworkTaskList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetNetworksError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworkTasks("net-abcdef12", gomock.Any()).Return([]ecloud.Task{}, errors.New("test error")).Times(1)

		err := ecloudNetworkTaskList(service, &cobra.Command{}, []string{"net-abcdef12"})

		assert.Equal(t, "Error retrieving network tasks: test error", err.Error())
	})
}

func Test_ecloudNetworkTaskWaitCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudNetworkTaskWaitCmd(nil).Args(nil, []string{"net-abcdef12", "task-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_MissingNetwork_Error", func(t *testing.T) {
		err := ecloudNetworkTaskWaitCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing network", err.Error())
	})

	t.Run("InvalidArgs_MissingTask_Error", func(t *testing.T) {
		err := ecloudNetworkTaskWaitCmd(nil).Args(nil, []string{"net-abcdef12"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing task", err.Error())
	})
}

func Test_ecloudNetworkTaskWait(t *testing.T) {
	t.Run("ValidFlags_ExpectedCalls", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworkTasks("net-abcdef12", gomock.Any()).Return([]ecloud.Task{}, nil).Times(1)

		ecloudNetworkTaskWait(service, &cobra.Command{}, []string{"net-abcdef12", "task-abcdef12"})
	})

	t.Run("InvalidStatus_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudNetworkTaskWaitCmd(nil)
		cmd.ParseFlags([]string{"--status=invalid"})

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudNetworkTaskWait(service, cmd, []string{"net-abcdef12", "task-abcdef12"})

		assert.Equal(t, "Failed to parse status: Invalid ecloud.TaskStatus. Valid values: complete, failed, in-progress", err.Error())
	})

	t.Run("GetNetworkTasksError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetNetworkTasks("net-abcdef12", gomock.Any()).Return([]ecloud.Task{}, errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error waiting for network task [task-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudNetworkTaskWait(service, &cobra.Command{}, []string{"net-abcdef12", "task-abcdef12"})
		})
	})
}
