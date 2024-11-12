package ecloud

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_ecloudBackupGatewayList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetBackupGateways(gomock.Any()).Return([]ecloud.BackupGateway{}, nil).Times(1)

		ecloudBackupGatewayList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudBackupGatewayList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetBackupGatewaysError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetBackupGateways(gomock.Any()).Return([]ecloud.BackupGateway{}, errors.New("test error")).Times(1)

		err := ecloudBackupGatewayList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving backup gateways: test error", err.Error())
	})
}

func Test_ecloudBackupGatewayShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudBackupGatewayShowCmd(nil).Args(nil, []string{"bgw-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudBackupGatewayShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing backup gateway ID", err.Error())
	})
}

func Test_ecloudBackupGatewayShow(t *testing.T) {
	t.Run("SingleBackupGateway", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetBackupGateway("bgw-abcdef12").Return(ecloud.BackupGateway{}, nil).Times(1)

		ecloudBackupGatewayShow(service, &cobra.Command{}, []string{"bgw-abcdef12"})
	})

	t.Run("MultipleBackupGateways", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetBackupGateway("bgw-abcdef12").Return(ecloud.BackupGateway{}, nil),
			service.EXPECT().GetBackupGateway("bgw-abcdef23").Return(ecloud.BackupGateway{}, nil),
		)

		ecloudBackupGatewayShow(service, &cobra.Command{}, []string{"bgw-abcdef12", "bgw-abcdef23"})
	})

	t.Run("GetBackupGatewayError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetBackupGateway("bgw-abcdef12").Return(ecloud.BackupGateway{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving backup gateway [bgw-abcdef12]: test error\n", func() {
			ecloudBackupGatewayShow(service, &cobra.Command{}, []string{"bgw-abcdef12"})
		})
	})
}

func Test_ecloudBackupGatewayCreate(t *testing.T) {
	t.Run("DefaultCreate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudBackupGatewayCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgateway", "--vpc=vpc-abcdef12", "--router=rtr-abcdef12", "--specification=bgws-abcdef12"})

		req := ecloud.CreateBackupGatewayRequest{
			Name:          "testgateway",
			VPCID:         "vpc-abcdef12",
			RouterID:      "rtr-abcdef12",
			GatewaySpecID: "bgws-abcdef12",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "bgw-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateBackupGateway(req).Return(resp, nil),
			service.EXPECT().GetBackupGateway("bgw-abcdef12").Return(ecloud.BackupGateway{}, nil),
		)

		ecloudBackupGatewayCreate(service, cmd, []string{})
	})

	t.Run("CreateWithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudBackupGatewayCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgateway", "--vpc=vpc-abcdef12", "--router=rtr-abcdef12", "--specification=bgws-abcdef12", "--wait"})

		req := ecloud.CreateBackupGatewayRequest{
			Name:          "testgateway",
			VPCID:         "vpc-abcdef12",
			RouterID:      "rtr-abcdef12",
			GatewaySpecID: "bgws-abcdef12",
		}

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "bgw-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateBackupGateway(req).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
			service.EXPECT().GetBackupGateway("bgw-abcdef12").Return(ecloud.BackupGateway{}, nil),
		)

		ecloudBackupGatewayCreate(service, cmd, []string{})
	})

	t.Run("WithWaitFlag_TaskError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudBackupGatewayCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgateway", "--vpc=vpc-abcdef12", "--router=rtr-abcdef12", "--specification=bgws-abcdef12", "--wait"})

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "bgw-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateBackupGateway(gomock.Any()).Return(resp, nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{}, errors.New("test error")),
		)

		err := ecloudBackupGatewayCreate(service, cmd, []string{})

		assert.Equal(t, "Error waiting for backup gateway task to complete: Error waiting for command: Failed to retrieve task status: test error", err.Error())
	})

	t.Run("CreateError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudBackupGatewayCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgateway"})

		service.EXPECT().CreateBackupGateway(gomock.Any()).Return(ecloud.TaskReference{}, errors.New("test error"))

		err := ecloudBackupGatewayCreate(service, cmd, []string{})

		assert.Equal(t, "Error creating backup gateway: test error", err.Error())
	})

	t.Run("GetBackupGatewayError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := ecloudBackupGatewayCreateCmd(nil)
		cmd.ParseFlags([]string{"--name=testgateway", "--vpc=vpc-abcdef12", "--router=rtr-abcdef12", "--specification=bgws-abcdef12"})

		resp := ecloud.TaskReference{
			TaskID:     "task-abcdef12",
			ResourceID: "bgw-abcdef12",
		}

		gomock.InOrder(
			service.EXPECT().CreateBackupGateway(gomock.Any()).Return(resp, nil),
			service.EXPECT().GetBackupGateway("bgw-abcdef12").Return(ecloud.BackupGateway{}, errors.New("test error")),
		)

		err := ecloudBackupGatewayCreate(service, cmd, []string{})

		assert.Equal(t, "Error retrieving new backup gateway: test error", err.Error())
	})
}

func Test_ecloudBackupGatewayUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudBackupGatewayUpdateCmd(nil).Args(nil, []string{"bgw-abcdef12"})
		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudBackupGatewayUpdateCmd(nil).Args(nil, []string{})
		assert.NotNil(t, err)
		assert.Equal(t, "Missing backup gateway", err.Error())
	})
}

func Test_ecloudBackupGatewayDelete(t *testing.T) {
	t.Run("SingleBackupGateway", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteBackupGateway("bgw-abcdef12").Return("task-abcdef12", nil)

		ecloudBackupGatewayDelete(service, &cobra.Command{}, []string{"bgw-abcdef12"})
	})

	t.Run("WithWaitFlag_NoError_Succeeds", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudBackupGatewayDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteBackupGateway("bgw-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{Status: ecloud.TaskStatusComplete}, nil),
		)

		ecloudBackupGatewayDelete(service, cmd, []string{"bgw-abcdef12"})
	})

	t.Run("DeleteError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().DeleteBackupGateway("bgw-abcdef12").Return("", errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error removing backup gateway [bgw-abcdef12]: test error\n", func() {
			ecloudBackupGatewayDelete(service, &cobra.Command{}, []string{"bgw-abcdef12"})
		})
	})

	t.Run("WaitFlag_GetTaskError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cmd := ecloudBackupGatewayDeleteCmd(nil)
		cmd.ParseFlags([]string{"--wait"})

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().DeleteBackupGateway("bgw-abcdef12").Return("task-abcdef12", nil),
			service.EXPECT().GetTask("task-abcdef12").Return(ecloud.Task{}, errors.New("test error")),
		)

		test_output.AssertErrorOutput(t, "Error waiting for task to complete for backup gateway [bgw-abcdef12]: Error waiting for command: Failed to retrieve task status: test error\n", func() {
			ecloudBackupGatewayDelete(service, cmd, []string{"bgw-abcdef12"})
		})
	})
}
