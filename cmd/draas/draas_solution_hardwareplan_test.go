package draas

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/internal/pkg/clierrors"
	"github.com/ans-group/cli/test/mocks"
	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/service/draas"
	gomock "github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_draasSolutionHardwarePlanListCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := draasSolutionHardwarePlanListCmd(nil)
		err := cmd.Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Nil(t, err)
	})

	t.Run("MissingSolution_Error", func(t *testing.T) {
		cmd := draasSolutionHardwarePlanListCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing solution", err.Error())
	})
}

func Test_draasSolutionHardwarePlanList(t *testing.T) {
	t.Run("DefaultRetrieve", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetSolutionHardwarePlans("00000000-0000-0000-0000-000000000000", gomock.Any()).Return([]draas.HardwarePlan{}, nil).Times(1)

		draasSolutionHardwarePlanList(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})
	})

	t.Run("MalformedFlag_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{"invalidfilter"}, "")

		err := draasSolutionHardwarePlanList(service, cmd, []string{"00000000-0000-0000-0000-000000000000"})

		assert.IsType(t, &clierrors.ErrInvalidFlagValue{}, err)
	})

	t.Run("GetSolutionHardwarePlansError_ReturnsError", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetSolutionHardwarePlans("00000000-0000-0000-0000-000000000000", gomock.Any()).Return([]draas.HardwarePlan{}, errors.New("test error")).Times(1)

		err := draasSolutionHardwarePlanList(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000"})

		assert.Equal(t, "error retrieving solution hardware plans: test error", err.Error())
	})
}

func Test_draasSolutionHardwarePlanShowCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		cmd := draasSolutionHardwarePlanShowCmd(nil)
		err := cmd.Args(nil, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})

		assert.Nil(t, err)
	})

	t.Run("MissingSolution_Error", func(t *testing.T) {
		cmd := draasSolutionHardwarePlanShowCmd(nil)
		err := cmd.Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing solution", err.Error())
	})

	t.Run("MissingHardwarePlan_Error", func(t *testing.T) {
		cmd := draasSolutionHardwarePlanShowCmd(nil)
		err := cmd.Args(nil, []string{"00000000-0000-0000-0000-000000000000"})

		assert.NotNil(t, err)
		assert.Equal(t, "missing hardware plan", err.Error())
	})
}

func Test_draasSolutionHardwarePlanShow(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetSolutionHardwarePlan("00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001").Return(draas.HardwarePlan{}, nil)

		draasSolutionHardwarePlanShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
	})

	t.Run("GetSolutionHardwarePlan_OutputsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		service := mocks.NewMockDRaaSService(mockCtrl)

		service.EXPECT().GetSolutionHardwarePlan("00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001").Return(draas.HardwarePlan{}, errors.New("test error"))

		test_output.AssertErrorOutput(t, "Error retrieving solution hardware plan [00000000-0000-0000-0000-000000000001]: test error\n", func() {
			draasSolutionHardwarePlanShow(service, &cobra.Command{}, []string{"00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000001"})
		})
	})
}
