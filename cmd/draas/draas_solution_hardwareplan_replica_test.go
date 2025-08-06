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

func Test_draasSolutionHardwarePlanReplicaListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := draasSolutionHardwarePlanReplicaListCmd(nil)
		err := cmd.Args(nil, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})

		assert.Nil(t, err)
	})

	t.Run("MissingSolution_Error", func(t *testing.T) {
		cmd := draasSolutionHardwarePlanReplicaListCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing solution", err.Error())
	})

	t.Run("MissingHardwarePlan_Error", func(t *testing.T) {
		cmd := draasSolutionHardwarePlanReplicaListCmd(nil)
		err := cmd.Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.NotNil(t, err)
		assert.Equal(t, "missing hardware plan", err.Error())
	})
}

func Test_draasSolutionHardwarePlanReplicaList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetSolutionHardwarePlanReplicas("00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001", gomock.Any()).Return([]draas.Replica{}, nil).Times(1)

		draasSolutionHardwarePlanReplicaList(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := draasSolutionHardwarePlanReplicaList(service, cmd, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetSolutionHardwarePlansError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetSolutionHardwarePlanReplicas("00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001", gomock.Any()).Return([]draas.Replica{}, errors.New("test error")).Times(1)

		err := draasSolutionHardwarePlanReplicaList(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})

		assert.Equal(t, "error retrieving solution hardware plan replicas: test error", err.Error())
	})
}
