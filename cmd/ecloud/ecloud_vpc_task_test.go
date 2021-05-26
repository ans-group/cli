package ecloud

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/clierrors"
	"github.com/ukfast/cli/test/mocks"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func Test_ecloudVPCTaskListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPCTaskListCmd(nil).Args(nil, []string{"vpc-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudVPCTaskListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing VPC", err.Error())
	})
}

func Test_ecloudVPCTaskList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPCTasks("vpc-abcdef12", gomock.Any()).Return([]ecloud.Task{}, nil).Times(1)

		ecloudVPCTaskList(service, &cobra.Command{}, []string{"vpc-abcdef12"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudVPCTaskList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetVPCsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPCTasks("vpc-abcdef12", gomock.Any()).Return([]ecloud.Task{}, errors.New("test error")).Times(1)

		err := ecloudVPCTaskList(service, &cobra.Command{}, []string{"vpc-abcdef12"})

		assert.Equal(t, "Error retrieving VPC tasks: test error", err.Error())
	})
}

func Test_ecloudVPCTaskWaitCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudVPCTaskWaitCmd(nil).Args(nil, []string{"vpc-abcdef12", "task-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_MissingVPC_Error", func(t *testing.T) {
		err := ecloudVPCTaskWaitCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing VPC", err.Error())
	})

	t.Run("InvalidArgs_MissingTask_Error", func(t *testing.T) {
		err := ecloudVPCTaskWaitCmd(nil).Args(nil, []string{"vpc-abcdef12"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing task", err.Error())
	})
}

func Test_ecloudVPCTaskWait(t *testing.T) {
	t.Run("ValidFlags_ExpectedCalls", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPCTasks("vpc-abcdef12", gomock.Any()).Return([]ecloud.Task{}, nil).Times(1)

		ecloudVPCTaskWait(service, &cobra.Command{}, []string{"vpc-abcdef12", "task-abcdef12"})
	})

	t.Run("InvalidStatus_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudVPCTaskWaitCmd(nil)
		cmd.ParseFlags([]string{"--status=invalid"})

		service := mocks.NewMockECloudService(mockCtrl)

		err := ecloudVPCTaskWait(service, cmd, []string{"vpc-abcdef12", "task-abcdef12"})

		assert.Equal(t, "Failed to parse status: Invalid ecloud.TaskStatus. Valid values: complete, failed, in-progress", err.Error())
	})

	t.Run("GetVPCTasksError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetVPCTasks("vpc-abcdef12", gomock.Any()).Return([]ecloud.Task{}, errors.New("test error")).Times(1)

		err := ecloudVPCTaskWait(service, &cobra.Command{}, []string{"vpc-abcdef12", "task-abcdef12"})

		assert.Equal(t, "Error waiting for VPC task: Error waiting for command: Failed to retrieve task status: test error", err.Error())
	})
}
