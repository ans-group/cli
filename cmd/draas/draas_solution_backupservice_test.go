package draas

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/sdk-go/pkg/service/draas"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_draasSolutionBackupServiceShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := draasSolutionBackupServiceShowCmd(nil).Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := draasSolutionBackupServiceShowCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing solution", err.Error())
	})
}

func Test_draasSolutionBackupServiceShow(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetSolutionBackupService("00000000-0000-0000-0000-000000000000").Return(draas.BackupService{}, nil).Times(1)

		draasSolutionBackupServiceShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("GetSolutionBackupService_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetSolutionBackupService("00000000-0000-0000-0000-000000000000").Return(draas.BackupService{}, errors.New("test error")).Times(1)

		err := draasSolutionBackupServiceShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
		assert.Equal(t, "Error retrieving solution backup service: test error", err.Error())
	})
}

func Test_draasSolutionBackupServiceResetCredentialsCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := draasSolutionBackupServiceResetCredentialsCmd(nil).Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := draasSolutionBackupServiceResetCredentialsCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing solution", err.Error())
	})
}

func Test_draasSolutionBackupServiceResetCredentials(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		cmd := &cobra.Command{}
		cmd.Flags().String("password", "testpassword", "")

		req := draas.ResetBackupServiceCredentialsRequest{
			Password: "testpassword",
		}

		service.EXPECT().ResetSolutionBackupServiceCredentials("00000000-0000-0000-0000-000000000000", req).Return(nil).Times(1)

		draasSolutionBackupServiceResetCredentials(service, cmd, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("ResetSolutionBackupServiceCredentialsError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().ResetSolutionBackupServiceCredentials("00000000-0000-0000-0000-000000000000", gomock.Any()).Return(errors.New("test error")).Times(1)

		err := draasSolutionBackupServiceResetCredentials(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
		assert.Equal(t, "Error resetting credentials for solution backup service: test error", err.Error())
	})
}
