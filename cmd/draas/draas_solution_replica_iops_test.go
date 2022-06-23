package draas

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/sdk-go/pkg/service/draas"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_draasSolutionReplicaIOPSUpdateCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := draasSolutionReplicaIOPSUpdateCmd(nil).Args(nil, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := draasSolutionReplicaIOPSUpdateCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing solution", err.Error())
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := draasSolutionReplicaIOPSUpdateCmd(nil).Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.NotNil(t, err)
		assert.Equal(t, "Missing replica", err.Error())
	})
}

func Test_draasSolutionReplicaIOPSUpdate(t *testing.T) {
	t.Run("DefaultUpdate", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)
		cmd := draasSolutionReplicaIOPSUpdateCmd(nil)
		cmd.Flags().Set("iops-tier", "testtier")

		req := draas.UpdateReplicaIOPSRequest{
			IOPSTierID: "testtier",
		}

		service.EXPECT().UpdateSolutionReplicaIOPS("00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001", req).Return(nil)

		draasSolutionReplicaIOPSUpdate(service, cmd, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
	})

	t.Run("UpdateSolutionReplicaIOPSError_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)
		cmd := draasSolutionReplicaIOPSUpdateCmd(nil)

		service.EXPECT().UpdateSolutionReplicaIOPS("00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001", gomock.Any()).Return(errors.New("test error")).Times(1)

		err := draasSolutionReplicaIOPSUpdate(service, cmd, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})

		assert.Equal(t, "Error updating replica [00000000-0000-0000-0000-000000000001]: test error", err.Error())
	})
}
