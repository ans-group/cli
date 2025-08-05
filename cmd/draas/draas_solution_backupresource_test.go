package draas

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/sdk-go/pkg/service/draas"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_draasSolutionBackupResourceListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := draasSolutionBackupResourceListCmd(nil).Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := draasSolutionBackupResourceListCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing solution", err.Error())
	})
}

func Test_draasSolutionBackupResourceList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetSolutionBackupResources("00000000-0000-0000-0000-000000000000", gomock.Any()).Return([]draas.BackupResource{}, nil).Times(1)

		draasSolutionBackupResourceList(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := draasSolutionBackupResourceList(service, cmd, []string{"00000000-0000-0000-0000-000000000000"})
		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetSolutionBackupResourcesError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetSolutionBackupResources("00000000-0000-0000-0000-000000000000", gomock.Any()).Return([]draas.BackupResource{}, errors.New("test error")).Times(1)

		err := draasSolutionBackupResourceList(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
		assert.Equal(t, "error retrieving solution backup resources: test error", err.Error())
	})
}
