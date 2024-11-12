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

func Test_ecloudBackupGatewaySpecificationList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetBackupGatewaySpecifications(gomock.Any()).Return([]ecloud.BackupGatewaySpecification{}, nil).Times(1)

		ecloudBackupGatewaySpecificationList(service, &cobra.Command{}, []string{})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := ecloudBackupGatewaySpecificationList(service, cmd, []string{})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetBackupGatewaySpecificationsError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetBackupGatewaySpecifications(gomock.Any()).Return([]ecloud.BackupGatewaySpecification{}, errors.New("test error")).Times(1)

		err := ecloudBackupGatewaySpecificationList(service, &cobra.Command{}, []string{})

		assert.Equal(t, "Error retrieving backup gateway specifications: test error", err.Error())
	})
}

func Test_ecloudBackupGatewaySpecificationShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := ecloudBackupGatewaySpecificationShowCmd(nil).Args(nil, []string{"bgws-abcdef12"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := ecloudBackupGatewaySpecificationShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing backup gateway specification", err.Error())
	})
}

func Test_ecloudBackupGatewaySpecificationShow(t *testing.T) {
	t.Run("SingleSpec", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetBackupGatewaySpecification("bgws-abcdef12").Return(ecloud.BackupGatewaySpecification{}, nil).Times(1)

		ecloudBackupGatewaySpecificationShow(service, &cobra.Command{}, []string{"bgws-abcdef12"})
	})

	t.Run("MultipleSpecs", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		gomock.InOrder(
			service.EXPECT().GetBackupGatewaySpecification("bgws-abcdef12").Return(ecloud.BackupGatewaySpecification{}, nil),
			service.EXPECT().GetBackupGatewaySpecification("bgws-abcdef23").Return(ecloud.BackupGatewaySpecification{}, nil),
		)

		ecloudBackupGatewaySpecificationShow(service, &cobra.Command{}, []string{"bgws-abcdef12", "bgws-abcdef23"})
	})

	t.Run("GetBackupGatewaySpecificationError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetBackupGatewaySpecification("bgws-abcdef12").Return(ecloud.BackupGatewaySpecification{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving backup gateway specification [bgws-abcdef12]: test error\n", func() {
			ecloudBackupGatewaySpecificationShow(service, &cobra.Command{}, []string{"bgws-abcdef12"})
		})
	})

	t.Run("NotFound_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockECloudService(mockCtrl)

		service.EXPECT().GetBackupGatewaySpecification("bgws-abcdef12").Return(ecloud.BackupGatewaySpecification{}, &ecloud.BackupGatewaySpecificationNotFoundError{ID: "bgws-abcdef12"})

		test_output.AssertErrorOutput(t, "Error retrieving backup gateway specification [bgws-abcdef12]: Backup gateway specification not found with ID [bgws-abcdef12]\n", func() {
			ecloudBackupGatewaySpecificationShow(service, &cobra.Command{}, []string{"bgws-abcdef12"})
		})
	})
}
