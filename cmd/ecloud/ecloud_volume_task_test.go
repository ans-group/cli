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

func Test_ecloudVolumeTaskListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVolumeTaskListCmd(nil).Args(nil, []string{"vol-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVolumeTaskListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing volume", err.Error())
	})
}

func Test_ecloudVolumeTaskList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVolumeTasks("vol-abcdef12", gomock.Any()).Return([]ecloud.Task{}, nil).Times(1)

		ecloudVolumeTaskList(service, &cobra.Command{}, []string{"vol-abcdef12"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudVolumeTaskList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetVolumesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVolumeTasks("vol-abcdef12", gomock.Any()).Return([]ecloud.Task{}, errors.New("test error")).Times(1)

		err := ecloudVolumeTaskList(service, &cobra.Command{}, []string{"vol-abcdef12"})

		assert.Equal(t, "Error retrieving volume tasks: test error", err.Error())
	})
}

func Test_ecloudVolumeTaskWaitCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVolumeTaskWaitCmd(nil).Args(nil, []string{"vol-abcdef12", "task-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_MissingVolume_Error", func(t *testing.T) {
		err := ecloudVolumeTaskWaitCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing volume", err.Error())
	})

	t.Run("InvalidArgs_MissingTask_Error", func(t *testing.T) {
		err := ecloudVolumeTaskWaitCmd(nil).Args(nil, []string{"vol-abcdef12"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing task", err.Error())
	})
}

func Test_ecloudVolumeTaskWait(t *testing.T) {
	t.Run("ValidFlags_ExpectedCalls", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVolumeTasks("vol-abcdef12", gomock.Any()).Return([]ecloud.Task{}, nil).Times(1)

		ecloudVolumeTaskWait(service, &cobra.Command{}, []string{"vol-abcdef12", "task-abcdef12"})
	})

	t.Run("InvalidStatus_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudVolumeTaskWaitCmd(nil)
		cmd.ParseFlags([]string{"--status=invalid"})

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudVolumeTaskWait(service, cmd, []string{"vol-abcdef12", "task-abcdef12"})

		assert.Equal(t, "Failed to parse status: Invalid ecloud.TaskStatus. Valid values: complete, failed, in-progress", err.Error())
	})

	t.Run("GetVolumeTasksError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVolumeTasks("vol-abcdef12", gomock.Any()).Return([]ecloud.Task{}, errors.New("test error")).Times(1)

		test_output.AssertErrorOutput(t, "Error waiting for volume task [task-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudVolumeTaskWait(service, &cobra.Command{}, []string{"vol-abcdef12", "task-abcdef12"})
		})
	})
}
